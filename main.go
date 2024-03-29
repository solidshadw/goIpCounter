package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Define command-line flags
	ipFile := flag.String("ci", "", "Path to the file containing IP addresses")
	subnetFile := flag.String("s", "", "Path to the file containing subnets")
	foundFlag := flag.Bool("f", false, "Flag to specify IPs found in subnets")
	notFoundFlag := flag.Bool("nf", false, "Flag to specify IPs not found in subnets")
	flag.Parse()

	if *ipFile == "" || *subnetFile == "" || !(*foundFlag != *notFoundFlag) {
		fmt.Println("Usage: go run main.go -ci <ip_file> -s <subnet_file> (-f | -nf)")
		os.Exit(1)
	}

	// Read IPs and subnets from the files
	ipsAndSubnets, err := readLines(*ipFile)
	if err != nil {
		fmt.Println("Error reading IP file:", err)
		os.Exit(1)
	}

	// Read subnets from the subnet file
	subnets, err := readLines(*subnetFile)
	if err != nil {
		fmt.Println("Error reading subnet file:", err)
		os.Exit(1)
	}

	// Store subnets in a map for efficient lookup
	subnetMap := make(map[string]bool)
	for _, subnet := range subnets {
		subnetMap[subnet] = true
	}

	// Check each IP or subnet against each subnet
	for _, ipOrSubnet := range ipsAndSubnets {
		ipOrSubnet = strings.TrimSpace(ipOrSubnet) // Trim whitespace
		if strings.Contains(ipOrSubnet, "/") {
			// It's a subnet, check if it's in any of the subnets
			found := false
			for _, subnet := range subnets {
				if ipInSubnet(ipOrSubnet, subnet) {
					found = true
					break
				}
			}
			if (*foundFlag && found) || (*notFoundFlag && !found) {
				fmt.Println(ipOrSubnet)
			}
		}
	}

	for _, ipOrSubnet := range ipsAndSubnets {
		ipOrSubnet = strings.TrimSpace(ipOrSubnet) // Trim whitespace
		if !strings.Contains(ipOrSubnet, "/") {
			// Check if it's an IP address
			if isValidIP(ipOrSubnet) {
				found := false
				for _, subnet := range subnets {
					if ipInSubnet(ipOrSubnet, subnet) {
						found = true
						break
					}
				}
				if (*foundFlag && found) || (*notFoundFlag && !found) {
					fmt.Println(ipOrSubnet)
				}
			} else {
				fmt.Printf("%s is not a valid IP address\n", ipOrSubnet)
			}
		}
	}
}

// Read lines from a file and return them as a slice of strings
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, scanner.Err()
}

// Check if an IP address or subnet is within a subnet
func ipInSubnet(ipOrSubnetStr, subnetStr string) bool {
	// Try parsing as a subnet
	_, ipSubnet, err := net.ParseCIDR(ipOrSubnetStr)
	if err != nil {
		// Not a valid subnet, try parsing as an IP
		ip := net.ParseIP(ipOrSubnetStr)
		if ip == nil {
			// Not a valid IP or subnet
			return false
		}

		// Parse subnetStr as a subnet
		_, subnet, err := net.ParseCIDR(subnetStr)
		if err != nil {
			return false
		}

		// Check if IP is in subnet
		return subnet.Contains(ip)
	}

	// Parse subnetStr as a subnet
	_, subnet, err := net.ParseCIDR(subnetStr)
	if err != nil {
		return false
	}

	// Check if ipSubnet is a subset of subnet or vice versa
	return subnet.Contains(ipSubnet.IP) || ipSubnet.Contains(subnet.IP)
}

// Check if a string represents a valid IP address (including CIDR notation)
func isValidIP(ipStr string) bool {
	// Regular expression to match IPv4 address format (with optional CIDR notation)
	ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(/\d{1,2})?\s*$`
	match, err := regexp.MatchString(ipPattern, ipStr)
	return match && err == nil
}
