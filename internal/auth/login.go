package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LoginRequest representa o corpo da requisição de login
type LoginRequest struct {
	Email               string `json:"email"`
	SMSVerificationCode string `json:"sms_verification_code"`
	Source              string `json:"source"`
}

// LoginResponse representa a resposta da API de login
type LoginResponse struct {
	Token      string `json:"token"`
	EmployeeID string `json:"_id"`
}

// DoLogin envia o código de verificação e retorna o token e ID do funcionário
func DoLogin(ctx context.Context, email, code string) (*LoginResponse, error) {
	url := "https://api.tiquetaque.com/employees/code/verify"

	body := LoginRequest{
		Email:               email,
		SMSVerificationCode: code,
		Source:              "web",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar corpo: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Headers
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("origin", "https://tiquetaque.app")
	req.Header.Set("referer", "https://tiquetaque.app/")
	req.Header.Set("user-agent", "Go-http-client/1.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar requisição: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao logar: %s", string(respBody))
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &loginResp, nil
}
