#!/bin/bash
# Applies the standard branch protection to ONE repository.
# This is the single-repo version of the owner's gh-repo-baseline.sh:
#   - production branch (default branch) protected: PRs required, admins included,
#     no force pushes, no deletions;
#   - a `dev` integration branch with light protection (no force push / deletion);
#   - merged branches deleted automatically;
#   - Copilot code review ruleset on the production branch.
# Idempotent. Usage: scripts/apply-baseline.sh [owner/repo]
set -euo pipefail
REPO="${1:-AnderCMD/TaskGlass}"
T=$(mktemp -d); trap 'rm -rf "$T"' EXIT
DEF=$(gh api "repos/$REPO" -q .default_branch)

echo "protecting $REPO ($DEF)"
cat > "$T/strict.json" <<'JSON'
{"required_status_checks":null,"enforce_admins":true,
 "required_pull_request_reviews":{"required_approving_review_count":0,"dismiss_stale_reviews":true},
 "restrictions":null,"allow_force_pushes":false,"allow_deletions":false,"required_conversation_resolution":false}
JSON
cat > "$T/light.json" <<'JSON'
{"required_status_checks":null,"enforce_admins":false,"required_pull_request_reviews":null,
 "restrictions":null,"allow_force_pushes":false,"allow_deletions":false,"required_conversation_resolution":false}
JSON
gh api -X PUT "repos/$REPO/branches/$DEF/protection" --input "$T/strict.json" >/dev/null

if ! gh api "repos/$REPO/branches/dev" >/dev/null 2>&1; then
  sha=$(gh api "repos/$REPO/git/ref/heads/$DEF" -q .object.sha)
  gh api -X POST "repos/$REPO/git/refs" -f ref=refs/heads/dev -f sha="$sha" >/dev/null
fi
gh api -X PUT "repos/$REPO/branches/dev/protection" --input "$T/light.json" >/dev/null

gh api -X PATCH "repos/$REPO" -F delete_branch_on_merge=true >/dev/null

if [ "$(gh api "repos/$REPO/rulesets" -q '[.[]|select(.name=="Copilot code review")]|length')" = 0 ]; then
  python3 - "$DEF" > "$T/copilot.json" <<'PY'
import sys, json
print(json.dumps({"name": "Copilot code review", "target": "branch", "enforcement": "active",
  "conditions": {"ref_name": {"include": ["refs/heads/" + sys.argv[1]], "exclude": []}},
  "rules": [{"type": "copilot_code_review", "parameters": {"review_on_push": False, "review_draft_pull_requests": False}}]}))
PY
  gh api -X POST "repos/$REPO/rulesets" --input "$T/copilot.json" >/dev/null
fi
echo "done"
