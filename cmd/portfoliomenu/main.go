package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/antedoro/PortfolioMenu/internal/assets"
	"github.com/antedoro/PortfolioMenu/internal/config"
	"github.com/antedoro/PortfolioMenu/internal/server"
	"github.com/antedoro/PortfolioMenu/internal/services"
	"github.com/antedoro/PortfolioMenu/internal/tray"
)

// baseDir restituisce la directory di riferimento:
// se l'eseguibile è dentro un .app bundle usa la sua cartella
// MacOS (la working directory sarebbe "/"), altrimenti la cwd.
func baseDir() string {

	if exe, err := os.Executable(); err == nil {

		if strings.Contains(
			exe,
			".app/Contents/MacOS",
		) {
			return filepath.Dir(exe)
		}

	}

	return "."
}

// locate trova un file esistente (config).
func locate(rel string) string {

	candidates := []string{
		filepath.Join(baseDir(), rel),
		rel,
	}

	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}

	return candidates[0]
}

// locateNew risolve il percorso di un file da creare, ancorandolo
// alla base se la sua cartella esiste già.
func locateNew(rel string) string {

	bd := filepath.Join(baseDir(), rel)

	if st, err := os.Stat(filepath.Dir(bd)); err == nil &&
		st.IsDir() {
		return bd
	}

	return rel
}

func main() {

	cfgPath := locate(
		"configs/portfoliomenu.toml",
	)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	cfg.DataFile = locateNew(cfg.DataFile)

	fmt.Println()
	fmt.Println("PortfolioMenu")
	fmt.Println("========================================")
	fmt.Printf(
		"Refresh ogni %d minuti\n",
		cfg.RefreshMinutes,
	)
	fmt.Printf("File dati: %s\n", cfg.DataFile)

	svc := services.New(cfg)
	if err := svc.Load(); err != nil {
		log.Fatal(err)
	}

	svc.Start(cfg.RefreshMinutes)

	fmt.Println("Portfolio updater avviato")

	webServer := server.New(
		svc,
		assets.Templates,
	)
	webServer.Start("localhost:8080")

	fmt.Println(
		"Dashboard disponibile su http://localhost:8080",
	)

	appTray := tray.New(svc)
	appTray.Run()
}
