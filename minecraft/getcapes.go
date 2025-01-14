package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Cape struct {
	Type    string `json:"type"`
	Removed bool   `json:"removed"`
}

type CapesResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Capes    []Cape `json:"capes"`
}

func GetCapes(name string) []string {
	url := fmt.Sprintf("https://capes.me/api/user/%s", name)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur lors de la requête HTTP : %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur : statut HTTP invalide %d\n", resp.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture de la réponse : %v\n", err)
		return nil
	}

	var response CapesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Erreur de décodage JSON : %v\n", err)
		return nil
	}

	var capeTypes []string
	for _, cape := range response.Capes {
		if !cape.Removed {
			capeTypes = append(capeTypes, cape.Type)
		}
	}

	return capeTypes
}
