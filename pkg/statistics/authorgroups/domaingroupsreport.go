package authorgroups

import (
	"fmt"
	"sort"
	"strings"

	"github.com/claucambra/commit-analysis-tool/pkg/common"
)

const fallbackGroupName = "unknown"

type DomainGroup struct {
	AuthorCount  int
	DomainCounts map[string]int
}

type DomainGroupsReport struct {
	AuthorCount  int
	DomainGroups map[string]*DomainGroup

	authors       map[string]bool
	domainToGroup map[string]string
}

func NewDomainGroup() *DomainGroup {
	return &DomainGroup{
		AuthorCount:  0,
		DomainCounts: make(map[string]int),
	}
}

func NewDomainGroupsReport(domainGroups map[string]string) *DomainGroupsReport {
	report := &DomainGroupsReport{
		AuthorCount:  0,
		DomainGroups: make(map[string]*DomainGroup),

		authors:       make(map[string]bool),
		domainToGroup: make(map[string]string),
	}

	for group, domainName := range domainGroups {
		report.domainToGroup[domainName] = group
	}

	return report
}

func (report *DomainGroupsReport) ParseCommits(commits []*common.CommitData) {
	if len(commits) == 0 {
		return
	}

	for _, commit := range commits {
		report.AddCommit(*commit)
	}
}

func (report *DomainGroupsReport) AddCommit(commit common.CommitData) {
	authorString := commit.AuthorEmail
	if authorString == "" {
		authorString = commit.AuthorName
	}

	if report.authors[authorString] { // Already counted, skip
		return
	} else if authorString != "" {
		report.authors[authorString] = true
		report.AuthorCount += 1
	}

	groupString := fallbackGroupName
	emailDomain := "unknown"

	if splitAuthorEmail := strings.Split(commit.AuthorEmail, "@"); len(splitAuthorEmail) == 2 {
		emailDomain = splitAuthorEmail[1]
		groupString = report.domainToGroup[emailDomain]

		if groupString == "" {
			groupString = fallbackGroupName
		}
	}

	group := report.DomainGroups[groupString]
	if group == nil {
		group = NewDomainGroup()
		report.DomainGroups[groupString] = group
	}

	group.AuthorCount += 1
	group.DomainCounts[emailDomain] += 1
}

func (report *DomainGroupsReport) GroupPercentageOfTotal(group string) float32 {
	DomainGroup := report.DomainGroups[group]
	if DomainGroup == nil {
		return 0
	}

	return (float32(DomainGroup.AuthorCount) / float32(report.AuthorCount)) * 100
}

func (report *DomainGroupsReport) String() string {
	// Get sorted domain group names
	sortedDomainGroupNames := make([]string, 0, len(report.DomainGroups))
	for domainGroupName := range report.DomainGroups {
		sortedDomainGroupNames = append(sortedDomainGroupNames, domainGroupName)
	}

	sort.SliceStable(sortedDomainGroupNames, func(i, j int) bool {
		groupNameA := sortedDomainGroupNames[i]
		groupNameB := sortedDomainGroupNames[j]

		groupACount := report.DomainGroups[groupNameA].AuthorCount
		groupBCount := report.DomainGroups[groupNameB].AuthorCount

		if groupACount == groupBCount {
			return groupNameA < groupNameB
		}

		return groupACount > groupBCount
	})

	reportString := "Author domain groups report\n"
	reportString += fmt.Sprintf("Total repository authors: %d\n", report.AuthorCount)
	reportString += "Number of authors by group:\n"

	for _, groupName := range sortedDomainGroupNames {
		groupStruct := report.DomainGroups[groupName]
		reportString += fmt.Sprintf("\t\"%s\":\t%d (%f%%)\n", groupName, groupStruct.AuthorCount, report.GroupPercentageOfTotal(groupName))

		// Get sorted domains
		sortedDomainNames := make([]string, 0, len(groupStruct.DomainCounts))
		for domainName := range groupStruct.DomainCounts {
			sortedDomainNames = append(sortedDomainNames, domainName)
		}

		sort.SliceStable(sortedDomainNames, func(i, j int) bool {
			domainA := sortedDomainNames[i]
			domainB := sortedDomainNames[j]

			domainACount := groupStruct.DomainCounts[domainA]
			domainBCount := groupStruct.DomainCounts[domainB]

			if domainACount == domainBCount {
				return domainA < domainB
			}

			return domainACount > domainBCount
		})

		for _, domainName := range sortedDomainNames {
			print(domainName)
			domainCount := groupStruct.DomainCounts[domainName]
			reportString += fmt.Sprintf("\t\t%s:\t%d\n", domainName, domainCount)
		}
	}

	return reportString
}