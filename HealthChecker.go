package main

import (
	"log"
)

func (sp *ServerPool) HealthCheck() {
	for _, b := range sp.Backends {
		status := "up"
		alive:=b.GetRealStatus()
		b.SetAlive(alive)
		if !b.GetRealStatus() {
			status = "down"
		}
		log.Printf("%s [%s]", b.URL, status)

	}
}
