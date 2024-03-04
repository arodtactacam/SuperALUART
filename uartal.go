package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ going to put path here"
)

var appName string

func printLogo() {
	logo := `
 ____                  _____       _   _     
|  _ \ ___  __ _ _ __|_   _|__   | |_| |__  
| |_) / _ \/ _' | '_ \| |/ _ \  | __| '_ \ 
|  _ <  __/ (_| | | | | |  __/  | |_| | | |
|_| \_\___|\__,_|_| |_|_|\___|   \__|_| |_|`
	fmt.Println(logo)
	fmt.Println("\n" + appName + "\n")
}

func openTerminal(port string, baudRate int) {
	cmd := exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("cu -l %s -s %d", port, baudRate))
	err := cmd.Run()
	if err != nil {
		log.Printf("Error opening terminal for port %s: %v\n", port, err)
	}
}

func readFromSerial(port string) {
	cmd := exec.Command("gnome-terminal", "--", "bash", "-c", fmt.Sprintf("cu -l %s", port))
	err := cmd.Run()
	if err != nil {
		log.Printf("Error reading from port %s: %v\n", port, err)
	}
}

func main() {
	flag.StringVar(&appName, "name", "My UART Terminal", "Name of the application")
	flag.Parse()

	printLogo()

	// Enumerate serial ports
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	// Display available ports
	fmt.Println("Available serial ports:")
	for i, port := range ports {
		fmt.Printf("%d: %s\n", i+1, port)
	}

	// Prompt user to enter baud rate
	var baudRate int
	fmt.Print("Enter the baud rate: ")
	reader := bufio.NewReader(os.Stdin)
	baudRateStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading baud rate:", err)
	}
	baudRateStr = strings.TrimSpace(baudRateStr)
	baudRate, err = strconv.Atoi(baudRateStr)
	if err != nil {
		log.Fatal("Invalid baud rate:", err)
	}

	// Open terminals for each port
	for _, port := range ports {
		go openTerminal(port, baudRate)
		go readFromSerial(port)
	}

	// Create a channel to listen for OS signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop to keep the program running (press Ctrl+C to exit)
	for {
		select {
		case <-sigc:
			return
		default:
			time.Sleep(time.Second) // Sleep to keep the program running
		}
	}
}
