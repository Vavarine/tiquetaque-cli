package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/app"
	client "github.com/vavarine/ttq/internal/httpclient"
)

var (
	email string
	code  string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Faz login no sistema",
	RunE: func(cmd *cobra.Command, args []string) error {
		email, _ := cmd.Flags().GetString("email")
		code, _ := cmd.Flags().GetString("code")

		if email == "" || code == "" {
			return fmt.Errorf("use --email e --code")
		}

		c := client.NewClient("https://api.tiquetaque.com")

		if err := app.DoLogin(context.Background(), c, email, code); err != nil {
			return err
		}

		fmt.Println("Login realizado com sucesso!")
		return nil
	},
}

func init() {
	loginCmd.Flags().String("email", "", "Email do usuário")
	loginCmd.Flags().String("code", "", "Código de verificação")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("code")

	rootCmd.AddCommand(loginCmd)
}
