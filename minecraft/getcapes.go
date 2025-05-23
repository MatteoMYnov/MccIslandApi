package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Structures pour représenter les capes
type CapeGroups struct {
	Capes []CapeForced `json:"capes"`
}

type Cape struct {
	Type    string `json:"type"`
	Removed bool   `json:"removed"`
}

type CapeForced struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Title string   `json:"title"`
	UUID  []string `json:"UUID"`
	Score int      `json:"score"` // Ajout du champ Score
}

type CapesResponse struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Capes    []Cape `json:"capes"`
}

// Fonction pour charger les capes par nom
func LoadCapesByName(name string) []map[string]interface{} {
	capesList, err := LoadCapesFromFile("./site/infos/capes.json")
	if err != nil {
		fmt.Println("Erreur lors du chargement des capes :", err)
		return nil
	}
	return GetCapes(name, capesList)
}

// Fonction pour charger les capes à partir d'un fichier JSON
func LoadCapesFromFile(filePath string) (CapeGroups, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return CapeGroups{}, fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return CapeGroups{}, fmt.Errorf("erreur lors de la lecture du fichier : %v", err)
	}

	var capes CapeGroups
	err = json.Unmarshal(data, &capes)
	if err != nil {
		return CapeGroups{}, fmt.Errorf("erreur lors du décodage JSON : %v", err)
	}

	return capes, nil
}

func GetCapes(name string, capeGroups CapeGroups) []map[string]interface{} {
	url := fmt.Sprintf("https://capes.me/api/user/%s", name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erreur HTTP : %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erreur lecture body : %v\n", err)
		return nil
	}

	var response CapesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Erreur JSON : %v\n", err)
		return nil
	}

	var capesList []map[string]interface{}

	// Ajouter les capes récupérées depuis l'API
	for _, cape := range response.Capes {
		capesList = append(capesList, map[string]interface{}{
			"cape":    cape.Type,
			"removed": cape.Removed,
		})
	}

	// Charger les capes depuis le fichier JSON
	capeGroups, err = LoadCapesFromFile("./site/infos/capes.json")
	if err != nil {
		fmt.Printf("Erreur lors du chargement des capes : %v\n", err)
		return capesList
	}

	// Ajouter les capes associées au joueur par UUID
	normalizedPlayerName := strings.ToLower(name)
	for _, cape := range capeGroups.Capes {
		for _, playerUUID := range cape.UUID {
			if strings.ToLower(GetName(playerUUID)) == normalizedPlayerName {
				capesList = append(capesList, map[string]interface{}{
					"cape":    cape.Name,
					"removed": false,
				})
				break
			}
		}
	}
	return capesList
}
