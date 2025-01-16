package main

import (
	"html/template"
	"hypixel-info/load"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var IGN string = "Leroidesafk"

//var HypixelAPIKey string = "e6fbfd94-d9c9-4638-8503-7e9248bb26d1"

type Infos struct {
	Name string
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

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Charger les capes par défaut pour Leroidesafk
	infos := DataMenuPage{"Leroidesafk", load.Load("Leroidesafk"), []string{}}

	// Si un joueur est spécifié, charger ses capes
	IGN := r.FormValue("playername")
	if IGN != "" {
		listCapes := load.Load(IGN)
		imageURLs := []string{}
		if len(listCapes) == 0 {
			infos.ImageURLs = []string{}                    // Pas de capes
			infos.Name = "Ce Joueur ne Possède aucune cape" // Message alternatif
		} else {
			for _, cape := range listCapes {
				imageURLs = append(imageURLs, "/img/capes/"+cape+".png")
			}
			infos.ImageURLs = imageURLs
		}
		infos.Name = IGN
	} else {
		// Si aucun joueur spécifié, afficher les capes par défaut de Leroidesafk
		listCapes := load.Load("Leroidesafk")
		imageURLs := []string{}
		for _, cape := range listCapes {
			imageURLs = append(imageURLs, "/img/capes/"+cape+".png")
		}
		infos.ImageURLs = imageURLs
	}

	// Passer la fonction contains au template
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
