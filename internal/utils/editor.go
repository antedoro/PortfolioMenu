package utils

import (
	"os/exec"
	"runtime"
)

func OpenEditor(
	file string,
) error {

	switch runtime.GOOS {

	case "darwin":

		return exec.Command(
			"open",
			file,
		).Start()

	case "linux":

		return exec.Command(
			"xdg-open",
			file,
		).Start()

	case "windows":

		return exec.Command(
			"notepad",
			file,
		).Start()

	default:

		return nil

	}

}
