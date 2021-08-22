package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const times = 5
const delay = 5

func main() {

	for {
		showMenu()
		command := getCommand()

		switch command {
		case 1:
			sites := getSitesByFile()
			execMonitor(sites)

		case 2:
			fmt.Println("Printing Logs..,")
			printLogs()

		case 0:
			os.Exit(0)

		default:
			fmt.Println("Invalid option!")
			fmt.Println("")
			main()
		}
	}
}

func showMenu() {
	fmt.Println("MENU:")
	fmt.Println("1 - Start monitor:")
	fmt.Println("2 - Show logs:")
	fmt.Println("0 - Quit")
}

func getCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func execMonitor(sites []string) {
	fmt.Println("Monitor started!")

	for i := 0; i < times; i++ {
		for _, site := range sites {
			testSite(site)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func testSite(site string) {
	response, err := http.Get(site)

	if err != nil {
		fmt.Println("erro! message:", err)
	} else {
		if response.StatusCode == 200 {
			fmt.Println(site, " - OK")
			fmt.Printf("")
			recordLog(site, true)
		} else {
			fmt.Println(site, " - OFF - Status code: ", response.StatusCode)
			fmt.Printf("")
			recordLog(site, false)
		}
	}
}

func getSitesByFile() []string {
	file, err := os.Open("sites.txt")
	reader := bufio.NewReader(file)

	if err != nil {
		fmt.Println("erro! ", err)
	}

	var sites []string

	for {
		site, err := reader.ReadString('\n')
		site = strings.TrimSpace(site)

		sites = append(sites, site)

		if err == io.EOF {
			break
		}
	}
	file.Close()
	return sites
}

func recordLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		file.WriteString(time.Now().Format("02/01/2006 15:04:05 - ") + site + " online: " + strconv.FormatBool(status) + "\n")
	}

	file.Close()
}

func printLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(file))
	}
}
