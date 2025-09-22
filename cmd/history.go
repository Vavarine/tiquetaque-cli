package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vavarine/ttq/internal/app"
	client "github.com/vavarine/ttq/internal/httpclient"
)

var historyCmd = &cobra.Command{
	Use:     "history",
	Short:   "Exibe o histórico de batidas de pontos",
	Aliases: []string{"h"},
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.NewClient("https://api.tiquetaque.com")

		resp, err := app.GetPunchHistory(context.Background(), c)
		if err != nil {
			return err
		}

		// Print the punch history
		for _, dayRecord := range resp.Items {
			if len(dayRecord.TimeEntries) == 0 {
				continue
			}
			fmt.Printf("Data: %s\n", strings.Split(dayRecord.Date, " ")[0])
			for _, punch := range dayRecord.TimeEntries {
				fmt.Printf("  - Horário: %s\n", punch.Time)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
