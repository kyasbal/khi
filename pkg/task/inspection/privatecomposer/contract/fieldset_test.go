// Copyright 2026 Google LLC
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

package privatecomposer_contract

import (
	"testing"

	"github.com/GoogleCloudPlatform/khi/pkg/common/structured"
	"github.com/google/go-cmp/cmp"
)

func TestCloudSQLFieldSetReader_Read(t *testing.T) {
	testCases := []struct {
		name      string
		inputYAML string
		wantFS    *CloudSQLFieldSet
	}{
		{
			name: "sample log with LOG prefix, database_id, and logName with URL encoded slash",
			inputYAML: `
logName: "projects/u7b17c24abaedf8c8p-tp/logs/cloudsql.googleapis.com%2Fpostgres.log"
textPayload: "2026-08-04 14:52:17.232 UTC [115855]: [1-1] db=[unknown],user=[unknown] LOG:  connection received: host=127.0.0.1 port=60716"
resource:
  labels:
    database_id: "u7b17c24abaedf8c8p-tp:us-central1-composer-2-170061d7-sql"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "us-central1-composer-2-170061d7-sql",
				LogFileName: "postgres.log",
				Summary:     "connection received: host=127.0.0.1 port=60716",
			},
		},
		{
			name: "FATAL prefix summary extraction with postgres-audit.log",
			inputYAML: `
logName: "projects/test-proj/logs/cloudsql.googleapis.com%2Fpostgres-audit.log"
textPayload: "2026-08-04 14:52:17 UTC [1]: FATAL:  the database system is starting up"
resource:
  labels:
    database_id: "test-db"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "test-db",
				LogFileName: "postgres-audit.log",
				Summary:     "the database system is starting up",
			},
		},
		{
			name: "ERROR prefix summary extraction with regular slash logName",
			inputYAML: `
logName: "projects/test-proj/logs/cloudsql.googleapis.com/postgres-upgrade.log"
textPayload: "ERROR:  syntax error at or near \"SELECT\""
resource:
  labels:
    database_id: "test-db-2"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "test-db-2",
				LogFileName: "postgres-upgrade.log",
				Summary:     "syntax error at or near \"SELECT\"",
			},
		},
		{
			name: "missing database_id and logName defaults to unknown",
			inputYAML: `
textPayload: "INFO:  system check ok"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "",
				LogFileName: "unknown",
				Summary:     "system check ok",
			},
		},
		{
			name: "database_id with empty instance part defaults to unknown",
			inputYAML: `
resource:
  labels:
    database_id: "test-proj:"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "",
				LogFileName: "unknown",
				Summary:     "",
			},
		},
		{
			name: "audit log with methodName and principalEmail",
			inputYAML: `
logName: "projects/test-proj/logs/cloudaudit.googleapis.com%2Factivity"
protoPayload:
  methodName: "cloudsql.instances.connect"
  authenticationInfo:
    principalEmail: "sa@test-proj.iam.gserviceaccount.com"
resource:
  labels:
    database_id: "test-proj:test-db"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "test-db",
				LogFileName: "activity",
				Summary:     "cloudsql.instances.connect from sa@test-proj.iam.gserviceaccount.com",
			},
		},
		{
			name: "audit log with methodName only",
			inputYAML: `
logName: "projects/test-proj/logs/cloudaudit.googleapis.com%2Factivity"
protoPayload:
  methodName: "cloudsql.instances.create"
resource:
  labels:
    database_id: "test-proj:test-db"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "test-db",
				LogFileName: "activity",
				Summary:     "cloudsql.instances.create",
			},
		},
		{
			name: "audit log with unknown principalEmail",
			inputYAML: `
logName: "projects/test-proj/logs/cloudaudit.googleapis.com/activity"
protoPayload:
  methodName: "cloudsql.instances.delete"
  authenticationInfo:
    principalEmail: "unknown"
resource:
  labels:
    database_id: "test-proj:test-db"
`,
			wantFS: &CloudSQLFieldSet{
				DatabaseID:  "test-db",
				LogFileName: "activity",
				Summary:     "cloudsql.instances.delete",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node, err := structured.FromYAML(tc.inputYAML)
			if err != nil {
				t.Fatalf("failed to parse test YAML: %v", err)
			}
			reader := structured.NewNodeReader(node)
			fsReader := &CloudSQLFieldSetReader{}

			gotFS, err := fsReader.Read(reader)
			if err != nil {
				t.Fatalf("Read() returned unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.wantFS, gotFS); diff != "" {
				t.Errorf("Read() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
