package controllers

import (
	"exemple/services"
	temp "exemple/templates"
	"fmt"
	"net/http"
)

// Contrôleur qui va permettre de lister les différents albums de SCH avec les données provenant de l'API de spotify
func PageListAlbums(w http.ResponseWriter, r *http.Request) {
	// Appel du service qui va émettre la requête
	listAlbums, listAlbumsCode, listAlbumsErr := services.RequestAlbums()
	// Vérification d'une erreur dans la réponse de l'api
	if listAlbumsErr != nil {
		// Redirection vers la page d'erreur en cas d'erreur dans la réponse
		http.Redirect(w, r, fmt.Sprintf("/error?code=%d&message=Erreur lors de la récupération des albums", listAlbumsCode), http.StatusPermanentRedirect)
		return
	}
	// Chargement et rendue du template "albums" avec les données de l'API
	// Envoie de la réponse au client
	temp.Temp.ExecuteTemplate(w, "albums", listAlbums)
}
