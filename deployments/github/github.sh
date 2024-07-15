#!/bin/bash
set -euo pipefail

gh=$1
github_files=$2

$gh release create --generate-notes --prerelease $github_files # --draft # --latest ${STABLE_GIT_TAG}