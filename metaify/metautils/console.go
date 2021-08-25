package metautils

import (
	"fmt"
	"log"
)

func AskNewInput(name string) string {
	var data string = ""

	fmt.Print("[+] Type " + name + ": ")
	_, err := fmt.Scanln(&data)
	if err != nil {
		if err.Error() == "unexpected newline" {
			fmt.Println("--Skipped--")
		} else {
			log.Fatalln(err)
		}
	}

	return data
}
