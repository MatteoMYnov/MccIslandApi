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
	TotalGames       int `json:"total_games"`
	SBGames          int `json:"sb_games"`
	BBGames          int `json:"bb_games"`
	TGTTOSGames      int `json:"tgttos_games"`
	HITWGames        int `json:"hitw_games"`
	RSGames          int `json:"rs_games"`
	DBGames          int `json:"db_games"`
	PWSGames         int `json:"pws_games"`
	SB_Wins          int `json:"sb_wins"`
	SB_Loses         int
	SB_WLR           float32
	SB_Kills         int `json:"sb_kills"`
	SB_Subdeaths     int `json:"sb_subdeaths"`
	SB_Deaths        int
	SB_KDR           float32
	BB_Wins          int `json:"bb_wins"`
	BB_Loses         int
	BB_WLR           float32
	BB_Kills         int `json:"bb_kills"`
	BB_Deaths        int `json:"bb_deaths"`
	BB_KDR           float32
	TGTTOS_Wins      int `json:"tgttos_wins"`
	TGTTOS_Loses     int
	TGTTOS_WLR       float32
	TGTTOS_Chicks    int `json:"tgttos_chicks"`
	TGTTOS_CGR       float32
	HITW_Wins        int `json:"hitw_wins"`
	HITW_Loses       int
	HITW_WLR         float32
	HITW_Wallsdodged int `json:"hitw_wallsdodged"`
	HITW_WGR         float32
	RS_Wins          int `json:"rs_wins"`
	RS_Loses         int
	RS_WLR           float32
	RS_Kills         int `json:"rs_kills"`
	RS_Deaths        int `json:"rs_deaths"`
	RS_KDR           float32
	DB_Wins          int `json:"db_wins"`
	DB_Loses         int
	DB_WLR           float32
	DB_Kills         int `json:"db_kills"`
	DB_Deaths        int
	DB_KDR           float32
}

type EquippedCosmetic struct {
	Category    string `json:"category"`
	Name        string `json:"name"`
	Rarity      string `json:"rarity"`
	Description string `json:"description"`
}

type InvCos struct {
	Owned           bool
	Name            string
	RealName        string
	Rarity          string
	IsBonusTrophies bool
	Trophies        int
	Description     string
}

type CosmeticInfos struct {
	Name            string `json:"name"`
	Rarity          string `json:"rarity"`
	IsBonusTrophies bool   `json:"isBonusTrophies"`
	Trophies        int    `json:"trophies"`
	Description     string `json:"description"`
}

