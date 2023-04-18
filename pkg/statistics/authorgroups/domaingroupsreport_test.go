package authorgroups

import (
	"reflect"
	"testing"

	"github.com/claucambra/commit-analysis-tool/internal/db"
)

const testNumAuthors = 31
const testNumInsertions = 78388
const testNumDeletions = 42629
const testNumFilesChanged = 3915
const testNumGroupAuthors = 5
const testGroupName = "VideoLAN"
const testGroupDomain = "videolan.org"
const testGroupInsertions = 660
const testGroupDeletions = 685
const testGroupFilesChanged = 120

var testGroupAuthorsPercent = (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100
var testCommitsFile = "../../../test/data/log.txt"

var testEmailGroups = map[string][]string{
	testGroupName: {testGroupDomain},
}

func TestNewDomainGroupsReport(t *testing.T) {
	db.TestLogFilePath = testCommitsFile
	sqlb := db.InitTestDB(t)
	cleanup := func() { db.CleanupTestDB(sqlb) }
	t.Cleanup(cleanup)

	db.IngestTestCommits(sqlb, t)

	report := NewDomainGroupsReport(testEmailGroups)
	report.Generate(sqlb)

	if authorCount := report.TotalAuthors; authorCount != testNumAuthors {
		t.Fatalf("Unexpected number of authors: received %d, expected %d", authorCount, testNumAuthors)
	} else if numGroupAuthors := report.DomainNumAuthors[testGroupDomain]; numGroupAuthors != testNumGroupAuthors {
		t.Fatalf("Unexpected number of domain authors: received %d, expected %d", numGroupAuthors, testNumGroupAuthors)
	}

	print(report.TotalInsertions, "b", report.TotalDeletions, "c", report.TotalNumFilesChanged)

	testGroupData := &GroupData{
		NumAuthors:          testNumGroupAuthors,
		NumInsertions:       testGroupInsertions,
		NumDeletions:        testGroupDeletions,
		NumFilesChanged:     testGroupFilesChanged,
		AuthorsPercent:      (float32(testNumGroupAuthors) / float32(testNumAuthors)) * 100,
		InsertionsPercent:   (float32(testGroupInsertions) / float32(testNumInsertions)) * 100,
		DeletionsPercent:    (float32(testGroupDeletions) / float32(testNumDeletions)) * 100,
		FilesChangedPercent: (float32(testGroupFilesChanged) / float32(testNumFilesChanged)) * 100,
	}
	groupData := report.GroupData(testGroupName)

	if !reflect.DeepEqual(testGroupData, groupData) {
		t.Fatalf(`Retrieved group data does not match test group data: 
			Expected %+v
			Received %+v`, testGroupData, groupData)
	}

}
