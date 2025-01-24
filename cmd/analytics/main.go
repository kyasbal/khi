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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/khi/pkg/private/analytics/types"
	"github.com/gin-gonic/gin"
)

// Analytics server deployed on CloudRun in `kubernetes-history-inspector` GCP project.
// Send received metadata to BigQuery

var PROJECT_ID = "kubernetes-history-inspector"
var DATASET_ID = "usage"
var TABLE_ID = "khi-usage"

func main() {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, PROJECT_ID)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to initialize bigquery client\n%s", err))
		os.Exit(1)
		return
	}
	defer client.Close()

	router := gin.Default()
	router.POST("/", func(ctx *gin.Context) {
		var request = types.RecordAnalyticsDataRequest{}
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("failed to bind the json request\n%s", err))
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		metadataInJson, err := json.Marshal(request.Metadata)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("failed to marshal the metadata field of the request\n%s", err))
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		inserter := client.Dataset(DATASET_ID).Table(TABLE_ID).Inserter()
		item := types.RecordAnalyticsData{
			Time:     time.Now(),
			Debug:    request.Debug,
			Event:    request.Event,
			Metadata: string(metadataInJson),
		}
		err = inserter.Put(ctx, item)
		if err != nil {
			slog.WarnContext(ctx, fmt.Sprintf("failed to insert bigquery record\n%s", err))
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.String(http.StatusOK, "OK")
	})
	router.Run()
}
