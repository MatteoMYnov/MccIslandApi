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
	Level             int               `json:"level"`
	Evolution         int               `json:"evolution"`
	NextLevelProgress NextLevelProgress `json:"nextLevelProgress"`
	FishingData       FishingData       `json:"fishingLevelData"`
	Trophies          Trophies          `json:"trophies"`
	TrophiesSKILL     Trophies          `json:"trophiesSKILL,omitempty"`
	TrophiesSTYLE     Trophies          `json:"trophiesSTYLE,omitempty"`
	TrophiesANGLER    Trophies          `json:"trophiesANGLER,omitempty"`
}

type Currency struct {
	Coins           int `json:"coins"`
	RoyalReputation int `json:"royalReputation"`
	Silver          int `json:"silver"`
	MaterialDust    int `json:"materialDust"`
	AnglrTokens     int `json:"anglrTokens"`
}

type Friend struct {
	Username   string   `json:"username"`
	Ranks      []string `json:"ranks"`
	CrownLevel struct {
		Evolution int `json:"evolution"`
	} `json:"crownLevel"`
}

type Player struct {
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
				ranks
				crownLevel {
					level
					evolution
					nextLevelProgress {
						obtained
						obtainable
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
						crownLevel {
							evolution
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
		Ranks:           response.Data.Player.Ranks,
		CrownLevel:      response.Data.Player.CrownLevel.Level,
		Evolution:       response.Data.Player.CrownLevel.Evolution,
		CrownObtained:   response.Data.Player.CrownLevel.NextLevelProgress.CrownObtained,
		CrownObtainable: response.Data.Player.CrownLevel.NextLevelProgress.CrownObtainable,
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
