# cloud-platform-po-linter
### Description
---
`cloud-platform-po-linter` is a github action, that will check Prometheus Rules YAML for any incorrect syntax.
### Github Action Status
---
[![Release Go project](https://github.com/ministryofjustice/cloud-platform-po-linter/actions/workflows/go-release.yaml/badge.svg)](https://github.com/ministryofjustice/cloud-platform-po-linter/actions/workflows/go-release.yaml)
### How to publish a new binary
---

Within this respoistory goreleaser tool is used to create a Go binary and push a image to Docker Hub. To publish a new binary once the Pull Request is approved and merged into the `main` branch, create a new release and the github action with automatically start and publish the new binary. 

After the new binary has been published, no changes will be needed to the workflow that is calling the image aslong as the `latest` image is being pulled.

### Github action example 
---
```
name: Prometheus Operator Linter

on:
  pull_request:
    paths:
    - *prometheus.yaml ## change to location and file name of the prometheus rule yaml
  workflow_dispatch:

jobs:
  po-lint:
    runs-on: ${{ matrix.os }}

    strategy:
      matrix:
        os: [ubuntu-latest]

    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Run po-linter
        uses: docker://ministryofjustice/cloud-platform-po-linter:latest
        continue-on-error: true
        id: po-linter
        env:
          GITHUB_OAUTH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          
      - name: Result
        uses: actions/github-script@v6
        env:
            summary: "Result:\n${{ steps.po-linter.outputs.po-linter }}"
        with:
            github-token: ${{ secrets.GITHUB_TOKEN }}
            script: |
                const output = `#### Prometheus Operator Linting Results \`${{ steps.po-linter.outcome }}\`
                <details><summary>Show</summary>
                ${process.env.summary}
                </details>`
                github.rest.issues.createComment({
                  issue_number: context.issue.number,
                  owner: context.repo.owner,
                  repo: context.repo.repo,
                  body: output
                })
      - name: Exitcode
        if: steps.po-linter.outcome != 'success'
        run: exit 1
```
