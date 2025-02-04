package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"hypixel-info/mcc"
	"hypixel-info/minecraft"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

// Nouveau format JSON pour les capes
type Cape struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Title string   `json:"title"`
	UUID  []string `json:"UUID"`
}

type CapeGroups struct {
	Capes []Cape `json:"capes"`
}

type CapeInfo struct {
	URL      string
	Class    string
	CapeName string
	Title    string
	Removed  bool
}

type BadgeInfo struct {
	URL       string
	Class     string
	BadgeName string
}

type Infos struct {
	Name        string
	ListCapes   []string
	ImageURLs   []CapeInfo
	BadgeURLs   []BadgeInfo
	PlayerClass string
}

type DataMenuPage struct {
	Name                 string
	ListCapes            []string
	ImageURLs            []CapeInfo
	BadgeURLs            []BadgeInfo
	PlayerClass          string
	MccRank              string
	Evolution            string
	CrownLevel           string
	CrownLevelplus1      string
	CrownObtained        int
	CrownObtainable      int
	CrownPourcentage     int
	FishingEvolution     string
	FishingLevel         string
	FishingLevelplus1    string
	FishingObtained      int
	FishingObtainable    int
	FishingPourcentage   int
	CurrencyCoins        string
	CurrencyRoyalRep     string
	CurrencySilver       string
	CurrencyMaterialDust string
	CurrencyAnglrTokens  string
	TotalTrophies        string
	MaxTotalTrophies     string
	SKILLTrophies        string
	MaxSKILLTrophies     string
	STYLETrophies        string
	MaxSTYLETrophies     string
	ANGLERTrophies       string
	MaxANGLERTrophies    string
	BonusTrophies        string
	Friends              []FriendInfo
	// Player Rank
	PlayerRank int
}

type FriendInfo struct {
	Username   string
	Ranks      string
	CrownLevel struct {
		Evolution int
		Level     int
	}
}

type PlayerRank struct {
	UUID  string `json:"uuid"`
	Capes int    `json:"capes"`
	Score int    `json:"score"` // Ajout du champ Score
}

