package main

// Simple wkhtmltopdf as service
// TODO: tests
// TODO: add version

import (
	"log"
	"log/syslog"
	"net/http"
	"os/exec"
	"strings"
)

var wkCmd string

func main() {
	// enable syslog
	syslogger, err := syslog.New(syslog.LOG_INFO, "wkservice")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(syslogger)

	// check wkhtmltopdf command
	wkCmd = checkWkExists()

	// set endpoiunt
	http.HandleFunc("/", LaunchCmd)

	// launch service
	log.Println("Server listeng at :8080")
	http.ListenAndServe(":8080", nil)

}

func LaunchCmd(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := query.Get("params")
	// check pipe
	if strings.Contains(params, "|") {
		log.Println("Aborting command. Pipe detected.")
		return
	}

	log.Printf("Launching wkhtmltopdf with params: '%s'", params)

	cmd := wkCmd + " " + params
	res := exec.Command("bash", "-c", cmd)

	// configure `Stdout` and `Stderr`
	// res.Stdout = os.Stdout
	// res.Stderr = os.Stdout

	// run command
	if err := res.Run(); err != nil {
		log.Println("Error:", err)
	}

}

func checkWkExists() string {
	path, err := exec.LookPath("_wkhtmltopdf")
	if err != nil {
		log.Fatal("didn't find '_wkhtmltopdf' executable.\n")
	} else {
		log.Printf("'wkhtmltopdf' executable is in '%s'\n", path)
	}
	return path
}
