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

type Player struct {
	UUID  string `json:"uuid"`
	Capes int    `json:"capes"`
	Score int    `json:"score"`
}

type Joueur struct {
	UUID       string `json:"uuid"`
	ActualName string `json:"actualname"`
	Rank       int
	Capes      int `json:"capes"`
	Score      int `json:"score"`
}

type Classement struct {
	Joueurs []Joueur `json:"classement"`
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
	BackgroundUUID       string
	BackgroundNavbar     string
	BackgroundDark       string
	BackgroundLight      string
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
	GameStats            mcc.Statistics
	// Player Rank
	PlayerRank     int
	PlayerRankPage int
}

type FriendInfo struct {
	Username     string
	Ranks        string
	Online       bool
	OnlineStatus string
	CrownLevel   struct {
		Evolution int
		Level     int
	}
}

// Structure principale du fichier JSON

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

func convertToFriendInfo(friends []mcc.Friend) []FriendInfo {
	var friendInfo []FriendInfo
	for _, friend := range friends {
		rank := "PLAYER"
		if len(friend.Ranks) > 0 {
			rank = friend.Ranks[0] // Récupère le premier rang
		}

		onlineStatus := "offline"
		if friend.Status.Online {
			onlineStatus = "online"
		}

		friendInfo = append(friendInfo, FriendInfo{
			Username:     friend.Username,
			Ranks:        rank,
			Online:       friend.Status.Online,
			OnlineStatus: onlineStatus,
			CrownLevel: struct {
				Evolution int
				Level     int
			}{
				Evolution: friend.CrownLevel.Evolution,
				Level:     friend.CrownLevel.Level,
			},
		})
	}

	// Trier d'abord par statut en ligne (Online en premier), puis par CrownLevel.Level décroissant
	sort.Slice(friendInfo, func(i, j int) bool {
		if friendInfo[i].Online == friendInfo[j].Online {
			return friendInfo[i].CrownLevel.Level > friendInfo[j].CrownLevel.Level
		}
		return friendInfo[i].Online
	})

	return friendInfo
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
	if len(IGN) == 36 {
		IGN = minecraft.GetNameFast(IGN)
	}
	fmt.Println(IGN)

	playerClass := ""
	if !minecraft.IsValidIGN(IGN) {
		playerClass = "badName"
	} else {
		playerClass = "playerName"
	}

	playerUUID, IGN := minecraft.GetUUID(IGN) // Utilisation de playerUUID pour obtenir les infos MCC
	playerCapesJSON := minecraft.LoadCapesByName(IGN)
	playerBadgesJSON := minecraft.LoadBadgesByName(IGN)

	stylesFile := "./site/infos/styles.json"
	var stylesData struct {
		Players []struct {
			UUID   string `json:"uuid"`
			Navbar string `json:"navbar"`
			Dark   string `json:"dark"`
			Light  string `json:"light"`
		} `json:"players"`
	}

	var backgroundUUID, backgroundNavbar, backgroundDark, backgroundLight string
	if fileContent, err := ioutil.ReadFile(stylesFile); err == nil {
		if json.Unmarshal(fileContent, &stylesData) == nil {
			// Parcourez les joueurs et recherchez celui qui correspond à l'UUID du joueur
			for _, player := range stylesData.Players {
				if player.UUID == playerUUID {
					backgroundUUID = player.UUID
					backgroundNavbar = player.Navbar
					backgroundDark = player.Dark
					backgroundLight = player.Light
					break
				}
			}
		}
	}

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

	var capeDetails []struct {
		Name    string
		Removed bool
	}
	for _, cape := range capeInfos {
		capeDetails = append(capeDetails, struct {
			Name    string
			Removed bool
		}{
			Name:    cape.CapeName,
			Removed: cape.Removed,
		})
	}

	actualname := IGN
	playerRank := minecraft.UpdateClassement(playerUUID, capeDetails, actualname)
	playerRankPage := ((playerRank - 1) / 50) + 1 //f(x)=⌊(x−1)/50⌋+1

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
		//Player
		Name:             IGN,
		ListCapes:        prioritizedCapes,
		ImageURLs:        prioritizedCapeInfos,
		BadgeURLs:        badgeInfos,
		PlayerClass:      playerClass,
		BackgroundUUID:   backgroundUUID,
		BackgroundNavbar: backgroundNavbar,
		BackgroundDark:   backgroundDark,
		BackgroundLight:  backgroundLight,
		MccRank:          MccRank,
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
		GameStats:            mccInfos.Statistics,
		// Player Rank
		PlayerRank:     playerRank,
		PlayerRankPage: playerRankPage,
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

func classementHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire le numéro de page depuis l'URL
	var page int
	_, err := fmt.Sscanf(r.URL.Path, "/classement/%d", &page)
	if err != nil || page < 1 {
		http.Error(w, "Page invalide", http.StatusBadRequest)
		return
	}

	// Lire le fichier JSON du classement
	filePath := "./site/infos/z_db_classement.json"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Erreur de lecture du fichier de classement", http.StatusInternalServerError)
		return
	}

	// Décoder le fichier JSON
	var classement Classement
	err = json.Unmarshal(file, &classement)
	if err != nil {
		http.Error(w, "Erreur lors du décodage du JSON", http.StatusInternalServerError)
		return
	}

	// Définir la taille du groupe de joueurs par page
	const pageSize = 50
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Vérifier que l'index ne dépasse pas la taille du classement
	if startIndex >= len(classement.Joueurs) {
		http.Error(w, "Page inexistante", http.StatusNotFound)
		return
	}

	if endIndex > len(classement.Joueurs) {
		endIndex = len(classement.Joueurs)
	}

	// Extraire les joueurs pour la page demandée
	joueursPage := classement.Joueurs[startIndex:endIndex]

	// Mettre à jour le rang des joueurs pour refléter leur position globale
	for i := range joueursPage {
		joueursPage[i].Rank = startIndex + i + 1
	}

	// Renvoyer la page HTML avec les données du classement
	tmplPath := filepath.Join("site", "template", "classement.html")
	tmpl, err := template.New("classement.html").ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Classement []Joueur
		Page       int
		HasNext    bool
		HasPrev    bool
	}{
		Classement: joueursPage,
		Page:       page,
		HasPrev:    page > 1,
		HasNext:    endIndex < len(classement.Joueurs),
	}

	// Exécuter le template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

	http.HandleFunc("/dbdl", downloadFileHandler) // pour Download la database des joueurs
	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/capes", capesHandler)

	// Redirection de /classement vers /classement/1
	http.HandleFunc("/classement/", classementHandler)

	if err := http.ListenAndServe(":1606", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
