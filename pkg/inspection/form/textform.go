// Copyright 2024 Google LLC
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

package form

import (
	"context"
	"fmt"

	form_metadata "github.com/GoogleCloudPlatform/khi/pkg/inspection/metadata/form"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/task"
	common_task "github.com/GoogleCloudPlatform/khi/pkg/task"
)

// TextFormValidator is a function to check if the given value is valid or not.
// Returns "" as the result when it has no error, otherwise the returned value is used as an error message on frontend.
// Returning an error as the 2nd returning value is only when the validator detects an unrecoverble error.
type TextFormValidator = func(ctx context.Context, value string, variables *common_task.VariableSet) (string, error)

// TextFormDefaultValueGenerator is a function type to generate the default value.
type TextFormDefaultValueGenerator = func(ctx context.Context, variables *common_task.VariableSet, previousValues []string) (string, error)

// TextFormReadonlyProvider is a function type to compute if the field is allowed edit or not.
type TextFormReadonlyProvider = func(ctx context.Context, variables *common_task.VariableSet) (bool, error)

// TextFormSuggestionsProvider is a function to return the list of strings shown in the autocomplete.
// Return nil instead of emptry string array means the autocomplete is disabled for the field.
type TextFormSuggestionsProvider = func(ctx context.Context, value string, variables *common_task.VariableSet, previousValues []string) ([]string, error)

// TextFormValueConverter is a function type to convert the given string value to another type stored in the variable set.
type TextFormValueConverter = func(ctx context.Context, value string, variables *common_task.VariableSet) (any, error)

// TextFormHintGenerator is a function type to generate a hint string
type TextFormHintGenerator = func(ctx context.Context, value string, convertedValue any, variables *common_task.VariableSet) (string, form_metadata.ParameterHintType, error)

// TextFormDefinitionBuilder is an utility to construct an instance of Definition for input form field.
// This will generate the Definition instance with `Build()` method call after chaining several configuration methods.
type TextFormDefinitionBuilder struct {
	id                  string
	label               string
	priority            int
	dependencies        []string
	description         string
	defaultValue        TextFormDefaultValueGenerator
	validator           TextFormValidator
	readonlyProvider    TextFormReadonlyProvider
	suggestionsProvider TextFormSuggestionsProvider
	hintGenerator       TextFormHintGenerator
	converter           TextFormValueConverter
}

// NewInputFormDefinitionBuilder constructs an instace of TextFormDefinitionBuilder.
// id,prioirity and label will be initialized with the value given in the argument. The other values are initialized with the following values.
// dependencies : Initialized with an empty string array indicating this definition is not depending on anything.
// description: Initialized with an empty string.
// defaultValue: Initialized with a function to return empty string.
// validator: Initialized with a function to return empty string that indicates the validation is always passing.
// allowEditProvider: Initialized with a function to return true.
// suggestionsProvider: Initialized with a function to return nil.
// converter: Initialized with a function to return the given value. This means no conversion applied and treated as a string.
func NewInputFormDefinitionBuilder(id string, priority int, fieldLabel string) *TextFormDefinitionBuilder {
	return &TextFormDefinitionBuilder{
		id:           id,
		priority:     priority,
		label:        fieldLabel,
		dependencies: []string{},
		defaultValue: func(ctx context.Context, variables *common_task.VariableSet, previousValues []string) (string, error) {
			return "", nil
		},
		validator: func(ctx context.Context, value string, variables *common_task.VariableSet) (string, error) {
			return "", nil
		},
		readonlyProvider: func(ctx context.Context, variables *common_task.VariableSet) (bool, error) {
			return false, nil
		},
		suggestionsProvider: func(ctx context.Context, value string, variables *common_task.VariableSet, previousValues []string) ([]string, error) {
			return nil, nil
		},
		converter: func(ctx context.Context, value string, variables *common_task.VariableSet) (any, error) {
			return value, nil
		},
		hintGenerator: func(ctx context.Context, value string, convertedValue any, variables *common_task.VariableSet) (string, form_metadata.ParameterHintType, error) {
			return "", form_metadata.Info, nil
		},
	}
}

func (b *TextFormDefinitionBuilder) WithDependencies(dependencies []string) *TextFormDefinitionBuilder {
	b.dependencies = dependencies
	return b
}

func (b *TextFormDefinitionBuilder) WithDescription(description string) *TextFormDefinitionBuilder {
	b.description = description
	return b
}

func (b *TextFormDefinitionBuilder) WithValidator(validator TextFormValidator) *TextFormDefinitionBuilder {
	b.validator = validator
	return b
}

func (b *TextFormDefinitionBuilder) WithDefaultValueFunc(defFunc TextFormDefaultValueGenerator) *TextFormDefinitionBuilder {
	b.defaultValue = defFunc
	return b
}

func (b *TextFormDefinitionBuilder) WithDefaultValueConstant(defValue string, preferPrevValue bool) *TextFormDefinitionBuilder {
	return b.WithDefaultValueFunc(func(ctx context.Context, variables *common_task.VariableSet, previousValues []string) (string, error) {
		if preferPrevValue {
			if len(previousValues) > 0 {
				return previousValues[0], nil
			}
		}
		return defValue, nil
	})
}

