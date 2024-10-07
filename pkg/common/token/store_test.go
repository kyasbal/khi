package token

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBasicTokenStore_GetToken(t *testing.T) {
	t.Parallel()
	type fields struct {
		resolver      TokenResolver
		expiredTokens map[string]interface{}
		lastToken     string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "GetToken should return the token from the resolver if lastToken is empty",
			fields: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{},
				lastToken:     "",
			},
			args: args{
				ctx: context.Background(),
			},
			want:    "test-token",
			wantErr: false,
		},
		{
			name: "GetToken should pass the expiredTokens and the expired token should be ignored",
			fields: fields{
				resolver: NewMultiTokenResolver(
					newSpyTokenResolver("test-token-1"),
					newSpyTokenResolver("test-token-2"),
				),
				expiredTokens: map[string]interface{}{
					"test-token-1": struct{}{},
				},
				lastToken: "test-token-2",
			},
			args: args{
				ctx: context.Background(),
			},
			want:    "test-token-2",
			wantErr: false,
		},
		{
			name: "GetToken should return the lastToken if it is not empty",
			fields: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{},
				lastToken:     "cached-token",
			},
			args: args{
				ctx: context.Background(),
			},
			want:    "cached-token",
			wantErr: false,
		},
		{
			name: "GetToken should return error if resolver returns error",
			fields: fields{
				resolver:      newMockErrorTokenResolver(),
				expiredTokens: map[string]interface{}{},
				lastToken:     "",
			},
			args: args{
				ctx: context.Background(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBasicTokenStore("test", tt.fields.resolver)
			b.lastToken = tt.fields.lastToken
			b.expiredTokens = tt.fields.expiredTokens
			got, err := b.GetToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicTokenStore.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("BasicTokenStore.GetToken() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestBasicTokenStore_RefreshToken(t *testing.T) {
	t.Parallel()
	type fields struct {
		resolver      TokenResolver
		expiredTokens map[string]interface{}
		lastToken     string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "RefreshToken should clear lastToken and add it to expiredTokens",
			fields: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{},
				lastToken:     "expired-token",
			},
			args: args{
				ctx: context.Background(),
			},
			want: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{"expired-token": struct{}{}},
				lastToken:     "test-token",
			},
		},
		{
			name: "RefreshToken should update the token if lastToken is empty",
			fields: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{},
				lastToken:     "",
			},
			args: args{
				ctx: context.Background(),
			},
			want: fields{
				resolver:      newSpyTokenResolver("test-token"),
				expiredTokens: map[string]interface{}{},
				lastToken:     "test-token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBasicTokenStore("test", tt.fields.resolver)
			b.lastToken = tt.fields.lastToken
			b.expiredTokens = tt.fields.expiredTokens
			b.RefreshToken(tt.args.ctx)
			if diff := cmp.Diff(tt.want.expiredTokens, b.expiredTokens); diff != "" {
				t.Errorf("BasicTokenStore.RefreshToken() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.want.lastToken, b.lastToken); diff != "" {
				t.Errorf("BasicTokenStore.RefreshToken() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
