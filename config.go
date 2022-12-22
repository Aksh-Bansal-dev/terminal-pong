package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Rows      int `json:"rows"`
	Cols      int `json:"cols"`
	BatLength int `json:"bat-length"`
	BallSpeed int `json:"ball-speed"` // 1-4
	BatSpeed  int `json:"bat-speed"`  // 1-3
}

// defaults
const (
	n         = 25
	m         = 100
	batLength = 8
	batSpeed  = 1
	ballSpeed = 2
)

func getConfig(configPath string) Config {
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("Error when opening file: ", err)
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Println("Error when opening file: ", err)
	}

	// rows
	if payload.Rows == 0 {
		payload.Rows = n
	}
	// cols
	if payload.Cols == 0 {
		payload.Cols = m
	}
	// bat length
	if payload.BatLength == 0 {
		payload.BatLength = batLength
	}
	// bat speed
	if payload.BatSpeed == 0 {
		payload.BatSpeed = batSpeed
	}
	// ball speed
	if payload.BallSpeed == 0 {
		payload.BallSpeed = ballSpeed
	}
	return payload
}
