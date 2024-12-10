package routes

import (
	"exemple/controllers"
	"net/http"
)

// Méthode permettant d'initialiser les routes lié à l'affichage des personnages
func albumRoutes() {
	// Route permettant d'afficher sous forme de liste les albums de SCH
	// La route /albums est associé au contrôleur PageListAlbums
	http.HandleFunc("/albums", controllers.PageListAlbums)
}
