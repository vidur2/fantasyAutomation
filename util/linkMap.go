package util

var LinkMap = map[string]string{
	"players":        "https://api.sleeper.app/v1/players/nfl",
	"rosters":        "https://api.sleeper.app/v1/league/%v/rosters",
	"player_leagues": "https://api.sleeper.app/v1/user/%v/leagues/nfl/2022",
	"league":         "https://api.sleeper.app/v1/league/%v",
	"player_id":      "https://api.sleeper.app/v1/user/%v",
	"player_stats":   "https://www.pro-football-reference.com/teams/%v/2022.htm",
	"transactions":   "https://api.sleeper.app/v1/league/%v/transactions/%v",
}

var RUSH_RECV_STATS = map[string]int{"targets": 1, "rec": 1, "rec_yds": 1, "rec_td": 1, "rec_drops": -1, "rec_pass_rating": 1, "rush_att": 1, "rush_yds": 1}
var PASS_STATS = map[string]int{"pass_cmp": 1, "pass_att": 1, "pass_yds": 1, "pass_drops": -1, "pass_drop_pct": -1, "pass_poor_throw_pct": -1, "pass_sacked": -1, "pass_hits": -1}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// var TeamArr = map[string]string{
// 	"BUF": "buf",
// 	"NE":  "nwe",
// 	"MIA": "mia",
// 	"NYJ": "nyj",
// 	"TEN": "oti",
// 	"IND": "clt",
// 	"HOU": "htx",
// 	"JAX": "jax",
// 	"CIN": "cin",
// 	"PIT": "pit",
// 	"CLE": "cle",
// 	"BAL": "rav",
// 	"KC":  "kan",
// 	"LV":  "rai",
// 	"LAC": "sdg",
// 	"DEN": "den",
// 	"DAL": "dal",
// 	"PHI": "phi",
// 	"WAS": "was",
// 	"NYG": "nyg",
// 	"TB":  "tam",
// 	"NO":  "nor",
// 	"ATL": "atl",
// 	"CAR": "car",
// 	"GB":  "gnb",
// 	"MIN": "min",
// 	"CHI": "chi",
// 	"DET": "det",
// 	"LAR": "lar",
// 	"ARI": "crd",
// 	"SF":  "sfo",
// 	"SEA": "sea",
// }
