package main

import (
	"github.com/joho/godotenv"
	"github.com/karthxk07/sol-III/internals/server"
)

func main() {
	//Loading the environment variables
	godotenv.Load()

	//init a new router and listen at port 3030
	r := server.New()
	r.Run(":3030")
}
