package metrika_sdk

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	Tr        *http.Client
	Token     string
	CounterId int64
}

func (c *Client) buildHeaders(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+c.Token)
}

// Returns a list of log requests.
func (c *Client) LogsList() ([]LogRequest, error) {
	createdLogs := []LogRequest{}
	req, err := http.NewRequest("GET", fmt.Sprintf(LogsListUrl, c.CounterId), nil)
	if err != nil {
		return nil, err
	}
	c.buildHeaders(req)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	res := make(map[string][]LogRequest)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	for _, v := range res["requests"] {
		createdLogs = append(createdLogs, v)
	}
	return createdLogs, nil
}

// Returns information about the log request and parts of it.
func (c *Client) GetParts(reqId int) ([]Part, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(LogsStatusUrl, c.CounterId, reqId), nil)
	if err != nil {
		panic(err)
	}
	c.buildHeaders(req)
	var resp *http.Response
	var res MetrikaResponse
	for {
		resp, err = c.Tr.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		if res.LogReq.Status == "processed" {
			fmt.Println("Отчет готов")
			return res.LogReq.Parts, nil
		} else if res.LogReq.Status == "cleaned_by_user" {
			return nil, fmt.Errorf("Такого отчета не существует")
		}
		time.Sleep(10 * time.Second)
	}
}

// Upload the log.
func (c *Client) CollectAllParts(reqId int, parts []Part) ([]string, error) {
	files := []string{}
	for _, p := range parts {
		file, err := c.DownloadLogPart(reqId, p.PartNumber)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// Uploads a part of the prepared log.
func (c *Client) DownloadLogPart(reqId int, partNumber int) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(DownloadLogUrl, c.CounterId, reqId, partNumber), nil)
	if err != nil {
		return "", err
	}
	c.buildHeaders(req)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	filename := fmt.Sprintf("%v_%v_%v-*.csv", c.CounterId, reqId, partNumber)
	file, err := os.CreateTemp("/Users/kostya/Gowork/metrika-sdk", filename)
	if err != nil {
		return "", err
	}
	writer := csv.NewWriter(file)
	writer.Comma = '|'
	defer file.Close()
	reader := bufio.NewReader(resp.Body)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\"", "")
		//line = strings.ReplaceAll(line, "--", "")
		// Split the line by tabs
		fields := strings.Split(line, "\t")
		// Write the fields to the CSV file
		err = writer.Write(fields)
		if err != nil {
			return "", fmt.Errorf("failed to write to CSV file: %v", err)
		}

	}
	//body, _ := io.ReadAll(resp.Body)
	//err = os.WriteFile(file.Name(), body, 0644)
	//if err != nil {
	//	return "", err
	//}
	return filename, nil

}

// Delete the log.
func (c *Client) DeleteLog(counterId int, reqId int) (bool, error) {
	res := MetrikaResponse{}
	req, err := http.NewRequest("POST", fmt.Sprintf(DeleteLogUrl, counterId, reqId), nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, err
	}
	return res.LogReq.Status == "cleaned_by_user", nil

}

// Get all counters.
func (c *Client) GetCounters() ([]Counter, error) {
	req, err := http.NewRequest("GET", CountersUrl, nil)
	if err != nil {
		panic(err)
	}
	c.buildHeaders(req)
	resp, err := c.Tr.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var res CounterResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res.Counters, err
	}
	return res.Counters, nil
}

// This function returns reques_id for the log.
func (c *Client) CreateLog(dateFrom string, dateTo string, fields string, source string) (int, error) {
	d := url.Values{
		"date1":  []string{dateFrom},
		"date2":  []string{dateTo},
		"fields": []string{fields},
		"source": []string{source},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf(CreateLogUrl, c.CounterId), nil)
	req.URL.RawQuery = d.Encode()

	if err != nil {
		panic(err)
	}
	c.buildHeaders(req)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var res MetrikaResponse

	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return 0, err
	}
	return res.LogReq.RequestID, nil
}
