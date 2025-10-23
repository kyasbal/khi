package server

import (
	"github.com/GoogleCloudPlatform/khi/pkg/api/googlecloud"
	"github.com/GoogleCloudPlatform/khi/pkg/private/api/iamtoken"
	"github.com/GoogleCloudPlatform/khi/pkg/server/option"
	"github.com/gin-gonic/gin"
)

// IAMTokenPostRequest is the payload JSON type for POST /api/v3/iam_token
type IAMTokenPostRequest struct {
	ProjectID string `json:"projectID"`
	IAMToken  string `json:"iamToken"`
}

// configureRoute configures the given gin.Engine instance to handle private endpoints.
func configureRoute(engine *gin.Engine, basePath string, iamTokenInjector *iamtoken.IAMTokenCallOptionInjectorOption) {
	router := engine.Group(basePath)

	router.POST("/api/v3/iam_token", func(ctx *gin.Context) {
		var reqBody IAMTokenPostRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			ctx.String(400, err.Error())
			return
		}
		iamTokenInjector.SetTokenFor(googlecloud.Project(reqBody.ProjectID), reqBody.IAMToken)
		ctx.String(200, "")
	})
}

type PrivateServerOption struct {
	basePath         string
	iamTokenInjector *iamtoken.IAMTokenCallOptionInjectorOption
}

// NewPrivateServerOption returns a server option to configure routing for private API.
func NewPrivateServerOption(basePath string, iamTokenHandler *iamtoken.IAMTokenCallOptionInjectorOption) *PrivateServerOption {
	return &PrivateServerOption{
		basePath: basePath, iamTokenInjector: iamTokenHandler,
	}
}

// Apply implements option.Option.
func (p *PrivateServerOption) Apply(engine *gin.Engine) error {
	configureRoute(engine, p.basePath, p.iamTokenInjector)
	return nil
}

// ID implements option.Option.
func (p *PrivateServerOption) ID() string {
	return "private-api"
}

// Order implements option.Option.
func (p *PrivateServerOption) Order() int {
	return 1000
}

var _ option.Option = (*PrivateServerOption)(nil)
