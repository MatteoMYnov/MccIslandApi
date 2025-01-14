package main

import (
	"html/template"
	"hypixel-info/load"
	"log"
	"net/http"
	"path/filepath"
)

var IGN string = "Leroidesafk"
var HypixelAPIKey string = "6031e81e-677f-4922-b46c-5149f06b0f9c"

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

func menuHandler(w http.ResponseWriter, r *http.Request) {
	// Charger les capes par défaut pour Leroidesafk
	infos := DataMenuPage{"Leroidesafk", load.Load("Leroidesafk", HypixelAPIKey), []string{}}

	// Si un joueur est spécifié, charger ses capes
	IGN := r.FormValue("playername")
	if IGN != "" {
		listCapes := load.Load(IGN, HypixelAPIKey)
		imageURLs := []string{}
		for _, cape := range listCapes {
			imageURLs = append(imageURLs, "/img/capes/"+cape+".png")
		}
		infos = DataMenuPage{IGN, listCapes, imageURLs}
	} else {
		// Si aucun joueur spécifié, afficher les capes par défaut de Leroidesafk
		listCapes := load.Load("Leroidesafk", HypixelAPIKey)
		imageURLs := []string{}
		for _, cape := range listCapes {
			imageURLs = append(imageURLs, "/img/capes/"+cape+".png")
		}
		infos.ImageURLs = imageURLs
	}

	tmplPath := filepath.Join("site", "template", "menu.html")
	tmpl, err := template.ParseFiles(tmplPath)
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

	load.Load(IGN, HypixelAPIKey)

	http.HandleFunc("/", rootHandler)

	setupFileServer("./site/styles", "/styles/")
	setupFileServer("./site/img", "/img/")
	setupFileServer("./site/js", "/js/")

	http.HandleFunc("/menu", menuHandler)

	log.Println("Démarrage du serveur sur le port 1551...")
	if err := http.ListenAndServe(":1551", nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}
