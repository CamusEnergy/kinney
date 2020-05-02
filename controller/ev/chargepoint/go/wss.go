package chargepoint

import "encoding/xml"

// securityHeader is a WS-Security ("WSS") "Security" message, to be used as a
// SOAP Header, that embeds a "UsernameToken" as the means of authentication.
type securityHeader struct {
	XMLName xml.Name `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd Security"`

	// Must always be "1".
	// https://www.w3.org/TR/2000/NOTE-SOAP-20000508/#_Toc478383500
	MustUnderstand uint8 `xml:"http://schemas.xmlsoap.org/soap/envelope/ mustUnderstand,attr"`

	Username string `xml:"UsernameToken>Username"`
	Password struct {
		Type  string `xml:"Type,attr"`
		Value string `xml:",chardata"`
	} `xml:"UsernameToken>Password"`
}

const passwordTypeText = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"

func newSecurityHeader(username, passwordValue string) *securityHeader {
	out := &securityHeader{
		MustUnderstand: 1,
		Username:       username,
	}
	// Use field assignment instead of a struct literal here to avoid
	// repeating the type definition for the `Password` field.
	out.Password.Type = passwordTypeText
	out.Password.Value = passwordValue
	return out
}
