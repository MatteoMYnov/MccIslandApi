package main

import (
	"encoding/json"
	"html/template"
	"hypixel-info/load"
	"hypixel-info/minecraft"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var IGN string = "Leroidesafk"

// Structure pour contenir les groupes de capes
type CapeGroups struct {
	Special []string `json:"special"`
	Normal  []string `json:"normal"`
	Common  []string `json:"common"`
}

// Structure pour contenir les informations d'une cape avec URL et classe CSS
type CapeInfo struct {
	URL      string
	Class    string
	CapeName string
	Removed  bool
}

// Structure pour contenir les informations d'un badge avec URL et classe CSS
type BadgeInfo struct {
	URL       string
	Class     string
	BadgeName string
}

type Infos struct {
	Name      string
	ListCapes []string
	ImageURLs []CapeInfo
	BadgeURLs []BadgeInfo // Nouvelle section pour inclure les badges
}

// Fonction pour vérifier si une chaîne est présente dans un tableau
func contains(sub, str string) bool {
	return strings.Contains(str, sub)
}

// Charger les groupes de capes depuis le fichier JSON
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

// Trier les capes selon les groupes définis
func prioritizeCapes(allCapes []string, capeGroups CapeGroups) []string {
	var prioritizedCapes []string
	// Prioriser les capes par groupe
	for _, group := range [][]string{capeGroups.Special, capeGroups.Normal, capeGroups.Common} {
		for _, cape := range group {
			if containsAny(allCapes, cape) {
				prioritizedCapes = append(prioritizedCapes, cape)
			}
		}
	}
	// Ajouter le reste des capes non prioritaires
	for _, cape := range allCapes {
		if !containsAny(prioritizedCapes, cape) {
			prioritizedCapes = append(prioritizedCapes, cape)
		}
	}
	return prioritizedCapes
}

// Vérifier si une cape existe dans un tableau
func containsAny(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// Fonction pour obtenir la classe CSS en fonction du groupe de la cape
func getCapeClass(cape string, capeGroups CapeGroups) string {
	if containsAny(capeGroups.Special, cape) {
		return "special-cape"
	} else if containsAny(capeGroups.Normal, cape) {
		return "normal-cape"
	} else if containsAny(capeGroups.Common, cape) {
		return "common-cape"
	}
	return "" // Pas de classe spécifique si la cape n'est dans aucun groupe
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
}

type DataMenuPage struct {
	Name      string
	ListCapes []string
	ImageURLs []CapeInfo
	BadgeURLs []BadgeInfo // Informations pour les badges
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Charger les groupes de capes depuis le fichier JSON
	capeGroups, err := loadCapeGroups()
	if err != nil {
		log.Fatal("Erreur de chargement des groupes de capes:", err)
	}

	// Si aucun nom de joueur n'est spécifié dans la requête, utiliser "Leroidesafk" par défaut
	IGN := r.FormValue("playername")
	if IGN == "" {
		IGN = "Leroidesafk"
	}

	// Charger les capes du joueur (c'est une liste de JSON)
	playerCapesJSON := load.Load(IGN)
	playerBadgesJSON := minecraft.LoadBadgesByName(IGN) // Charger les badges du joueur

	// Extraire les noms des capes et leur état "removed"
	var listCapes []string
	capeInfos := []CapeInfo{}
	for _, cape := range playerCapesJSON {
		capeName := cape["cape"].(string)
		removed := cape["removed"].(bool)

		// Ajouter la cape à la liste
		listCapes = append(listCapes, capeName)

		// Obtenir la classe CSS de base pour la cape
		class := getCapeClass(capeName, capeGroups)

		// Si la cape est supprimée, ajouter la classe "removed-cape"
		if removed {
			class = class + " removed-cape"
		}

		// Ajouter l'info de la cape avec URL et classe
		capeInfos = append(capeInfos, CapeInfo{
			URL:      "/img/capes/" + capeName + ".png",
			Class:    class,
			CapeName: capeName,
			Removed:  removed,
		})
	}

	// Gérer les badges du joueur
	badgeInfos := []BadgeInfo{}
	for _, badgeName := range playerBadgesJSON { // "playerBadgesJSON" est directement un tableau de strings
		// Ajouter l'info du badge avec URL et classe CSS
		badgeInfos = append(badgeInfos, BadgeInfo{
			URL:       "/img/badges/" + badgeName + ".png",
			Class:     "badge-class", // Ajouter une classe par défaut ou dynamique
			BadgeName: badgeName,
		})
	}

	// Prioriser les capes en fonction des groupes
	prioritizedCapes := prioritizeCapes(listCapes, capeGroups)

	// Réorganiser les informations de capes en fonction de la priorité
	var prioritizedCapeInfos []CapeInfo
	for _, cape := range prioritizedCapes {
		for _, capeInfo := range capeInfos {
			if capeInfo.URL == "/img/capes/"+cape+".png" {
				prioritizedCapeInfos = append(prioritizedCapeInfos, capeInfo)
			}
		}
	}

	// Passer les informations au template
	infos := DataMenuPage{
		Name:      IGN,
		ListCapes: prioritizedCapes,
		ImageURLs: prioritizedCapeInfos,
		BadgeURLs: badgeInfos,
	}

	// Charger et exécuter le template
	tmplPath := filepath.Join("site", "template", "menu.html")
	tmpl, err := template.New("menu.html").Funcs(template.FuncMap{
		"contains": contains, // Ajouter la fonction personnalisée
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

	if err := http.ListenAndServe(":1556", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
