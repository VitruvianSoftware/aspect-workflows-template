package github_repos

import (
	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ManageMainRepo(ctx *pulumi.Context) (*github.Repository, error) {
	repo, err := github.NewRepository(ctx, "aspect-workflows-template", &github.RepositoryArgs{
		Name:        pulumi.String("aspect-workflows-template"),
		Description: pulumi.String("Scaffolding to create an Aspect-flavored Bazel workspace"),
		Visibility:  pulumi.String("public"),
		HasIssues:   pulumi.Bool(true),
		HasProjects: pulumi.Bool(true),
		HasWiki:     pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return repo, nil
}
