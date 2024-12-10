package controllers

import (
	"exemple/services"
	temp "exemple/templates"
	"fmt"
	"net/http"
)

// Contrôleur qui va permettre de lister les différents albums de SCH avec les données provenant de l'API de spotify
func PageTrack(w http.ResponseWriter, r *http.Request) {
	// Appel du service qui va émettre la requête
	track, trackCode, trackErr := services.RequestTrack()
	// Vérification d'une erreur dans la réponse de l'api
	if trackErr != nil {
		// Redirection vers la page d'erreur en cas d'erreur dans la réponse
		http.Redirect(w, r, fmt.Sprintf("/error?code=%d&message=Erreur lors de la récupération des détails de la musique", trackCode), http.StatusPermanentRedirect)
		return
	}
	// Chargement et rendue du template "tracks" avec les données de l'API
	// Envoie de la réponse au client
	temp.Temp.ExecuteTemplate(w, "tracks", track)
}
