package load

import (
	"fmt"
	"time"

	"hypixel-info/hypixel"
	"hypixel-info/minecraft"
)

func Load(Name string) []map[string]interface{} {
	// PlayerUUID := minecraft.GetUUID(Name)
	PlayerCapes := minecraft.GetCapes(Name)
	PlayerBadges := minecraft.LoadBadges(PlayerCapes)
	fmt.Print(PlayerBadges)
	// client := hypixel.NewClient(HypixelAPIKey)

	// playerInfo, err := client.GetPlayerInfo(PlayerUUID)
	// if err != nil {
	// 	log.Fatalf("Erreur lors de la récupération des données : %v", err)
	// }

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
	var capesList []map[string]interface{}
	for _, cape := range PlayerCapes {
		capeObj := map[string]interface{}{
			"cape":    cape["cape"],    // Nom de la cape
			"removed": cape["removed"], // Statut de la cape
		}
		capesList = append(capesList, capeObj)
	}
	return capesList
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
