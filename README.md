# cloud-platform-po-linter
### Description
---
This linter will check Prometheus Rules YAML for any incorrect syntax.

### Building binary for release 
---
From within the repository run the following command to create the correct binary to be ran in the docker image 
```
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" .
```

### Github action example 
---
```
name: Prometheus Operator Linter

on:
  pull_request:
    paths:
    - *prometheus.yaml
    - *prometheus.yml

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
        uses: ministryofjustice/cloud-platform-po-linter@p1.0.0
        continue-on-error: true
        id: po-linter
        env:
          GITHUB_OAUTH_TOKEN: ${{ secrets.ACTION_TOKEN }}
          
      - name: Result
        uses: actions/github-script@v6
        env:
            summary: "Result:\n${{ steps.po-linter.outputs.po-linter }}"
        with:
            github-token: ${{ secrets.ACTION_TOKEN }}
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
