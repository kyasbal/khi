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

package server

import (
	"net/http"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	"github.com/gin-gonic/gin"
)

// IAMTokenPostRequest is the payload JSON type for POST /api/v3/iam_token
type IAMTokenPostRequest struct {
	ProjectID string `json:"projectID"`
	IAMToken  string `json:"iamToken"`
}

// ConfigureRoute configures the given gin.IRouter instance to handle private endpoints.
func ConfigureRoute(router gin.IRouter, iamTokenInjector *iamtoken.IAMTokenCallOptionInjectorOption) {
	router.POST("/api/v3/iam_token", func(ctx *gin.Context) {
		var reqBody IAMTokenPostRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		iamTokenInjector.SetTokenFor(googlecloud.Project(reqBody.ProjectID), reqBody.IAMToken)
		ctx.String(http.StatusOK, "")
	})
}

// configureRoute configures the given gin.Engine instance to handle private endpoints.
func configureRoute(engine *gin.Engine, basePath string, iamTokenInjector *iamtoken.IAMTokenCallOptionInjectorOption) {
	ConfigureRoute(engine.Group(basePath), iamTokenInjector)
}
