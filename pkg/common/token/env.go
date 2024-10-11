package token

import (
	"context"
	"fmt"
	"os"
)

// EnvironmentVariableTokenResolver resolves the token from environment variable.
type EnvironmentVariableTokenResolver struct {
	envVariableName string
	// tokenResolved is set to true once this resolver returns a token.
	tokenResolved bool
}

func NewEnvironmentVariableTokenResolver(envVariableName string) *EnvironmentVariableTokenResolver {
	return &EnvironmentVariableTokenResolver{
		envVariableName: envVariableName,
		tokenResolved:   false,
	}
}

// Resolve implements TokenResolver.
func (e *EnvironmentVariableTokenResolver) Resolve(ctx context.Context) (*Token, error) {
	token, found := os.LookupEnv(e.envVariableName)
	if found && !e.tokenResolved {
		e.tokenResolved = true
		return New(token), nil
	} else {
		return nil, fmt.Errorf("token not found from environment variable `%s`", e.envVariableName)
	}
}

var _ TokenResolver = (*EnvironmentVariableTokenResolver)(nil)
