package metautils

import (
	"bufio"
	"fmt"
	"os"
)

func AskNewInput(name string) string {
	var data string = ""

	fmt.Print("[+] Type " + name + ": ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()

		if text == "" {
			fmt.Println("--Skipped--")
		}

		data = text
		break
	}

	return data
}
