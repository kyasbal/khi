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

package inspection_cached_task

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/khi/pkg/common/khictx"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	inspection_task_contextkey "github.com/GoogleCloudPlatform/khi/pkg/inspection/contextkey"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
	"github.com/GoogleCloudPlatform/khi/pkg/task/taskid"
)

type CachableResult[T any] struct {
	Value            T
	DependencyDigest string
}

func NewCachedTask[T any](taskID taskid.TaskImplementationID[T], depdendencies []taskid.UntypedTaskReference, f func(ctx context.Context, prevValue CachableResult[T]) (CachableResult[T], error), labelOpt ...task.LabelOpt) task.Definition[T] {
	return task.NewTask(taskID, depdendencies, func(ctx context.Context) (T, error) {
		inspectionSharedMap := khictx.MustGetValue(ctx, inspection_task_contextkey.GlobalSharedMap)
		cacheKey := typedmap.NewTypedKey[CachableResult[T]](fmt.Sprintf("cached_result-%s", taskID.String()))
		cachedResult := typedmap.GetOrDefault(inspectionSharedMap, cacheKey, CachableResult[T]{
			Value:            *new(T),
			DependencyDigest: "",
		})

		nextCache, err := f(ctx, cachedResult)
		if err != nil {
			return *new(T), err
		}

		typedmap.Set(inspectionSharedMap, cacheKey, nextCache)
		return nextCache.Value, nil
	}, labelOpt...)
}
