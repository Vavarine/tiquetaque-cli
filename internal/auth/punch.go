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

type PunchRequest struct {
	Times []struct {
		Source string `json:"source"`
		Date   string `json:"date"`
		Time   string `json:"time"`
	} `json:"times"`
	EmployerID string `json:"employerId"`
}

type PunchResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DoPunch envia um registro de ponto para a API do TiqueTaQue
func DoPunch(ctx context.Context, token, employerID, date, timeStr string) (*PunchResponse, error) {
	url := "https://api.tiquetaque.com/employees/day-records/add-times"

	// Monta o corpo da requisição
	body := PunchRequest{
		Times: []struct {
			Source string `json:"source"`
			Date   string `json:"date"`
			Time   string `json:"time"`
		}{
			{Source: "web", Date: date, Time: timeStr},
		},
		EmployerID: employerID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar corpo: %w", err)
	}

	// Cria cliente HTTP com timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Cria requisição com contexto
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Define headers
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("origin", "https://tiquetaque.app")
	req.Header.Set("referer", "https://tiquetaque.app/")
	req.Header.Set("user-agent", "Go-http-client/1.1")

	// Executa requisição
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
		return nil, fmt.Errorf("erro ao registrar ponto: %s", string(respBody))
	}

	var punchResp PunchResponse
	if err := json.Unmarshal(respBody, &punchResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &punchResp, nil
}
