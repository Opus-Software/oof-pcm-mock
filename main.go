package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Report struct {
	ClientOrgId       string `json:"clientOrgId"`
	CorrelationId     string `json:"correlationId"`
	Endpoint          string `json:"endpoint"`
	FapiInteractionId string `json:"fapiInteractionId"`
	HttpMethod        string `json:"httpMethod"`
	ServerOrgId       string `json:"serverOrgId"`
	Timestamp         string `json:"timestamp"`
	ProcessTimespam   int    `json:"processTimespan"`
	StatusCode        int    `json:"statusCode"`
}

type ReportResponseSuccess struct {
	ReportId      string `json:"reportId"`
	Status        string `json:"status"`
	CorrelationId string `json:"correlationId,omitempty"`
	Message       string `json:"message,omitempty"`
}

type ReportResponseError struct {
	Message string `json:"message"`
}

type TokenResponseSuccess struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func main() {
	address := "localhost:8090"
	fmt.Printf("Started Mock PCM Server Tool on %s...\n", address)

	gin.DefaultWriter = io.Discard
	r := gin.Default()
	r.POST("/report-api/v1/private/report", handlePostPrivateReport)
	r.POST("/report-api/v1/opendata/report", handlePostOpenDataReport)
	r.POST("/token/", handlePostToken)
	r.Run(address)
}

func handlePostPrivateReport(c *gin.Context) {
	fmt.Println("Handling new post private report")

	var requestReports []Report
	var responseReports []ReportResponseSuccess

	expectedResponseStatus := "200"
	possibleStatuses := [5]string{"ACCEPTED", "DISCARDED", "PAIRED_INCONSISTENT", "UNPAIRED", "SINGLE"}

	if len(os.Args) > 1 {
		expectedResponseStatus = os.Args[1]
	}

	_ = c.Bind(&requestReports)
	responseStatus, _ := strconv.Atoi(expectedResponseStatus)

	if responseStatus == 200 || responseStatus == 207 {
		fmt.Printf("New post with %d incoming events\n", len(requestReports))

		fmt.Printf("Request headers received:\n")
		for k, vals := range c.Request.Header {
			fmt.Printf("\n%s", k)
			for _, v := range vals {
				fmt.Printf("\t%s ", v)
			}
		}
		fmt.Printf("\n\n")

		for _, reqReport := range requestReports {
			resReport := &ReportResponseSuccess{
				ReportId:      uuid.New().String(),
				Status:        possibleStatuses[rand.Intn(len(possibleStatuses))],
				CorrelationId: reqReport.CorrelationId,
			}

			if expectedResponseStatus == "207" {
				shouldSetMessageField := rand.Intn(2)
				if shouldSetMessageField == 1 {
					resReport.Message = "Missing fields message example"
				}
			}

			responseReports = append(responseReports, *resReport)
		}

		out, _ := json.MarshalIndent(responseReports, "", "\t")
		fmt.Printf("Response:\n\n%s\n\n", string(out))

		c.JSON(responseStatus, responseReports)
	} else {
		errorResponse := &ReportResponseError{
			Message: buildErrorMessage(responseStatus),
		}
		out, _ := json.MarshalIndent(errorResponse, "", "\t")
		fmt.Printf("Response:\n\n%s\n\n", string(out))
		c.JSON(responseStatus, errorResponse)
	}
}

func handlePostOpenDataReport(c *gin.Context) {
	fmt.Println("Handling new post open data report")

	var requestReports []Report
	var responseReports []ReportResponseSuccess

	expectedResponseStatus := "200"
	possibleStatuses := [5]string{"ACCEPTED", "DISCARDED", "PAIRED_INCONSISTENT", "UNPAIRED", "SINGLE"}

	if len(os.Args) > 1 {
		expectedResponseStatus = os.Args[1]
	}

	_ = c.Bind(&requestReports)
	responseStatus, _ := strconv.Atoi(expectedResponseStatus)

	if responseStatus == 200 || responseStatus == 207 {
		fmt.Printf("New post with %d incoming events\n", len(requestReports))

		fmt.Printf("Request headers received:\n")
		for k, vals := range c.Request.Header {
			fmt.Printf("\n%s", k)
			for _, v := range vals {
				fmt.Printf("\t%s ", v)
			}
		}
		fmt.Printf("\n\n")

		for _, _ = range requestReports {
			resReport := &ReportResponseSuccess{
				ReportId: uuid.New().String(),
				Status:   possibleStatuses[rand.Intn(len(possibleStatuses))],
			}

			if expectedResponseStatus == "207" {
				shouldSetMessageField := rand.Intn(2)
				if shouldSetMessageField == 1 {
					resReport.Message = "Missing fields message example"
				}
			}

			responseReports = append(responseReports, *resReport)
		}

		out, _ := json.MarshalIndent(responseReports, "", "\t")
		fmt.Printf("Response:\n\n%s\n\n", string(out))

		c.JSON(responseStatus, responseReports)
	} else {
		errorResponse := &ReportResponseError{
			Message: buildErrorMessage(responseStatus),
		}
		out, _ := json.MarshalIndent(errorResponse, "", "\t")
		fmt.Printf("Response:\n\n%s\n\n", string(out))
		c.JSON(responseStatus, errorResponse)
	}
}

func handlePostToken(c *gin.Context) {
	tokenResponse := &TokenResponseSuccess{
		AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkFiT",
		ExpiresIn:   21600,
		TokenType:   "Bearer",
	}
	out, _ := json.MarshalIndent(tokenResponse, "", "\t")
	fmt.Printf("Response:\n\n%s\n\n", string(out))
	c.JSON(200, tokenResponse)
}

func buildErrorMessage(responseStatus int) string {
	switch responseStatus {
	case 400:
		return "Invalid payload format: MUST be an array"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 406:
		return "Content type not accepted"
	case 413:
		return "Record limit exceeded"
	case 415:
		return "Unsupported Media Type"
	case 429:
		return "Unsupported Media Type"
	default:
		return "Internal Server Error"
	}
}
