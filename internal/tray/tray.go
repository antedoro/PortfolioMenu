package tray

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/getlantern/systray"

	"github.com/antedoro/PortfolioMenu/internal/portfolio"
)

type Tray struct {
	Updater *portfolio.Updater

	index int
}

func New(
	updater *portfolio.Updater,
) *Tray {

	return &Tray{
		Updater: updater,
	}

}

func (t *Tray) Run() {

	systray.Run(
		t.onReady,
		t.onExit,
	)

}

func (t *Tray) onReady() {

	systray.SetTitle(
		"PortfolioMenu",
	)

	systray.SetTooltip(
		"PortfolioMenu",
	)

	openDashboard :=
		systray.AddMenuItem(
			"Open Dashboard",
			"Apri dashboard web",
		)

	editConfig :=
		systray.AddMenuItem(
			"Edit config",
			"Modifica configurazione",
		)

	systray.AddSeparator()

	checkUpdate :=
		systray.AddMenuItem(
			"Check for Update...",
			"Aggiornamenti",
		)

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

				openURL(
					"http://localhost:8080",
				)

			case <-editConfig.ClickedCh:

				openFile(
					"configs/portfoliomenu.toml",
				)

			case <-checkUpdate.ClickedCh:

				fmt.Println(
					"Check update",
				)

			case <-about.ClickedCh:

				fmt.Println(
					"PortfolioMenu",
				)

			case <-quit.ClickedCh:

				systray.Quit()

				return

			}

		}

	}()

}

func (t *Tray) updateTitle() {

	ticker :=
		time.NewTicker(
			10 * time.Second,
		)

	defer ticker.Stop()

	for range ticker.C {

		p :=
			t.Updater.Get()

		if len(p.Assets) == 0 {

			continue

		}

		if t.index >= len(p.Assets) {

			t.index = 0

		}

		a :=
			p.Assets[t.index]

		text :=
			fmt.Sprintf(
				"%s: %.2f %s",
				a.Name,
				a.LastPrice,
				a.CurrencySymbol,
			)

		systray.SetTitle(
			text,
		)

		t.index++

	}

}

func openURL(url string) {

	if runtime.GOOS == "darwin" {

		exec.Command(
			"open",
			url,
		).Start()

	}

}

func openFile(file string) {

	if runtime.GOOS == "darwin" {

		exec.Command(
			"open",
			file,
		).Start()

	}

}

func (t *Tray) onExit() {

}
