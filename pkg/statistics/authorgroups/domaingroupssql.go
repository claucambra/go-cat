package authorgroups

import (
	"database/sql"
	"log"
	"time"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

func domainCommitsRows(sqlb *db.SQLiteBackend, domain string) (*sql.Rows, error) {
	stmt := "SELECT * FROM commits WHERE instr(author_email, ?) > 0"
	accStmt, err := sqlb.Db.Prepare(stmt)
	if err != nil {
		log.Fatalf("Encountered error preparing commits retrieval statement: %s", err)
		return nil, err
	}

	defer accStmt.Close()
	return accStmt.Query(domain)
}

func domainCommits(sqlb *db.SQLiteBackend, domain string) ([]*common.Commit, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	commits := []*common.Commit{}

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)
		commits = append(commits, commit)
	}

	return commits, nil
}

func domainLineChanges(sqlb *db.SQLiteBackend, domain string) (*common.LineChanges, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	numInsertions := 0
	numDeletions := 0
	numFilesChanged := 0

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)

		numInsertions += commit.NumInsertions
		numDeletions += commit.NumDeletions
		numFilesChanged += commit.NumFilesChanged
	}

	return &common.LineChanges{
		NumInsertions: numInsertions,
		NumDeletions:  numDeletions,
	}, nil
}

func domainYearlyLineChanges(sqlb *db.SQLiteBackend, domain string) (common.YearlyLineChangeMap, error) {
	rows, err := domainCommitsRows(sqlb, domain)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	yearBuckets := common.YearlyLineChangeMap{}

	for rows.Next() {
		commit := sqlb.ScanRowInRowsToCommits(rows)
		commitYear := time.Unix(commit.AuthorTime, 0).UTC().Year()
		yearBuckets.AddLineChanges(&(commit.LineChanges), commitYear)
	}

	return yearBuckets, nil
}

// The years in which an author has contributed
func authorYears(sqlb *db.SQLiteBackend, authorEmail string) ([]int, error) {
	authorCommits, err := sqlb.AuthorCommits(authorEmail)
	if err != nil {
		log.Fatalf("Error retrieving rows: %s", err)
		return nil, err
	}

	yearsMap := map[int]bool{}
	for _, commit := range authorCommits {
		commitYear := time.Unix(commit.AuthorTime, 0).UTC().Year()
		yearsMap[commitYear] = true
	}

	return common.SortedMapKeys(yearsMap), nil
}
