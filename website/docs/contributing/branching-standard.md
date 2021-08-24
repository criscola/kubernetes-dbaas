---
sidebar_position: 2
---

# Branching standard

Make sure you have read [How to contribute](/docs/contributing/how-to-contribute) first.

## Quick legend

<table>
  <thead>
    <tr>
      <th>Instance</th>
      <th>Branch</th>
      <th>Description, Instructions, Notes.</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>Stable</td>
      <td>stable</td>
      <td>Accepts merges from Working and Hotfixes.</td>
    </tr>
    <tr>
      <td>Working</td>
      <td>main</td>
      <td>Accepts merges from Features/Issues and Hotfixes.</td>
    </tr>
    <tr>
      <td>Features</td>
      <td>topic/*</td>
      <td>Always branch off HEAD of Working.</td>
    </tr>
    <tr>
      <td>Bug fixes</td>
      <td>bugfix/*</td>
      <td>Always branch off HEAD of Working.</td>
    </tr>
    <tr>
      <td>Hotfix</td>
      <td>hotfix/*</td>
      <td>Always branch off Stable.</td>.
    </tr>
  </tbody>
</table>

## Main branches

The main repository will always hold two evergreen branches:

* `main`
* `stable`

The main branch should be considered `origin/main` and will be the main branch where the source code of `HEAD` always 
reflects a state with the latest delivered development changes for the next release. As a developer, you will be 
branching and merging from `main`.

Consider `origin/stable` to always represent the latest code deployed to production. During day to day development, 
the `stable` branch will not be interacted with.

When the source code in the `main` branch is stable and has been deployed, all the changes will be merged into `stable` 
and tagged with a release number. Tagging and releasing is explained in detail in [CI pipeline](/docs/contributing/ci).

Squash and rebase as desired; but strive to present a consistent and descriptive commit history when doing so.

Do not fast-forward commits to the `main` branch; make sure to create a commit with `--no-ff` during merge.

## Supporting branches

Supporting branches are used to aid parallel development between team members, ease tracking of features, and to assist 
in quickly fixing live production problems. Unlike the main branches, these branches always have a limited lifetime, 
since they will be removed eventually.

The different types of supporting branches we use are:

- Feature branches
- Bug branches
- Hotfix branches

### Feature branches

No matter when the feature branch will be finished, it will always be merged back into 
the main branch.

* Must branch from: `main`
* Must merge back into: `main`
* Branch naming convention: `topic/<short descriptive name>`

Periodically, changes made to `main` (if any) should be merged back into your feature branch.

### Bug branches

Bug branches differ have the same lifecycle of feature branches. Bug branches will be created when there is a bug 
that should be fixed and merged into the next release. For that reason, a bug branch typically will not last longer 
than one deployment cycle (where a new release is produced). No matter when the bug branch will be finished, 
it will always be merged back into `main`.

* Must branch from: `main`
* Must merge back into: `main`
* Branch naming convention: `bugfix/<short descriptive name>`

Periodically, changes made to `main` (if any) should be merged back into your bug branch.

### Hotfix Branches

A hotfix branch comes from the need to act immediately upon an undesired state of a live release version. 
Additionally, because of the urgency, a hotfix is not required to be pushed during a scheduled release. 
Due to these requirements, a hotfix branch is always branched from a tagged `stable` branch. This is done 
for two reasons:

First, development on the `main` branch can continue while the hotfix is being addressed. Second, a tagged `stable` 
branch still represents what is in production. At the point in time when a hotfix is needed, there could have been
multiple commits to `main` which would then no longer represent production.

* Must branch from: tagged `stable`
* Must merge back into: `main` and `stable`
* Branch naming convention: `hotfix/<short descriptive name>`

#### Working with a hotfix branch

When development on the hotfix is complete, a maintainer should merge changes into `stable` and then the patch version
bumped (following semantic versioning).

Merge changes into `main` so not to lose the hotfix and then delete the remote hotfix branch.

## Additional information

The branching standard was based on [this gist](https://gist.github.com/digitaljhelms/4287848).

