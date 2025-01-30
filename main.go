package main

import (
	"encoding/json"
	"html/template"
	"hypixel-info/minecraft"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 14719ec2-7760-4b7c-a85f-f8d3775a2bb5

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

func contains(sub, str string) bool {
	return strings.Contains(str, sub)
}

func loadCapeGroups() (CapeGroups, error) {
	var capeGroups CapeGroups
	file, err := ioutil.ReadFile("./site/infos/capes.json")
	if err != nil {
		return capeGroups, err
	}

	err = json.Unmarshal(file, &capeGroups)
	if err != nil {
		return capeGroups, err
	}

	return capeGroups, nil
}

func prioritizeCapes(allCapes []string, capeGroups CapeGroups) []string {
	var prioritizedCapes []string

	for _, cape := range capeGroups.Capes {
		if containsAny(allCapes, cape.Name) {
			prioritizedCapes = append(prioritizedCapes, cape.Name)
		}
	}

	for _, cape := range allCapes {
		if !containsAny(prioritizedCapes, cape) {
			prioritizedCapes = append(prioritizedCapes, cape)
		}
	}

	return prioritizedCapes
}

func containsAny(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func getCapeClass(cape string, capeGroups CapeGroups) string {
	for _, group := range capeGroups.Capes {
		if group.Name == cape {
			return group.Type + "-cape"
		}
	}
	return ""
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
}

type DataMenuPage struct {
	Name        string
	ListCapes   []string
	ImageURLs   []CapeInfo
	BadgeURLs   []BadgeInfo
	PlayerClass string
}

func isValidIGN(name string) bool {
	validIGNPattern := "^[a-zA-Z0-9_]+$"
	matched, _ := regexp.MatchString(validIGNPattern, name)
	return matched
}

func getRandomName() string {
	var names struct {
		Name []string `json:"name"`
	}

	file, err := ioutil.ReadFile("./site/infos/names.json")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier names.json: %v", err)
	}

	err = json.Unmarshal(file, &names)
	if err != nil {
		log.Fatalf("Erreur lors du décodage du JSON: %v", err)
	}

	if len(names.Name) == 0 {
		log.Fatalf("La liste des noms est vide dans names.json")
	}

	rand.Seed(time.Now().UnixNano())
	return names.Name[rand.Intn(len(names.Name))]
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	capeGroups, err := loadCapeGroups()
	if err != nil {
		log.Fatal("Erreur de chargement des groupes de capes:", err)
	}

	IGN := r.FormValue("playername")
	if IGN == "" {
		IGN = getRandomName()
	}

	playerClass := ""
	if !isValidIGN(IGN) {
		playerClass = "badName"
	} else {
		playerClass = "playerName"
	}

	playerCapesJSON := minecraft.LoadCapesByName(IGN)
	playerBadgesJSON := minecraft.LoadBadgesByName(IGN)

	var listCapes []string
	capeInfos := []CapeInfo{}
	for _, cape := range playerCapesJSON {
		capeName := cape["cape"].(string)
		removed := cape["removed"].(bool)

		listCapes = append(listCapes, capeName)

		class := getCapeClass(capeName, capeGroups)
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

	prioritizedCapes := prioritizeCapes(listCapes, capeGroups)

	var prioritizedCapeInfos []CapeInfo
	for _, cape := range prioritizedCapes {
		for _, capeInfo := range capeInfos {
			if capeInfo.URL == "/img/capes/"+cape+".png" {
				prioritizedCapeInfos = append(prioritizedCapeInfos, capeInfo)
			}
		}
	}

	infos := DataMenuPage{
		Name:        IGN,
		ListCapes:   prioritizedCapes,
		ImageURLs:   prioritizedCapeInfos,
		BadgeURLs:   badgeInfos,
		PlayerClass: playerClass,
	}

	tmplPath := filepath.Join("site", "template", "menu.html")
	tmpl, err := template.New("menu.html").Funcs(template.FuncMap{
		"contains": contains,
	}).ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, infos)
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

	if err := http.ListenAndServe(":1512", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
