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
	CounterID int64
}

// buildHeaders adds authorization headers to the request.
func (c *Client) buildHeaders(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+c.Token)
}

// LogsList Returns a list of log requests.
func (c *Client) LogsList(ctx context.Context) ([]LogRequest, error) {
	URL := fmt.Sprintf(LogsListURL, c.CounterID)

	resp, err := c.doRequest(ctx, http.MethodGet, URL)
	if err != nil {
		return nil, err
	}

	defer closeBody(resp.Body)

	var res map[string][]LogRequest
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newInternalError(err, UnmarshalResponseFailedMsg)
	}

	return res["requests"], nil
}

// GetParts Returns information about the log request and parts of it.
func (c *Client) GetParts(ctx context.Context, reqID int) ([]Part, error) {
	URL := fmt.Sprintf(LogsStatusURL, c.CounterID, reqID)

	for {
		resp, err := c.doRequest(ctx, http.MethodGet, URL)
		if err != nil {
			return nil, err
		}

		var res MetrikaResponse

		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			closeBody(resp.Body)

			return nil, newInternalError(err, UnmarshalResponseFailedMsg)
		}

		closeBody(resp.Body)

		switch res.LogReq.Status {
		case string(Processed):
			return res.LogReq.Parts, nil
		case string(CleanedByUser), string(CleanedAutomaticallyAsTooOld):
			return nil, &APIError{
				Errors:  nil,
				Code:    http.StatusNotFound,
				Message: "отчет удалён пользователем или автоматически",
			}
		case string(Canceled):
			return nil, &APIError{
				Errors:  nil,
				Code:    http.StatusNotFound,
				Message: "отчет отменен",
			}
		case string(ProcessingFailed):
			return nil, &APIError{
				Errors:  nil,
				Code:    http.StatusInternalServerError,
				Message: "Ошибка при обработке отчета",
			}
		case string(AwaitingRetry), string(Created):
			time.Sleep(10 * time.Second)
		default:
			return nil, &InternalError{msg: fmt.Sprintf("неизвестный статус %s", res.LogReq.Status)}
		}
	}
}

// CollectAllParts Upload the log.
func (c *Client) CollectAllParts(ctx context.Context, reqID int, parts []Part, directory string) ([]string, error) {
	files := make([]string, 0, len(parts))

	for _, p := range parts {
		file, err := c.DownloadLogPart(ctx, reqID, p.PartNumber, directory)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

// DownloadLogPart Uploads a part of the prepared log.
func (c *Client) DownloadLogPart(ctx context.Context, reqID, partNumber int, directory string) (string, error) {
	URL := fmt.Sprintf(DownloadLogURL, c.CounterID, reqID, partNumber)

	resp, err := c.doRequest(ctx, http.MethodGet, URL)
	if err != nil {
		return "", err
	}

	defer closeBody(resp.Body)

	filename := fmt.Sprintf("%v_%v_%v-*.csv", c.CounterID, reqID, partNumber)

	f, err := os.CreateTemp(directory, filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

// DeleteLog Delete the log.
func (c *Client) DeleteLog(ctx context.Context, counterID, reqID int) (bool, error) {
	URL := fmt.Sprintf(DeleteLogURL, counterID, reqID)

	resp, err := c.doRequest(ctx, http.MethodPost, URL)
	if err != nil {
		return false, err
	}

	defer closeBody(resp.Body)

	var res MetrikaResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, newInternalError(err, UnmarshalResponseFailedMsg)
	}

	return res.LogReq.Status == string(CleanedByUser), nil
}

// GetCounters Get all counters.
func (c *Client) GetCounters(ctx context.Context) ([]Counter, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, CountersURL)
	if err != nil {
		return nil, err
	}

	defer closeBody(resp.Body)

	var res CounterResponse

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, newInternalError(err, UnmarshalResponseFailedMsg)
	}

	return res.Counters, nil
}

// CreateLog This function returns request_id for the log.
func (c *Client) CreateLog(ctx context.Context, dateFrom, dateTo, fields, source string) (int, error) {
	d := url.Values{
		"date1":  []string{dateFrom},
		"date2":  []string{dateTo},
		"fields": []string{fields},
		"source": []string{source},
	}

	URL := fmt.Sprintf(CreateLogURL, c.CounterID) + "?" + d.Encode()

	resp, err := c.doRequest(ctx, http.MethodPost, URL)
	if err != nil {
		return 0, err
	}

	defer closeBody(resp.Body)

	var res MetrikaResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, newInternalError(err, UnmarshalResponseFailedMsg)
	}

	return res.LogReq.RequestID, nil
}

func statusCodeHandler(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	defer closeBody(resp.Body)

	defaultErr := &APIError{
		Code:    resp.StatusCode,
		Message: http.StatusText(resp.StatusCode),
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %s\n", err)

		return defaultErr
	}

	var apiError APIError

	if err := json.Unmarshal(body, &apiError); err != nil {
		fmt.Printf("Failed to unmarshal API error: %s\n", err)

		return defaultErr
	}

	return &apiError
}

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		fmt.Printf("Failed to close response body: %s\n", err)
	}
}

func (c *Client) doRequest(ctx context.Context, method, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, newInternalError(err, CreateRequestFailedMsg)
	}

	c.buildHeaders(req)

	resp, err := c.Tr.Do(req)
	if err != nil {
		return nil, newInternalError(err, RequestFailedMsg)
	}

	if err := statusCodeHandler(resp); err != nil {
		return nil, err
	}

	return resp, nil
}
