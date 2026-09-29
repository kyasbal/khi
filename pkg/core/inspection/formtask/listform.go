// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package formtask

import (
	"context"
	"fmt"
	"strconv"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspectionmetadata "github.com/GoogleCloudPlatform/khi/pkg/core/inspection/metadata"
	coretask "github.com/GoogleCloudPlatform/khi/pkg/core/task"
	"github.com/GoogleCloudPlatform/khi/pkg/core/task/taskid"
	"github.com/GoogleCloudPlatform/khi/pkg/task/inspection/inspectioncore"
)

// ListFormItemResolver resolves a child form field and its evaluated value for a single item in the list.
type ListFormItemResolver[T any] = func(ctx context.Context, index int, itemKey string, itemFieldID string) (inspectionmetadata.ParameterFormField, T, error)

// ListFormValidator validates the aggregated slice of resolved item values.
// Returns an empty string when valid, or an error message to display on the frontend.
type ListFormValidator[T any] = func(ctx context.Context, values []T) (string, error)

// ListFormHintGenerator generates a hint message and severity type for the overall list field based on item values.
type ListFormHintGenerator[T any] = func(ctx context.Context, values []T) (string, inspectionmetadata.ParameterHintType, error)

// ListFormTaskBuilder builds a DAG task that renders and processes a repeatable list form field.
type ListFormTaskBuilder[T any] struct {
	FormTaskBuilderBase[[]T]
	itemResolver   ListFormItemResolver[T]
	defaultCount   int
	minCount       int
	maxCount       int
	addButtonLabel string
	validator      ListFormValidator[T]
	hintGenerator  ListFormHintGenerator[T]
}

// NewListFormTaskBuilder constructs a new ListFormTaskBuilder instance.
func NewListFormTaskBuilder[T any](
	id taskid.TaskImplementationID[[]T],
	priority int,
	fieldLabel string,
	itemResolver ListFormItemResolver[T],
) *ListFormTaskBuilder[T] {
	return &ListFormTaskBuilder[T]{
		FormTaskBuilderBase: NewFormTaskBuilderBase(id, priority, fieldLabel),
		itemResolver:        itemResolver,
		defaultCount:        0,
		minCount:            0,
		maxCount:            0,
		addButtonLabel:      "Add",
		validator: func(ctx context.Context, values []T) (string, error) {
			return "", nil
		},
		hintGenerator: func(ctx context.Context, values []T) (string, inspectionmetadata.ParameterHintType, error) {
			return "", inspectionmetadata.Info, nil
		},
	}
}

// WithDependencies sets the task dependencies.
func (b *ListFormTaskBuilder[T]) WithDependencies(dependencies []coretask.Dependency) *ListFormTaskBuilder[T] {
	b.FormTaskBuilderBase.WithDependencies(dependencies)
	return b
}

// WithDescription sets the description for the list form field.
func (b *ListFormTaskBuilder[T]) WithDescription(description string) *ListFormTaskBuilder[T] {
	b.FormTaskBuilderBase.WithDescription(description)
	return b
}

// WithDefaultCount sets a static default item count for the list.
func (b *ListFormTaskBuilder[T]) WithDefaultCount(count int) *ListFormTaskBuilder[T] {
	b.defaultCount = count
	return b
}

// WithMinCount sets a static minimum item count constraint.
func (b *ListFormTaskBuilder[T]) WithMinCount(count int) *ListFormTaskBuilder[T] {
	b.minCount = count
	return b
}

// WithMaxCount sets a static maximum item count constraint. Zero means unlimited.
func (b *ListFormTaskBuilder[T]) WithMaxCount(count int) *ListFormTaskBuilder[T] {
	b.maxCount = count
	return b
}

// WithAddButtonLabel sets the label text displayed on the button to add a new item.
func (b *ListFormTaskBuilder[T]) WithAddButtonLabel(label string) *ListFormTaskBuilder[T] {
	b.addButtonLabel = label
	return b
}

// WithValidator sets the validator function for the resolved item values.
func (b *ListFormTaskBuilder[T]) WithValidator(validator ListFormValidator[T]) *ListFormTaskBuilder[T] {
	b.validator = validator
	return b
}

