package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/auth"
)

var (
	punchDate  string
	punchTime  string
	employerID string
)

var punchCmd = &cobra.Command{
	Use:   "punch",
	Short: "Bate o ponto usando o token salvo",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Carrega token do keyring/fallback
		token, err := auth.LoadToken()
		if err != nil {
			return fmt.Errorf("token não encontrado: %w", err)
		}

		// Define employerID
		empID := employerID
		if empID == "" {
			empID, err = auth.LoadEmployeeID()
			if err != nil {
				return fmt.Errorf("employer id não encontrado: %w", err)
			}
		}

		if empID == "" || punchDate == "" || punchTime == "" {
			return fmt.Errorf("use --employerid, --date e --time")
		}

		// Cria contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Chama DoPunch refatorado
		resp, err := auth.DoPunch(ctx, token, empID, punchDate, punchTime)
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
	punchCmd.Flags().StringVar(&employerID, "employerid", "", "ID do empregador")
	punchCmd.Flags().StringVar(&punchDate, "date", "", "Data (dd/mm/yyyy)")
	punchCmd.Flags().StringVar(&punchTime, "time", "", "Hora (hh:mm)")
}
