# :recycle: Contributing

Contributions to this project are very much welcome, please make sure that your code changes are tested and that they follow Go best-practices.

## Getting started

### Required environmental variables

The tests require at least two environment variables: `STREAM_API_KEY` and `STREAM_API_SECRET`. There are multiple ways to provide that:
- simply set it in your current shell (`export STREAM_API_KEY=xyz`)
- you could use [direnv](https://direnv.net/)
- if you debug the tests in VS Code, you can set up an env file there as well: `"go.testEnvFile": "${workspaceFolder}/.env"`.

### Code formatting & linter

We enforce code formatting with [`gofumpt`](https://github.com/mvdan/gofumpt) (a stricter `gofmt`). If you use VS Code, it's recommended to set this setting there for auto-formatting:

```json
{
    "editor.formatOnSave": true,
    "gopls": {
        "formatting.gofumpt": true
    },
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ]
}
```

Gofumpt will mostly take care of your linting issues as well.

## Commit message convention

Since we're autogenerating our [CHANGELOG](./CHANGELOG.md), we need to follow a specific commit message convention.
You can read about conventional commits [here](https://www.conventionalcommits.org/). Here's how a usual commit message looks like for a new feature: `feat: allow provided config object to extend other configs`. A bugfix: `fix: prevent racing of requests`.

## Release (for Stream developers)

Releasing this package involves two GitHub Action steps:

- Kick off a job called `initiate_release` ([link](https://github.com/GetStream/stream-chat-go/actions/workflows/initiate_release.yml)).

The job creates a pull request with the changelog. Check if it looks good.

- Merge the pull request.

Once the PR is merged, it automatically kicks off another job which will create the tag and created a GitHub release.
