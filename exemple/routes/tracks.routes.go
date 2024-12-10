package routes

import (
	"exemple/controllers"
	"net/http"
)

// Méthode permettant d'initialiser les routes lié à l'affichage des personnages
func tracksRoutes() {
	// Route permettant d'afficher les détails de la musique "CARTIER SANTOS" de SDM
	// La route /tracks est associé au contrôleur PageTrack
	http.HandleFunc("/tracks", controllers.PageTrack)
}