// WithHintFunc sets the hint generator function for the list form field.
func (b *ListFormTaskBuilder[T]) WithHintFunc(hintFunc ListFormHintGenerator[T]) *ListFormTaskBuilder[T] {
	b.hintGenerator = hintFunc
	return b
}

// Build constructs the coretask.Task that manages list form evaluation.
func (b *ListFormTaskBuilder[T]) Build(labelOpts ...coretask.LabelOpt) coretask.Task[[]T] {
	return coretask.NewTask(b.FormTaskBuilderBase.id, b.FormTaskBuilderBase.dependencies, func(ctx context.Context) ([]T, error) {
		metadata := khictx.MustGetValue(ctx, inspectionmetadata.MapContextKey)
		req := khictx.MustGetValue(ctx, inspectioncore.InspectionTaskInput)
		taskMode := khictx.MustGetValue(ctx, inspectioncore.InspectionTaskMode)
		taskID := b.FormTaskBuilderBase.id.ReferenceIDString()

		if err := validateCountConstraints(b.defaultCount, b.minCount, b.maxCount, taskID); err != nil {
			return nil, err
		}

		defaultKeys := make([]string, b.defaultCount)
		for i := 0; i < b.defaultCount; i++ {
			defaultKeys[i] = strconv.Itoa(i)
		}

		currentKeys, err := parseAndValidateItemKeys(req[taskID], defaultKeys, taskID)
		if err != nil {
			return nil, err
		}

		items, resolvedValues, itemErr := resolveListItems(ctx, currentKeys, taskID, b.itemResolver)
		if itemErr != nil && items == nil {
			return nil, itemErr
		}

		var validationErr string
		if itemErr == nil {
			validationErr, err = b.validateItems(ctx, currentKeys, resolvedValues, taskID)
			if err != nil {
				return nil, err
			}
		}

		field, err := b.buildFormField(ctx, items, defaultKeys, resolvedValues, validationErr, itemErr != nil)
		if err != nil {
			return nil, err
		}

		formFields, found := typedmap.Get(metadata, inspectionmetadata.FormFieldSetMetadataKey)
		if !found {
			return nil, fmt.Errorf("failed to get form fields from metadata in task %s", taskID)
		}
		if err := formFields.SetField(field); err != nil {
			return nil, fmt.Errorf("failed to set list form field in task %s: %w", taskID, err)
		}

		if itemErr != nil {
			return nil, itemErr
		}
		if validationErr != "" && taskMode == inspectioncore.TaskModeRun {
			return nil, fmt.Errorf("validator for task %s returned a validation error in Run mode: %s", taskID, validationErr)
		}

		return resolvedValues, nil
	}, append(labelOpts, inspectioncore.NewFormTaskLabelOpt(b.label, b.description))...)
}

// validateCountConstraints checks whether the item count bounds are valid.
func validateCountConstraints(defaultCount, minCount, maxCount int, taskID string) error {
	if defaultCount < 0 {
		return fmt.Errorf("default count must be non-negative, got %d in task %s", defaultCount, taskID)
	}
	if minCount < 0 {
		return fmt.Errorf("min count must be non-negative, got %d in task %s", minCount, taskID)
	}
	if maxCount < 0 {
		return fmt.Errorf("max count must be non-negative, got %d in task %s", maxCount, taskID)
	}
	if maxCount > 0 && minCount > maxCount {
		return fmt.Errorf("min count (%d) cannot be greater than max count (%d) in task %s", minCount, maxCount, taskID)
	}
	if defaultCount < minCount {
		return fmt.Errorf("default count (%d) cannot be less than min count (%d) in task %s", defaultCount, minCount, taskID)
	}
	if maxCount > 0 && defaultCount > maxCount {
		return fmt.Errorf("default count (%d) cannot be greater than max count (%d) in task %s", defaultCount, maxCount, taskID)
	}
	return nil
}

