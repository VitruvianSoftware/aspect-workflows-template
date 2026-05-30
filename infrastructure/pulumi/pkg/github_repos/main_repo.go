package github_repos

import (
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ManageMainRepo(ctx *pulumi.Context) (*github.Repository, error) {
	repo, err := github.NewRepository(ctx, "aspect-workflows-template", &github.RepositoryArgs{
		Name:                pulumi.String("aspect-workflows-template"),
		Description:         pulumi.String("Scaffolding to create an Aspect-flavored Bazel workspace"),
		Visibility:          pulumi.String("public"),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(true),
		DeleteBranchOnMerge: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	// Best-practice protection on the default branch: require a PR (and an
	// up-to-date branch) before merging, and block force-pushes and branch
	// deletion. NOTE: starter repos are intentionally NOT protected — the
	// deliver pipeline force-pushes generated content to their main branch
	// (see starter_repos.go).
	_, err = github.NewBranchProtection(ctx, "aspect-workflows-template-default-protection", &github.BranchProtectionArgs{
		RepositoryId: repo.NodeId,
		Pattern:      pulumi.String("main"),
		RequiredPullRequestReviews: github.BranchProtectionRequiredPullRequestReviewArray{
			&github.BranchProtectionRequiredPullRequestReviewArgs{
				RequiredApprovingReviewCount: pulumi.Int(0),
			},
		},
		RequiredStatusChecks: github.BranchProtectionRequiredStatusCheckArray{
			&github.BranchProtectionRequiredStatusCheckArgs{
				Strict: pulumi.Bool(true),
			},
		},
		AllowsForcePushes: pulumi.Bool(false),
		AllowsDeletions:   pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
