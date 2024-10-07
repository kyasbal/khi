package httpclient

import (
	"context"
	"net/http"
	"time"
)

// TokenApplyResult contains the operation result of ApplyCurrentToken.
// This contains TokenApplyTime to purge related token only when the token was expired.
// Because the time to purge token and obtaining token is different and these are used in parallel, it can purge newly refreshed token without checking the time to issue the token.
type TokenApplyResult struct {
	TokenObtainedAt time.Time
}

type TokenApplier interface {
	ApplyCurrentToken(ctx context.Context, req *http.Request) (*TokenApplyResult, error)
}

type NopTokenApplier struct{}

var _ TokenApplier = (*NopTokenApplier)(nil)

func (n *NopTokenApplier) ApplyCurrentToken(ctx context.Context, req *http.Request) (*TokenApplyResult, error) {
	return &TokenApplyResult{
		TokenObtainedAt: time.Now(),
	}, nil
}
