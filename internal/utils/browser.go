package utils

import (
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {

	switch runtime.GOOS {

	case "darwin":

		exec.Command(
			"open",
			url,
		).Start()

	case "linux":

		exec.Command(
			"xdg-open",
			url,
		).Start()

	case "windows":

		exec.Command(
			"rundll32",
			"url.dll,FileProtocolHandler",
			url,
		).Start()

	}

}
