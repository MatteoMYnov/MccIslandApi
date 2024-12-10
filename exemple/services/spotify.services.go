package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Déclaration du client qui va émettre les requêtes
var _httpClient = http.Client{
	Timeout: 5 * time.Second,
}

// Initialisation des variables pour s'authentifier à l'API
var _clientId string = ""
var _clientSecret string = ""

// Valeur par défaut pour éviter l'erreur 400 due à un problème de format, type...
var _token string = "Bearer "

type DataToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func requestToken() error {
	body := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", _clientId, _clientSecret))

	req, reqErr := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", body)
	if reqErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors de l'envois de la réquête")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("RequestToken - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var data DataToken

	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return fmt.Errorf("RequestToken - Erreur lors du décodage des données")
	}

	_token = fmt.Sprintf("%s %s", data.TokenType, data.AccessToken)
	return nil
}

func requestSpotifyGet(url string) (*http.Response, error) {
	fmt.Printf("valeur token : %s\n", _token)
	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("Erreur lors de l'initialisation de la réquête")
	}

	req.Header.Set("Authorization", _token)

	res, resErr := _httpClient.Do(req)
	if resErr != nil {
		return nil, fmt.Errorf("Erreur lors de l'envois de la réquête")
	}
	fmt.Println(res.Status, res.StatusCode, url)
	if res.StatusCode == 401 {
		errToken := requestToken()
		fmt.Println(errToken)
		if errToken != nil {
			return nil, fmt.Errorf("Erreur lors de la récupération du token")
		}
		return requestSpotifyGet(url)
	}
	return res, nil
}

type DataAlbums struct {
	Items []struct {
		TotalTracks int `json:"total_tracks"`
		Images      []struct {
			Url string `json:"url"`
		} `json:"images"`
		Name       string `json:"name"`
		ReleasDate string `json:"release_date"`
	} `json:"items"`
}

func RequestAlbums() (DataAlbums, int, error) {
	res, resErr := requestSpotifyGet("https://api.spotify.com/v1/artists/2kXKa3aAFngGz2P4GjG5w2/albums")
	if resErr != nil {
		return DataAlbums{}, 500, fmt.Errorf("RequestAlbums - %s", resErr.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return DataAlbums{}, res.StatusCode, fmt.Errorf("RequestAlbums - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var data DataAlbums

	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return DataAlbums{}, 500, fmt.Errorf("RequestAlbums - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return data, res.StatusCode, nil
}

type DataTrack struct {
	NameTack  string `json:"name"`
	AlbumData struct {
		ArtistsData []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Images []struct {
			Url string `json:"url"`
		}
		Name       string `json:"name"`
		ReleasDate string `json:"release_date"`
	} `json:"album"`
	Url struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}

func RequestTrack() (DataTrack, int, error) {
	res, resErr := requestSpotifyGet("https://api.spotify.com/v1/tracks/1Mzg6bu3hkCwJKEf7v49MN")
	if resErr != nil {
		return DataTrack{}, 500, fmt.Errorf("RequestTrack - %s", resErr.Error())
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return DataTrack{}, res.StatusCode, fmt.Errorf("RequestTrack - Erreur dans la réponse de la requête code : %d", res.StatusCode)
	}

	var data DataTrack

	decodeErr := json.NewDecoder(res.Body).Decode(&data)
	if decodeErr != nil {
		return DataTrack{}, 500, fmt.Errorf("RequestTrack - Erreur lors du décodage des données : %s", decodeErr.Error())
	}

	return data, res.StatusCode, nil
}
