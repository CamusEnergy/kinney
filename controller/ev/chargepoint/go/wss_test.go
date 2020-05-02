package chargepoint

import (
	"encoding/xml"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const headerXML = `<Security xmlns="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" xmlns:envelope="http://schemas.xmlsoap.org/soap/envelope/" envelope:mustUnderstand="1"><UsernameToken><Username>username</Username><Password Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">password</Password></UsernameToken></Security>`

func TestMarshalSecurityHeader(t *testing.T) {
	header := newSecurityHeader("username", "password")
	if b, err := xml.Marshal(header); err != nil {
		t.Errorf("xml.Marshal(header) = %q; want nil", err)
	} else if diff := cmp.Diff(headerXML, string(b)); diff != "" {
		t.Errorf("newSecurityHeader(\"username\", \"password\") mismatch (-want +got):\n%s", diff)
	}
}
