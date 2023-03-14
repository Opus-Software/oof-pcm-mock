package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OPINRequestPayload struct {
	OrganizationId string   `json:"organizationId"`
	Events         []Report `json:"events"`
}

type Report struct {
	ClientOrgId       string         `json:"clientOrgId"`
	ClientSSId        string         `json:"clientSSId,omitempty"`
	CorrelationId     string         `json:"correlationId"`
	Endpoint          string         `json:"endpoint"`
	Url               string         `json:"url"`
	FapiInteractionId string         `json:"fapiInteractionId"`
	HttpMethod        string         `json:"httpMethod"`
	ServerOrgId       string         `json:"serverOrgId"`
	ServerASId        string         `json:"serverASId,omitempty"`
	Timestamp         string         `json:"timestamp"`
	EndpointUriPrefix string         `json:"endpointUriPrefix,omitempty"`
	ProcessTimespam   int            `json:"processTimespan"`
	StatusCode        int            `json:"statusCode"`
	AdditionalInfo    AdditionalInfo `json:"additionalInfo,omitempty"`
}

type AdditionalInfo struct {
	ConsentId       string `json:"consentId,omitempty"`
	PersonType      string `json:"personType,omitempty"`
	LocalInstrument string `json:"localInstrument,omitempty"`
	Status          string `json:"status,omitempty"`
	RejectReason    string `json:"rejectReason,omitempty"`
	ClientIp        string `json:"clientIp,omitempty"`
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
	r.POST("/report-api/v1/private/report", handleOFBPostPrivateReport)
	r.POST("/report-api/v1/opendata/report", handleOFBPostOpenDataReport)
	r.POST("/token/", handleOFBPostToken)
	r.POST("/report-api/v1/server-batch", handleOPINPostServerBatch)
	r.Run(address)
}

func handleOFBPostPrivateReport(c *gin.Context) {
	fmt.Println("Handling new OFB post private report")

	var requestReports []Report
	var responseReports []ReportResponseSuccess

	expectedResponseStatus := "200"
	possibleStatuses := [5]string{"ACCEPTED", "DISCARDED", "PAIRED_INCONSISTENT", "UNPAIRED", "SINGLE"}

	if len(os.Args) > 1 {
		expectedResponseStatus = os.Args[1]
	}

	fmt.Printf("Request headers received:\n")
	for k, vals := range c.Request.Header {
		fmt.Printf("\n%s", k)
		for _, v := range vals {
			fmt.Printf("\t%s ", v)
		}
	}
	fmt.Printf("\n\n")

	_ = c.Bind(&requestReports)
	responseStatus, _ := strconv.Atoi(expectedResponseStatus)

	if responseStatus == 200 || responseStatus == 207 {
		fmt.Printf("New post with %d incoming events\n", len(requestReports))

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

func handleOFBPostOpenDataReport(c *gin.Context) {
	fmt.Println("Handling new OFB post open data report")

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

func handleOFBPostToken(c *gin.Context) {
	tokenResponse := &TokenResponseSuccess{
		AccessToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkFiT",
		ExpiresIn:   21600,
		TokenType:   "Bearer",
	}
	out, _ := json.MarshalIndent(tokenResponse, "", "\t")
	fmt.Printf("Response:\n\n%s\n\n", string(out))
	c.JSON(200, tokenResponse)
}

func handleOPINPostServerBatch(c *gin.Context) {
	fmt.Println("Handling new OPIN post private report")

	var opinRequestPayload OPINRequestPayload
	var responseReports []ReportResponseSuccess

	expectedResponseStatus := "200"
	possibleStatuses := [5]string{"ACCEPTED", "DISCARDED", "PAIRED_INCONSISTENT", "UNPAIRED", "SINGLE"}

	if len(os.Args) > 1 {
		expectedResponseStatus = os.Args[1]
	}

	fmt.Printf("Request headers received:\n")
	for k, vals := range c.Request.Header {
		fmt.Printf("\n%s", k)
		for _, v := range vals {
			fmt.Printf("\t%s ", v)
		}
	}
	fmt.Printf("\n\n")

	jwtBody, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Received JWT payload: %s\n\n", string(jwtBody))

	jwtParts := strings.Split(string(jwtBody), ".")
	decodedBody, _ := base64.RawURLEncoding.DecodeString(jwtParts[1])

	json.Unmarshal(decodedBody, &opinRequestPayload)
	requestPayload, _ := json.MarshalIndent(opinRequestPayload, "", "\t")
	fmt.Printf("Decoded payload: %s\n", requestPayload)

	responseStatus, _ := strconv.Atoi(expectedResponseStatus)

	if responseStatus == 200 || responseStatus == 207 {
		requestReports := opinRequestPayload.Events
		fmt.Printf("New post with %d incoming events\n", len(requestReports))

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
