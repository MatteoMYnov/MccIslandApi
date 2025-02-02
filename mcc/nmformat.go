package mcc

import (
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
