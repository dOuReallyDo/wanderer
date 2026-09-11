package models

import "strings"

// AirportCodes — dizionario compatto dei principali aeroporti europei/globali.
// Codice IATA → nome città + paese.
var AirportCodes = []City{
	// Italia
	{"MIL", "Milano (tutti)", "airport", "Italia"},
	{"LIN", "Milano Linate", "airport", "Italia"},
	{"MXP", "Milano Malpensa", "airport", "Italia"},
	{"BGY", "Milano Bergamo", "airport", "Italia"},
	{"ROM", "Roma (tutti)", "airport", "Italia"},
	{"FCO", "Roma Fiumicino", "airport", "Italia"},
	{"CIA", "Roma Ciampino", "airport", "Italia"},
	{"NAP", "Napoli", "airport", "Italia"},
	{"VCE", "Venezia", "airport", "Italia"},
	{"FLR", "Firenze", "airport", "Italia"},
	{"BLQ", "Bologna", "airport", "Italia"},
	{"TRN", "Torino", "airport", "Italia"},
	{"CTA", "Catania", "airport", "Italia"},
	{"PMO", "Palermo", "airport", "Italia"},
	{"BRI", "Bari", "airport", "Italia"},
	{"GOA", "Genova", "airport", "Italia"},
	{"PSA", "Pisa", "airport", "Italia"},
	// Europa
	{"PAR", "Parigi (tutti)", "airport", "Francia"},
	{"CDG", "Parigi Charles de Gaulle", "airport", "Francia"},
	{"ORY", "Parigi Orly", "airport", "Francia"},
	{"LHR", "Londra Heathrow", "airport", "UK"},
	{"LGW", "Londra Gatwick", "airport", "UK"},
	{"STN", "Londra Stansted", "airport", "UK"},
	{"LTN", "Londra Luton", "airport", "UK"},
	{"AMS", "Amsterdam", "airport", "Olanda"},
	{"FRA", "Francoforte", "airport", "Germania"},
	{"MUC", "Monaco di Baviera", "airport", "Germania"},
	{"BER", "Berlino", "airport", "Germania"},
	{"MAD", "Madrid", "airport", "Spagna"},
	{"BCN", "Barcellona", "airport", "Spagna"},
	{"ZRH", "Zurigo", "airport", "Svizzera"},
	{"GVA", "Ginevra", "airport", "Svizzera"},
	{"BRU", "Bruxelles", "airport", "Belgio"},
	{"VIE", "Vienna", "airport", "Austria"},
	{"ATH", "Atene", "airport", "Grecia"},
	{"LIS", "Lisbona", "airport", "Portogallo"},
	{"CPH", "Copenaghen", "airport", "Danimarca"},
	{"ARN", "Stoccolma", "airport", "Svezia"},
	{"HEL", "Helsinki", "airport", "Finlandia"},
	{"DUB", "Dublino", "airport", "Irlanda"},
	{"PRG", "Praga", "airport", "Cechia"},
	{"WAW", "Varsavia", "airport", "Polonia"},
	{"BUD", "Budapest", "airport", "Ungheria"},
	{"IST", "Istanbul", "airport", "Turchia"},
	// Extra-Europa
	{"NYC", "New York (tutti)", "airport", "USA"},
	{"JFK", "New York JFK", "airport", "USA"},
	{"LAX", "Los Angeles", "airport", "USA"},
	{"SFO", "San Francisco", "airport", "USA"},
	{"DXB", "Dubai", "airport", "UAE"},
	{"TYO", "Tokyo (tutti)", "airport", "Giappone"},
	{"HND", "Tokyo Haneda", "airport", "Giappone"},
	{"NRT", "Tokyo Narita", "airport", "Giappone"},
	{"SYD", "Sydney", "airport", "Australia"},
	{"SIN", "Singapore", "airport", "Singapore"},
	{"HKG", "Hong Kong", "airport", "Cina"},
	{"BKK", "Bangkok", "airport", "Tailandia"},
	{"DEL", "Nuova Delhi", "airport", "India"},
	{"BOM", "Mumbai", "airport", "India"},
	{"GRU", "Sao Paulo", "airport", "Brasile"},
	{"EZE", "Buenos Aires", "airport", "Argentina"},
	{"CPT", "Città del Capo", "airport", "Sudafrica"},
	// Stazioni treni principali (Europa)
	{"Berlin", "Berlino (stazione)", "station", "Germania"},
	{"Munich", "Monaco (stazione)", "station", "Germania"},
	{"Hamburg", "Amburgo (stazione)", "station", "Germania"},
	{"Frankfurt", "Francoforte (stazione)", "station", "Germania"},
	{"Paris", "Parigi (stazione)", "station", "Francia"},
	{"London", "Londra (stazione)", "station", "UK"},
	{"Amsterdam", "Amsterdam (stazione)", "station", "Olanda"},
	{"Zurich", "Zurigo (stazione)", "station", "Svizzera"},
	{"Vienna", "Vienna (stazione)", "station", "Austria"},
	{"Milan", "Milano (stazione)", "station", "Italia"},
	{"Rome", "Roma (stazione)", "station", "Italia"},
	{"Venice", "Venezia (stazione)", "station", "Italia"},
	{"Bologna", "Bologna (stazione)", "station", "Italia"},
	{"Florence", "Firenze (stazione)", "station", "Italia"},
	{"Naples", "Napoli (stazione)", "station", "Italia"},
	{"Madrid", "Madrid (stazione)", "station", "Spagna"},
	{"Barcelona", "Barcellona (stazione)", "station", "Spagna"},
	{"Brussels", "Bruxelles (stazione)", "station", "Belgio"},
	{"Prague", "Praga (stazione)", "station", "Cechia"},
	{"Budapest", "Budapest (stazione)", "station", "Ungheria"},
}

// SearchCities filtra per query (case-insensitive, match su code o name).
func SearchCities(query string) []City {
	if query == "" {
		return AirportCodes[:20] // prime 20 se vuoto
	}
	var out []City
	q := strings.ToLower(query)
	for _, c := range AirportCodes {
		if strings.Contains(strings.ToLower(c.Code), q) ||
			strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.Country), q) {
			out = append(out, c)
		}
	}
	if len(out) > 30 {
		out = out[:30]
	}
	return out
}