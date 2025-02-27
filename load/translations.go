package load

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Structure pour stocker les traductions
type Translations map[string]map[string]string

// Charger les traductions à partir d'un fichier JSON
func LoadTranslations(lang string) Translations {
	var translations Translations
	fileContent, err := ioutil.ReadFile("./site/infos/lang.json")
	if err != nil {
		log.Fatalf("Erreur de lecture du fichier de traductions: %v", err)
	}

	err = json.Unmarshal(fileContent, &translations)
	if err != nil {
		log.Fatalf("Erreur de parsing du JSON: %v", err)
	}

	// Si la langue spécifiée n'existe pas, utiliser "en-US" par défaut
	if _, ok := translations[lang]; !ok {
		lang = "en-US"
	}

	return translations
}
