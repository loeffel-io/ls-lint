#!/bin/bash
set -euo pipefail

pnpm=$1
package=$2

# cd deployments/npm/ls_lint
# echo $pnpm
ls -al
$pnpm publish --no-git-checks --tag beta --dry-run