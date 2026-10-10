package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func getHostName() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostName, nil
}

func main() {

	fmt.Println("System Check")
	fmt.Println("------------")
	//get hostname
	hostName, err := getHostName()
	if err != nil {
		fmt.Println("Error retrieving hostname:", err)
		return
	}
	// Get CPU information using PowerShell command
	output, err := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_Processor).Name").Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	cpu := strings.TrimSpace(string(output))

	fmt.Println("CPU:", cpu)

	// Get uptime using GetTickCount64 from kernel32.dll
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getTickCount64 := kernel32.NewProc("GetTickCount64")

	ret, _, _ := getTickCount64.Call()

	uptime := time.Duration(ret) * time.Millisecond
	fmt.Println("Hostname:", hostName)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	fmt.Printf("Uptime: %dh %dm %ds\n", hours, minutes, seconds)
}
