# Subov88r
Simple Go tool for alnalyzing the subdomains for subdomain takeover vulnerability 
# Usage
- `subov88r -f subdomains.txt`
- `subov88r -f subdomains.txt | cat cname.txt | grep -E 'cloudapp.net|azurewebsites.net|cloudapp.azure.com'`
# install 
- `go install github.com/h0tak88r/subov88r@latest`
