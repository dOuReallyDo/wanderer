#!/bin/bash
# Installa Wanderer come servizio launchd su Mac mini
set -e

PLIST=~/Library/LaunchAgents/com.wanderer.app.plist
SOURCE=/Volumes/HD_esterno/Progetti/Wanderer/deploy/com.wanderer.app.plist

echo "📦 Installazione Wanderer service..."

# Copia plist
cp "$SOURCE" "$PLIST"

# Sostituisci env vars placeholder con valori reali se presenti
if [ -f /Volumes/HD_esterno/Progetti/Wanderer/.env ]; then
    echo "📝 Caricando .env..."
    set -a
    source /Volumes/HD_esterno/Progetti/Wanderer/.env
    set +a
    # Rigenera plist con valori reali
    sed -i '' "s|\${AMADEUS_API_KEY}|${AMADEUS_API_KEY:-}|g" "$PLIST"
    sed -i '' "s|\${AMADEUS_API_SECRET}|${AMADEUS_API_SECRET:-}|g" "$PLIST"
    sed -i '' "s|\${TEQUILA_API_KEY}|${TEQUILA_API_KEY:-}|g" "$PLIST"
fi

# Load
launchctl unload "$PLIST" 2>/dev/null || true
launchctl load "$PLIST"

sleep 2

# Verify
if curl -s --max-time 3 http://localhost:8899/api/health | grep -q "ok"; then
    echo "✅ Wanderer attivo su http://localhost:8899"
else
    echo "⚠️  Verifica: tail /tmp/wanderer.err"
fi