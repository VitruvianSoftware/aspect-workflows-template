package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/VitruvianSoftware/aspect-workflows-template/infrastructure/pulumi/pkg/github_repos"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define the main template repository
		mainRepo, err := github_repos.ManageMainRepo(ctx)
		if err != nil {
			return err
		}

		// Define all the external starter repositories and their deploy keys
		err = github_repos.ManageStarterRepos(ctx, mainRepo.Name)
		if err != nil {
			return err
		}

		return nil
	})
}
