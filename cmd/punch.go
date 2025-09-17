package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/app"
)

var (
	punchDate string
	punchTime string
)

var punchCmd = &cobra.Command{
	Use:   "punch",
	Short: "Bate o ponto usando o token salvo",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Carrega token do keyring/fallback
		token, err := app.LoadToken()
		if err != nil || token == "" {
			return fmt.Errorf("token não encontrado. Faça login primeiro com: ttq-cli login --email seu@email --code 123456")
		}

		employeeID, err := app.LoadEmployeeID()
		if err != nil || employeeID == "" {
			return fmt.Errorf("employeeID não encontrado. Refazer login pode corrigir: ttq-cli login --email ... --code ...")
		}

		// Se data ou hora não forem informadas, usa a data/hora atual
		now := time.Now()
		if punchDate == "" {
			punchDate = now.Format("02-01-2006")
		}
		if punchTime == "" {
			punchTime = now.Format("15:04")
		}

		fmt.Printf("Registrando ponto às %s do dia %s\n", punchTime, punchDate)

		// Cria contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		fmt.Println("Usando employeeID:", employeeID)
		fmt.Println("id length", len(employeeID))

		fullName, err := app.LoadEmployeeName()
		if err != nil || fullName == "" {
			return fmt.Errorf("nome do funcionário não encontrado. Refazer login pode corrigir: ttq-cli login --email ... --code ...")
		}

		xCheck := app.GetCheck(employeeID, fullName, punchDate, punchTime)

		resp, err := app.DoPunch(ctx, token, xCheck, now.Format("02/01/2006"), punchTime)
		if err != nil {
			return fmt.Errorf("erro ao registrar ponto: %w", err)
		}

		fmt.Println("Ponto registrado com sucesso!")
		fmt.Printf("Resultado: %+v\n", resp)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(punchCmd)
	punchCmd.Flags().StringVar(&punchDate, "date", "", "Data (dd/mm/yyyy)")
	punchCmd.Flags().StringVar(&punchTime, "time", "", "Hora (hh:mm)")
}
