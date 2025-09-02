package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "minha-cli",
	Short: "Uma CLI de exemplo com Cobra",
	Long:  `Essa é uma CLI simples feita em Go usando o framework Cobra.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bem-vindo à CLI do Evailson 🚀")
	},
}

// Função que o main.go chama
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
