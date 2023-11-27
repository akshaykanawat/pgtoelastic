/*
Version 1.00
Date Created: 2023-11-28
Copyright (c) 2023, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"time"
)

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	//TODO: Fetch it from environment or AWS secrets manager
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	channels := []string{"crud_operations"} // Replace with your channel name
	_, err = db.Exec("LISTEN crud_operations")
	if err != nil {
		panic(err)
	}

	for _, channel := range channels {
		_, err = db.Exec("LISTEN " + channel)
		if err != nil {
			panic(err)
		}
	}

	notificationHandler := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			panic(err)
		}
	}

	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, notificationHandler)
	defer listener.Close()

	err = listener.Listen("crud_operations")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case notification := <-listener.Notify:
			fmt.Println("Received notification:", notification.Extra, notification.BePid)
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				listener.Ping()
			}()
		}
	}

}
