// Copyright 2025 Google LLC
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

package task_test

import (
	"context"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
	task_contextkey "github.com/GoogleCloudPlatform/khi/pkg/task/contextkey"
	"github.com/GoogleCloudPlatform/khi/pkg/task/taskid"
)

type TaskDependencyValuePair interface {
	Register(resultMap *typedmap.TypedMap)
}

type taskDependencyValuePair[T any] struct {
	Value T
	Key   taskid.TaskReference[T]
}

// Register implements TaskDependencyValuePair.
func (t *taskDependencyValuePair[T]) Register(resultMap *typedmap.TypedMap) {
	typedmap.Set(resultMap, typedmap.NewTypedKey[T](t.Key.ReferenceIDString()), t.Value)
}

var _ TaskDependencyValuePair = (*taskDependencyValuePair[any])(nil)

func NewTaskDependencyValuePair[T any](key taskid.TaskReference[T], value T) TaskDependencyValuePair {
	return &taskDependencyValuePair[T]{
		Value: value,
		Key:   key,
	}
}

func RunTask[T any](baseContext context.Context, task task.Definition[T], taskDependencyValues ...TaskDependencyValuePair) (T, error) {
	taskCtx := khictx.WithValue(baseContext, task_contextkey.TaskImplementationIDContextKey, task.UntypedID())

	resultMap := typedmap.NewTypedMap()
	for _, taskDependencyValue := range taskDependencyValues {
		taskDependencyValue.Register(resultMap)
	}
	taskCtx = khictx.WithValue(taskCtx, task_contextkey.TaskResultMapContextKey, resultMap)
	return task.Run(taskCtx)
}
