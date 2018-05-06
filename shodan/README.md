# Shodan Plugin

## Details

* This plugin will return the folloiwng information:
  * Hostnames
  * Organization
  * ASN
  * Ports 
  * You'll need to set the API key with the following command
    * ```export CHATBOT_SHODAN={APIKEY}```

```
CMD: /shodan
Help: /shodan {ip}
```

### Example

```
/shodan 8.8.8.8
---------------
HostNames: google-public-dns-a.google.com 
Organization: Google
ASN: AS15169
Ports: 53/udp
```