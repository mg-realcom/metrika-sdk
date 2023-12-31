package metrika_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
func (c *Client) LogsList(ctx context.Context) ([]LogRequest, error) {
	createdLogs := []LogRequest{}
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(LogsListUrl, c.CounterId), nil)
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
func (c *Client) GetParts(ctx context.Context, reqId int) ([]Part, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(LogsStatusUrl, c.CounterId, reqId), nil)
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
func (c *Client) CollectAllParts(ctx context.Context, reqId int, parts []Part, directory string) ([]string, error) {
	files := []string{}
	for _, p := range parts {
		file, err := c.DownloadLogPart(ctx, reqId, p.PartNumber, directory)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// Uploads a part of the prepared log.
func (c *Client) DownloadLogPart(ctx context.Context, reqId, partNumber int, directory string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(DownloadLogUrl, c.CounterId, reqId, partNumber), nil)
	if err != nil {
		return "", err
	}
	c.buildHeaders(req)
	resp, err := c.Tr.Do(req)
	if err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%v_%v_%v-*.csv", c.CounterId, reqId, partNumber)
	f, err := os.CreateTemp(directory, filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return f.Name(), nil

}

// Delete the log.
func (c *Client) DeleteLog(ctx context.Context, counterId, reqId int) (bool, error) {
	res := MetrikaResponse{}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(DeleteLogUrl, counterId, reqId), nil)
	if err != nil {
		return false, err
	}
	c.buildHeaders(req)
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
func (c *Client) GetCounters(ctx context.Context) ([]Counter, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", CountersUrl, nil)
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
func (c *Client) CreateLog(ctx context.Context, dateFrom, dateTo, fields, source string) (int, error) {
	d := url.Values{
		"date1":  []string{dateFrom},
		"date2":  []string{dateTo},
		"fields": []string{fields},
		"source": []string{source},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(CreateLogUrl, c.CounterId), nil)
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
