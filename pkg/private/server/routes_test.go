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
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	"github.com/gin-gonic/gin"
)

func TestConfigureRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful IAM Token Set", func(t *testing.T) {
		w := httptest.NewRecorder()
		engine := gin.New()
		injector := iamtoken.NewInjector()

		configureRoute(engine, "/private", injector)

		reqBody := IAMTokenPostRequest{
			ProjectID: "test-project",
			IAMToken:  "test-token",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/private/api/v3/iam_token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
		}

		header := http.Header{}
		injector.ApplyToRawHTTPHeader(header, googlecloud.Project("test-project"))
		gotToken := header.Get("x-goog-iam-authorization-token")

		if gotToken != "test-token" {
			t.Errorf("expected token 'test-token', got %q", gotToken)
		}
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		engine := gin.New()
		injector := iamtoken.NewInjector()

		configureRoute(engine, "/private", injector)

		req, _ := http.NewRequest("POST", "/private/api/v3/iam_token", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		engine.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
