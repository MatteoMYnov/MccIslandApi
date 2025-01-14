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

// GetUUID prend un nom de joueur et retourne son UUID ou termine le programme en cas d'erreur
func GetUUID(name string) string {
	// Construire l'URL
	url := fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s", name)

	// Effectuer la requête HTTP
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur lors de la requête HTTP : %v\n", err)
		return "" // Renvoie une chaîne vide en cas d'erreur
	}
	defer resp.Body.Close()

	// Vérifier le statut HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur : statut HTTP invalide %d\n", resp.StatusCode)
		return "" // Renvoie une chaîne vide si le statut HTTP n'est pas 200 OK
	}

	// Lire le corps de la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la réponse : %v\n", err)
		return "" // Renvoie une chaîne vide en cas d'erreur
	}

	// Décoder la réponse JSON
	var response MojangResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Erreur de décodage JSON : %v\n", err)
		return "" // Renvoie une chaîne vide si le décodage échoue
	}

	// Retourner l'UUID du joueur
	return response.ID
}
