# CI pipeline

## Overview

The CI pipeline is set up using GitHub Actions and Workflows. Workflows
can be found under the `.github/workflows` directory.

There are currently 4 different workflows defined:
- CodeQL
- Test suite
- Docker image build & release
- Docs generation 

## CodeQL

This workflow will run whenever a push or pull request is submitted to the main branch, except when there are 
commits involving files contained exclusively in `website/`.

This is the default code scanning workflow scaffolded by GitHub. Detected security vulnerabilities in the codebase are
logged into the "Security" tab of the repository.

## Test suite

This workflow will run whenever a push or pull request is submitted to the main branch, except when there are commits 
involving files contained exclusively in `website/`.

The full test suite including integration and e2e tests is executed. If any test case fails, the workflow will fail.

This workflow automatically pushes the updated coverage report to the repo.

## Docker image build & release

This workflow will run whenever a push is submitted to the main branch or when a new release is published.

This workflow will build & push the Docker image defined in the `Dockerfile`. If the build succeeded, the Trivy
vulnerability scanner will scan the codebase and the Docker image in search of vulnerabilities.

### Version tags

The versioning convention follow semantic versioning. When either a push to the main branch or a release is submitted, 
the image tagged with `latest` is updated.

The tags pushed in case of a release are: `{{major}}`, `{{major}}.{{minor}}`, and `{{major}}.{{minor}}.{{patch}}`, so
that users can decide to receive updates of minor or minor+patch versions without having to bump version manually.

Currently, releases are done manually.

## Docs generation

This workflow will run whenever a push or pull request is submitted to the main branch. Commit pushes
to the main branch will trigger this workflow only when at least one file is modified in the `website/` directory.
Pull requests trigger a job checking whether a build is successful, it does not trigger deployment to GitHub pages.

We are using Docusaurus to generate the documentation.

## Protected branches and status checks

The `stable` branch is a protected branch:
- Only a linear commit history is allowed.
- All commits must be signed and verified.
- All workflows must pass, and the main branch must be up to date before merging.

## Additonal information

You can also run the workflows locally by setting up [act](https://github.com/nektos/act) on your machine.