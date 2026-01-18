package main

import (
	"log"
)

func (sp *ServerPool) HealthCheck() {
	for _, b := range sp.Backends {
		status := "UP"
		alive:=b.GetRealStatus()
		b.SetAlive(alive)
		if !b.GetRealStatus() {
			status = "DOWN"
		}
		log.Printf("%s is %s", b.URL, status)

	}
}
