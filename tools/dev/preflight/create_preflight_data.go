package main

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"log"
	"os"
)

var filename = "preflightdata.zlib"

// ClientMode is a placeholder type
type ClientMode string

// PreflightPayload represents the payload
type PreflightPayload struct {
	SerialNumber         string     `json:"serial_num"`
	Hostname             string     `json:"hostname"`
	OSVersion            string     `json:"os_version"`
	OSBuild              string     `json:"os_build"`
	ModelIdentifier      string     `json:"model_identifier"`
	SantaVersion         string     `json:"santa_version"`
	PrimaryUser          string     `json:"primary_user"`
	BinaryRuleCount      int        `json:"binary_rule_count"`
	CertificateRuleCount int        `json:"certificate_rule_count"`
	CompilerRuleCount    int        `json:"compiler_rule_count"`
	TransitiveRuleCount  int        `json:"transitive_rule_count"`
	TeamIDRuleCount      int        `json:"teamid_rule_count"`
	SigningIDRuleCount   int        `json:"signingid_rule_count"`
	CdHashRuleCount      int        `json:"cdhash_rule_count"`
	ClientMode           ClientMode `json:"client_mode"`
	RequestCleanSync     bool       `json:"request_clean_sync"`
}

func main() {
	payload := PreflightPayload{
		SerialNumber:         "KDQL6DT633",
		Hostname:             "MacBook-Pro",
		OSVersion:            "14.2.1",
		OSBuild:              "23C71",
		ModelIdentifier:      "MacBookPro18,2",
		SantaVersion:         "2023.9.579822660",
		PrimaryUser:          "foobar",
		BinaryRuleCount:      243,
		CertificateRuleCount: 619,
		CompilerRuleCount:    2,
		TransitiveRuleCount:  4,
		TeamIDRuleCount:      1,
		SigningIDRuleCount:   3,
		CdHashRuleCount:      6,
		ClientMode:           "LOCKDOWN",
		RequestCleanSync:     true,
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

	log.Printf("Preflight data written to %s\n", filename)
}
