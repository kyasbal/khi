package token

import (
	"context"
	"fmt"
	"os"
)

// EnvironmentVariableTokenResolver resolves the token from environment variable.
type EnvironmentVariableTokenResolver struct {
	envVariableName string
}

func NewEnvironmentVariableTokenResolver(envVariableName string) *EnvironmentVariableTokenResolver {
	return &EnvironmentVariableTokenResolver{
		envVariableName: envVariableName,
	}
}

// Resolve implements TokenResolver.
func (e *EnvironmentVariableTokenResolver) Resolve(ctx context.Context, expiredTokens map[string]interface{}) (string, error) {
	token, found := os.LookupEnv(e.envVariableName)
	if found {
		if _, found := expiredTokens[token]; found {
			// If the token is received from env varialble, it could be the expired token already.
			// This resolver will ignore if the token is already in the expired token set.
			return "", fmt.Errorf("token found from environment variable `%s` is already expired", e.envVariableName)
		}
		return token, nil
	} else {
		return "", fmt.Errorf("token not found from environment variable `%s`", e.envVariableName)
	}
}

var _ TokenResolver = (*EnvironmentVariableTokenResolver)(nil)
