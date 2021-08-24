---
sidebar_position: 1
---

# How to contribute

Thank you for taking your valuable time to go through the contributor documentation. 

We'd love to review your PRs. Please submit a GitHub issue describing your enhancement or bugfix before 
starting to work on your proposed changes. Make sure to read the following instructions as well.

## Pull request process

1. Fork the project repository, develop and test your changes on that fork
2. Commit changes
3. Submit a PR from your fork to this project.

## Commits

Please have meaningful and well-formed commit messages. [This article](https://chris.beams.io/posts/git-commit/) 
summarizes nicely the interpretation of "meaningful" and "well-formed".

### Commit signature verification

Each commit signature **must be signed** (using the uppercase `-S` option of `git commit`) and **verified** (by uploading your public
key to your own GitHub account). 

- [About commit signature verification](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/about-commit-signature-verification)
- You can read [this excellent article](https://mikegerwitz.com/2012/05/a-git-horror-story-repository-integrity-with-signed-commits)
  for an overview of git commit signing (if you have read it, we are using option #3).
### Sign off your changes

The Developer Certificate of Origin (DCO) is a lightweight way for contributors to certify that they wrote or otherwise 
have the right to submit the code they are contributing to the project. Here is the full text of the [DCO](https://github.com/bedag/kubernetes-dbaas/blob/main/DCO). 

Contributors **must sign-off** each commit by adding a Signed-off-by line to commit messages (using the lowercase `-s` option of `git commit`).

```
Signed-off-by: Random J Developer <random@developer.example.org>
See git help commit:

-s, --signoff
    Add Signed-off-by line by the committer at the end of the commit log
    message. The meaning of a signoff depends on the project, but it typically
    certifies that committer has the rights to submit this work under the same
    license and agrees to a Developer Certificate of Origin (see
    http://developercertificate.org/ for more information).
```

By making a contribution to this project and signing-off your commits, 
you agree to and comply with the Developer's Certificate of Origin.

## Pull requests requirements

All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. Consult [GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on using pull requests. See the above stated requirements for PR on this project.

Your PR has to fulfill the following points, to be considered:

- Workflows must pass (see [CI pipeline](/docs/contributing/ci)).
- All commits correspond to the requirements (see [Commits](/docs/contributing/how-to-contribute)).

### Versioning

This project follows the [semver standard](https://semver.org/) for everything. Versions are bumped by maintainers.
If your contribution contains breaking changes, please explicitly state it in your GitHub issue.

### Review

If your PR does not have any activity after certain days, feel free to comment a reminder. Your PR requires approval 
to be mergeable.

### Testing

Your PR should enclose updated or new test cases. See [Testing](/docs/contributing/testing). 

## Documentation

### Website

The documentation is done using [Docusaurus](https://github.com/facebook/docusaurus). Install Docusaurus, 
run `yarn start` in the `/website` directory and start modifying the Markdown
files contained in `/website/docs`. Refer directly to the official Docusaurus documentation to learn more. 

### Helm Chart

Any change to the Helm Chart's documentation requires you to run [helm-docs](https://github.com/norwoodj/helm-docs),
which will autogenerate a README.md for the Chart.

### Godocs

Any change to the codebase requires you to update the relative godocs as well.