// Structure principale du fichier JSON
type Classement struct {
	Classement []PlayerRank `json:"classement"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Définir le chemin du fichier que tu veux télécharger
	filePath := "./site/infos/z_db_classement.json"

	// Vérifie si le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	// Définir les en-têtes pour indiquer que le fichier est à télécharger
	w.Header().Set("Content-Disposition", "attachment; filename=z_db_classement.json")
	w.Header().Set("Content-Type", "application/json")

	// Utilise http.ServeFile pour envoyer le fichier
	http.ServeFile(w, r, filePath)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	capeGroups, err := minecraft.LoadCapeGroups()
	if err != nil {
		log.Fatal("Erreur de chargement des groupes de capes:", err)
	}

	IGN := r.FormValue("q")
	if IGN == "" {
		IGN = minecraft.GetRandomName()
	}

	playerClass := ""
	if !minecraft.IsValidIGN(IGN) {
		playerClass = "badName"
	} else {
		playerClass = "playerName"
	}

	playerUUID := minecraft.GetUUID(IGN) // Utilisation de playerUUID pour obtenir les infos MCC
	playerCapesJSON := minecraft.LoadCapesByName(IGN)
	playerBadgesJSON := minecraft.LoadBadgesByName(IGN)

	var listCapes []string
	capeInfos := []CapeInfo{}
	for _, cape := range playerCapesJSON {
		capeName := cape["cape"].(string)
		removed := cape["removed"].(bool)

		listCapes = append(listCapes, capeName)

		class := minecraft.GetCapeClass(capeName, capeGroups)
		title := ""

		if removed {
			class += " removed-cape"
		}

		for _, capeGroup := range capeGroups.Capes {
			if capeGroup.Name == capeName {
				title = capeGroup.Title
				break
			}
		}

		capeInfos = append(capeInfos, CapeInfo{
			URL:      "/img/capes/" + capeName + ".png",
			Class:    class,
			CapeName: capeName,
			Title:    title,
			Removed:  removed,
		})
	}

	playerRank := UpdateClassement(playerUUID, len(listCapes), listCapes)

	badgeInfos := []BadgeInfo{}
	for _, badgeName := range playerBadgesJSON {
		badgeInfos = append(badgeInfos, BadgeInfo{
			URL:       "/img/badges/" + badgeName + ".png",
			Class:     "badge-class",
			BadgeName: badgeName,
		})
	}

	if len(badgeInfos) == 0 {
		badgeInfos = nil
	}

	prioritizedCapes := minecraft.PrioritizeCapes(listCapes, capeGroups)

	var prioritizedCapeInfos []CapeInfo
	for _, cape := range prioritizedCapes {
		for _, capeInfo := range capeInfos {
			if capeInfo.URL == "/img/capes/"+cape+".png" {
				prioritizedCapeInfos = append(prioritizedCapeInfos, capeInfo)
			}
		}
	}

	// Obtenez les informations du MCC pour le joueur
	mccInfos := mcc.GetInfos(playerUUID)

	// Définir la valeur de MccRank
	MccRank := "PLAYER"
	if mccInfos != nil && len(mccInfos.Ranks) > 0 {
		MccRank = mccInfos.Ranks[0]
	}

	// Calculer l'évolution et le niveau de couronne pour le joueur
	crownLevelPlus1 := mccInfos.CrownLevel + 1
	fishingLevelPlus1 := mccInfos.FishingData.Level + 1
	fishingEvolution := mccInfos.FishingData.Level / 10

	// Obtenez les informations de la couronne
	crownObtained := 0
	crownObtainable := 0
	fishingObtained := 0
	fishingObtainable := 0
	if mccInfos != nil {
		crownObtained = mccInfos.CrownObtained
		crownObtainable = mccInfos.CrownObtainable
		fishingObtained = mccInfos.FishingData.NextLevelProgress.CrownObtained
		fishingObtainable = mccInfos.FishingData.NextLevelProgress.CrownObtainable
	}
	var crowncalculatedPercent int
	var fishingcalculatedPercent int
	if mccInfos.CrownObtainable > 0 {
		crowncalculatedPercent = (mccInfos.CrownObtained * 100) / mccInfos.CrownObtainable
	} else {
		crowncalculatedPercent = 0
	}
	if mccInfos.FishingData.NextLevelProgress.CrownObtainable > 0 {
		fishingcalculatedPercent = (mccInfos.FishingData.NextLevelProgress.CrownObtained * 100) / mccInfos.FishingData.NextLevelProgress.CrownObtainable
	} else {
		fishingcalculatedPercent = 0
	}

	infos := DataMenuPage{
		Name:        IGN,
		ListCapes:   prioritizedCapes,
		ImageURLs:   prioritizedCapeInfos,
		BadgeURLs:   badgeInfos,
		PlayerClass: playerClass,
		MccRank:     MccRank,
		// Crown
		Evolution:        fmt.Sprintf("%d", mccInfos.Evolution),
		CrownLevel:       fmt.Sprintf("%d", mccInfos.CrownLevel),
		CrownLevelplus1:  fmt.Sprintf("%d", crownLevelPlus1),
		CrownObtained:    crownObtained,
		CrownObtainable:  crownObtainable,
		CrownPourcentage: crowncalculatedPercent,
		// Crown Fishing
		FishingEvolution:   fmt.Sprintf("%d", fishingEvolution),
		FishingLevel:       fmt.Sprintf("%d", mccInfos.FishingData.Level),
		FishingLevelplus1:  fmt.Sprintf("%d", fishingLevelPlus1),
		FishingObtained:    fishingObtained,
		FishingObtainable:  fishingObtainable,
		FishingPourcentage: fishingcalculatedPercent,
		// Ajout des informations de currency en tant qu'entiers
		CurrencyCoins:        mcc.FormatNumberWithSpaces(mccInfos.Currency.Coins),
		CurrencyRoyalRep:     mcc.FormatNumberWithSpaces(mccInfos.Currency.RoyalReputation),
		CurrencySilver:       mcc.FormatNumberWithSpaces(mccInfos.Currency.Silver),
		CurrencyMaterialDust: mcc.FormatNumberWithSpaces(mccInfos.Currency.MaterialDust),
		CurrencyAnglrTokens:  mcc.FormatNumberWithSpaces(mccInfos.Currency.AnglrTokens),
		TotalTrophies:        mcc.FormatNumberWithSpaces(mccInfos.Trophies.Obtained),
		MaxTotalTrophies:     mcc.FormatNumberWithSpaces(mccInfos.Trophies.Obtainable),
		SKILLTrophies:        mcc.FormatNumberWithSpaces(mccInfos.TrophiesSKILL.Obtained),
		MaxSKILLTrophies:     mcc.FormatNumberWithSpaces(mccInfos.TrophiesSKILL.Obtainable),
		STYLETrophies:        mcc.FormatNumberWithSpaces(mccInfos.TrophiesSTYLE.Obtained),
		MaxSTYLETrophies:     mcc.FormatNumberWithSpaces(mccInfos.TrophiesSTYLE.Obtainable),
		ANGLERTrophies:       mcc.FormatNumberWithSpaces(mccInfos.TrophiesANGLER.Obtained),
		MaxANGLERTrophies:    mcc.FormatNumberWithSpaces(mccInfos.TrophiesANGLER.Obtainable),
		BonusTrophies:        mcc.FormatNumberWithSpaces(mccInfos.Trophies.Bonus),
		Friends:              convertToFriendInfo(mccInfos.Friends),
		// Player Rank
		PlayerRank: playerRank,
	}

	tmplPath := filepath.Join("site", "template", "menu.html")
	tmpl, err := template.New("menu.html").Funcs(template.FuncMap{
		"contains": minecraft.Contains,
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, infos)
}

func convertToFriendInfo(friends []mcc.Friend) []FriendInfo {
	var friendInfo []FriendInfo
	for _, friend := range friends {
		rank := "PLAYER"
		if len(friend.Ranks) > 0 {
			rank = friend.Ranks[0] // Récupère le premier rang
		}
		friendInfo = append(friendInfo, FriendInfo{
			Username: friend.Username,
			Ranks:    rank,
			CrownLevel: struct {
				Evolution int
				Level     int
			}{
				Evolution: friend.CrownLevel.Evolution,
				Level:     friend.CrownLevel.Level,
			},
		})
	}
	return friendInfo
}

func capesHandler(w http.ResponseWriter, r *http.Request) {
	capeGroups, err := minecraft.LoadCapeGroups()
	if err != nil {
		http.Error(w, "Erreur lors du chargement des capes", http.StatusInternalServerError)
		return
	}

	var capeInfos []CapeInfo
	for _, cape := range capeGroups.Capes {
		capeInfos = append(capeInfos, CapeInfo{
			URL:      "/img/capes/" + cape.Name + ".png",
			Class:    cape.Type + "-cape",
			CapeName: cape.Name,
			Title:    cape.Title,
			Removed:  false,
		})
	}

	data := struct {
		ImageURLs []CapeInfo
	}{
		ImageURLs: capeInfos,
	}

	tmplPath := filepath.Join("site", "template", "capes.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func UpdateClassement(uuid string, capesCount int, listCapes []string) int {
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

	// Si le joueur n'a aucune cape, ne pas l'ajouter au classement
	if capesCount == 0 {
		return -1
	}

	// Récupérer les capes et leur score
	capeGroups, err := minecraft.LoadCapeGroups()
	if err != nil {
		log.Printf("Erreur lors du chargement des groupes de capes : %v", err)
		return -1
	}

	// Calculer le score total du joueur en fonction des capes
	totalScore := 0
	for _, cape := range capeGroups.Capes {
		for _, capeName := range listCapes {
			if capeName == cape.Name {
				totalScore += cape.Score
			}
		}
	}

	// Vérifier si le joueur est déjà dans la liste
	found := false
	for i, player := range classement.Classement {
		if player.UUID == uuid {
			// Mettre à jour le nombre de capes et le score du joueur
			classement.Classement[i].Capes = capesCount
			classement.Classement[i].Score = totalScore
			found = true
			break
		}
	}

	// Ajouter le joueur s'il n'existe pas encore
	if !found {
		classement.Classement = append(classement.Classement, PlayerRank{UUID: uuid, Capes: capesCount, Score: totalScore})
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

func setupFileServer(path, route string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle(route, http.StripPrefix(route, fs))
}

func main() {
	http.HandleFunc("/", rootHandler)

	setupFileServer("./site/styles", "/styles/")
	setupFileServer("./site/img", "/img/")
	setupFileServer("./site/js", "/js/")

	http.HandleFunc("/download", downloadFileHandler) // pour Download la database des joueurs

	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/capes", capesHandler)

	if err := http.ListenAndServe(":1618", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
