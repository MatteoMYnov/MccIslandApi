package minecraft

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type PlayerRank struct {
	UUID       string `json:"uuid"`
	Capes      int    `json:"capes"`
	Score      int    `json:"score"`
	ActualName string `json:"actualname"` // Nouveau champ pour le pseudonyme
	Badge      string `json:"badge"`      // Nouveau champ pour le premier badge
}

type Classement struct {
	Classement []PlayerRank `json:"classement"`
}

func UpdateClassement(uuid string, listCapes []struct {
	Name    string
	Removed bool
}, actualName string, badge string) int {
	filePath := "./site/infos/z_db_classement.json"

	// Lire le fichier
	file, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Erreur lors de la lecture du fichier classement : %v", err)
		return -1
	}

	var classement Classement

	// Charger les données existantes
	if len(file) > 0 {
		err = json.Unmarshal(file, &classement)
		if err != nil {
			log.Printf("Erreur lors de l'analyse du JSON : %v", err)
			return -1
		}
	}

	// Récupérer les capes et leur score
	capeGroups, err := LoadCapeGroups()
	if err != nil {
		log.Printf("Erreur lors du chargement des groupes de capes : %v", err)
		return -1
	}

	// Création d'une map pour un accès rapide aux scores des capes
	capeScores := make(map[string]int)
	for _, cape := range capeGroups.Capes {
		capeScores[cape.Name] = cape.Score
	}

	// Calculer le score total du joueur
	totalScore := 0
	capesCount := 0

	for _, cape := range listCapes {
		if cape.Removed {
			totalScore += 1 // Si la cape est supprimée, elle vaut 1 point
		} else if score, exists := capeScores[cape.Name]; exists {
			totalScore += score // Sinon, on prend son score normal
		}

		// On compte toutes les capes, qu'elles soient supprimées ou non
		capesCount++
	}

	// Si le joueur n'a aucune cape, ne pas l'ajouter au classement
	if capesCount == 0 {
		return -1
	}

	// Vérifier si le joueur est déjà dans la liste
	found := false
	for i, player := range classement.Classement {
		if player.UUID == uuid {
			// Mettre à jour le nombre de capes, le score, le pseudonyme et le badge du joueur
			classement.Classement[i].Capes = capesCount
			classement.Classement[i].Score = totalScore
			classement.Classement[i].ActualName = actualName
			classement.Classement[i].Badge = badge
			found = true
			break
		}
	}

	// Ajouter le joueur s'il n'existe pas encore
	if !found {
		classement.Classement = append(classement.Classement, PlayerRank{
			UUID:       uuid,
			Capes:      capesCount,
			Score:      totalScore,
			ActualName: actualName,
			Badge:      badge,
		})
	}

	// Trier le classement : priorité au nombre de capes, puis au score
	sort.Slice(classement.Classement, func(i, j int) bool {
		if classement.Classement[i].Capes == classement.Classement[j].Capes {
			return classement.Classement[i].Score > classement.Classement[j].Score
		}
		return classement.Classement[i].Capes > classement.Classement[j].Capes
	})

	// Trouver la position du joueur (1-based)
	playerPosition := -1
	for i, player := range classement.Classement {
		if player.UUID == uuid {
			playerPosition = i + 1
			break
		}
	}

	// Sauvegarder les modifications
	updatedData, err := json.MarshalIndent(classement, "", "    ")
	if err != nil {
		log.Printf("Erreur lors de la conversion en JSON : %v", err)
		return -1
	}

	err = ioutil.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du fichier classement : %v", err)
		return -1
	}

	return playerPosition
}
