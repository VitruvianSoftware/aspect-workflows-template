package github_repos

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi-tls/sdk/v5/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var presets = []string{
	"cpp",
	"go",
	"java",
	"js",
	"kitchen-sink",
	"kotlin",
	"minimal",
	"py",
	"ruby",
	"rust",
	"swift",
	"scala",
	"shell",
	"backstage-cpp",
	"backstage-go",
	"backstage-java",
	"backstage-js",
	"backstage-kitchen-sink",
	"backstage-kotlin",
	"backstage-minimal",
	"backstage-py",
	"backstage-ruby",
	"backstage-rust",
	"backstage-swift",
	"backstage-scala",
	"backstage-shell",
}

func ManageStarterRepos(ctx *pulumi.Context, mainRepoName pulumi.StringInput) error {
	for _, preset := range presets {
		// 1. Create the repository
		repoArgs := &github.RepositoryArgs{
			Name:        pulumi.String(preset),
			Description: pulumi.String(fmt.Sprintf("Aspect Workflows Template for %s", preset)),
			Visibility:  pulumi.String("public"),
			IsTemplate:  pulumi.Bool(true),
			HasIssues:   pulumi.Bool(true),
			HasProjects: pulumi.Bool(true),
			HasWiki:     pulumi.Bool(true),
			// Auto-delete head branches after merge. No branch protection here:
			// the deliver pipeline force-pushes generated content to main.
			DeleteBranchOnMerge: pulumi.Bool(true),
		}
		repo, err := github.NewRepository(ctx, preset, repoArgs)
		if err != nil {
			return err
		}

		// 2. Create the TLS Private Key for SSH Deploy Key
		privateKey, err := tls.NewPrivateKey(ctx, fmt.Sprintf("%s-deploy-key", preset), &tls.PrivateKeyArgs{
			Algorithm: pulumi.String("ED25519"),
		})
		if err != nil {
			return err
		}

		// 3. Add the Deploy Key to the Starter Repository
		_, err = github.NewRepositoryDeployKey(ctx, fmt.Sprintf("%s-repo-deploy-key", preset), &github.RepositoryDeployKeyArgs{
			Title:      pulumi.String("Delivery Pipeline Key"),
			Repository: repo.Name,
			Key:        privateKey.PublicKeyOpenssh,
			ReadOnly:   pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// 4. Add the Private Key as a GitHub Secret in the Main Repository
		secretSuffix := strings.ReplaceAll(preset, "-", "_")
		secretName := fmt.Sprintf("STARTER_DEPLOY_%s", strings.ToUpper(secretSuffix))

		_, err = github.NewActionsSecret(ctx, fmt.Sprintf("%s-actions-secret", preset), &github.ActionsSecretArgs{
			Repository:     mainRepoName,
			SecretName:     pulumi.String(secretName),
			PlaintextValue: privateKey.PrivateKeyOpenssh,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
