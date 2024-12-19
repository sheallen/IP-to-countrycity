package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ip2location/ip2location-go/v9"
)

// 載入虛擬IP對應實體IP
func loadIPMapping(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var ipMapping map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ipMapping); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return ipMapping, nil
}

func main() {
	defer func() {
		fmt.Println("\nPress Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}()

	ipMapping, err := loadIPMapping("ip_mapping.dat")
	if err != nil {
		fmt.Println("Error loading IP mapping:", err)
		return
	}

	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB3.BIN")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Virtual IP Address: ")
	inputIP, _ := reader.ReadString('\n')
	virtualIP := strings.TrimSpace(inputIP)

	//convert IPv4 ip to real IP
	realIP, exists := ipMapping[virtualIP]
	if !exists {
		fmt.Println("----------------------------")
		fmt.Printf("No mapping found for Virtual IP: %s\n", virtualIP)
		fmt.Println("----------------------------")
		return
	}

	fmt.Printf("Virtual IP: %s maps to Real IP: %s\n", virtualIP, realIP)

	if net.ParseIP(realIP) == nil {
		fmt.Println("Invalid Real IP address format in mapping.")
		return
	}

	results, err := db.Get_all(realIP)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nReal IP Address Information:")
	fmt.Println("----------------------------")
	fmt.Printf("IP Address: %s\n", realIP)
	fmt.Printf("Country Code: %s\n", results.Country_short)
	fmt.Printf("Country Name: %s\n", results.Country_long)
	fmt.Printf("Region: %s\n", results.Region)
	fmt.Printf("City: %s\n", results.City)
}
