package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ANSI color codes
const (
	Red   = "\033[0;31m"
	Blue  = "\033[0;34m"
	Green = "\033[0;32m"
	NC    = "\033[0m" // No Color
)

func main() {
	filepath := flag.String("f", "", "Path to the subdomains file")
	flag.Parse()

	// check for provided subdomains file
	if *filepath == "" {
		fmt.Println("Usage: subov88r -filepath subdomains.txt")
		os.Exit(88)
	}

	// open subdomains file
	file, err := os.Open(*filepath)
	if err != nil {
		fmt.Println("Error While Opening a file:", err)
		return
	}
	defer file.Close()

	// loop over the list of the subdomains
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		subdomain := scanner.Text()
		cname, err := getCname(subdomain)
		if err != nil {
			fmt.Printf("Error getting CNAME for %s: %v\n", subdomain, err)
			continue
		}

		status, err := getStatus(subdomain)
		if err != nil {
			fmt.Printf("Error getting status for %s: %v\n", subdomain, err)
			continue
		}

		fmt.Printf("%sSubdomain: %s %s, %s CNAME: %s %s, %sStatus: %s%s\n", Red, subdomain, NC, Blue, cname, NC, Green, status, NC)
	}
}

// queries the CNAME record for a given subdomain
func getCname(subdomain string) (string, error) {
	cmd := exec.Command("dig", "+short", subdomain, "CNAME")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Error while running $ dig +short %s CNAME: %v", subdomain, err)
	}
	return strings.TrimSpace(string(output)), nil
}

// Get the status from the dig output
func getStatus(subdomain string) (string, error) {
	cmd := exec.Command("dig", subdomain)
	digResult, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error while running dig for %s: %v", subdomain, err)
	}

	digOutput := string(digResult)
	status := ""
	lines := strings.Split(digOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "status:") {
			fields := strings.Fields(line)
			if len(fields) >= 4 {
				status = fields[5]
				break
			}
		}
	}
	return status, nil
}
