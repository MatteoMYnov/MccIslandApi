package mcc

import (
	"math"
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
		return 0.00 // Évite la division par zéro
	}
	return roundFloat(float32(numerator)/float32(denominator), 2)
}

// Fonction d'arrondi à 2 chiffres après la virgule
func roundFloat(value float32, precision int) float32 {
	multiplier := float32(math.Pow(10, float64(precision)))
	return float32(math.Round(float64(value)*float64(multiplier))) / multiplier
}
