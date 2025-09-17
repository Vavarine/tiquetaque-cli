package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetFullName(token string) (string, error) {
	req, err := c.newRequest("GET", "/employees", nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf("request failed: %s", string(data))
	}

	var result struct {
		Items []struct {
			FullName string `json:"full_name"`
		} `json:"_items"`
	}

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Resposta da API:", string(body))

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("error decoding response:", err)
		return "", err
	}

	return result.Items[0].FullName, nil
}
