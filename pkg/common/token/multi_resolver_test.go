package token

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMultiTokenResolver_Resolve(t *testing.T) {
	t.Parallel()
	type fields struct {
		resolvers []TokenResolver
	}
	type args struct {
		ctx           context.Context
		expiredTokens map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Resolve should return the first token successfully resolved",
			fields: fields{
				resolvers: []TokenResolver{
					newSpyTokenResolver("token1"),
					newSpyTokenResolver("token2"),
				},
			},
			args: args{
				ctx:           context.Background(),
				expiredTokens: map[string]interface{}{},
			},
			want:    "token1",
			wantErr: false,
		},
		{
			name: "Resolve should skip resolvers that return an error",
			fields: fields{
				resolvers: []TokenResolver{
					newMockErrorTokenResolver(),
					newSpyTokenResolver("token2"),
				},
			},
			args: args{
				ctx:           context.Background(),
				expiredTokens: map[string]interface{}{},
			},
			want:    "token2",
			wantErr: false,
		},
		{
			name: "Resolve should return ErrNoNewTokenResolved if every resolvers return error",
			fields: fields{
				resolvers: []TokenResolver{
					newMockErrorTokenResolver(),
					newMockErrorTokenResolver(),
				},
			},
			args: args{
				ctx:           context.Background(),
				expiredTokens: map[string]interface{}{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MultiTokenResolver{
				resolvers: tt.fields.resolvers,
			}
			got, err := m.Resolve(tt.args.ctx, tt.args.expiredTokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiTokenResolver.Resolve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("MultiTokenResolver.Resolve() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
