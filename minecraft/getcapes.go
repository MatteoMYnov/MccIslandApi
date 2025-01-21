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

func GetCapeNames(name string) []string {
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

	var capesList []string
	for _, cape := range response.Capes {
		if !cape.Removed {
			capesList = append(capesList, cape.Type)
		}
	}
	return capesList
}

func GetCapes(name string) []map[string]interface{} {
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

	var capesList []map[string]interface{}
	for _, cape := range response.Capes {
		capeObj := map[string]interface{}{
			"cape":    cape.Type,
			"removed": cape.Removed,
		}
		capesList = append(capesList, capeObj)
	}
	return capesList
}
