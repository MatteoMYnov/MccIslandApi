package main

import (
	"fmt"
	"html/template"
	"hypixel-info/mcc"
	"hypixel-info/minecraft"
	"log"
	"net/http"
	"path/filepath"
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
	Evolutionplus1       string
	CrownLevelplus1      string
	CrownObtained        int
	CrownObtainable      int
	CrownPourcentage     int
	CurrencyCoins        int
	CurrencyGems         int
	CurrencyRoyalRep     int
	CurrencySilver       int
	CurrencyMaterialDust int
	CurrencyAnglrTokens  int
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
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

	// Obtenez les informations de la couronne
	crownObtained := 0
	crownObtainable := 0
	if mccInfos != nil {
		crownObtained = mccInfos.CrownObtained
		crownObtainable = mccInfos.CrownObtainable
	}
	var calculatedPercent int
	if mccInfos.CrownObtainable > 0 {
		calculatedPercent = (mccInfos.CrownObtained * 100) / mccInfos.CrownObtainable
	} else {
		calculatedPercent = 0 // Ou une valeur par défaut si tu préfères
	}

	fmt.Println(mccInfos.Currency.Coins, mccInfos.Currency.Gems, mccInfos.Currency.RoyalReputation, mccInfos.Currency.Silver, mccInfos.Currency.MaterialDust, mccInfos.Currency.AnglrTokens)

	infos := DataMenuPage{
		Name:             IGN,
		ListCapes:        prioritizedCapes,
		ImageURLs:        prioritizedCapeInfos,
		BadgeURLs:        badgeInfos,
		PlayerClass:      playerClass,
		MccRank:          MccRank,
		Evolution:        fmt.Sprintf("%d", mccInfos.Evolution),
		CrownLevel:       fmt.Sprintf("%d", mccInfos.CrownLevel),
		CrownLevelplus1:  fmt.Sprintf("%d", crownLevelPlus1),
		CrownObtained:    crownObtained,
		CrownObtainable:  crownObtainable,
		CrownPourcentage: calculatedPercent,
		// Ajout des informations de currency en tant qu'entiers
		CurrencyCoins:        mccInfos.Currency.Coins,
		CurrencyGems:         mccInfos.Currency.Gems,
		CurrencyRoyalRep:     mccInfos.Currency.RoyalReputation,
		CurrencySilver:       mccInfos.Currency.Silver,
		CurrencyMaterialDust: mccInfos.Currency.MaterialDust,
		CurrencyAnglrTokens:  mccInfos.Currency.AnglrTokens,
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

func setupFileServer(path, route string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle(route, http.StripPrefix(route, fs))
}

func main() {
	http.HandleFunc("/", rootHandler)

	setupFileServer("./site/styles", "/styles/")
	setupFileServer("./site/img", "/img/")
	setupFileServer("./site/js", "/js/")

	http.HandleFunc("/menu", menuHandler)
	http.HandleFunc("/capes", capesHandler)

	if err := http.ListenAndServe(":1608", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
