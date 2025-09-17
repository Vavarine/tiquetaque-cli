package app

import (
	"context"
	"fmt"
	"time"

	client "github.com/vavarine/ttq/internal/httpclient"
)

func DoPunch(ctx context.Context, cli *client.Client) (*client.PunchResponse, error) {
	// Loads the saved token, employeeID, and fullName from secure storage
	token, err := LoadToken()
	if err != nil || token == "" {
		return nil, fmt.Errorf("token não encontrado. Faça login primeiro")
	}

	employeeID, err := LoadEmployeeID()
	if err != nil || employeeID == "" {
		return nil, fmt.Errorf("employeeID não encontrado. Refazer login")
	}

	fullName, err := LoadEmployeeName()
	if err != nil || fullName == "" {
		return nil, fmt.Errorf("employeeName não encontrado. Refazer login")
	}

	now := time.Now()

	checkDate := now.Format("02-01-2006")
	jsonDate := now.Format("02/01/2006")
	punchTime := now.Format("15:04")
	xCheck := GetCheck(employeeID, fullName, checkDate, punchTime)

	resp, err := cli.Punch(ctx, token, xCheck, jsonDate, punchTime)

	if err != nil {
		return nil, fmt.Errorf("error when punching: %w", err)
	}

	return resp, nil
}
