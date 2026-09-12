package models

import (
	"sort"
	"strings"
)

// City — un luogo prenotabile (aeroporto, stazione, porto o "metro" = tutti).
// Metro raggruppa i luoghi della stessa area urbana per la coerenza ricerca.
type City struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Kind    string `json:"kind"`    // "metro" | "airport" | "station" | "port"
	Country string `json:"country"`
	Metro   string `json:"metro"` // chiave area metropolitana, es "milano"
}

// AirportCodes — dizionario luoghi (aeroporti reali IATA, stazioni per DB API, porti UN/LOCODE).
var AirportCodes = []City{
	// ── ITALIA ──
	{"MIL", "Milano — tutti", "metro", "Italia", "milano"},
	{"LIN", "Milano Linate", "airport", "Italia", "milano"},
	{"MXP", "Milano Malpensa", "airport", "Italia", "milano"},
	{"BGY", "Milano Bergamo (Orio al Serio)", "airport", "Italia", "milano"},
	{"Milan", "Milano Centrale (stazione)", "station", "Italia", "milano"},
	{"ROM", "Roma — tutti", "metro", "Italia", "roma"},
	{"FCO", "Roma Fiumicino", "airport", "Italia", "roma"},
	{"CIA", "Roma Ciampino", "airport", "Italia", "roma"},
	{"Rome", "Roma Termini (stazione)", "station", "Italia", "roma"},
	{"NAP", "Napoli", "airport", "Italia", "napoli"},
	{"Naples", "Napoli Centrale (stazione)", "station", "Italia", "napoli"},
	{"ITNAP", "Napoli (porto)", "port", "Italia", "napoli"},
	{"VCE", "Venezia Marco Polo", "airport", "Italia", "venezia"},
	{"Venice", "Venezia S. Lucia (stazione)", "station", "Italia", "venezia"},
	{"ITVCE", "Venezia (porto)", "port", "Italia", "venezia"},
	{"FLR", "Firenze Peretola", "airport", "Italia", "firenze"},
	{"Florence", "Firenze S.M.N. (stazione)", "station", "Italia", "firenze"},
	{"BLQ", "Bologna Guglielmo Marconi", "airport", "Italia", "bologna"},
	{"Bologna", "Bologna Centrale (stazione)", "station", "Italia", "bologna"},
	{"TRN", "Torino Caselle", "airport", "Italia", "torino"},
	{"Turin", "Torino Porta Nuova (stazione)", "station", "Italia", "torino"},
	{"GOA", "Genova Cristoforo Colombo", "airport", "Italia", "genova"},
	{"ITGOA", "Genova (porto)", "port", "Italia", "genova"},
	{"CTA", "Catania Fontanarossa", "airport", "Italia", "catania"},
	{"PMO", "Palermo Falcone-Borsellino", "airport", "Italia", "palermo"},
	{"ITPMO", "Palermo (porto)", "port", "Italia", "palermo"},
	{"BRI", "Bari Karol Wojtyła", "airport", "Italia", "bari"},
	{"ITBRI", "Bari (porto)", "port", "Italia", "bari"},
	{"PSA", "Pisa Galilei", "airport", "Italia", "pisa"},
	{"CAG", "Cagliari Elmas", "airport", "Italia", "cagliari"},
	{"ITCAG", "Cagliari (porto)", "port", "Italia", "cagliari"},
	{"ITLIV", "Livorno (porto)", "port", "Italia", "livorno"},
	{"OLB", "Olbia Costa Smeralda", "airport", "Italia", "olbia"},
	{"VRN", "Verona Villafranca", "airport", "Italia", "verona"},
	{"Verona", "Verona Porta Nuova (stazione)", "station", "Italia", "verona"},
	// ── EUROPA ──
	{"PAR", "Parigi — tutti", "metro", "Francia", "parigi"},
	{"CDG", "Parigi Charles de Gaulle", "airport", "Francia", "parigi"},
	{"ORY", "Parigi Orly", "airport", "Francia", "parigi"},
	{"Paris", "Parigi (stazione)", "station", "Francia", "parigi"},
	{"LON", "Londra — tutti", "metro", "UK", "londra"},
	{"LHR", "Londra Heathrow", "airport", "UK", "londra"},
	{"LGW", "Londra Gatwick", "airport", "UK", "londra"},
	{"STN", "Londra Stansted", "airport", "UK", "londra"},
	{"LTN", "Londra Luton", "airport", "UK", "londra"},
	{"London", "Londra St Pancras (stazione)", "station", "UK", "londra"},
	{"BHX", "Birmingham (UK)", "airport", "UK", "birmingham-uk"},
	{"BHM", "Birmingham (USA)", "airport", "USA", "birmingham-us"},
	{"AMS", "Amsterdam Schiphol", "airport", "Olanda", "amsterdam"},
	{"Amsterdam", "Amsterdam Centraal (stazione)", "station", "Olanda", "amsterdam"},
	{"FRA", "Francoforte am Main", "airport", "Germania", "francoforte"},
	{"Frankfurt", "Francoforte (stazione)", "station", "Germania", "francoforte"},
	{"MUC", "Monaco di Baviera", "airport", "Germania", "monaco"},
	{"Munich", "Monaco (stazione)", "station", "Germania", "monaco"},
	{"BER", "Berlino Brandenburg", "airport", "Germania", "berlino"},
	{"Berlin", "Berlino (stazione)", "station", "Germania", "berlino"},
	{"HAM", "Amburgo", "airport", "Germania", "amburgo"},
	{"Hamburg", "Amburgo (stazione)", "station", "Germania", "amburgo"},
	{"MAD", "Madrid Barajas", "airport", "Spagna", "madrid"},
	{"Madrid", "Madrid Atocha (stazione)", "station", "Spagna", "madrid"},
	{"BCN", "Barcellona El Prat", "airport", "Spagna", "barcellona"},
	{"Barcelona", "Barcellona Sants (stazione)", "station", "Spagna", "barcellona"},
	{"ESBCN", "Barcellona (porto)", "port", "Spagna", "barcellona"},
	{"VLC", "Valencia (Spagna)", "airport", "Spagna", "valencia-es"},
	{"VLN", "Valencia (Venezuela)", "airport", "Venezuela", "valencia-ve"},
	{"SCQ", "Santiago di Compostela (Spagna)", "airport", "Spagna", "santiago-es"},
	{"SCL", "Santiago del Cile", "airport", "Cile", "santiago-cl"},
	{"SJO", "San José (Costa Rica)", "airport", "Costa Rica", "sanjose-cr"},
	{"SJC", "San Jose (California, USA)", "airport", "USA", "sanjose-us"},
	{"ZRH", "Zurigo", "airport", "Svizzera", "zurigo"},
	{"Zurich", "Zurigo (stazione)", "station", "Svizzera", "zurigo"},
	{"GVA", "Ginevra", "airport", "Svizzera", "ginevra"},
	{"BRU", "Bruxelles", "airport", "Belgio", "bruxelles"},
	{"Brussels", "Bruxelles Midi (stazione)", "station", "Belgio", "bruxelles"},
	{"VIE", "Vienna", "airport", "Austria", "vienna"},
	{"Vienna", "Vienna (stazione)", "station", "Austria", "vienna"},
	{"ATH", "Atene", "airport", "Grecia", "atene"},
	{"GRPIR", "Il Pireo (porto, Atene)", "port", "Grecia", "atene"},
	{"LIS", "Lisbona", "airport", "Portogallo", "lisbona"},
	{"CPH", "Copenaghen", "airport", "Danimarca", "copenaghen"},
	{"ARN", "Stoccolma Arlanda", "airport", "Svezia", "stoccolma"},
	{"HEL", "Helsinki", "airport", "Finlandia", "helsinki"},
	{"DUB", "Dublino", "airport", "Irlanda", "dublino"},
	{"PRG", "Praga", "airport", "Cechia", "praga"},
	{"Prague", "Praga (stazione)", "station", "Cechia", "praga"},
	{"WAW", "Varsavia", "airport", "Polonia", "varsavia"},
	{"BUD", "Budapest", "airport", "Ungheria", "budapest"},
	{"Budapest", "Budapest (stazione)", "station", "Ungheria", "budapest"},
	{"IST", "Istanbul", "airport", "Turchia", "istanbul"},
	{"REK", "Reykjavík Keflavík", "airport", "Islanda", "reykjavik"},
	// ── EXTRA-EUROPA ──
	{"NYC", "New York — tutti", "metro", "USA", "newyork"},
	{"JFK", "New York JFK", "airport", "USA", "newyork"},
	{"EWR", "New York Newark", "airport", "USA", "newyork"},
	{"LAX", "Los Angeles", "airport", "USA", "losangeles"},
	{"SFO", "San Francisco", "airport", "USA", "sanfrancisco"},
	{"MIA", "Miami", "airport", "USA", "miami"},
	{"DXB", "Dubai", "airport", "UAE", "dubai"},
	{"TYO", "Tokyo — tutti", "metro", "Giappone", "tokyo"},
	{"HND", "Tokyo Haneda", "airport", "Giappone", "tokyo"},
	{"NRT", "Tokyo Narita", "airport", "Giappone", "tokyo"},
	{"SYD", "Sydney", "airport", "Australia", "sydney"},
	{"SIN", "Singapore", "airport", "Singapore", "singapore"},
	{"HKG", "Hong Kong", "airport", "Cina", "hongkong"},
	{"BKK", "Bangkok", "airport", "Tailandia", "bangkok"},
	{"DEL", "Nuova Delhi", "airport", "India", "delhi"},
	{"BOM", "Mumbai", "airport", "India", "mumbai"},
	{"GRU", "São Paulo", "airport", "Brasile", "saopaulo"},
	{"EZE", "Buenos Aires", "airport", "Argentina", "buenosaires"},
	{"CPT", "Città del Capo", "airport", "Sudafrica", "capetown"},
	{"CAI", "Il Cairo", "airport", "Egitto", "cairo"},
	{"MEX", "Città del Messico", "airport", "Messico", "mexico"},
	{"YYZ", "Toronto", "airport", "Canada", "toronto"},
}

