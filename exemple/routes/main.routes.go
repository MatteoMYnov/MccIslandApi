package routes

import (
	"fmt"
	"net/http"
)

// Méthode permettant d'd'initialiser le serveur, mais également d'associer l'ensemble des routes
// gerer par le serveur. Le serveur écoute sur le port :8080 du localhost de la machine
func InitServe() {
	// Récuperation des routes
	albumRoutes()
	tracksRoutes()
	errorRoutes()

	// Initialisation du serveur HTTP
	// à l'écoute sur le port 8080 du localhost de la machine
	fmt.Println("Le serveur est opérationel : http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
