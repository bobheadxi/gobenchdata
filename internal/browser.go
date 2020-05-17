package internal

import (
	"os/exec"
	"runtime"
)

// OpenBrowser opens a given URL in browser - sourced from https://gist.github.com/nanmu42/4fbaf26c771da58095fa7a9f14f23d27
func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		return
	}
	if err != nil {
		panic(err)
	}
}
