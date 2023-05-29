package db

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
)

const testDbFileName = "test.db"
const testDirName = "sqlite_test"

var testDir = ""
var TestLogFilePath = "../../test/data/log.txt"

func InitTestDB(t *testing.T) *SQLiteBackend {
	testDir, err := os.MkdirTemp("", testDirName)
	if err != nil {
		t.Fatalf("Could not create temp test dir, received error: %s", err)
		return nil
	}

	testDbPath := filepath.Join(testDir, testDbFileName)

	log.Printf("Setting up test database at %s\n", testDir)
	var sqlb = new(SQLiteBackend)
	err = sqlb.Open(testDbPath)

	if err != nil {
		t.Fatalf("Could not open database: %s", err)
		return nil
	}

	err = sqlb.Setup()
	if err != nil {
		t.Fatalf("Could not setup database: %s", err)
		return nil
	}

	return sqlb
}

func ReadTestLogFile(t *testing.T) string {
	testCommitLogBytes, err := os.ReadFile(TestLogFilePath)
	if err != nil {
		t.Fatalf("Could not read test commits file")
	}

	return string(testCommitLogBytes)
}

func ParsedTestCommitLog(t *testing.T) []*common.Commit {
	testCommitLog := ReadTestLogFile(t)
	testCommits, err := logread.ParseCommitLog(testCommitLog)
	if err != nil {
		t.Fatalf("Could not parse test commit log")
	}

	return testCommits
}

func IngestTestCommits(sqlb *SQLiteBackend, t *testing.T) {
	parsedCommitLog := ParsedTestCommitLog(t)

	err := sqlb.AddCommits(parsedCommitLog)
	if err != nil {
		t.Fatalf("Error during test log file ingest: %s", err)
	}

	parsedCommits, err := sqlb.Commits()
	if err != nil {
		t.Fatalf("Error checking ingested commits: %s", err)
	}

	if len(parsedCommits) != 1000 {
		t.Fatalf("Missing commits.")
	}
}

func CleanupTestDB(sqlb *SQLiteBackend) {
	if sqlb == nil {
		return
	}

	sqlb.Close()

	if testDir != "" {
		os.RemoveAll(testDir)
	}
}
