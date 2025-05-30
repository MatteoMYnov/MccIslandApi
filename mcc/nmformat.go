package mcc

import (
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Formater un nombre avec des espaces comme séparateurs de milliers
func FormatNumberWithSpaces(number int) string {
	// Convertir le nombre en chaîne de caractères
	numStr := strconv.Itoa(number)

	// Créer un slice pour stocker les groupes de chiffres
	var groups []string
	for i := len(numStr); i > 0; i -= 3 {
		// Prendre les 3 derniers chiffres (ou moins si la longueur du nombre est inférieure à 3)
		start := i - 3
		if start < 0 {
			start = 0
		}
		groups = append([]string{numStr[start:i]}, groups...)
	}

	// Joindre les groupes avec un espace
	return strings.Join(groups, " ")
}

func safeDivide(numerator, denominator int) float32 {
	if denominator == 0 {
		return roundFloat(float32(numerator)/1, 2)
	}
	return roundFloat(float32(numerator)/float32(denominator), 2)
}

// Fonction d'arrondi à 2 chiffres après la virgule
func roundFloat(value float32, precision int) float32 {
	multiplier := float32(math.Pow(10, float64(precision)))
	return float32(math.Round(float64(value)*float64(multiplier))) / multiplier
}

func CleanCosmeticName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`) // Expression régulière pour tout sauf lettres et chiffres
	return re.ReplaceAllString(name, "")     // Remplacer tout ce qui n'est pas une lettre ou un chiffre par rien
}

func SortCosmetics(cosmetics []InvCos) []InvCos {
	// Définir l'ordre des rarités
	rarityOrder := map[string]int{
		"COMMON":    0,
		"UNCOMMON":  1,
		"RARE":      2,
		"EPIC":      3,
		"LEGENDARY": 4,
		"MYTHIC":    5,
	}

	// Trier en utilisant sort.SliceStable
	sort.SliceStable(cosmetics, func(i, j int) bool {
		// Trier d'abord par owned (true avant false)
		if cosmetics[i].Owned != cosmetics[j].Owned {
			return cosmetics[i].Owned
		}
		// Puis par rareté en utilisant rarityOrder
		if rarityOrder[cosmetics[i].Rarity] != rarityOrder[cosmetics[j].Rarity] {
			return rarityOrder[cosmetics[i].Rarity] < rarityOrder[cosmetics[j].Rarity]
		}
		// Enfin par nom (ordre alphabétique)
		return cosmetics[i].Name < cosmetics[j].Name
	})

	return cosmetics
}
