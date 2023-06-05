package corpimpact

import (
	"github.com/claucambra/commit-analysis-tool/internal/db"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/authorgroups"
	"github.com/claucambra/commit-analysis-tool/pkg/statistics/commitimpact"
)

type CorporateReport struct {
	CorporateGroupName string
	GroupsOfDomains    map[string][]string

	CorporateGroup *authorgroups.GroupData
	CommunityGroup *authorgroups.GroupData

	// Correlations based upon year-by-year aggregated figures for both groups
	InsertionsCorrel float64
	DeletionsCorrel  float64
	AuthorsCorrel    float64

	DomainGroupsReport           *authorgroups.DomainGroupsReport
	CorporateGroupSurvivalReport *authorgroups.GroupSurvivalReport
	CommunityGroupSurvivalReport *authorgroups.GroupSurvivalReport

	CorporateCommitImpactReport *commitimpact.CommitImpactReport
	CommunityCommitImpactReport *commitimpact.CommitImpactReport

	sqlb *db.SQLiteBackend
}

func NewCorporateReport(groupsOfDomains map[string][]string, sqlb *db.SQLiteBackend, corporateGroupName string) *CorporateReport {
	if corporateGroupName == "" {
		corporateGroupName = "Corporate"
	}

	return &CorporateReport{
		CorporateGroupName: corporateGroupName,
		GroupsOfDomains:    groupsOfDomains,
		sqlb:               sqlb,
	}
}

func (corpReport *CorporateReport) Generate() {
	domainGroupsReport := authorgroups.NewDomainGroupsReport(corpReport.GroupsOfDomains, corpReport.sqlb)
	domainGroupsReport.Generate()
	corpReport.DomainGroupsReport = domainGroupsReport

	corpGroup := domainGroupsReport.GroupData(corpReport.CorporateGroupName)
	corpReport.CorporateGroup = corpGroup

	commGroup := domainGroupsReport.UnknownGroupData()
	corpReport.CommunityGroup = commGroup

	insertionsCorrel, deletionsCorrel, authorsCorrel := corpGroup.Correlation(commGroup)
	corpReport.InsertionsCorrel = insertionsCorrel
	corpReport.DeletionsCorrel = deletionsCorrel
	corpReport.AuthorsCorrel = authorsCorrel

	corpGroupSurvival := authorgroups.NewGroupSurvivalReport(corpReport.sqlb, corpGroup.Authors)
	corpGroupSurvival.Generate()
	corpReport.CorporateGroupSurvivalReport = corpGroupSurvival

	commGroupSurvival := authorgroups.NewGroupSurvivalReport(corpReport.sqlb, commGroup.Authors)
	commGroupSurvival.Generate()
	corpReport.CommunityGroupSurvivalReport = commGroupSurvival

	corpGroupImpact := commitimpact.NewCommitImpactReport(corpGroup.Commits)
	corpGroupImpact.Generate()
	corpReport.CorporateCommitImpactReport = corpGroupImpact

	commGroupImpact := commitimpact.NewCommitImpactReport(commGroup.Commits)
	commGroupImpact.Generate()
	corpReport.CommunityCommitImpactReport = commGroupImpact
}