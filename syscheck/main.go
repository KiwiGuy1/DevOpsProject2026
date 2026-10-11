package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

func getCPUInfo() (string, error) {
	output, err := exec.Command("powershell", "-Command", "(Get-CimInstance Win32_Processor).Name").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
func getRAMInfo() (string, error) {
	output, err := exec.Command(
		"powershell",
		"-Command",
		"(Get-CimInstance Win32_ComputerSystem).TotalPhysicalMemory",
	).Output()
	if err != nil {
		return "", err
	}

	memoryBytes, err := strconv.ParseUint(
		strings.TrimSpace(string(output)),
		10,
		64,
	)
	if err != nil {
		return "", err
	}

	memoryGB := float64(memoryBytes) / (1024 * 1024 * 1024)

	return fmt.Sprintf("%.2f GiB", memoryGB), nil
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
	cpu, err := getCPUInfo()
	if err != nil {
		fmt.Println("Error retrieving CPU info:", err)
		return
	}

	fmt.Println("CPU:", cpu)

	// Get RAM information using PowerShell command
	ram, err := getRAMInfo()
	if err != nil {
		fmt.Println("Error retrieving RAM info:", err)
		return
	}
	fmt.Println("RAM:", ram)

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
