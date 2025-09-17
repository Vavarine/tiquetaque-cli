package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

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

func saveErrorToFile(data []byte) error {
	file, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("erro ao escrever no arquivo de log: %w", err)
	}

	return nil
}

// DoPunch envia um registro de ponto para a API do TiqueTaQue
func DoPunch(ctx context.Context, token, xCheck, date, timeStr string) (*PunchResponse, error) {
	url := "https://api.tiquetaque.com/employees/day-records/add-times"

	body := PunchRequest{
		Times: []struct {
			Source string `json:"source"`
			Date   string `json:"date"`
			Time   string `json:"time"`
		}{
			{Source: "web", Date: date, Time: timeStr},
		},
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

	// Headers based on your curl
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "pt-BR,pt;q=0.9")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("origin", "https://tiquetaque.app")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tiquetaque.app/")
	req.Header.Set("sec-ch-ua", `"Not;A=Brand";v="99", "Brave";v="139", "Chromium";v="139"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	req.Header.Set("X-Check", xCheck)

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
		if err := saveErrorToFile(respBody); err != nil {
			return nil, fmt.Errorf("erro ao salvar erro em arquivo: %w", err)
		}
		return nil, fmt.Errorf("erro ao registrar ponto: %s", string(respBody))
	}

	var punchResp PunchResponse
	if err := json.Unmarshal(respBody, &punchResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &punchResp, nil
}
