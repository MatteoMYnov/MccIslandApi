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
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
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

func GetNameFast(uuid string) string {
	url := fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/lookup/%s", uuid)

	resp, err := http.Get(url)
	if err != nil {
		return "" // Retourne une chaîne vide si une erreur survient
	}
	defer resp.Body.Close()

	var data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "" // Retourne une chaîne vide en cas d'erreur de décodage
	}

	if data.Name != "" {
		return data.Name // Retourne le nom si trouvé
	}

	return "" // Retourne une chaîne vide si aucun nom n'est trouvé
}
