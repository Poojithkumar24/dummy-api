package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Struct to parse the API response
type LogEntry struct {
	UserID       string `json:"user_id"`
	CreatedAt    string `json:"createdAt"`
	RequestBody  string `json:"requestBody"`
	ResponseBody string `json:"responseBody"`
	Service      string `json:"service"`
}

type APIResponse struct {
	Data struct {
		Logs []LogEntry `json:"logs"`
	} `json:"data"`
}

// Read user IDs from an Excel (.xlsx) file
func readUserIDsFromXLSX(filePath string) ([]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening Excel file: %v", err)
	}
	defer f.Close()

	var userIDs []string

	// Assuming the user IDs are in the first column (A) of the first sheet
	rows, err := f.GetRows(f.GetSheetName(0)) // Get the first sheet by index (0)
	if err != nil {
		return nil, fmt.Errorf("error reading rows: %v", err)
	}

	for _, row := range rows {
		if len(row) > 0 {
			userIDs = append(userIDs, strings.TrimSpace(row[0])) // Read the first column (user_id)
		}
	}
	return userIDs, nil
}

// Write logs to an Excel (.xlsx) file
func writeLogsToXLSX(filePath string, data [][]string) error {
	f := excelize.NewFile()

	// Write header row
	headers := []string{"user_id", "created_at", "request_body", "response_body", "service"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i) // A1, B1, C1, etc.
		f.SetCellValue("Sheet1", cell, header)
	}

	// Write log entries
	for i, row := range data {
		for j, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+j, i+2) // A2, B2, C2, etc.
			f.SetCellValue("Sheet1", cell, value)
		}
	}

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("error saving Excel file: %v", err)
	}
	return nil
}

// Fetch logs from the API for a given user ID
func fetchLogsForUser(userID, authToken string) ([]LogEntry, error) {
	url := fmt.Sprintf("https://gateway.finbox.in/lisa/getUserLogs?user_id=%s&collection_name=lender_request_response&service=abflplHunter", userID)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", authToken)

	client := &http.Client{Timeout: 30 * time.Second} // Increased timeout to 30 seconds
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Print the raw API response for debugging purposes
	fmt.Printf("API response for user %s: %s\n", userID, string(body))

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Print number of logs fetched
	fmt.Printf("Fetched %d logs for user %s\n", len(apiResponse.Data.Logs), userID)

	return apiResponse.Data.Logs, nil
}

// Main function to read user IDs, fetch logs, and write to Excel
func main() {
	// Input and output Excel paths
	inputXLSX := "/Users/poojithkumar/Downloads/4000.xlsx"
	outputXLSX := "/Users/poojithkumar/Downloads/abfl_output.xlsx"
	authToken := "Bearer ory_at_oCxRXn9jxFDnB1ZwpzvYQWeb-FQKp_pOqjtH6nLNbYg.U9WNtGpiOmENDuU9xO5bc8nO_EFxyX616xR11VHz5JM" // Provide authorization token here

	// Read user IDs
	userIDs, err := readUserIDsFromXLSX(inputXLSX)
	if err != nil {
		log.Fatalf("Error reading user IDs: %v", err)
	}
	fmt.Printf("Found %d user IDs to process\n", len(userIDs))

	var allLogs [][]string

	for _, userID := range userIDs {
		// Add 200ms delay before each API call
		time.Sleep(200 * time.Millisecond)

		// Fetch logs for the current user
		logs, err := fetchLogsForUser(userID, authToken)
		if err != nil {
			log.Printf("Error fetching logs for user %s: %v", userID, err)
			continue
		}

		// Append the fetched logs
		for _, logEntry := range logs {
			allLogs = append(allLogs, []string{
				userID,
				logEntry.CreatedAt,
				logEntry.RequestBody,
				logEntry.ResponseBody,
				logEntry.Service,
			})
		}
	}

	// Write logs to the output Excel file
	if err := writeLogsToXLSX(outputXLSX, allLogs); err != nil {
		log.Fatalf("Error writing logs to Excel: %v", err)
	}

	fmt.Printf("Successfully wrote logs to %s\n", outputXLSX)
}
