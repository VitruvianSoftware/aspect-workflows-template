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

	// Best-practice protection on the default branch (`main`) and the active
	// integration branch (`platform-v2.0`): require a PR (and an up-to-date
	// branch) before merging, and block force-pushes and branch deletion.
	// enforceAdmins is left false so repo owners can still sync `main` from the
	// upstream template. NOTE: starter repos are intentionally NOT protected —
	// the deliver pipeline force-pushes generated content to their main branch
	// (see starter_repos.go).
	for _, branch := range []string{"main", "platform-v2.0"} {
		if err := protectBranch(ctx, repo, branch); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// protectBranch applies the standard best-practice protection rule to a single
// branch pattern on repo: require a PR with an up-to-date branch, and block
// force-pushes and deletion. enforceAdmins is intentionally left at its default
// (false) so owners retain an escape hatch (e.g. mirroring `main` upstream).
func protectBranch(ctx *pulumi.Context, repo *github.Repository, branch string) error {
	_, err := github.NewBranchProtection(ctx, "aspect-workflows-template-"+branch+"-protection", &github.BranchProtectionArgs{
		RepositoryId: repo.NodeId,
		Pattern:      pulumi.String(branch),
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
	return err
}
