package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/temperature", temperatureHandler)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		panic(err)
	}
}

func temperatureHandler(writer http.ResponseWriter, _ *http.Request) {
	now := getCurrentDate()
	temp := getTemperature()
	response := buildJsonResponse(temp, now)
	writeResponse(writer, response)
}

func getCurrentDate() []byte {
	now, err := time.Now().MarshalText()
	if err != nil {
		panic(err.Error())
	}
	return now
}

func buildJsonResponse(temp string, now []byte) []byte {
	temperature := Temperature{
		Value: temp,
		Date:  string(now),
	}
	response, err := json.Marshal(temperature)
	if err != nil {
		panic(err.Error())
	}
	return response
}

func getTemperature() string {
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "vcgencmd measure_temp | egrep -o '[0-9]*\\.[0-9]*'")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		panic(err.Error())
	}
	return strings.TrimRight(stdout.String(), "\n")
}

func writeResponse(writer http.ResponseWriter, response []byte) {
	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err := writer.Write(response)
	if err != nil {
		panic(err.Error())
	}
}

type Temperature struct {
	Value string `json:"value"`
	Date  string `json:"date"`
}

