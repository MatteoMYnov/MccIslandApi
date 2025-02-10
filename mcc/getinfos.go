package mcc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type GraphQLRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables,omitempty"`
}

type Statistics struct {
	TotalGames   int `json:"total_games"`
	SBGames      int `json:"sb_games"`
	BBGames      int `json:"bb_games"`
	TGTTOSGames  int `json:"tgttos_games"`
	HITWGames    int `json:"hitw_games"`
	RSGames      int `json:"rs_games"`
	DBGames      int `json:"db_games"`
	PWSGames     int `json:"pws_games"`
	SB_Wins      int `json:"sb_wins"`
	SB_Loses     int
	SB_WLR       float32
	SB_Kills     int `json:"sb_kills"`
	SB_Subdeaths int `json:"sb_subdeaths"`
	SB_Deaths    int
	SB_KDR       float32
	BB_Wins      int `json:"bb_wins"`
	BB_Loses     int
	BB_WLR       float32
	BB_Kills     int `json:"bb_kills"`
	BB_Deaths    int `json:"bb_deaths"`
	BB_KDR       float32
	RS_Wins      int `json:"rs_wins"`
	RS_Loses     int
	RS_WLR       float32
	RS_Kills     int `json:"rs_kills"`
	RS_Deaths    int `json:"rs_deaths"`
	RS_KDR       float32
	DB_Wins      int `json:"db_wins"`
	DB_Loses     int
	DB_WLR       float32
	DB_Kills     int `json:"db_kills"`
	DB_Deaths    int
	DB_KDR       float32
}

type NextLevelProgress struct {
	CrownObtained   int `json:"obtained"`
	CrownObtainable int `json:"obtainable"`
}

type Trophies struct {
	Obtained   int `json:"obtained"`
	Obtainable int `json:"obtainable"`
	Bonus      int `json:"bonus,omitempty"`
}

type FishingData struct {
	Level             int               `json:"level"`
	NextLevelProgress NextLevelProgress `json:"nextLevelProgress"`
}

type CrownLevel struct {
	LevelData struct {
		Level             int               `json:"level"`
		Evolution         int               `json:"evolution"`
		NextLevelProgress NextLevelProgress `json:"nextLevelProgress"`
	} `json:"levelData"`
	FishingData    FishingData `json:"fishingLevelData"`
	Trophies       Trophies    `json:"trophies"`
	TrophiesSKILL  Trophies    `json:"trophiesSKILL,omitempty"`
	TrophiesSTYLE  Trophies    `json:"trophiesSTYLE,omitempty"`
	TrophiesANGLER Trophies    `json:"trophiesANGLER,omitempty"`
}

type Currency struct {
	Coins           int `json:"coins"`
	RoyalReputation int `json:"royalReputation"`
	Silver          int `json:"silver"`
	MaterialDust    int `json:"materialDust"`
	AnglrTokens     int `json:"anglrTokens"`
}

type Friend struct {
	Username string   `json:"username"`
	Ranks    []string `json:"ranks"`
	Status   struct {
		Online bool `json:"online"`
	} `json:"status"`
	CrownLevel struct {
		LevelData struct {
			Level     int `json:"level"`
			Evolution int `json:"evolution"`
		} `json:"levelData"`
	} `json:"crownLevel"`
}

type Player struct {
	Statistics  Statistics `json:"statistics"`
	UUID        string     `json:"uuid"`
	Username    string     `json:"username"`
	Ranks       []string   `json:"ranks"`
	CrownLevel  CrownLevel `json:"crownLevel"`
	Collections struct {
		Currency Currency `json:"currency"`
	} `json:"collections"`
	Social struct {
		Friends []Friend `json:"friends"`
	} `json:"social"`
}

type Response struct {
	Data struct {
		Player Player `json:"player"`
	} `json:"data"`
}

type APIConfig struct {
	Mcctoken string `json:"mcctoken"`
}

type MccInfos struct {
	Statistics      Statistics  `json:"statistics"`
	Ranks           []string    `json:"ranks"`
	CrownLevel      int         `json:"crownLevel"`
	Evolution       int         `json:"evolution"`
	CrownObtained   int         `json:"crownObtained"`
	CrownObtainable int         `json:"crownObtainable"`
	FishingData     FishingData `json:"fishingData"`
	Currency        Currency    `json:"currency"`
	Trophies        Trophies    `json:"trophies"`
	TrophiesSKILL   Trophies    `json:"trophiesSKILL"`
	TrophiesSTYLE   Trophies    `json:"trophiesSTYLE"`
	TrophiesANGLER  Trophies    `json:"trophiesANGLER"`
	Friends         []Friend    `json:"friends"` // Liste des amis avec leurs informations
}

