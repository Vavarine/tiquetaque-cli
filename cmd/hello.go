package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var nome string

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Diz olá para alguém",
	Run: func(cmd *cobra.Command, args []string) {
		if nome == "" {
			fmt.Println("Olá, mundo!")
		} else {
			fmt.Printf("Olá, %s!\n", nome)
		}
	},
}

func init() {
	// adiciona o comando hello ao root
	rootCmd.AddCommand(helloCmd)

	// adiciona uma flag --nome
	helloCmd.Flags().StringVarP(&nome, "nome", "n", "", "Nome da pessoa")
}