// parseAndValidateItemKeys extracts and validates list item keys from request input.
func parseAndValidateItemKeys(valueRaw any, defaultKeys []string, taskID string) ([]string, error) {
	if valueRaw == nil {
		return defaultKeys, nil
	}
	var parsedKeys []string
	switch v := valueRaw.(type) {
	case []string:
		parsedKeys = make([]string, len(v))
		copy(parsedKeys, v)
	case []any:
		parsedKeys = make([]string, len(v))
		for i, elem := range v {
			strElem, ok := elem.(string)
			if !ok {
				return nil, fmt.Errorf("item key at index %d in task %s is not a string (got %T)", i, taskID, elem)
			}
			parsedKeys[i] = strElem
		}
	default:
		return nil, fmt.Errorf("expected []string or []any for task %s, got %T", taskID, valueRaw)
	}

	seenKeys := make(map[string]bool, len(parsedKeys))
	for i, key := range parsedKeys {
		if key == "" {
			return nil, fmt.Errorf("empty item key found at index %d in task %s", i, taskID)
		}
		if seenKeys[key] {
			return nil, fmt.Errorf("duplicate item key %q found at index %d in task %s", key, i, taskID)
		}
		seenKeys[key] = true
	}
	return parsedKeys, nil
}

// resolveListItems resolves form fields and values for each item in the list.
func resolveListItems[T any](
	ctx context.Context,
	keys []string,
	taskID string,
	resolver ListFormItemResolver[T],
) ([]inspectionmetadata.ListParameterFormFieldItem, []T, error) {
	items := make([]inspectionmetadata.ListParameterFormFieldItem, 0, len(keys))
	resolvedValues := make([]T, 0, len(keys))
	var firstItemErr error

	for i, key := range keys {
		itemFieldID := fmt.Sprintf("%s/%s", taskID, key)
		itemField, itemVal, itemErr := resolver(ctx, i, key, itemFieldID)
		if itemErr != nil && itemField == nil {
			return nil, nil, fmt.Errorf("failed to resolve item %q in task %s: %w", key, taskID, itemErr)
		}
		items = append(items, inspectionmetadata.ListParameterFormFieldItem{
			Key:   key,
			Field: itemField,
		})
		resolvedValues = append(resolvedValues, itemVal)
		if itemErr != nil && firstItemErr == nil {
			firstItemErr = itemErr
		}
	}

	return items, resolvedValues, firstItemErr
}

// validateItems validates the item count constraints and delegates to the custom validator.
func (b *ListFormTaskBuilder[T]) validateItems(ctx context.Context, keys []string, values []T, taskID string) (string, error) {
	switch {
	case len(keys) < b.minCount:
		return fmt.Sprintf("at least %d item(s) are required", b.minCount), nil
	case b.maxCount > 0 && len(keys) > b.maxCount:
		return fmt.Sprintf("at most %d item(s) are allowed", b.maxCount), nil
	default:
		validationErr, err := b.validator(ctx, values)
		if err != nil {
			return "", fmt.Errorf("validator in task %s failed: %w", taskID, err)
		}
		return validationErr, nil
	}
}

// buildFormField constructs the metadata parameter form field with validation or hints.
func (b *ListFormTaskBuilder[T]) buildFormField(
	ctx context.Context,
	items []inspectionmetadata.ListParameterFormFieldItem,
	defaultKeys []string,
	values []T,
	validationErr string,
	hasItemErr bool,
) (inspectionmetadata.ListParameterFormField, error) {
	field := inspectionmetadata.ListParameterFormField{
		ParameterFormFieldBase: inspectionmetadata.ParameterFormFieldBase{
			Type:     inspectionmetadata.List,
			HintType: inspectionmetadata.None,
			Hint:     "",
		},
		Items:          items,
		Default:        defaultKeys,
		MinCount:       b.minCount,
		MaxCount:       b.maxCount,
		AddButtonLabel: b.addButtonLabel,
	}
	b.FormTaskBuilderBase.SetupBaseFormField(&field.ParameterFormFieldBase)

	if validationErr != "" {
		field.HintType = inspectionmetadata.Error
		field.Hint = validationErr
		return field, nil
	}

	if hasItemErr {
		return field, nil
	}

	hint, hintType, err := b.hintGenerator(ctx, values)
	if err != nil {
		return field, fmt.Errorf("hint generator in task %s failed: %w", b.FormTaskBuilderBase.id, err)
	}
	if hint == "" {
		hintType = inspectionmetadata.None
	}
	field.Hint = hint
	field.HintType = hintType
	return field, nil
}
