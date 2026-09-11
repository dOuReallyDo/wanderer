# 🌍 Wanderer

Ricerca ottimizzata di itinerari reali — volo e treno — con ranking intelligente.

## Caratteristiche

- **Multi-provider in parallelo**: Amadeus (voli), Kiwi/Tequila (low-cost), Deutsche Bahn (treni)
- **Wander Score**: ranking 0-100 che bilancia costo, durata, orario e "serendipity"
- **Filtri**: prezzo max, durata max, vincoli orario partenza, tipo trasporto
- **Web UI**: dark mode, mobile-first, embedded nel binary (zero dipendenze)
- **REST API**: `/api/search` per integrazione esterna
- **Binary singolo**: Go compilato, deploy su Mac mini o VPS con un comando

## Quick Start

```bash
# Build
go build -o bin/wanderer ./cmd/wanderer/

# Run (porta default 8899)
./bin/wanderer

# Con provider reali
export AMADEUS_API_KEY=xxx
export AMADEUS_API_SECRET=xxx
export TEQUILA_API_KEY=xxx
./bin/wanderer -port 8899
```

Apri `http://localhost:8899`.

## API

### Ricerca

```bash
# GET
curl "http://localhost:8899/api/search?origin=MIL&destination=PAR&date=2026-10-15&max_price=300&adults=1"

# POST
curl -X POST http://localhost:8899/api/search \
  -H "Content-Type: application/json" \
  -d '{"origin":"MIL","destination":"PAR","date":"2026-10-15","max_price":300,"adults":1}'
```

### Response

```json
{
  "count": 7,
  "trips": [
    {
      "rank": 1,
      "score": 92.6,
      "total_price": 47,
      "currency": "EUR",
      "total_duration_min": 283,
      "transport_types": ["train"],
      "segments": [...],
      "booking_url": "https://...",
      "source": "demo"
    }
  ],
  "best": { ... },
  "providers_queried": ["amadeus", "kiwi-tequila", "deutsche-bahn", "demo"],
  "took_ms": 1200
}
```

## Wander Score

Il Wander Score (0-100) combina 4 dimensioni con pesi configurabili:

| Dimensione | Peso default | Descrizione |
|---|---|---|
| **Cost** | 40% | Più basso è il prezzo, più alto il punteggio |
| **Duration** | 25% | Più breve è il viaggio, più alto il punteggio |
| **Schedule** | 20% | Partenze 8-20h punteggio massimo, vincoli orario rispettati |
| **Serendipity** | 15% | Offerte sorprendentemente buone (bottom 25% prezzo) |

## Provider

| Provider | Tipo | API Key | Note |
|---|---|---|---|
| **Amadeus** | Voli | `AMADEUS_API_KEY` + `AMADEUS_API_SECRET` | [developers.amadeus.com](https://developers.amadeus.com) |
| **Kiwi/Tequila** | Voli low-cost | `TEQUILA_API_KEY` | [tequila.kiwi.com](https://tequila.kiwi.com/portal) |
| **Deutsche Bahn** | Treni | Nessuna | API pubblica, può essere instabile |
| **Demo** | Simulato | Nessuna | Offerte realistiche per testing |

## Deploy

### Mac mini / VPS

```bash
# Build per la piattaforma target
GOOS=linux GOARCH=amd64 go build -o bin/wanderer-linux ./cmd/wanderer/
scp bin/wanderer-linux vps:/usr/local/bin/wanderer

# Run come servizio (systemd)
# o direttamente
./wanderer -port 8899
```

### Docker

```bash
docker build -t wanderer .
docker run -p 8899:8899 --env-file .env wanderer
```

## Stack

- **Go 1.26** — compilato, goroutine per ricerche parallele
- **Zero dipendenze esterne** — solo standard library
- **UI vanilla JS** — no build step, embedded nel binary
- **~15MB binary** — deploy ovunque

## License

MIT