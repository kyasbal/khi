package token

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestEnvironmentVariableTokenResolver_Resolve(t *testing.T) {
	t.Parallel()
	type fields struct {
		envVariableName string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        string
		wantErr     bool
		before      func()
		after       func()
		envVarValue string
	}{
		{
			name: "Resolve should return the token from the environment variable if it is set",
			fields: fields{
				envVariableName: "TEST_TOKEN",
			},
			args: args{
				ctx: context.Background(),
			},
			want:        "test-token",
			wantErr:     false,
			before:      func() { os.Setenv("TEST_TOKEN", "test-token") },
			after:       func() { os.Unsetenv("TEST_TOKEN") },
			envVarValue: "test-token",
		},
		{
			name: "Resolve should return an error if the environment variable is not set",
			fields: fields{
				envVariableName: "TEST_TOKEN_NOT_SET",
			},
			args: args{
				ctx: context.Background(),
			},
			want:        "",
			wantErr:     true,
			before:      func() {},
			after:       func() {},
			envVarValue: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			defer tt.after()
			e := &EnvironmentVariableTokenResolver{
				envVariableName: tt.fields.envVariableName,
			}
			got, err := e.Resolve(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvironmentVariableTokenResolver.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.want != got.RawToken {
					t.Errorf("EnvironmentVariableTokenResolver.Resolve() didn't return the expected token. want=%s,got=%s", tt.want, got.RawToken)
				}
				if !(time.Time{}).Equal(got.ValidAtLeastUntil) {
					t.Errorf("EnvironmentVariableTokenResolver.Resolve() should return a token without expiry, but non default time was given:%v", got.ValidAtLeastUntil)
				}
			}
		})
	}
}
