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
	Bonus      int `json:"bonus"`
}

type CrownLevel struct {
	Level             int               `json:"level"`
	Evolution         int               `json:"evolution"`
	NextLevelProgress NextLevelProgress `json:"nextLevelProgress"`
	Trophies          Trophies          `json:"trophies"`
}

type Currency struct {
	Coins           int `json:"coins"`
	RoyalReputation int `json:"royalReputation"`
	Silver          int `json:"silver"`
	MaterialDust    int `json:"materialDust"`
	AnglrTokens     int `json:"anglrTokens"`
}

type Player struct {
	UUID        string     `json:"uuid"`
	Username    string     `json:"username"`
	Ranks       []string   `json:"ranks"`
	CrownLevel  CrownLevel `json:"crownLevel"`
	Collections struct {
		Currency Currency `json:"currency"`
	} `json:"collections"`
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
	Ranks           []string `json:"ranks"`
	CrownLevel      int      `json:"crownLevel"`
	Evolution       int      `json:"evolution"`
	CrownObtained   int      `json:"crownObtained"`
	CrownObtainable int      `json:"crownObtainable"`
	Currency        Currency `json:"currency"`
	Trophies        Trophies `json:"trophies"`
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
					trophies {
						obtained
						obtainable
						bonus
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

	Infos := &MccInfos{
		Ranks:           response.Data.Player.Ranks,
		CrownLevel:      response.Data.Player.CrownLevel.Level,
		Evolution:       response.Data.Player.CrownLevel.Evolution,
		CrownObtained:   response.Data.Player.CrownLevel.NextLevelProgress.CrownObtained,
		CrownObtainable: response.Data.Player.CrownLevel.NextLevelProgress.CrownObtainable,
		Currency: Currency{
			Coins:           response.Data.Player.Collections.Currency.Coins,
			RoyalReputation: response.Data.Player.Collections.Currency.RoyalReputation,
			Silver:          response.Data.Player.Collections.Currency.Silver,
			MaterialDust:    response.Data.Player.Collections.Currency.MaterialDust,
			AnglrTokens:     response.Data.Player.Collections.Currency.AnglrTokens,
		},
		Trophies: Trophies{
			Obtained:   response.Data.Player.CrownLevel.Trophies.Obtained,
			Obtainable: response.Data.Player.CrownLevel.Trophies.Obtainable,
			Bonus:      response.Data.Player.CrownLevel.Trophies.Bonus,
		},
	}
	return Infos
}
