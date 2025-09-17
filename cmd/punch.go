package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/app"
	client "github.com/vavarine/ttq/internal/httpclient"
)

var punchCmd = &cobra.Command{
	Use:   "punch",
	Short: "Bate o ponto usando o token salvo",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient("https://api.tiquetaque.com")

		resp, err := app.DoPunch(context.Background(), c)
		if err != nil {
			return err
		}

		fmt.Println("Ponto registrado com sucesso!")
		fmt.Printf("Resultado: %+v\n", resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(punchCmd)
}
