package mcc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GraphQLRequest struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables,omitempty"`
}

type CrownLevel struct {
	Level int `json:"level"`
}

type Player struct {
	UUID       string     `json:"uuid"`
	Username   string     `json:"username"`
	Ranks      []string   `json:"ranks"`
	CrownLevel CrownLevel `json:"crownLevel"`
}

type Response struct {
	Data struct {
		Player Player `json:"player"`
	} `json:"data"`
}

type APIConfig struct {
	Mcctoken string `json:"mcctoken"`
}

func GetInfos(UUID string) []string {
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

	// DEBUG: Vérifier si le token est bien extrait
	fmt.Printf("Token décodé: %s\n", apiConfig.Mcctoken)

	if apiConfig.Mcctoken == "" {
		log.Println("Le token API est vide")
		return nil
	}

	// Définir la requête GraphQL
	query := `
		query player($uuid: UUID!) {
			player(uuid: $uuid) {
				uuid
				username
				ranks
				crownLevel {
      				level
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

	return response.Data.Player.Ranks
}
