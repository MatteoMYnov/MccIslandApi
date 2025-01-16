package main

import (
	"encoding/json"
	"html/template"
	"hypixel-info/load"
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

type Infos struct {
	Name      string
	ListCapes []string
	ImageURLs []string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
}

type DataMenuPage struct {
	Name      string
	ListCapes []string
	ImageURLs []string
}

func contains(sub, str string) bool {
	return strings.Contains(str, sub)
}

// Charger les capes du fichier JSON
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

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Charger les groupes de capes depuis le JSON
	capeGroups, err := loadCapeGroups()
	if err != nil {
		log.Fatal("Erreur de chargement des groupes de capes:", err)
	}

	// Charger les capes du joueur
	IGN := r.FormValue("playername")
	var listCapes []string
	if IGN != "" {
		listCapes = load.Load(IGN)
	} else {
		listCapes = load.Load("Leroidesafk")
	}

	// Prioriser les capes en fonction des groupes
	prioritizedCapes := prioritizeCapes(listCapes, capeGroups)

	// Générer les URLs des images de capes
	imageURLs := []string{}
	for _, cape := range prioritizedCapes {
		imageURLs = append(imageURLs, "/img/capes/"+cape+".png")
	}

	// Passer les informations au template
	infos := DataMenuPage{
		Name:      IGN,
		ListCapes: prioritizedCapes,
		ImageURLs: imageURLs,
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
	load.Load(IGN)

	http.HandleFunc("/", rootHandler)

	setupFileServer("./site/styles", "/styles/")
	setupFileServer("./site/img", "/img/")
	setupFileServer("./site/js", "/js/")

	http.HandleFunc("/menu", menuHandler)

	if err := http.ListenAndServe(":1551", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
