package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"hypixel-info/load"
	"hypixel-info/mcc"
	"hypixel-info/minecraft"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ImageInfo struct {
	URL string `json:"URL"`
	Alt string `json:"Alt"`
}

type ModalEntry struct {
	Name        string      `json:"name"`
	Date        string      `json:"date"`
	Images      []ImageInfo `json:"images"`
	Description string      `json:"description"`
}

type ModalData map[string]ModalEntry

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
	Capes      int           `json:"capes"`
	Score      int           `json:"score"`
	Badge      string        `json:"badge"`
	Capelist   []interface{} `json:"capelist"`
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
	//Cos
	Equipped    EquippedCosmetics
	Hats        []mcc.InvCos
	Accessories []mcc.InvCos
	Auras       []mcc.InvCos
	Trails      []mcc.InvCos
	// Player Rank
	PlayerRank     int
	PlayerRankPage int
	//lang
	ModalData    ModalData
	Lang         string
	Translations load.Translations
}

type EquippedCosmetics struct {
	Hats        Cosmetic
	Accessories Cosmetic
	Auras       Cosmetic
	Trails      Cosmetic
	Backs       Cosmetic
	Rods        Cosmetic
}

type Cosmetic struct {
	Name     string
	RealName string
	Rarity   string
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
	http.Redirect(w, r, "/en-US/menu", http.StatusFound)
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
				Evolution: friend.CrownLevel.LevelData.Evolution,
				Level:     friend.CrownLevel.LevelData.Level,
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
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	lang := parts[1]

	supportedLangs := map[string]bool{"fr-FR": true, "en-US": true, "es-ES": true, "de-DE": true}

	if _, ok := supportedLangs[lang]; !ok {
		lang = "en-US"
	}

	translations := load.LoadTranslations(lang)

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

	firstBadge := ""
	if len(playerBadgesJSON) > 0 {
		firstBadge = playerBadgesJSON[0]
	}

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
	playerRank := minecraft.UpdateClassement(playerUUID, capeDetails, actualname, firstBadge)
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

	modalFile := fmt.Sprintf("./site/infos/modals/modal_infos_%s.json", lang)
	var modalData ModalData

	if fileContent, err := ioutil.ReadFile(modalFile); err == nil {
		if err := json.Unmarshal(fileContent, &modalData); err != nil {
			log.Println("Erreur de parsing du JSON modal_infos:", err)
		}
	} else {
		log.Println("Erreur de lecture du fichier modal_infos.json:", err)
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

	var hat, accessory, aura, trail, cloak, rod Cosmetic
	if mccInfos != nil && len(mccInfos.EquippedCosmetics) > 0 {
		for _, cosmetic := range mccInfos.EquippedCosmetics {
			cleanedCosmeticName := mcc.CleanCosmeticName(cosmetic.Name)
			switch cosmetic.Category {
			case "HAT":
				hat = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			case "ACCESSORY":
				accessory = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			case "AURA":
				aura = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			case "TRAIL":
				trail = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			case "CLOAK":
				cloak = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			case "ROD":
				rod = Cosmetic{Name: cleanedCosmeticName, RealName: cosmetic.Name, Rarity: cosmetic.Rarity}
			}
		}
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
		// Cosmetics
		Equipped: EquippedCosmetics{
			Hats:        hat,
			Accessories: accessory,
			Auras:       aura,
			Trails:      trail,
			Backs:       cloak,
			Rods:        rod,
		},
		Hats:        mccInfos.Hats,
		Accessories: mccInfos.Accessories,
		Auras:       mccInfos.Auras,
		Trails:      mccInfos.Trails,
		// Player Rank
		PlayerRank:     playerRank,
		PlayerRankPage: playerRankPage,
		//lang
		ModalData:    modalData,
		Lang:         lang,         // Ajoutez la langue ici
		Translations: translations, // Ajoutez les traductions ici
	}

	tmplPath := filepath.Join("site", "template", "menu.html")
	tmpl, err := template.New("menu.html").Funcs(template.FuncMap{
		"contains": minecraft.Contains,
		"mod":      minecraft.Mod,
		"seq":      minecraft.Seq,
		"sub":      minecraft.Sub,
		"add":      minecraft.Add,
		"mul":      minecraft.Mul,
		"toJson":   minecraft.ToJSON,
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, infos)
}

func classementHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire les parties de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "URL mal formée", http.StatusBadRequest)
		return
	}

	// Extraire la langue et le numéro de page
	lang := parts[1] // Langue (par exemple 'fr-FR' ou 'en-US')
	var page int
	_, err := fmt.Sscanf(parts[3], "%d", &page)
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

	// Charger les traductions en fonction de la langue
	supportedLangs := map[string]bool{"fr-FR": true, "en-US": true, "es-ES": true, "de-DE": true}
	if _, ok := supportedLangs[lang]; !ok {
		lang = "en-US" // Langue par défaut
	}

	translations := load.LoadTranslations(lang)

	// Charger les capes
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

	// Extraire les capes demandées dans l'URL
	capesParam := r.URL.Query().Get("capes")
	var requestedCapes []string
	if capesParam != "" {
		requestedCapes = strings.Split(capesParam, ",")
	}

	// Avant de filtrer, mettez à jour les rangs de tous les joueurs
	for i := range classement.Joueurs {
		classement.Joueurs[i].Rank = i + 1 // Récupère le rang global
	}

	// Filtrer tous les joueurs en fonction des capes avant de paginer
	var filteredPlayers []Joueur
	for _, joueur := range classement.Joueurs {
		// Vérifier si le joueur possède toutes les capes demandées
		hasAllRequestedCapes := true
		for _, cape := range requestedCapes {
			if !playerHasCape(joueur, cape) {
				hasAllRequestedCapes = false
				break
			}
		}

		// Ajouter le joueur à la liste filtrée s'il possède toutes les capes
		if hasAllRequestedCapes {
			filteredPlayers = append(filteredPlayers, joueur)
		}
	}

	// Si aucun joueur ne correspond au filtre, rediriger vers la page 1 sans filtre
	if len(filteredPlayers) == 0 {
		http.Redirect(w, r, fmt.Sprintf("/%s/classement/1", lang), http.StatusFound)
		return
	}

	// Définir la taille du groupe de joueurs par page
	const pageSize = 50
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// Vérifier que l'index ne dépasse pas la taille du classement filtré
	if startIndex >= len(filteredPlayers) {
		http.Error(w, "Page inexistante", http.StatusNotFound)
		return
	}

	if endIndex > len(filteredPlayers) {
		endIndex = len(filteredPlayers)
	}

	// Extraire les joueurs pour la page demandée
	joueursPage := filteredPlayers[startIndex:endIndex]

	// Mettre à jour le rang des joueurs filtrés pour refléter leur position globale
	for i := range joueursPage {
		// Utilisez la position originale du joueur dans la liste complète
		originalPlayer := findPlayerByUUID(classement.Joueurs, joueursPage[i].UUID)
		joueursPage[i].Rank = originalPlayer.Rank
	}

	// Renvoyer la page HTML avec les données du classement
	tmplPath := filepath.Join("site", "template", "classement.html")
	tmpl, err := template.New("classement.html").Funcs(template.FuncMap{
		"contains": minecraft.Contains,
		"mod":      minecraft.Mod,
		"seq":      minecraft.Seq,
		"sub":      minecraft.Sub,
		"add":      minecraft.Add,
		"mul":      minecraft.Mul,
		"toJson":   minecraft.ToJSON,
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Classement   []Joueur
		Page         int
		HasNext      bool
		HasPrev      bool
		Lang         string
		Translations load.Translations // Ajouter ici les traductions
		ImageURLs    []CapeInfo        // Ajouter les capes
	}{
		Classement:   joueursPage,
		Page:         page,
		HasPrev:      page > 1,
		HasNext:      endIndex < len(filteredPlayers),
		Lang:         lang,
		Translations: translations, // Passer les traductions au template
		ImageURLs:    capeInfos,    // Passer les capes au template
	}

	// Exécuter le template avec les données
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findPlayerByUUID(players []Joueur, uuid string) *Joueur {
	for i := range players {
		if players[i].UUID == uuid {
			return &players[i]
		}
	}
	return nil
}

func playerHasCape(joueur Joueur, capeName string) bool {
	// Parcourir la liste des capes du joueur
	for _, cape := range joueur.Capelist {
		// Assumer que chaque élément de Capelist est une structure contenant un champ Name
		if capeMap, ok := cape.(map[string]interface{}); ok {
			if name, exists := capeMap["name"].(string); exists && name == capeName {
				return true
			}
		}
	}
	return false
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
	setupFileServer("./site/sounds", "/sounds/")

	// Définir les routes pour chaque langue
	http.HandleFunc("/en-US/menu", menuHandler)
	http.HandleFunc("/fr-FR/menu", menuHandler)
	http.HandleFunc("/es-ES/menu", menuHandler)
	http.HandleFunc("/de-DE/menu", menuHandler)

	http.HandleFunc("/en-US/classement/", classementHandler)
	http.HandleFunc("/fr-FR/classement/", classementHandler)
	http.HandleFunc("/es-ES/classement/", classementHandler)
	http.HandleFunc("/de-DE/classement/", classementHandler)

	http.HandleFunc("/dbdl", downloadFileHandler)

	if err := http.ListenAndServe(":1637", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
