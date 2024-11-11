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

package iamtoken

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_extractTokenFromGACommand(t *testing.T) {
	type args struct {
		command string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "extractTokenFromGACommand should return the token from the command (old command style)",
			args: args{
				command: "gcloud auth configure-docker --quiet\n" +
					"docker run --rm --pull always -p 8080:8080 -it -e KHI_FIXED_PROJECT_ID=\"XXXXX\" -e GCP_DEFAULT_PROJECT=\"XXXXX\" -e GCP_ACCESS_TOKEN=`gcloud auth print-access-token` -e GCP_IDENTITY_TOKEN=`gcloud auth print-identity-token` -e KHI_GA_LABELS=\"justification=vector/XXXXX,user=xxxxx\" -e IAM_TOKEN=\"AAAA\\\n" +
					"BBBB\\\n" +
					"CCCC\\\n" +
					"DDDD\\\n" +
					"EEEE\\\n" +
					"FFFF\\\n" +
					"GGGG\\\n" +
					"HHHH\\\n" +
					"IIII\" gcr.io/kubernetes-history-inspector/standalone:latest",
			},
			want:    "AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIII",
			wantErr: false,
		},
		{
			name: "extractTokenFromGACommand should return the token from the command",
			args: args{
				command: "gcloud auth configure-docker --quiet\n" +
					"docker run --rm --pull always -p 8080:8080 -it -e KHI_FIXED_PROJECT_ID=\"XXXXX\" -e GCP_DEFAULT_PROJECT=\"XXXXX\" -e GCP_ACCESS_TOKEN=`gcloud auth print-access-token` -e KHI_GA_LABELS=\"justification=vector/XXXXX,user=xxxxx\" -e IAM_TOKEN=\"AAAA\\\n" +
					"BBBB\\\n" +
					"CCCC\\\n" +
					"DDDD\\\n" +
					"EEEE\\\n" +
					"FFFF\\\n" +
					"GGGG\\\n" +
					"HHHH\\\n" +
					"IIII\" gcr.io/kubernetes-history-inspector/standalone:latest",
			},
			want:    "AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIII",
			wantErr: false,
		},
		{
			name: "extractTokenFromGACommand should return error when it can't find the IAM_TOKEN section",
			args: args{
				command: "gcloud auth configure-docker --quiet\n" +
					"docker run --rm --pull always -p 8080:8080 -it -e KHI_FIXED_PROJECT_ID=\"XXXXX\" -e GCP_DEFAULT_PROJECT=\"XXXXX\" -e GCP_ACCESS_TOKEN=`gcloud auth print-access-token` -e GCP_IDENTITY_TOKEN=`gcloud auth print-identity-token` -e KHI_GA_LABELS=\"justification=vector/XXXXX,user=xxxxx\" gcr.io/kubernetes-history-inspector/standalone:latest",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTokenFromGACommand(tt.args.command)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTokenFromGACommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("extractTokenFromGACommand() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