func (b *TextFormDefinitionBuilder) WithAllowEditFunc(readonlyFunc TextFormReadonlyProvider) *TextFormDefinitionBuilder {
	b.readonlyProvider = readonlyFunc
	return b
}

func (b *TextFormDefinitionBuilder) WithSuggestionsFunc(suggestionsFunc TextFormSuggestionsProvider) *TextFormDefinitionBuilder {
	b.suggestionsProvider = suggestionsFunc
	return b
}

func (b *TextFormDefinitionBuilder) WithSuggestionsConstant(suggestions []string) *TextFormDefinitionBuilder {
	return b.WithSuggestionsFunc(func(ctx context.Context, value string, variables *common_task.VariableSet, previousValues []string) ([]string, error) {
		return suggestions, nil
	})
}

func (b *TextFormDefinitionBuilder) WithHintFunc(hintFunc TextFormHintGenerator) *TextFormDefinitionBuilder {
	b.hintGenerator = hintFunc
	return b
}

func (b *TextFormDefinitionBuilder) WithConverter(converter TextFormValueConverter) *TextFormDefinitionBuilder {
	b.converter = converter
	return b
}

func (b *TextFormDefinitionBuilder) Build(labelOpts ...common_task.LabelOpt) common_task.Definition {
	return common_task.NewProcessorTask(b.id, b.dependencies, func(ctx context.Context, taskMode int, v *common_task.VariableSet) (any, error) {
		m, err := task.GetMetadataSetFromVariable(v)
		if err != nil {
			return nil, err
		}
		req, err := task.GetInspectionRequestFromVariable(v)
		if err != nil {
			return nil, err
		}
		cacheStore, err := common_task.GetCacheStoreFromTaskVariable(v)
		if err != nil {
			return nil, err
		}
		previousValueStoreKey := fmt.Sprintf("text-form-pv-%s", b.id)
		prevValueAny, _ := cacheStore.LoadOrStore(previousValueStoreKey, []string{})
		prevValue := prevValueAny.([]string)

		readonly, err := b.readonlyProvider(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("allowEdit provider for task `%s` returned an error\n%v", b.id, err)
		}
		field := form_metadata.TextParameterFormField{}
		field.Readonly = readonly

		// Compute the default value of the form
		var currentValue string
		currentValue, err = b.defaultValue(ctx, v, prevValue)
		if err != nil {
			return nil, fmt.Errorf("default value generator for task `%s` returned an error\n%v", b.id, err)
		}
		field.Default = currentValue
		if valueRaw, exist := req.Values[b.id]; exist && !readonly {
			valueString, isString := valueRaw.(string)
			if !isString {
				return nil, fmt.Errorf("request parameter `%s` was not given in string in task %s", b.id, b.id)
			}
			currentValue = valueString
		}

		field.ID = b.id
		field.Type = form_metadata.Text
		field.Priority = b.priority
		field.Label = b.label
		field.Description = b.description
		field.HintType = form_metadata.Info

		suggestions, err := b.suggestionsProvider(ctx, currentValue, v, prevValue)
		if err != nil {
			return nil, fmt.Errorf("suggesion provider for task `%s` returned an error\n%v", b.id, err)
		}
		field.Suggestions = suggestions

		validationErr, err := b.validator(ctx, currentValue, v)
		if err != nil {
			return nil, fmt.Errorf("validator for task `%s` returned an unrecovable error\n%v", b.id, err)
		}
		if validationErr != "" {
			// When the given string is invalid, it should be the default value.
			currentValue, err = b.defaultValue(ctx, v, prevValue)
			if err != nil {
				return nil, fmt.Errorf("default value generator for task `%s` returned an error\n%v", b.id, err)
			}
		}
		if validationErr != "" && taskMode == task.TaskModeRun {
			return nil, fmt.Errorf("validator for task `%s` returned a validation error. But this task was executed as a Run mode not in DryRun. All validations must be resolved before running.\n%v", b.id, validationErr)
		}

		convertedValue, err := b.converter(ctx, currentValue, v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert the value `%s` to the dedicated value in task %s\n%v", currentValue, b.id, err)
		}
		if validationErr != "" {
			field.HintType = form_metadata.Error
			field.Hint = validationErr
		} else {
			hint, hintType, err := b.hintGenerator(ctx, currentValue, convertedValue, v)
			if err != nil {
				return nil, fmt.Errorf("failed to generate a hint for task %s\n%v", b.id, err)
			}
			if hint == "" {
				hintType = form_metadata.None
			}
			field.Hint = hint
			field.HintType = hintType
			if taskMode == task.TaskModeRun {
				newValueHistory := append([]string{currentValue}, prevValue...)
				cacheStore.Store(previousValueStoreKey, newValueHistory)
			}
		}

		formFields := m.LoadOrStore(form_metadata.FormFieldSetMetadataKey, &form_metadata.FormFieldSetMetadataFactory{}).(*form_metadata.FormFieldSet)
		err = formFields.SetField(field)
		if err != nil {
			return nil, fmt.Errorf("failed to configure the form metadata in task `%s`\n%v", b.id, err)
		}
		return convertedValue, nil
	}, labelOpts...)
}
