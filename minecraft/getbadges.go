package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type BadgeGroup struct {
	Restriction string   `json:"restriction"`
	Capes       []string `json:"capes"`
}

type BadgeGroups map[string]BadgeGroup

func LoadBadgesByName(name string) []string {
	capes := GetCapeNames(name)
	if capes == nil {
		fmt.Println("Impossible de récupérer les capes pour l'utilisateur :", name)
		return nil
	}
	return LoadBadges(capes)
}

func LoadBadges(capes []string) []string {
	// Charger les badges depuis un fichier JSON
	badges, err := LoadBadgesFromFile("./site/infos/badges.json")
	if err != nil {
		fmt.Println("Erreur lors du chargement des badges :", err)
		return nil
	}

	// Retourner les badges du joueur en fonction de ses capes et des badges définis
	return GetBadges(capes, badges)
}

// Charger le fichier JSON contenant les groupes de badges
func LoadBadgesFromFile(filePath string) (BadgeGroups, error) {
	// Lire le fichier JSON
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Lire le contenu du fichier
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier : %v", err)
	}

	// Décoder le JSON dans une structure BadgeGroups
	var badges BadgeGroups
	err = json.Unmarshal(data, &badges)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage JSON : %v", err)
	}

	return badges, nil
}

// Récupère les badges du joueur en fonction de ses capes et des badges chargés
func GetBadges(capes []string, badges BadgeGroups) []string {
	// Créer une carte des capes possédées par le joueur pour un accès rapide
	capeTypes := make(map[string]bool)
	for _, cape := range capes {
		capeTypes[cape] = true
	}

	// Créer une liste des badges du joueur
	var playerBadges []string

	// Vérifier chaque groupe de badges
	for badgeName, badge := range badges {
		if badge.Restriction == "one" {
			// Vérifier si le joueur possède au moins une cape du groupe
			if containsAnyCape(badge.Capes, capeTypes) {
				playerBadges = append(playerBadges, badgeName)
			}
		} else if badge.Restriction == "all" {
			// Vérifier si le joueur possède toutes les capes du groupe
			if containsAllCapes(badge.Capes, capeTypes) {
				playerBadges = append(playerBadges, badgeName)
			}
		}
	}

	// Retourner la liste des badges du joueur
	return playerBadges
}

// Vérifie si le joueur possède au moins une cape du groupe
func containsAnyCape(badgeGroup []string, capeTypes map[string]bool) bool {
	for _, badge := range badgeGroup {
		if capeTypes[badge] {
			return true
		}
	}
	return false
}

// Vérifie si le joueur possède toutes les capes du groupe
func containsAllCapes(badgeGroup []string, capeTypes map[string]bool) bool {
	for _, badge := range badgeGroup {
		if !capeTypes[badge] {
			return false
		}
	}
	return true
}
