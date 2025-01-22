package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MinecraftResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

func GetName(UUID string) string {
	url := fmt.Sprintf("https://api.ashcon.app/mojang/v2/user/%s", UUID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur lors de la requête HTTP : %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur : statut HTTP invalide %d pour UUID %s\n", resp.StatusCode, UUID)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la réponse : %v\n", err)
		return ""
	}

	var response MinecraftResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Erreur de décodage JSON : %v\n", err)
		return ""
	}
	return response.Username
}
