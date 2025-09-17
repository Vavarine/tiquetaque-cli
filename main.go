package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vavarine/ttq/cmd"
)

var version string = "dev" // set via ldflags in build

type Release struct {
	TagName string `json:"tag_name"`
}

func checkForUpdateAsync() {
	go func() {
		client := http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("https://api.github.com/repos/Vavarine/tiquetaque-cli/releases/latest")
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var r Release
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			return
		}

		if r.TagName != version {
			fmt.Fprintf(os.Stderr, "\n⚠️ Nova versão da cli disponível: %s (atual: %s)\n", r.TagName, version)
			fmt.Fprintf(os.Stderr, "Atualize com: curl -fsSL https://raw.githubusercontent.com/Vavarine/tiquetaque-cli/main/install.sh | bash \n\n")
		}
	}()
}

func main() {
	fmt.Printf("running ttq version %s\n\n", version)
	checkForUpdateAsync()

	cmd.Execute()
}
