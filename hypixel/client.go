package hypixel

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Client structure
type Client struct {
	APIKey string
}

// PlayerInfo représente les informations d'un joueur
type PlayerInfo struct {
	Name               string `json:"displayname"`
	FirstLogin         int64  `json:"firstLogin"`
	LastLogin          int64  `json:"lastLogin"`
	NewPackageRank     string `json:"newpackagerank"`
	UserLanguage       string `json:"userLanguage"`
	MostRecentGameType string `json:"mostRecentGameType"`
	BedWarsStars       int    `json:"bedwars_level"`
}

// Structure pour les réalisations
type achievementsData struct {
	BedWarsLevel int `json:"bedwars_level"`
}

type playerData struct {
	Displayname        string           `json:"displayname"`
	FirstLogin         int64            `json:"firstLogin"`
	LastLogin          int64            `json:"lastLogin"`
	NewPackageRank     string           `json:"newpackagerank"`
	UserLanguage       string           `json:"userLanguage"`
	MostRecentGameType string           `json:"mostRecentGameType"`
	Achievements       achievementsData `json:"achievements"`
}

// Response struct pour l'API
type apiResponse struct {
	Success bool       `json:"success"`
	Player  playerData `json:"player"`
}

// NewClient crée une nouvelle instance du client
func NewClient(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

// GetPlayerInfo récupère les informations de base d'un joueur
func (c *Client) GetPlayerInfo(playerUUID string) (*PlayerInfo, error) {
	// Construire l'URL
	url := fmt.Sprintf("https://api.hypixel.net/v2/player?key=%s&uuid=%s", c.APIKey, playerUUID)

	// Effectuer la requête
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("échec de la requête : %w", err)
	}
	defer resp.Body.Close()

	// Vérifier le statut HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("statut HTTP invalide : %d", resp.StatusCode)
	}

	// Décoder la réponse JSON
	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage JSON : %w", err)
	}

	// Vérifier si la réponse est réussie
	if !result.Success {
		return nil, errors.New("échec de l'API")
	}

	// Mapper les données dans PlayerInfo
	player := &PlayerInfo{
		Name:               result.Player.Displayname,
		FirstLogin:         result.Player.FirstLogin,
		LastLogin:          result.Player.LastLogin,
		NewPackageRank:     result.Player.NewPackageRank,
		UserLanguage:       result.Player.UserLanguage,
		MostRecentGameType: result.Player.MostRecentGameType,
		BedWarsStars:       result.Player.Achievements.BedWarsLevel,
	}
	if player.NewPackageRank == "" {
		player.NewPackageRank = "Default"
	}

	return player, nil
}
