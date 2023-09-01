#!/bin/bash
set -euo pipefail

gh=$1
github_files=$2

$gh release create --generate-notes --latest ${STABLE_GIT_TAG} $github_files # --draft