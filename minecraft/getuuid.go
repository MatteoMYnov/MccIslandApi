package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Structure pour la réponse JSON de l'API Mojang
type MojangResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Fonction pour formater l'UUID en ajoutant les tirets
func formatUUID(uuid string) string {
	if len(uuid) != 32 {
		return uuid // Retourne tel quel si l'UUID n'est pas de longueur 32
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		uuid[0:8], uuid[8:12], uuid[12:16], uuid[16:20], uuid[20:32])
}

// GetUUID prend un nom de joueur et retourne son UUID formaté
func GetUUID(name string) (string, string) {
	// Construire l'URL
	url := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", name)

	// Effectuer la requête HTTP
	resp, err := http.Get(url)
	if err != nil {
		return "", name
	}
	defer resp.Body.Close()

	// Vérifier le statut HTTP
	if resp.StatusCode != http.StatusOK {
		return "", name
	}

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la réponse : %v\n", err)
		return "", name
	}

	// Décoder la réponse JSON
	var response MojangResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Erreur de décodage JSON : %v\n", err)
		return "", name
	}

	// Retourner l'UUID formaté
	return formatUUID(response.ID), response.Name
}
