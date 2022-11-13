package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

type APIResponse struct {
	WPOSInformationTime struct {
		ListTotalCount int `json:"list_total_count"`
		RESULT         struct {
			CODE    string `json:"CODE"`
			MESSAGE string `json:"MESSAGE"`
		} `json:"RESULT"`
		Row []struct {
			MSRDATE string `json:"MSR_DATE"`
			MSRTIME string `json:"MSR_TIME"`
			SITEID  string `json:"SITE_ID"`
			WTEMP   string `json:"W_TEMP"`
			WPH     string `json:"W_PH"`
			WDO     string `json:"W_DO"`
			WTN     string `json:"W_TN"`
			WTP     string `json:"W_TP"`
			WTOC    string `json:"W_TOC"`
			WPHEN   string `json:"W_PHEN"`
			WCN     string `json:"W_CN"`
		} `json:"row"`
	} `json:"WPOSInformationTime"`
}

type HangangTemperature struct {
	Temperature string `json:"temperature"`
	MeasuredAt  string `json:"measured_at"`
}

func getHangangTemperature() (*HangangTemperature, error) {
	var parsedResponse *APIResponse

	response, err := http.Get("http://openapi.seoul.go.kr:8088/" + os.Getenv("HANGANG_API_KEY") + "/json/WPOSInformationTime/1/5/")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return nil, err
	}

	for _, data := range parsedResponse.WPOSInformationTime.Row {
		if data.SITEID == "노량진" {
			return &HangangTemperature{
				Temperature: data.WTEMP,
				MeasuredAt:  data.MSRTIME,
			}, nil
		}
	}

	return nil, errors.New("cannot get result")
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var (
			statusCode int
			response   []byte
		)

		writer.Header().Set("Content-Type", "application/json")

		hangangTemperature, err := getHangangTemperature()
		if err != nil {
			response, _ = json.Marshal(Response{
				Code: 500,
				Data: map[string]string{"message": err.Error()},
			})

			statusCode = http.StatusInternalServerError
		} else {
			response, _ = json.Marshal(Response{
				Code: 200,
				Data: hangangTemperature,
			})

			statusCode = http.StatusOK
		}
		writer.WriteHeader(statusCode)
		_, _ = writer.Write(response)

		log.Printf("%s %s %d", request.Method, request.Host, statusCode)
	})

	log.Println("On http://0.0.0.0:8000")

	panic(http.ListenAndServe(":8000", nil))
}
