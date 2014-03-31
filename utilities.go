package gadata

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

// BrowserOpen opens a URL is the OS' default web browser
// func BrowserOpen(url string) error {
// 	return exec.Command("open", url).Run()
// }

// WebCallback listens on a predefined port for a oauth response
// sends back via channel once it receives a response and shuts down.
func WebCallback(ch chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/?") {
			block := strings.SplitAfter(r.RequestURI, "/?")[1]
			if !strings.Contains(block, "code=") {
				ch <- block
			} else {
				ch <- strings.SplitAfter(block, "code=")[1]
				fmt.Fprintln(w, "Authentication completed, you can close this window.")
				close(ch)
				return
			}
		}
		fmt.Fprintln(w, "Error encountered during authentication.")
		return
	})

	log.Fatalf("Server exited: %v", http.ListenAndServe(ReturnURI, nil))
}
