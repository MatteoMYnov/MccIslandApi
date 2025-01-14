package main

import (
	"html/template"
	"hypixel-info/load"
	"log"
	"net/http"
	"path/filepath"
)

type Infos struct {
	Name string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/menu", http.StatusFound)
}

func menuHandler(w http.ResponseWriter, r *http.Request) {
	infos := Infos{Name: ""}

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
	IGN := "Leroidesafk"
	HypixelAPIKey := "74a12d16-fc4a-4446-8ef2-2c666606a01e"
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
