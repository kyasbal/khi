package server

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/metadata"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/popup"
	common_task "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/task"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func redirectMiddleware(exactPath string, redirectTo string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == exactPath {
			ctx.Redirect(302, redirectTo)
			return
		}
		ctx.Next()
	}
}

func CreateKHIServer(inspectionServer *inspection.InspectionTaskServer, viewerMode bool, staticFolderPath string, resourceMonitor ResourceMonitor) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	appHtmlPath := path.Join(staticFolderPath, "/index.html")
	indexHtml := generateIndexHtmlWithGALabels(appHtmlPath)
	engine.Use(redirectMiddleware("/", "/session/0")) // Request for `/` shouldn't be handled by `static.Serve`, redirect `/session/0` to be handled by patternToString
	engine.Use(static.Serve("/", static.LocalFile(staticFolderPath, false)))
	engine.Use(gin.Recovery())
	engine.Use(cors.New(corsConfig))

	// frontend uses Angular router. All frontend routing path should return the app html
	engine.GET("/session/*wild", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html")
		ctx.Writer.Write([]byte(indexHtml))
	})
	if !viewerMode {
		// GET /api/v2/inspection/types
		// Returns the list of inspection types available on the inspection server.
		engine.GET("/api/v2/inspection/types", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, &GetInspectionTypesResponse{
				Types: inspectionServer.GetAllInspectionTypes(),
			})
		})

		// GET /api/v2/inspection/tasks
		// Returns the all started inspections on the inspection server.
		engine.GET("/api/v2/inspection/tasks", func(ctx *gin.Context) {
			inspections := inspectionServer.GetAllRunners()
			responseInspections := map[string]SerializedMetadata{}
			for _, inspection := range inspections {
				if inspection.Started() {
					md, err := inspection.GetCurrentMetadata()
					if err != nil {
						ctx.String(http.StatusInternalServerError, err.Error())
						return
					}
					m, err := md.ToMap(common_task.EqualLabelFilter(metadata.LabelKeyIncludedInTaskListFlag, true, false))
					if err != nil {
						ctx.String(http.StatusInternalServerError, err.Error())
						return
					}
					responseInspections[inspection.ID] = m
				}
			}

			ctx.JSON(http.StatusOK, &GetInspectionTasksResponse{
				Tasks: responseInspections,
				ServerStat: &ServerStat{
					TotalMemoryAvailable: resourceMonitor.GetUsedMemory(),
				},
			})
		})

		// POST /api/v2/inspection/tasks
		engine.POST("/api/v2/inspection/types/:typeId", func(ctx *gin.Context) {
			typeId := ctx.Param("typeId")
			inspectionId, err := inspectionServer.CreateInspection(typeId)
			if err != nil {
				// only the not found error is expected here
				ctx.String(http.StatusNotFound, err.Error())
				return
			}
			ctx.JSON(http.StatusAccepted, &PostInspectionTaskResponse{InspectionId: inspectionId})
		})
		// PUT /api/v2/inspection/tasks/<task-id>/features
		engine.PUT("/api/v2/inspection/tasks/:taskId/features", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			task := inspectionServer.GetTask(taskId)
			if task == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			var reqBody PutInspectionTaskFeatureRequest
			if err := ctx.ShouldBindJSON(&reqBody); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			err := task.SetFeatureList(reqBody.Features)
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.String(http.StatusAccepted, "ok")
		})
		//GET /api/v2/inspection/tasks/<task-id>/features
		engine.GET("/api/v2/inspection/tasks/:taskId/features", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			task := inspectionServer.GetTask(taskId)
			if task == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			features, err := task.FeatureList()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, GetInspectionTaskFeatureResponse{
				Features: features,
			})
		})

		engine.POST("/api/v2/inspection/tasks/:taskId/dryrun", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			currentTask := inspectionServer.GetTask(taskId)
			if currentTask == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			var reqBody PostInspectionTaskDryRunRequest
			if err := ctx.ShouldBindJSON(&reqBody); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			result, err := currentTask.DryRun(ctx, &task.InspectionRequest{
				Values: reqBody,
			})
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, result)
		})

		engine.POST("/api/v2/inspection/tasks/:taskId/run", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			currentTask := inspectionServer.GetTask(taskId)
			if currentTask == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			var reqBody PostInspectionTaskDryRunRequest
			if err := ctx.ShouldBindJSON(&reqBody); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			err := currentTask.Run(ctx, &task.InspectionRequest{
				Values: reqBody,
			})
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.String(http.StatusAccepted, "ok")
		})

		engine.POST("/api/v2/inspection/tasks/:taskId/cancel", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			currentTask := inspectionServer.GetTask(taskId)
			if currentTask == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			err := currentTask.Cancel()
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.String(http.StatusOK, "ok")
		})

		engine.GET("/api/v2/inspection/tasks/:taskId/metadata", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			currentTask := inspectionServer.GetTask(taskId)
			if currentTask == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			result, err := currentTask.Metadata()
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, result)
		})

		engine.GET("/api/v2/inspection/tasks/:taskId/data", func(ctx *gin.Context) {
			taskId := ctx.Param("taskId")
			currentTask := inspectionServer.GetTask(taskId)
			if currentTask == nil {
				ctx.String(http.StatusNotFound, fmt.Sprintf("task %s was not found", taskId))
				return
			}
			result, err := currentTask.Result()
			if err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			inspectionDataReader, err := result.ResultStore.GetReader()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			fileSize, err := result.ResultStore.GetInspectionResultSizeInBytes()
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.DataFromReader(http.StatusOK, int64(fileSize), "application/octet-stream", inspectionDataReader, map[string]string{})
			result.ResultStore.Close()
		})

		engine.GET("/api/v2/popup", func(ctx *gin.Context) {
			currentPopup := popup.Instance.GetCurrentPopup()
			if currentPopup == nil {
				ctx.String(http.StatusOK, "")
				return
			}
			ctx.JSON(http.StatusOK, currentPopup)
		})

		engine.POST("/api/v2/popup/validate", func(ctx *gin.Context) {
			request := &popup.PopupAnswerResponse{}
			if err := ctx.ShouldBindJSON(request); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			result, err := popup.Instance.Validate(request)
			if errors.Is(err, popup.NoCurrentPopup) {
				ctx.String(http.StatusNotFound, err.Error())
				return
			}
			if errors.Is(err, popup.CurrentPopupIsntMatchingWithGivenId) {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, result)
		})

		engine.POST("/api/v2/popup/answer", func(ctx *gin.Context) {
			request := &popup.PopupAnswerResponse{}
			if err := ctx.ShouldBindJSON(request); err != nil {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			err := popup.Instance.Answer(request)
			if errors.Is(err, popup.NoCurrentPopup) {
				ctx.String(http.StatusNotFound, err.Error())
				return
			}
			if errors.Is(err, popup.CurrentPopupIsntMatchingWithGivenId) {
				ctx.String(http.StatusBadRequest, err.Error())
				return
			}
			if err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.String(http.StatusOK, "")
		})
	}
	return engine
}
