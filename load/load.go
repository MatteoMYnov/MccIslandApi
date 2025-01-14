package load

import (
	"fmt"
	"log"
	"time"

	"hypixel-info/hypixel"
	"hypixel-info/minecraft"
)

func Load(Name string, HypixelAPIKey string) []string {
	PlayerUUID := minecraft.GetUUID(Name)
	PlayerCapes := minecraft.GetCapes(Name)

	client := hypixel.NewClient(HypixelAPIKey)

	playerInfo, err := client.GetPlayerInfo(PlayerUUID)
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des données : %v", err)
	}
	fmt.Println(playerInfo.Name)

	/* fmt.Printf("        Nom du joueur : %s\n", playerInfo.Name)
	fmt.Println("---------' Minecraft '---------")

	capesnb := len(PlayerCapes)
	if capesnb > 0 {
		fmt.Printf("(%d) Capes possédées : \n%s\n", capesnb, PlayerCapes)
	} else {
		fmt.Println("Aucune cape disponible.")
	}

	fmt.Println("\n----------' Hypixel '----------")
	printLoginInfo(playerInfo)
	printGameInfo(playerInfo)
	fmt.Println("-------------------------------") */
	return PlayerCapes
}

func printLoginInfo(playerInfo *hypixel.PlayerInfo) {
	fmt.Printf("Première connexion : %s\n", formatUnixTime(playerInfo.FirstLogin))
	fmt.Printf("Dernière connexion : %s\n", formatUnixTime(playerInfo.LastLogin))
	fmt.Printf("Rang : %s\n", playerInfo.NewPackageRank)
	fmt.Printf("Langue : %s\n", playerInfo.UserLanguage)
}

func printGameInfo(playerInfo *hypixel.PlayerInfo) {
	fmt.Printf("Dernier jeu joué : %s\n", playerInfo.MostRecentGameType)
	fmt.Printf("Bedwars Stars : %d ✪\n", playerInfo.BedWarsStars)
}

func formatUnixTime(timestamp int64) string {
	return time.Unix(timestamp/1000, 0).Format("02/01/2006 15:04:05")
}
