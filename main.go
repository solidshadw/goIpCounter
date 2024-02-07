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
    flag.Parse()

    if *ipFile == "" || *subnetFile == "" {
        fmt.Println("Usage: go run main.go -ci <ip_file> -s <subnet_file>")
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
                if !found {
                    fmt.Printf("%s\n", ipOrSubnet)
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
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// Check if an IP address or subnet is within a subnet
func ipInSubnet(ipOrSubnetStr, subnetStr string) bool {
    ipOrSubnet := net.ParseIP(ipOrSubnetStr)
    if ipOrSubnet == nil {
        // Not an IP, try parsing as a subnet
        _, subnet, err := net.ParseCIDR(ipOrSubnetStr)
        if err != nil {
            // Not a valid IP or subnet
            return false
        }
        return subnet.Contains(subnet.IP)
    }

    _, subnet, err := net.ParseCIDR(subnetStr)
    if err != nil {
        return false
    }

    return subnet.Contains(ipOrSubnet)
}

// Check if a string represents a valid IP address (including CIDR notation)
func isValidIP(ipStr string) bool {
    // Regular expression to match IPv4 address format (with optional CIDR notation)
    ipPattern := `^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(/\d{1,2})?$`
    match, err := regexp.MatchString(ipPattern, ipStr)
    return match && err == nil
}
