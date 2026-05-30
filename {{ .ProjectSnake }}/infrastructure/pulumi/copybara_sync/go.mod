module github.com/YOUR_GITHUB_ORG/{{ .ProjectKebab }}/infrastructure/pulumi/copybara_sync

go 1.26.1

require (
	github.com/pulumi/pulumi-github/sdk/v6 v6.12.2
	github.com/pulumi/pulumi-tls/sdk/v5 v5.3.1
	github.com/pulumi/pulumi/sdk/v3 v3.231.0
)
