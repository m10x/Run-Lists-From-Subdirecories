package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
)

/*
go run C:\Users\max\Documents\git\Run-Lists-From-Subdirecories\main.go -pathlist "D:\aWCVS\tests\" -lines 0 -namelist "domains" -command 'C:\WCVS\wcvs-0.4.8.exe -ppath "C:\Users\max\Downloads\certificate.pem" -gc -gr -gp "$PATH" -u "file:$LIST" -r 5 -rl 5 -hw "C:/WCVS/wordlists/headers" -pw "C:/WCVS/wordlists/top-parameters" -useragentchrome -reqrate 1'
*/

const replaceString = "$LIST"

var sortedList []string

func main() {
	if runtime.GOOS == "windows" {
		reset = ""
		red = ""
		yellow = ""
		green = ""
	}

	pathList, nameList, _, command := parseFlags()

	if pathList == "" {
		fmt.Printf("%sError: -list wasn't specified%s\n", red, reset)
		os.Exit(1)
	}
	if nameList == "" {
		fmt.Printf("%sError: -output wasn't specified%s\n", red, reset)
		os.Exit(2)
	}
	if !strings.Contains(command, replaceString) {
		fmt.Printf("%sWarning: -command doesn't contain %s%s\n", yellow, replaceString, reset)
	}

	//sliceList := readLocalFile(pathList)
	started := 0
	finished := 1
	var wg sync.WaitGroup
	var m sync.Mutex

	filepath.Walk(pathList, func(path string, info os.FileInfo, err error) error {
		if info.Name() == nameList {
			wg.Add(1)
			started++

			newPath := strings.TrimSuffix(path, nameList)
			commandNew := strings.Replace(command, replaceString, path, -1)
			commandNew = strings.Replace(commandNew, "$PATH", newPath, -1)

			if strings.Contains(path, "roblox") {
				commandNew += " -sh ''User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 hackeronetest-m10xde''"
			} else if strings.Contains(path, "amazon") {
				commandNew += " -sh ''User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 amazonvrpresearcher_m10xde''"
			} else if strings.Contains(path, "jimdo") {
				commandNew += " -sh ''X-Bug-Bounty: HackerOne-m10xde''"
			} else if strings.Contains(path, "lazada") {
				commandNew += " -sh ''User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 LZD_YWH_BBP_PUBLIC''"
			} else if strings.Contains(path, "hilton") {
				commandNew += " -sh ''User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 HackerOne''"
			} else if strings.Contains(path, "upwork") {
				commandNew += " -sh ''User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 bugcrowd''"
			}

			fmt.Printf("started %d: %s\n", started, commandNew)

			go func(newPath string, commandNew string) {
				defer wg.Done()
				cmd := exec.Command("powershell", "start-process", "-Wait", "powershell.exe", "-argumentlist", "'(Get-Host).ui.RawUI.WindowTitle = ''"+newPath+"'' ; "+commandNew+"'")
				cmd.Start()
				cmd.Wait()
				m.Lock()
				fmt.Printf("finished (%d/%d): %s\n", finished, started, commandNew)
				finished++
				m.Unlock()
			}(newPath, commandNew)
			time.Sleep(100 * time.Millisecond)

		}
		return nil
	})
	wg.Wait()
}

func readLocalFile(path string) []string {

	w, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%sError while Reading list %s: %s%s\n", red, path, err.Error(), reset)
		os.Exit(3)
	}

	return strings.Split(string(w), "\n")
}

func parseFlags() (string, string, int, string) {
	var pathList string
	var nameList string
	var lines int
	var command string

	flag.StringVar(&pathList, "pathlist", "", "path to the list")
	flag.StringVar(&nameList, "namelist", "", "path to output folder")
	flag.IntVar(&lines, "lines", 0, "after how many lines should be splitted? Default is 0 (=unlimited)")
	flag.StringVar(&command, "command", "", "command to run. Use $LIST where the path of a splitted list shall be inserted")

	flag.Parse()

	return pathList, nameList, lines, command
}
