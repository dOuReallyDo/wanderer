# Wanderer install script — Mac mini
# launchd non può leggere binary da /Volumes (TCC restriction)
# Binary deve stare in ~/.local/bin/

set -e

BINARY_SRC=/Volumes/HD_esterno/Progetti/Wanderer/bin/wanderer
BINARY_DST="$HOME/.local/bin/wanderer"
WORKDIR="$HOME/.local/share/wanderer"
PLIST_SRC=/Volumes/HD_esterno/Progetti/Wanderer/deploy/com.wanderer.app.plist
PLIST_DST="$HOME/Library/LaunchAgents/com.wanderer.app.plist"

echo "📦 Installazione Wanderer service..."

# 1. Copia binary in home (launchd non legge /Volumes)
mkdir -p "$HOME/.local/bin"
cp "$BINARY_SRC" "$BINARY_DST"
chmod +x "$BINARY_DST"

# 2. Crea working directory
mkdir -p "$WORKDIR"

# 3. Copia plist
cp "$PLIST_SRC" "$PLIST_DST"

# 4. Sostituisci env vars se .env esiste
if [ -f /Volumes/HD_esterno/Progetti/Wanderer/.env ]; then
    echo "📝 Caricando .env..."
    set -a
    source /Volumes/HD_esterno/Progetti/Wanderer/.env
    set +a
    sed -i '' "s|\${AMADEUS_API_KEY}|${AMADEUS_API_KEY:-}|g" "$PLIST_DST"
    sed -i '' "s|\${AMADEUS_API_SECRET}|${AMADEUS_API_SECRET:-}|g" "$PLIST_DST"
    sed -i '' "s|\${TEQUILA_API_KEY}|${TEQUILA_API_KEY:-}|g" "$PLIST_DST"
fi

# 5. Unload + load
launchctl unload "$PLIST_DST" 2>/dev/null || true
launchctl load "$PLIST_DST"

sleep 2

# 6. Verify
if curl -s --max-time 3 http://localhost:8899/api/health | grep -q "ok"; then
    echo "✅ Wanderer attivo su http://localhost:8899"
else
    echo "⚠️  Verifica: tail /tmp/wanderer.err"
fi