func GetInfos(UUID string) *MccInfos {
	// Lire le fichier contenant le token de l'API
	file, err := ioutil.ReadFile("./site/infos/api.json")
	if err != nil {
		log.Printf("Erreur lors de la lecture du fichier api.json: %v", err)
		return nil
	}

	var apiConfig APIConfig
	err = json.Unmarshal(file, &apiConfig)
	if err != nil {
		log.Printf("Erreur lors du décodage du JSON: %v", err)
		return nil
	}

	if apiConfig.Mcctoken == "" {
		log.Println("Le token API est vide")
		return nil
	}

	// Définir la requête GraphQL
	query := `
		query player($uuid: UUID!) {	
			player(uuid: $uuid) {
				statistics {
					total_games: rotationValue(statisticKey: "games_played")
					sb_games: rotationValue(statisticKey: "sky_battle_quads_games_played")
					bb_games: rotationValue(statisticKey: "battle_box_quads_games_played")
					tgttos_games: rotationValue(statisticKey: "tgttos_games_played")
					hitw_games: rotationValue(statisticKey: "hole_in_the_wall_games_played")
					rs_games: rotationValue(statisticKey: "rocket_spleef_games_played")
					db_games: rotationValue(statisticKey: "dynaball_games_played")
					pws_games: rotationValue(statisticKey: "pw_survival_games_played")
					sb_wins:rotationValue(statisticKey: "sky_battle_quads_team_placement_1")
					sb_kills: rotationValue(statisticKey: "sky_battle_quads_players_killed")
					sb_subdeaths: rotationValue(statisticKey: "sky_battle_quads_survival_first_place")
					bb_wins:rotationValue(statisticKey: "battle_box_quads_team_first_place")
					bb_kills: rotationValue(statisticKey: "battle_box_quads_players_killed")
					bb_deaths: rotationValue(statisticKey: "battle_box_quads_times_eliminated")
					db_wins: rotationValue(statisticKey: "dynaball_wins")
					db_kills: rotationValue(statisticKey: "dynaball_players_eliminated")
					rs_wins: rotationValue(statisticKey: "rocket_spleef_first_place")
					rs_kills: rotationValue(statisticKey: "rocket_spleef_kills")
					rs_deaths: rotationValue(statisticKey: "rocket_spleef_deaths")
					}
				ranks
				crownLevel {
					levelData {
						level
						evolution
						nextLevelProgress {
							obtained
							obtainable
						}
					}
					fishingLevelData {
						level
						nextLevelProgress {
							obtained
							obtainable
						}
					}
					trophies {
						obtained
						obtainable
						bonus
					}
					trophiesSKILL: trophies(category: SKILL) { 
						obtained
						obtainable
					}
					trophiesSTYLE: trophies(category: STYLE) { 
						obtained
						obtainable
					}
					trophiesANGLER: trophies(category: ANGLER) { 
						obtained
						obtainable
					}
				}
				collections {
					currency {
						coins
						royalReputation
						silver
						materialDust
						anglrTokens
					}
				}
				social {
					friends {
						username
						ranks
						status {
							online
							}
						crownLevel {
							levelData {
								level
								evolution
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]string{
		"uuid": UUID,
	}

	requestBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Erreur lors de la conversion de la requête en JSON: %v", err)
		return nil
	}

	apiURL := "https://api.mccisland.net/graphql"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Erreur lors de la création de la requête: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiConfig.Mcctoken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la requête: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Erreur: Code de statut reçu %d", resp.StatusCode)
		return nil
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Erreur lors du décodage de la réponse: %v", err)
		return nil
	}

	// Mapper la réponse aux informations de MCC
	Infos := &MccInfos{
		Statistics: Statistics{
			TotalGames:  response.Data.Player.Statistics.TotalGames,
			SBGames:     response.Data.Player.Statistics.SBGames,
			BBGames:     response.Data.Player.Statistics.BBGames,
			TGTTOSGames: response.Data.Player.Statistics.TGTTOSGames,
			HITWGames:   response.Data.Player.Statistics.HITWGames,
			RSGames:     response.Data.Player.Statistics.RSGames,
			DBGames:     response.Data.Player.Statistics.DBGames,
			PWSGames:    response.Data.Player.Statistics.PWSGames,
			//Sky Battle
			SB_Wins:   response.Data.Player.Statistics.SB_Wins,
			SB_Loses:  response.Data.Player.Statistics.SBGames - response.Data.Player.Statistics.SB_Wins,
			SB_WLR:    safeDivide(response.Data.Player.Statistics.SB_Wins, response.Data.Player.Statistics.SBGames-response.Data.Player.Statistics.SB_Wins),
			SB_Kills:  response.Data.Player.Statistics.SB_Kills,
			SB_Deaths: response.Data.Player.Statistics.SBGames - response.Data.Player.Statistics.SB_Subdeaths,
			SB_KDR:    safeDivide(response.Data.Player.Statistics.SB_Kills, response.Data.Player.Statistics.SBGames-response.Data.Player.Statistics.SB_Subdeaths),
			//Battle Box
			BB_Wins:   response.Data.Player.Statistics.BB_Wins,
			BB_Loses:  response.Data.Player.Statistics.BBGames - response.Data.Player.Statistics.BB_Wins,
			BB_WLR:    safeDivide(response.Data.Player.Statistics.BB_Wins, response.Data.Player.Statistics.BBGames-response.Data.Player.Statistics.BB_Wins),
			BB_Kills:  response.Data.Player.Statistics.BB_Kills,
			BB_Deaths: response.Data.Player.Statistics.BB_Deaths,
			BB_KDR:    safeDivide(response.Data.Player.Statistics.BB_Kills, response.Data.Player.Statistics.BB_Deaths),
			//Rocket Spleef Rush
			RS_Wins:   response.Data.Player.Statistics.RS_Wins,
			RS_Loses:  response.Data.Player.Statistics.RSGames - response.Data.Player.Statistics.RS_Wins,
			RS_WLR:    safeDivide(response.Data.Player.Statistics.RS_Wins, response.Data.Player.Statistics.RSGames-response.Data.Player.Statistics.RS_Wins),
			RS_Kills:  response.Data.Player.Statistics.RS_Kills,
			RS_Deaths: response.Data.Player.Statistics.RS_Deaths,
			RS_KDR:    safeDivide(response.Data.Player.Statistics.RS_Kills, response.Data.Player.Statistics.RS_Deaths),
			//DynaBall
			DB_Wins:   response.Data.Player.Statistics.DB_Wins,
			DB_Loses:  response.Data.Player.Statistics.DBGames - response.Data.Player.Statistics.DB_Wins,
			DB_WLR:    safeDivide(response.Data.Player.Statistics.DB_Wins, response.Data.Player.Statistics.DBGames-response.Data.Player.Statistics.DB_Wins),
			DB_Kills:  response.Data.Player.Statistics.DB_Kills,
			DB_Deaths: response.Data.Player.Statistics.DBGames - response.Data.Player.Statistics.DB_Wins,
			DB_KDR:    safeDivide(response.Data.Player.Statistics.DB_Kills, response.Data.Player.Statistics.DBGames-response.Data.Player.Statistics.DB_Wins),
		},
		Ranks:           response.Data.Player.Ranks,
		CrownLevel:      response.Data.Player.CrownLevel.LevelData.Level,
		Evolution:       response.Data.Player.CrownLevel.LevelData.Evolution,
		CrownObtained:   response.Data.Player.CrownLevel.LevelData.NextLevelProgress.CrownObtained,
		CrownObtainable: response.Data.Player.CrownLevel.LevelData.NextLevelProgress.CrownObtainable,
		FishingData: FishingData{
			Level:             response.Data.Player.CrownLevel.FishingData.Level,
			NextLevelProgress: response.Data.Player.CrownLevel.FishingData.NextLevelProgress,
		},
		Currency:       response.Data.Player.Collections.Currency,
		Trophies:       response.Data.Player.CrownLevel.Trophies,
		TrophiesSKILL:  response.Data.Player.CrownLevel.TrophiesSKILL,
		TrophiesSTYLE:  response.Data.Player.CrownLevel.TrophiesSTYLE,
		TrophiesANGLER: response.Data.Player.CrownLevel.TrophiesANGLER,
		Friends:        response.Data.Player.Social.Friends, // Liste des amis avec leurs infos
	}
	return Infos
}
