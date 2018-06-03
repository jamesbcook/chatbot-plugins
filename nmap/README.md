# Nmap Plugins

* This plugin will query the Nmap api and return the scan results

```
CMD: /nmap
Help: /nmap {"apiIP:apiPort" "nmap args" }
```

### Example

```
/nmap "172.31.98.119:9292" "--open -p 80,8080,443 google.com"
-------------------------
Starting Nmap 7.70 ( https://nmap.org ) at 2018-06-01 08:17 MST
Nmap scan report for google.com (172.217.11.78)
Host is up (0.027s latency).
rDNS record for 172.217.11.78: lax17s34-in-f14.1e100.net
Not shown: 1 filtered port
Some closed ports may be reported as filtered due to --defeat-rst-ratelimit
PORT    STATE SERVICE
80/tcp  open  http
443/tcp open  https

Nmap done: 1 IP address (1 host up) scanned in 1.33 seconds
```