// PlaceQuery è la query di ricerca luoghi con filtri di coerenza.
type PlaceQuery struct {
	Query string // testo digitato
	Mode  string // "", "airport", "station", "port" — filtro esplicito
	From  string // codice luogo di partenza già scelto (per coerenza mezzo)
	Limit int
}

// cityKind mappa il codice alla sua City (indice interno).
var cityIndex = func() map[string]City {
	m := make(map[string]City, len(AirportCodes))
	for _, c := range AirportCodes {
		m[c.Code] = c
	}
	return m
}()

// KindOf ritorna il kind di un codice luogo ("" se sconosciuto).
func KindOf(code string) string {
	if c, ok := cityIndex[code]; ok {
		return c.Kind
	}
	return ""
}

// SearchPlaces ricerca con disambiguazione e coerenza.
// - From è aeroporto → destination propone solo aeroporti
// - From è stazione → solo stazioni; From è porto → solo porti
// - From è metro ("tutti") → tutto
// - Match "milano" → prima il cluster metro Milano (tutti + aeroporti + stazioni + porti),
//   poi gli altri luoghi che contengono il testo.
func SearchPlaces(q PlaceQuery) []City {
	limit := q.Limit
	if limit <= 0 {
		limit = 30
	}
	qq := strings.ToLower(strings.TrimSpace(q.Query))

	// coerenza: kind richiesto derivato da From
	wantedKind := q.Mode
	if wantedKind == "" && q.From != "" {
		switch KindOf(q.From) {
		case "airport":
			wantedKind = "airport"
		case "station":
			wantedKind = "station"
		case "port":
			wantedKind = "port"
		case "metro":
			wantedKind = "" // tutti i mezzi
		}
	}

	score := func(c City) int {
		n := strings.ToLower(c.Name)
		switch {
		case qq == "":
			return 2 // nessun testo: mostra top generale
		case strings.HasPrefix(n, qq):
			return 0
		case contains(n, qq):
			return 1
		case contains(strings.ToLower(c.Code), qq), contains(strings.ToLower(c.Country), qq), contains(strings.ToLower(c.Metro), qq):
			return 1
		default:
			return -1
		}
	}

	var exact, others []City
	var metroMatches []City
	for _, c := range AirportCodes {
		s := score(c)
		if s < 0 {
			continue
		}
		if wantedKind != "" && c.Kind != wantedKind && c.Kind != "metro" {
			continue // le voci metro restano sempre (es "Milano — tutti")
		}
		if wantedKind != "" && c.Kind == "metro" && q.Mode != "" {
			continue // filtro esplicito di tipo: niente cluster "tutti"
		}
		if s == 0 {
			exact = append(exact, c)
			if c.Metro != "" {
				metroMatches = append(metroMatches, c)
			}
		} else {
			others = append(others, c)
		}
	}

	// raggruppa: prima la metro della migliore corrispondenza esatta, poi il resto
	out := make([]City, 0, limit)
	seen := map[string]bool{}
	addIfNew := func(c City) {
		if !seen[c.Code] {
			seen[c.Code] = true
			out = append(out, c)
		}
	}

	metroKey := ""
	if len(exact) > 0 {
		metroKey = exact[0].Metro
		// tutte le voci della stessa metro (disambiguazione immediata)
		var cluster []City
		for _, c := range AirportCodes {
			if c.Metro == metroKey {
				cluster = append(cluster, c)
			}
		}
		sort.Slice(cluster, func(i, j int) bool {
			ord := map[string]int{"metro": 0, "airport": 1, "station": 2, "port": 3}
			if ord[cluster[i].Kind] != ord[cluster[j].Kind] {
				return ord[cluster[i].Kind] < ord[cluster[j].Kind]
			}
			return cluster[i].Name < cluster[j].Name
		})
		for _, c := range cluster {
			if wantedKind != "" && c.Kind != wantedKind && c.Kind != "metro" && q.Mode == "" {
				continue
			}
			if q.Mode != "" && c.Kind != wantedKind {
				continue
			}
			addIfNew(c)
		}
	}
	for _, c := range exact {
		addIfNew(c)
	}
	for _, c := range others {
		addIfNew(c)
	}
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

func contains(hay, needle string) bool {
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

// SearchCities — compatibilità con il vecchio formato (usato da server.go legacy).
func SearchCities(query string) []City {
	return SearchPlaces(PlaceQuery{Query: query, Limit: 30})
}