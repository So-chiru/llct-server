package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/joho/godotenv"

	"github.com/so-chiru/llct-server/dashboard"
	"github.com/so-chiru/llct-server/route"

	"net/http"
)

var GitCommit = "unknown"

func version() {
	address := os.Getenv("ADDRESS")

	fmt.Println("                                    ")
	fmt.Println("                      ****          ")
	fmt.Println("                     ******%        ")
	fmt.Println("                    ********        llct-server #" + GitCommit)
	fmt.Println("                   ********/        ")
	fmt.Println("                  /******/%         built with " + runtime.Compiler + "(" + runtime.Version() + ")")
	fmt.Println("                 &****/////         running on " + runtime.GOARCH)
	fmt.Println("                /**/*/////          listening on port: " + address + ", PID: " + fmt.Sprint(os.Getpid()))
	fmt.Println("               **///////            ")
	fmt.Println("             **///////#             ")
	fmt.Println("           *///////                 https://github.com/So-chiru/llct-server")
	fmt.Println("        *////                       ")
	fmt.Println("      */                            ")
	fmt.Println("                                    ")

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("- Error occurred while loading the .env file.")
	}

	version()

	if len(os.Args) > 1 {
		var command = os.Args[1]

		if command == "--version" || command == "version" || command == "-v" {
			return
		}
	}

	log.Println("# server is now running.")

	address := os.Getenv("ADDRESS")

	dashboard.SyncBirthdayCalendar()
	dashboard.SyncLiveCalendar()

	http.ListenAndServe(address, route.Router())
}
