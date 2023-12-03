/*
Version 1.00
Date Created: 2022-06-25
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"log"
	"pgsync/server/app/routes"
	"pgsync/server/config"
)

func main() {

	r := routes.SetupServer()
	log.Print("Server online. Enjoy the ride....")
	err := r.Run(config.PORT)
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
