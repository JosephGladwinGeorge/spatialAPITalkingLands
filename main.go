package main

import (
	"spatialDB/backend"
)

func main(){
	app:=backend.App{
		Port: "3000",
	}
	app.Initialize()
	app.Run()
}