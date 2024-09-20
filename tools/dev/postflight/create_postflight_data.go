package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"log"
	"os"
)

var filename = "postflightdata.zlib"

// PostflightPayload represents the payload
type PostflightPayload struct {
	RulesReceived  int `json:"rules_received"`
	RulesProcessed int `json:"rules_processed"`
}

func main() {
	payload := PostflightPayload{
		RulesReceived:  10,
		RulesProcessed: 8,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}

	var zlibBuffer bytes.Buffer
	zlibWriter := zlib.NewWriter(&zlibBuffer)
	_, err = zlibWriter.Write(jsonData)
	if err != nil {
		log.Fatalf("zlib writing failed: %s", err)
	}
	zlibWriter.Close()

	err = os.WriteFile(filename, zlibBuffer.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Writing compressed data to file failed: %s", err)
	}

	log.Printf("Postflight data written to %s\n", filename)
}
