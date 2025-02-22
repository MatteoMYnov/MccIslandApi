package minecraft

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func Contains(sub, str string) bool {
	return strings.Contains(str, sub)
}

func GetRandomName() string {
	var names struct {
		Name []string `json:"name"`
	}

	file, err := ioutil.ReadFile("./site/infos/names.json")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier names.json: %v", err)
	}

	err = json.Unmarshal(file, &names)
	if err != nil {
		log.Fatalf("Erreur lors du décodage du JSON: %v", err)
	}

	if len(names.Name) == 0 {
		log.Fatalf("La liste des noms est vide dans names.json")
	}

	rand.Seed(time.Now().UnixNano())
	return names.Name[rand.Intn(len(names.Name))]
}

func LoadCapeGroups() (CapeGroups, error) {
	var capeGroups CapeGroups
	file, err := ioutil.ReadFile("./site/infos/capes.json")
	if err != nil {
		return capeGroups, err
	}

	err = json.Unmarshal(file, &capeGroups)
	if err != nil {
		return capeGroups, err
	}

	return capeGroups, nil
}

func PrioritizeCapes(allCapes []string, capeGroups CapeGroups) []string {
	var prioritizedCapes []string

	for _, cape := range capeGroups.Capes {
		if containsAny(allCapes, cape.Name) {
			prioritizedCapes = append(prioritizedCapes, cape.Name)
		}
	}

	for _, cape := range allCapes {
		if !containsAny(prioritizedCapes, cape) {
			prioritizedCapes = append(prioritizedCapes, cape)
		}
	}

	return prioritizedCapes
}

func containsAny(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func GetCapeClass(cape string, capeGroups CapeGroups) string {
	for _, group := range capeGroups.Capes {
		if group.Name == cape {
			return group.Type + "-cape"
		}
	}
	return ""
}

func IsValidIGN(name string) bool {
	validIGNPattern := "^[a-zA-Z0-9_]+$"
	matched, _ := regexp.MatchString(validIGNPattern, name)
	return matched
}

func Mod(a, b int) int {
	return a % b
}

func Seq(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func Sub(a, b int) int {
	return a - b
}

func Add(x, y int) int {
	return x + y
}

func Mul(x, y int) int {
	return x * y
}

func ToJSON(data interface{}) template.JS {
	bytes, err := json.Marshal(data)
	if err != nil {
		return template.JS("[]") // Retourne un tableau vide si erreur
	}
	return template.JS(bytes)
}
