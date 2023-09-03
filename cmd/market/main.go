package main

import "market/internal/app"

const configDir = "./configs"

func main() {
	app.Run(configDir)
}
