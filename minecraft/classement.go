package minecraft

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

type PlayerRank struct {
	UUID       string `json:"uuid"`
	Capes      int    `json:"capes"`
	Score      int    `json:"score"`
	ActualName string `json:"actualname"`
	Badge      string `json:"badge"`
	CapeList   []struct {
		Name string `json:"name"`
	} `json:"capelist"`
}

type MccRank struct {
	UUID       string `json:"uuid"`
	ActualName string `json:"actualname"`
	CrownLevel int    `json:"crownlevel"`
	Trophies   int    `json:"trophies"`
	Rank       string `json:"rank"`
}

type Classement struct {
	Classement []PlayerRank `json:"classement"`
}

type MccClassement struct {
	Classement []MccRank `json:"classement"`
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

	// Calculer le score total du joueur et récupérer la liste des capes
	totalScore := 0
	capesCount := 0
	capeList := []struct {
		Name string `json:"name"`
	}{}

	for _, cape := range listCapes {
		if cape.Removed {
			totalScore += 0 // Si la cape est supprimée, elle vaut 1 point
		} else if score, exists := capeScores[cape.Name]; exists {
			totalScore += score // Sinon, on prend son score normal
		}

		// Ajouter à la liste des capes possédées
		capeList = append(capeList, struct {
			Name string `json:"name"`
		}{Name: cape.Name})

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
			// Mettre à jour le nombre de capes, le score, le pseudonyme, le badge et la liste des capes du joueur
			classement.Classement[i].Capes = capesCount
			classement.Classement[i].Score = totalScore
			classement.Classement[i].ActualName = actualName
			classement.Classement[i].Badge = badge
			classement.Classement[i].CapeList = capeList
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
			CapeList:   capeList,
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
	rawData, err := json.MarshalIndent(classement, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la conversion en JSON : %v", err)
		return -1
	}

	// Convertir en string pour manipulation
	formattedJSON := string(rawData)

	// Expression régulière pour détecter les capes sous forme de liste
	re := regexp.MustCompile(`"capelist": \[\s*([\s\S]*?)\s*\]`)

	// Reformater la liste pour qu'elle tienne sur une seule ligne
	formattedJSON = re.ReplaceAllStringFunc(formattedJSON, func(match string) string {
		// Supprimer les sauts de ligne et espaces inutiles
		compactList := strings.ReplaceAll(match, "\n", "")
		compactList = strings.ReplaceAll(compactList, "  ", "")  // Retire les espaces d'indentation
		compactList = strings.ReplaceAll(compactList, " ]", "]") // Corrige la fermeture
		compactList = strings.ReplaceAll(compactList, "[ ", "[") // Corrige l'ouverture
		return compactList
	})

	// Sauvegarder le JSON modifié
	err = ioutil.WriteFile(filePath, []byte(formattedJSON), 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du fichier classement : %v", err)
		return -1
	}

	return playerPosition
}

func UpdateMccClassement(uuid string, actualName string, crownLevel int, trophies int, rank string) int {
	filePath := "./site/infos/z_db_mccclassement.json"

	file, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Erreur lors de la lecture du fichier classement MCC : %v", err)
		return -1
	}

	var classement MccClassement

	if len(file) > 0 {
		err = json.Unmarshal(file, &classement)
		if err != nil {
			log.Printf("Erreur lors de l'analyse du JSON : %v", err)
			return -1
		}
	}

	found := false
	for i, player := range classement.Classement {
		if player.UUID == uuid {
			classement.Classement[i].ActualName = actualName
			classement.Classement[i].CrownLevel = crownLevel
			classement.Classement[i].Trophies = trophies
			classement.Classement[i].Rank = rank
			found = true
			break
		}
	}

	if !found {
		classement.Classement = append(classement.Classement, MccRank{
			UUID:       uuid,
			ActualName: actualName,
			CrownLevel: crownLevel,
			Trophies:   trophies,
			Rank:       rank,
		})
	}

	sort.Slice(classement.Classement, func(i, j int) bool {
		return classement.Classement[i].Trophies > classement.Classement[j].Trophies
	})

	playerPosition := -1
	for i, player := range classement.Classement {
		if player.UUID == uuid {
			playerPosition = i + 1
			break
		}
	}

	rawData, err := json.MarshalIndent(classement, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la conversion en JSON : %v", err)
		return -1
	}

	err = ioutil.WriteFile(filePath, []byte(rawData), 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du fichier classement MCC : %v", err)
		return -1
	}

	return playerPosition
}
