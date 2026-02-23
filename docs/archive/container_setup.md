# Status: Archived

Reason: Legacy devcontainer troubleshooting note; superseded by `.devcontainer/README.md` and `docs/agent-workflow.md`.

Replacement references:
- `.devcontainer/README.md`
- `docs/agent-workflow.md`

---


# 3 fixes to the devcontainer config

docker-outside-of-docker:1 → docker-in-docker:2 — The docker-outside-of-docker feature tries to bind-mount the host's Docker socket (/var/run/docker.sock), which frequently fails on Windows/Docker Desktop due to path translation issues between Windows ↔ WSL ↔ container. docker-in-docker runs its own Docker daemon inside the dev container — more isolated and reliable on Windows.

Removed workspaceMount and workspaceFolder overrides — VS Code was already passing --workspace-mount-consistency cached in the command line (visible in your error log). The explicit workspaceMount was conflicting with that. The default mount at /workspaces/card-generator works fine.

Made post-create.sh non-fatal — Removed set -e so a single failing step (like go install or npm install) won't prevent the container from starting. Each step now logs a warning and continues. Also tries both /workspaces/card-generator/ and /workspace/ paths.

Other cleanups:

Removed nodeGypDependencies: false (unnecessary option)
Removed ms-vscode.vscode-typescript-next and GitHub.copilot from extensions (both auto-install from your local VS Code)
Removed go.lintTool setting (golangci-lint install is best-effort now)
Changed postCreateCommand to bash .devcontainer/post-create.sh (explicit bash invocation)
Try reopening in container now. If it still fails, the error log will be more specific since we've eliminated the most common Windows failure points. Let me know what you see!
