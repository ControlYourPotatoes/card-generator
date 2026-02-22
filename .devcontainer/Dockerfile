FROM mcr.microsoft.com/devcontainers/go:1-1.23-bookworm

# The base image sometimes includes an APT source for Yarn (dl.yarnpkg.com).
# If the signing key isn't present, ANY `apt-get update` during feature install
# fails (including Docker features). Since this project uses npm, we remove the
# Yarn APT source to keep builds deterministic.
RUN set -eux; \
    if [ -d /etc/apt/sources.list.d ]; then \
      for f in /etc/apt/sources.list.d/*.list; do \
        [ -e "$f" ] || continue; \
        if grep -q "dl.yarnpkg.com" "$f"; then \
          rm -f "$f"; \
        fi; \
      done; \
    fi
