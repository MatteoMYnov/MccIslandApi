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

// Fonction pour récupérer les noms des capes d'un utilisateur, incluant celles par UUID
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

	// Ajouter les capes récupérées depuis l'API
	for _, cape := range response.Capes {
		if !cape.Removed {
			capesList = append(capesList, cape.Type)
		}
	}

	// Charger les capes depuis le fichier JSON
	capeGroups, err := LoadCapesFromFile("./site/infos/capes.json")
	if err != nil {
		fmt.Printf("Erreur lors du chargement des capes : %v\n", err)
		return capesList // Retourner uniquement les capes de l'API si une erreur survient
	}

	// Ajouter les capes associées au joueur par UUID
	normalizedPlayerName := strings.ToLower(name)
	for _, cape := range capeGroups.Capes {
		for _, playerUUID := range cape.UUID {
			if strings.ToLower(GetName(playerUUID)) == normalizedPlayerName {
				capesList = append(capesList, cape.Name)
				break
			}
		}
	}

	return capesList
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

// Fonction principale pour obtenir les capes
// Fonction pour obtenir les capes du joueur avec leur score associé
func GetCapes(name string, capeGroups CapeGroups) []map[string]interface{} {
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
	fmt.Println(capesList)
	return capesList
}
