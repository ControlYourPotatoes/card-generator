#!/bin/bash
set -e

echo "🔧 Setting up Card Generator development environment..."

# ── Go backend setup ──
echo "📦 Downloading Go modules..."
cd /workspace/backend
go mod download
echo "✅ Go modules ready"

# Install Go development tools
echo "🔧 Installing Go tools..."
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
echo "✅ Go tools installed"

# ── Frontend setup ──
echo "📦 Installing frontend dependencies..."
cd /workspace/frontend
# Use npm since the docker-compose uses npm (package.json declares yarn but npm works fine)
npm install
echo "✅ Frontend dependencies installed"

# Generate Prisma client
echo "🔧 Generating Prisma client..."
npx prisma generate || echo "⚠️  Prisma generate skipped (DB may not be running yet)"
echo "✅ Prisma client ready"

# ── Utilities ──
echo "🔧 Installing additional utilities..."
sudo apt-get update -qq && sudo apt-get install -y -qq jq entr postgresql-client > /dev/null 2>&1
echo "✅ Utilities installed (jq, entr, psql)"

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
