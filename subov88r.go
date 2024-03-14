package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
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
	// Parse command-line arguments
	filepath := flag.String("f", "", "Path to the subdomains file")
	flag.Parse()

	// Check for provided subdomains file
	if *filepath == "" {
		fmt.Println("Usage: subov88r -f subdomains.txt")
		os.Exit(88)
	}

	// Open subdomains file
	file, err := os.Open(*filepath)
	if err != nil {
		fmt.Println("Error while opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Loop over the list of subdomains
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subdomain := scanner.Text()

		// Get the CNAME record for the subdomain
		cname, _ := net.LookupCNAME(subdomain)

		// Get the status of the subdomain
		status, err := getStatus(subdomain)
		if err != nil {
			fmt.Printf("Error getting status for %s: %v\n", subdomain, err)
			continue
		}

		isVuln := azureSTO(cname, status)

		if isVuln {
			fmt.Printf("[%v,%v,%v] Possiply Vulnerable to subdomain takevover", subdomain, cname, status)
			continue
		}

		// Print results with ANSI colors
		fmt.Printf("%sSubdomain: %s %s, %sCNAME: %s %s, %sStatus: %s%s\n", Red, subdomain, NC, Blue, cname, NC, Green, status, NC)
	}
}

// getStatus gets the status from the dig output
func getStatus(subdomain string) (string, error) {
	cmd := exec.Command("dig", subdomain)
	digResult, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
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

// function that check for subdomain takeover in azure services
func azureSTO(cname string, status string) bool {
	azureRegex := regexp.MustCompile(`(?i)^(?:[a-z0-9-]+\.)?(?:cloudapp\.net|azurewebsites\.net|cloudapp\.azure\.com)$`)

	if strings.Contains(status, "NXDOMAIN") && azureRegex.MatchString(cname) {
		url := fmt.Sprintf("https://%s", cname)
		_, err := http.Get(url)
		if err != nil {
			return true // If there's an error, assume it's a possible subdomain takeover
		}
	}
	return false
}
