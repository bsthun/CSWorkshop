package main

import (
	"backend/common/config"
	"backend/common/fiber"
	"backend/common/mysql"
)

func main() {
	config.Init()
	mysql.Init()
	fiber.Init()
}
