package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/auth"
)

var (
	email string
	code  string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Faz login via código e email",
	RunE: func(cmd *cobra.Command, args []string) error {
		if email == "" || code == "" {
			return fmt.Errorf("use --email e --code para informar email e código")
		}

		// Cria contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Chama a função refatorada DoLogin
		resp, err := auth.DoLogin(ctx, email, code)
		if err != nil {
			return fmt.Errorf("falha no login: %w", err)
		}

		fmt.Println("Login feito com sucesso!")

		// Salva token no keyring ou fallback
		if resp.Token != "" {
			if err := auth.SaveToken(resp.Token); err != nil {
				return fmt.Errorf("erro ao salvar token: %w", err)
			}
			fmt.Println("Token salvo com sucesso 🔑")
		}

		// Salva EmployeeID
		if resp.EmployeeID != "" {
			if err := auth.SaveEmployeeID(resp.EmployeeID); err != nil {
				return fmt.Errorf("erro ao salvar EmployeeID: %w", err)
			}
			fmt.Println("ID do funcionário salvo com sucesso 🆔")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Email do usuário")
	loginCmd.Flags().StringVarP(&code, "code", "c", "", "Código de verificação")
}