type Cosmetic struct {
	Owned    bool          `json:"owned"`
	Cosmetic CosmeticInfos `json:"cosmetic"`
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
		Currency    Currency           `json:"currency"`
		Equipped    []EquippedCosmetic `json:"equippedCosmetics"`
		Hats        []Cosmetic         `json:"hats"`
		Accessories []Cosmetic         `json:"accessories"`
		Auras       []Cosmetic         `json:"auras"`
		Trails      []Cosmetic         `json:"trails"`
		Cloaks      []Cosmetic         `json:"cloaks"`
		Rods        []Cosmetic         `json:"rods"`
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
	Statistics        Statistics         `json:"statistics"`
	Ranks             []string           `json:"ranks"`
	CrownLevel        int                `json:"crownLevel"`
	Evolution         int                `json:"evolution"`
	CrownObtained     int                `json:"crownObtained"`
	CrownObtainable   int                `json:"crownObtainable"`
	FishingData       FishingData        `json:"fishingData"`
	Currency          Currency           `json:"currency"`
	Trophies          Trophies           `json:"trophies"`
	TrophiesSKILL     Trophies           `json:"trophiesSKILL"`
	TrophiesSTYLE     Trophies           `json:"trophiesSTYLE"`
	TrophiesANGLER    Trophies           `json:"trophiesANGLER"`
	Friends           []Friend           `json:"friends"`
	EquippedCosmetics []EquippedCosmetic `json:"equippedCosmetics"`
	Hats              []InvCos
	Accessories       []InvCos
	Auras             []InvCos
	Trails            []InvCos
	Cloaks            []InvCos
	Rods              []InvCos
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
					tgttos_wins: rotationValue(statisticKey: "tgttos_first_place")
					tgttos_chicks: rotationValue(statisticKey: "tgttos_chickens_punched")
					hitw_wins: rotationValue(statisticKey: "hole_in_the_wall_first_place")
      				hitw_wallsdodged: rotationValue(statisticKey: "hole_in_the_wall_walls_dodged")
					rs_wins: rotationValue(statisticKey: "rocket_spleef_first_place")
					rs_kills: rotationValue(statisticKey: "rocket_spleef_kills")
					rs_deaths: rotationValue(statisticKey: "rocket_spleef_deaths")
					db_wins: rotationValue(statisticKey: "dynaball_wins")
					db_kills: rotationValue(statisticKey: "dynaball_players_eliminated")
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
						materialDust
						royalReputation
						anglrTokens
						silver
					}
					equippedCosmetics {
						category
						name
						rarity
						description
					}
					hats: cosmetics(category: HAT) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
					}
					accessories: cosmetics(category: ACCESSORY) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
					}
					auras: cosmetics(category: AURA) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
					}
					trails: cosmetics(category: TRAIL) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
					}
					cloaks: cosmetics(category: CLOAK) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
					}
					rods: cosmetics(category: ROD) {
						owned
						cosmetic {
							rarity 
							name
							isBonusTrophies
							trophies
							description
						}
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

	equippedCosmetics := []EquippedCosmetic{}
	for _, cosmetic := range response.Data.Player.Collections.Equipped {
		equippedCosmetics = append(equippedCosmetics, EquippedCosmetic{
			Category:    cosmetic.Category,
			Name:        cosmetic.Name,
			Rarity:      cosmetic.Rarity,
			Description: cosmetic.Description,
		})
	}

	hats := []InvCos{}
	accessories := []InvCos{}
	auras := []InvCos{}
	trails := []InvCos{}
	cloaks := []InvCos{}
	rods := []InvCos{}
	for _, hat := range response.Data.Player.Collections.Hats {
		hats = append(hats, InvCos{
			Owned:           hat.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(hat.Cosmetic.Name),
			RealName:        hat.Cosmetic.Name,
			Rarity:          hat.Cosmetic.Rarity,
			IsBonusTrophies: hat.Cosmetic.IsBonusTrophies,
			Trophies:        hat.Cosmetic.Trophies,
			Description:     hat.Cosmetic.Description,
		})
	}
	for _, accessory := range response.Data.Player.Collections.Accessories {
		accessories = append(accessories, InvCos{
			Owned:           accessory.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(accessory.Cosmetic.Name),
			RealName:        accessory.Cosmetic.Name,
			Rarity:          accessory.Cosmetic.Rarity,
			IsBonusTrophies: accessory.Cosmetic.IsBonusTrophies,
			Trophies:        accessory.Cosmetic.Trophies,
			Description:     accessory.Cosmetic.Description,
		})
	}
	for _, aura := range response.Data.Player.Collections.Auras {
		auras = append(auras, InvCos{
			Owned:           aura.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(aura.Cosmetic.Name),
			RealName:        aura.Cosmetic.Name,
			Rarity:          aura.Cosmetic.Rarity,
			IsBonusTrophies: aura.Cosmetic.IsBonusTrophies,
			Trophies:        aura.Cosmetic.Trophies,
			Description:     aura.Cosmetic.Description,
		})
	}
	for _, trail := range response.Data.Player.Collections.Trails {
		trails = append(trails, InvCos{
			Owned:           trail.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(trail.Cosmetic.Name),
			RealName:        trail.Cosmetic.Name,
			Rarity:          trail.Cosmetic.Rarity,
			IsBonusTrophies: trail.Cosmetic.IsBonusTrophies,
			Trophies:        trail.Cosmetic.Trophies,
			Description:     trail.Cosmetic.Description,
		})
	}
	for _, cloak := range response.Data.Player.Collections.Cloaks {
		cloaks = append(cloaks, InvCos{
			Owned:           cloak.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(cloak.Cosmetic.Name),
			RealName:        cloak.Cosmetic.Name,
			Rarity:          cloak.Cosmetic.Rarity,
			IsBonusTrophies: cloak.Cosmetic.IsBonusTrophies,
			Trophies:        cloak.Cosmetic.Trophies,
			Description:     cloak.Cosmetic.Description,
		})
	}
	for _, rod := range response.Data.Player.Collections.Rods {
		rods = append(rods, InvCos{
			Owned:           rod.Owned, // Ajout du champ Owned
			Name:            CleanCosmeticName(rod.Cosmetic.Name),
			RealName:        rod.Cosmetic.Name,
			Rarity:          rod.Cosmetic.Rarity,
			IsBonusTrophies: rod.Cosmetic.IsBonusTrophies,
			Trophies:        rod.Cosmetic.Trophies,
			Description:     rod.Cosmetic.Description,
		})
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
			//To Go To The Other Side
			TGTTOS_Wins:   response.Data.Player.Statistics.TGTTOS_Wins,
			TGTTOS_Loses:  response.Data.Player.Statistics.TGTTOSGames - response.Data.Player.Statistics.TGTTOS_Wins,
			TGTTOS_WLR:    safeDivide(response.Data.Player.Statistics.TGTTOS_Wins, response.Data.Player.Statistics.TGTTOSGames-response.Data.Player.Statistics.TGTTOS_Wins),
			TGTTOS_Chicks: response.Data.Player.Statistics.TGTTOS_Chicks,
			TGTTOS_CGR:    safeDivide(response.Data.Player.Statistics.TGTTOS_Chicks, response.Data.Player.Statistics.TGTTOSGames),
			//Hole In The Wall
			HITW_Wins:        response.Data.Player.Statistics.HITW_Wins,
			HITW_Loses:       response.Data.Player.Statistics.HITWGames - response.Data.Player.Statistics.HITW_Wins,
			HITW_WLR:         safeDivide(response.Data.Player.Statistics.HITW_Wins, response.Data.Player.Statistics.HITWGames-response.Data.Player.Statistics.HITW_Wins),
			HITW_Wallsdodged: response.Data.Player.Statistics.HITW_Wallsdodged,
			HITW_WGR:         safeDivide(response.Data.Player.Statistics.HITW_Wallsdodged, response.Data.Player.Statistics.HITWGames),
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
		Currency:          response.Data.Player.Collections.Currency,
		Trophies:          response.Data.Player.CrownLevel.Trophies,
		TrophiesSKILL:     response.Data.Player.CrownLevel.TrophiesSKILL,
		TrophiesSTYLE:     response.Data.Player.CrownLevel.TrophiesSTYLE,
		TrophiesANGLER:    response.Data.Player.CrownLevel.TrophiesANGLER,
		Friends:           response.Data.Player.Social.Friends,
		EquippedCosmetics: equippedCosmetics,
		Hats:              SortCosmetics(hats),
		Accessories:       SortCosmetics(accessories),
		Auras:             SortCosmetics(auras),
		Trails:            SortCosmetics(trails),
		Cloaks:            SortCosmetics(cloaks),
		Rods:              SortCosmetics(rods),
	}
	return Infos
}
