import GitHubButton from "react-github-btn";
import React from "react";

export function HeaderButtons() {
  return (
    <>
      <GitHubButton href="https://github.com/bedag/kubernetes-dbaas" data-icon="octicon-star" data-size="large"
                    aria-label="Star bedag/kubernetes-dbaas on GitHub">Star</GitHubButton>
      <GitHubButton href="https://github.com/bedag/kubernetes-dbaas/subscription" data-icon="octicon-eye"
                    data-size="large" aria-label="Watch bedag/kubernetes-dbaas on GitHub">Watch</GitHubButton>
      <GitHubButton href="https://github.com/bedag/kubernetes-dbaas/fork" data-icon="octicon-repo-forked"
                    data-size="large" aria-label="Fork bedag/kubernetes-dbaas on GitHub">Fork</GitHubButton>
    </>
  )
}