package main

import (
	"exemple/routes"
	temp "exemple/templates"
)

func main() {
	temp.InitTemplates()
	routes.InitServe()
}
