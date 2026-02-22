#!/bin/bash
# Post-create script for dev container setup
# Non-fatal: individual steps can fail without blocking the container

echo "🔧 Setting up Card Generator development environment..."

# ── Go backend setup ──
echo "📦 Downloading Go modules..."
if cd /workspaces/card-generator/backend 2>/dev/null || cd /workspace/backend 2>/dev/null; then
  go mod download && echo "✅ Go modules ready" || echo "⚠️  Go mod download failed (non-fatal)"
else
  echo "⚠️  Backend directory not found, skipping Go setup"
fi

# Install Go development tools
echo "🔧 Installing Go tools..."
go install golang.org/x/tools/gopls@latest 2>/dev/null || true
go install github.com/go-delve/delve/cmd/dlv@latest 2>/dev/null || true
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest 2>/dev/null || true
echo "✅ Go tools installed"

# ── Frontend setup ──
echo "📦 Installing frontend dependencies..."
if cd /workspaces/card-generator/frontend 2>/dev/null || cd /workspace/frontend 2>/dev/null; then
  npm install && echo "✅ Frontend dependencies installed" || echo "⚠️  npm install failed (non-fatal)"
  # Generate Prisma client
  npx prisma generate 2>/dev/null && echo "✅ Prisma client ready" || echo "⚠️  Prisma generate skipped"
else
  echo "⚠️  Frontend directory not found, skipping Node setup"
fi

# ── Utilities ──
echo "🔧 Installing additional utilities..."
sudo apt-get update -qq && sudo apt-get install -y -qq jq entr postgresql-client > /dev/null 2>&1 || echo "⚠️  Some utilities failed to install"
echo "✅ Utilities installed"

# ── Summary ──
echo ""
echo "════════════════════════════════════════════════════════"
echo "  ✅ Dev container ready!"
echo ""
echo "  Quick start:"
echo "    Backend:  cd backend && go run ./services/api-gateway/"
echo "    Frontend: cd frontend && npm run dev"
echo ""
echo "  Infrastructure (run from host or inside container):"
echo "    docker compose --env-file .env.docker up -d postgres adminer"
echo ""
echo "  Ports:"
echo "    3000 → Next.js dev server"
echo "    8080 → API Gateway"
echo "    5433 → PostgreSQL"
echo "    8081 → Adminer (DB admin)"
echo "════════════════════════════════════════════════════════"
