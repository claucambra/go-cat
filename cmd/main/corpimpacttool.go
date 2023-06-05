package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/logread"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups/corpimpact"
)

func main() {
	var (
		batchRead            = flag.String("batch-read", "", "path to file of git clone urls to analyse")
		clonePath            = flag.String("clone-path", "", "path to store cloned repositories in")
		ingestDbPath         = flag.String("ingest-db-path", "", "path to database file")
		readDbPath           = flag.String("read-db-path", "", "path to database file")
		repoPath             = flag.String("repo-path", "", "path to git repository")
		domainGroupsFilePath = flag.String("domain-groups-file-path", "", "file containing email domain groups")
	)

	flag.Parse()

	if *batchRead != "" {

		if *clonePath == "" {
			log.Fatalf("Received empty clone path, don't know where to store cloned repos")
		} else if *domainGroupsFilePath == "" {
			log.Println("WARNING: No valid domain groupings file has been provided")
		}

		batchCloneAndRead(*batchRead, *clonePath, *domainGroupsFilePath)

	} else if *ingestDbPath != "" {

		if *repoPath == "" {
			log.Fatalf("Cannot ingest git repository commits to a database file without a path for said file.")
		}

		sqlb := newSql(*ingestDbPath)
		ingestRepoCommits(*ingestDbPath, *repoPath, sqlb)
		sqlb.Close()

	} else if *readDbPath != "" && *domainGroupsFilePath != "" {

		if *domainGroupsFilePath == "" {
			log.Println("WARNING: No valid domain groupings file has been provided")
		}

		sqlb := newSql(*readDbPath)
		printDomainGroups(*readDbPath, *domainGroupsFilePath, sqlb)
		sqlb.Close()

	} else {

		log.Fatalf("No valid individual repo or batch operation specified. Exiting.")

	}
}

func newSql(dbpath string) *db.SQLiteBackend {
	sqlb := new(db.SQLiteBackend)
	err := sqlb.Open(dbpath)
	if err != nil {
		log.Fatalf("Error opening sqlite database, received error: %s", err)
	}

	return sqlb
}

func ingestRepoCommits(ingestDbPath string, repoPath string, sqlb *db.SQLiteBackend) {
	err := sqlb.Setup()
	if err != nil {
		log.Fatalf("Error setting up sqlite database, received error: %s", err)
		os.Exit(0)
	}

	commits, err := logread.ReadCommits(repoPath)
	if err != nil {
		log.Fatalf("Error reading commits at %s: %s", repoPath, err)
	}

	log.Println("Starting commit ingest.")
	sqlb.AddCommits(commits)
	log.Println("Finished ingesting commits!")
}

func printDomainGroups(readDbPath string, domainGroupsFilePath string, sqlb *db.SQLiteBackend) {
	groupsJsonBytes, err := os.ReadFile(domainGroupsFilePath)
	if err != nil {
		log.Fatalf("Error opening domain groups json file: %s", err)
		sqlb.Close()
	}

	var groups map[string][]string
	err = json.Unmarshal(groupsJsonBytes, &groups)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		sqlb.Close()
	}

	corpReport := corpimpact.NewCorporateReport(groups, sqlb, "Corporate")
	corpReport.Generate()
	fmt.Printf("%+v", corpReport)
}

func batchCloneAndRead(urlsJsonFile string, clonePath string, domainGroupsFilePath string) {
	urlsJsonBytes, err := os.ReadFile(urlsJsonFile)
	if err != nil {
		log.Fatalf("Error opening batch fetch urls JSON file: %s", err)
	}

	var urls []string
	err = json.Unmarshal(urlsJsonBytes, &urls)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	for _, url := range urls {
		log.Printf("About to clone git repository: %s", url)

		splitUrl := strings.Split(url, "/")
		repoNameDotGit := splitUrl[len(splitUrl)-1]
		repoName := strings.TrimSuffix(repoNameDotGit, ".git")

		fullClonePath := filepath.Join(clonePath, repoName)
		ingestDbPath := filepath.Join(clonePath, repoName+".db")

		cmd := exec.Command("git",
			"clone",
			"--no-tags",
			"--filter=blob:none",
			"--filter=tree:0",
			url,
			fullClonePath)

		var stdBuffer bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		cmd.Stdout = mw
		cmd.Stderr = mw

		if err := cmd.Run(); err != nil {
			log.Panic(err)
		}

		log.Println(stdBuffer.String())

		sqlb := newSql(ingestDbPath)

		log.Printf("Clone of %s now complete. Beginning commit ingest at %s", repoName, ingestDbPath)
		ingestRepoCommits(ingestDbPath, fullClonePath, sqlb)

		log.Printf("Commit ingest for %s now complete. Beginning corporate impact analysis.", repoName)
		printDomainGroups(ingestDbPath, domainGroupsFilePath, sqlb)

		sqlb.Close()
	}
}
