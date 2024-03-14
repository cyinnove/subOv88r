# Subov88r
Simple Go tool for alnalyzing the subdomains for subdomain takeover vulnerability specially in azure services.  
# Usage
- `subov88r -f subdomains.txt`
- `subov88r -f subdomains.txt | cat cname.txt | grep -E 'cloudapp.net|azurewebsites.net|cloudapp.azure.com'`
# install 
- `go install github.com/h0tak88r/subov88r@latest`

![image](https://github.com/h0tak88r/subOv88r/assets/108616378/5bdaaf2d-ed34-40f4-91cc-67da4a088519)
