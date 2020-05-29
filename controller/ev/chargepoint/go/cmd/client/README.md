This tool allows issuing requests against the ChargePoint API from the command
line:

```bash
> go build ./cmd/client

> cat credentials.json
{
  "APIKey": "<SNIP>",
  "APIPassword": "<SNIP>"
}

> ./client \
    --credentials="credentials.json" \
    --http_log="http_log.jsonl" \
    --url="<API_URL>" \
    --method="GetCPNInstances" \
    --request="{}" \
    | jq
2020/05/11 22:31:20 Using request: &schema.GetCPNInstancesRequest{XMLName:xml.Name{Space:"", Local:""}}
{
  "XMLName": {
    "Space": "ns1",
    "Local": "getCPNInstancesResponse"
  },
  "ChargePointNetworks": [
    {
      "ChargePointNetworkID": "1",
      "ChargePointNetworkName": "NA",
      "ChargePointNetworkDescription": "ChargePoint Operations"
    },
    {
      "ChargePointNetworkID": "2",
      "ChargePointNetworkName": "EU",
      "ChargePointNetworkDescription": "ChargePoint Europe"
    },
    {
      "ChargePointNetworkID": "3",
      "ChargePointNetworkName": "AU",
      "ChargePointNetworkDescription": "ChargePoint Australia"
    },
    {
      "ChargePointNetworkID": "4",
      "ChargePointNetworkName": "CA",
      "ChargePointNetworkDescription": "ChargePoint Canada"
    }
  ]
}
```

The HTTP log is written in JSONL format (details are in `soap.go`), and so can
also be inspected with `jq`:

```bash
> jq '(.RequestBody, .ResponseBody) |= @base64d' < http_log.jsonl
{
  "RequestTimestamp": "2020-05-11T22:31:20.736844-07:00",
  "RequestMethod": "POST",
  "RequestURL": "https://webservices.chargepoint.com/webservices/chargepoint/services/5.0",
  "RequestHeaders": {
    "Content-Type": [
      "text/xml; charset=utf-8"
    ]
  },
  "RequestBody": "<Envelope xmlns=\"http://schemas.xmlsoap.org/soap/envelope/\"><Header><Security xmlns=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd\" xmlns:envelope=\"http://schemas.xmlsoap.org/soap/envelope/\" envelope:mustUnderstand=\"1\"><UsernameToken><Username>{{SNIP}}</Username><Password Type=\"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText\">{{SNIP}}</Password></UsernameToken></Security></Header><Body><getCPNInstances xmlns=\"urn:dictionary:com.chargepoint.webservices\"></getCPNInstances></Body></Envelope>",
  "ResponseTimestamp": "2020-05-11T22:31:21.140106-07:00",
  "ResponseStatusCode": 200,
  "ResponseBody": "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ns1=\"urn:dictionary:com.chargepoint.webservices\"><SOAP-ENV:Body><ns1:getCPNInstancesResponse><CPN><cpnID>1</cpnID><cpnName>NA</cpnName><cpnDescription>ChargePoint Operations</cpnDescription></CPN><CPN><cpnID>2</cpnID><cpnName>EU</cpnName><cpnDescription>ChargePoint Europe</cpnDescription></CPN><CPN><cpnID>3</cpnID><cpnName>AU</cpnName><cpnDescription>ChargePoint Australia</cpnDescription></CPN><CPN><cpnID>4</cpnID><cpnName>CA</cpnName><cpnDescription>ChargePoint Canada</cpnDescription></CPN></ns1:getCPNInstancesResponse></SOAP-ENV:Body></SOAP-ENV:Envelope>\n"
}
```

`xmllint` can be used to pretty-print the XML:

```console
$ jq '(.RequestBody, .ResponseBody) |= @base64d | .ResponseBody' --raw-output \
    < http_log.jsonl \
    | xmllint --format -
<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="urn:dictionary:com.chargepoint.webservices">
  <SOAP-ENV:Body>
    <ns1:getCPNInstancesResponse>
      <CPN>
        <cpnID>1</cpnID>
        <cpnName>NA</cpnName>
        <cpnDescription>ChargePoint Operations</cpnDescription>
      </CPN>
      <CPN>
        <cpnID>2</cpnID>
        <cpnName>EU</cpnName>
        <cpnDescription>ChargePoint Europe</cpnDescription>
      </CPN>
      <CPN>
        <cpnID>3</cpnID>
        <cpnName>AU</cpnName>
        <cpnDescription>ChargePoint Australia</cpnDescription>
      </CPN>
      <CPN>
        <cpnID>4</cpnID>
        <cpnName>CA</cpnName>
        <cpnDescription>ChargePoint Canada</cpnDescription>
      </CPN>
    </ns1:getCPNInstancesResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>
```
