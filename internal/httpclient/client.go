package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// Construtor
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

// cria request com headers padrões
func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return req, nil
}

// Login chama a API e retorna token e employeeID
func (c *Client) Login(email, code string) (string, string, error) {
	payload := map[string]string{
		"email":                 email,
		"sms_verification_code": code,
		"source":                "web",
	}
	body, _ := json.Marshal(payload)

	req, err := c.newRequest("POST", "/employees/code/verify", bytes.NewBuffer(body))

	if err != nil {
		return "", "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("login failed: %s", string(data))
	}

	var result struct {
		Token      string `json:"token"`
		EmployeeID string `json:"_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	c.token = result.Token
	return result.Token, result.EmployeeID, nil
}

type EmployeesResponse struct {
	Items []struct {
		FullName string `json:"full_name"`
	} `json:"_items"`
}

func (c *Client) GetFullName() (string, error) {
	req, err := c.newRequest("GET", "/employees", nil)

	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf("request failed: %s", string(data))
	}

	var employeesRes EmployeesResponse

	if err := json.NewDecoder(resp.Body).Decode(&employeesRes); err != nil {
		fmt.Println("error decoding response:", err)
		return "", err
	}

	return employeesRes.Items[0].FullName, nil
}

type PunchRequest struct {
	Times []struct {
		Source string `json:"source"`
		Date   string `json:"date"`
		Time   string `json:"time"`
	} `json:"times"`
}

type PunchResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (c *Client) Punch(ctx context.Context, xCheck, date, timeStr string) (*PunchResponse, error) {
	body := PunchRequest{
		Times: []struct {
			Source string `json:"source"`
			Date   string `json:"date"`
			Time   string `json:"time"`
		}{
			{Source: "web", Date: date, Time: timeStr},
		},
	}

	jsonBody, _ := json.Marshal(body)
	req, err := c.newRequest("POST", "/employees/day-records/add-times", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// adiciona header X-Check
	req.Header.Set("X-Check", xCheck)

	// timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao registrar ponto: %s", string(respBody))
	}

	var punchResp PunchResponse
	if err := json.Unmarshal(respBody, &punchResp); err != nil {
		return nil, err
	}

	return &punchResp, nil
}
