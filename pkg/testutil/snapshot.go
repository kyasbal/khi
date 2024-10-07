package testutil

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const SNAPSHOT_FOLDER = "test/snapshots/"

func ValidateSnapshot(t *testing.T, snapshotLabel string, currentResult string) {
	lastSnapshot := getLastSnapshotResult(snapshotLabel)
	if !willUpdateSnapshot() {
		assert.Equal(t, lastSnapshot, currentResult)
	} else {
		if lastSnapshot == "" {
			slog.Warn("Last snapshot seems not to be found!!")
		}
		if lastSnapshot != currentResult {
			slog.Warn(fmt.Sprintf("Updating snapshot for `%s", snapshotLabel))
		}
		updateLastSnapshotResult(snapshotLabel, currentResult)
	}
}

func willUpdateSnapshot() bool {
	return os.Getenv("UPDATE_SNAPSHOT") == "true"
}

func getLastSnapshotResult(snapshotLabel string) string {
	return MustReadText(SNAPSHOT_FOLDER+snapshotLabel, "")
}

func updateLastSnapshotResult(snapshotLabel string, content string) {
	err := os.WriteFile(SNAPSHOT_FOLDER+snapshotLabel, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
