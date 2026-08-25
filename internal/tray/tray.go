package tray

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/getlantern/systray"

	"github.com/antedoro/PortfolioMenu/internal/models"
	"github.com/antedoro/PortfolioMenu/internal/services"
	"github.com/antedoro/PortfolioMenu/internal/utils"
)

type Tray struct {
	svc   *services.Service
	index int
}

func New(
	svc *services.Service,
) *Tray {

	return &Tray{svc: svc}
}

func (t *Tray) Run() {

	systray.Run(
		t.onReady,
		t.onExit,
	)
}

func (t *Tray) onReady() {

	systray.SetTitle("PortfolioMenu")
	systray.SetTooltip("PortfolioMenu")

	openDashboard :=
		systray.AddMenuItem(
			"Open Dashboard",
			"Apri dashboard web",
		)

	refresh :=
		systray.AddMenuItem(
			"Refresh",
			"Aggiorna prezzi",
		)

	editConfig :=
		systray.AddMenuItem(
			"Edit config",
			"Modifica configurazione",
		)

	systray.AddSeparator()

	about :=
		systray.AddMenuItem(
			"About",
			"Informazioni",
		)

	systray.AddSeparator()

	quit :=
		systray.AddMenuItem(
			"Quit",
			"Chiudi",
		)

	go t.updateTitle()

	go func() {

		for {

			select {

			case <-openDashboard.ClickedCh:
				utils.OpenBrowser(
					"http://localhost:8080",
				)

			case <-refresh.ClickedCh:
				go t.svc.Refresh()

			case <-editConfig.ClickedCh:
				utils.OpenBrowser(
					"http://localhost:8080/config",
				)

			case <-about.ClickedCh:
				utils.OpenBrowser(
					"http://localhost:8080/about",
				)

			case <-quit.ClickedCh:
				t.svc.Stop()
				systray.Quit()
				return

			}

		}

	}()

}

func (t *Tray) updateTitle() {

	update := func() {
		p := t.svc.Get()
		if len(p.Assets) == 0 {
			return
		}

		if t.index >= len(p.Assets) {
			t.index = 0
		}

		a := p.Assets[t.index]
		systray.SetTitle(formatTitle(a, t.svc.MenubarFormat()))
	}

	update()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	ticks := 0
	for range ticker.C {
		ticks++
		if ticks % 10 == 0 {
			t.index++
		}
		update()
	}

}

func formatTitle(
	a models.Asset,
	format []string,
) string {

	name := a.Ticker
	if name == "" {
		name = a.Name
	}

	var result strings.Builder

	for i, f := range format {
		
		if i > 0 {
			if f == "percent" && format[i-1] == "value" {
				// No space between value and percent
			} else {
				result.WriteString(" ")
			}
		}

		switch f {
		case "ticker":
			result.WriteString(name + ":")
		case "value":
			result.WriteString(fmt.Sprintf("%d€", int(a.MarketValue)))
		case "percent":
			percentStr := fmt.Sprintf("%+.2f%%", a.GainPercent)
			percentStr = strings.Replace(percentStr, ".", ",", 1)
			result.WriteString("(" + percentStr + ")")
		case "gainloss":
			result.WriteString(fmt.Sprintf("%+d€", int(a.GainLoss)))
		}
	}

	return result.String()
}

func openURL(url string) {
	if runtime.GOOS == "darwin" {
		exec.Command("open", url).Start()
	}
}

func (t *Tray) onExit() {
}
