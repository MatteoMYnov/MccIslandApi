package minecraft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type BadgeGroup struct {
	Restriction string   `json:"restriction"`
	Capes       []string `json:"capes"`
	UUID        []string `json:"uuid"`
}

type BadgeGroups struct {
	Badges []map[string]BadgeGroup `json:"badges"`
}

// Fonction pour charger les badges en fonction des capes et du nom
func LoadBadges(name string, capes []string) []string {
	badges, err := LoadBadgesFromFile("./site/infos/badges.json")
	if err != nil {
		fmt.Println("Erreur lors du chargement des badges :", err)
		return nil
	}

	return GetBadges(name, capes, badges)
}

// Charger le fichier JSON contenant les groupes de badges
func LoadBadgesFromFile(filePath string) (BadgeGroups, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return BadgeGroups{}, fmt.Errorf("erreur lors de l'ouverture du fichier : %v", err)
	}
	defer file.Close()

	// Lire le contenu du fichier
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return BadgeGroups{}, fmt.Errorf("erreur lors de la lecture du fichier : %v", err)
	}

	// Décoder le JSON dans une structure BadgeGroups
	var badges BadgeGroups
	err = json.Unmarshal(data, &badges)
	if err != nil {
		return BadgeGroups{}, fmt.Errorf("erreur lors du décodage JSON : %v", err)
	}

	return badges, nil
}

// Récupère les badges du joueur en fonction de son nom, de ses capes et des badges chargés
func GetBadges(name string, capes []string, badges BadgeGroups) []string {
	capeTypes := make(map[string]bool)
	for _, cape := range capes {
		capeTypes[cape] = true
	}

	// Créer une liste des badges du joueur
	var playerBadges []string
	normalizedPlayerName := strings.ToLower(name)
	// Vérifier chaque groupe de badges
	for _, badgeGroup := range badges.Badges {
		for badgeName, badge := range badgeGroup {
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

			// Vérifier si le joueur est dans la liste des UUID
			if len(badge.UUID) > 0 {
				for _, playerUUID := range badge.UUID {
					// Normaliser le nom obtenu via l'UUID
					normalizedUUIDName := strings.ToLower(GetName(playerUUID))
					// Comparer les deux noms normalisés
					if normalizedUUIDName == normalizedPlayerName {
						playerBadges = append(playerBadges, badgeName)
						break
					}
				}
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
