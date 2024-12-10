package routes

import (
	"exemple/controllers"
	"net/http"
)

// Méthode permettant d'initialiser les routes lié à la gestion des erreurs
// soit aux différentes pages d'erreurs
func errorRoutes() {
	// route permettant d'd'accéder à une page d'erreur
	// la route /error est associé au contrôleur ErrorController
	http.HandleFunc("/error", controllers.ErrorController)
}
