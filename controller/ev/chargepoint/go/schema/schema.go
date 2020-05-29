// Package schema defines Go structs that map to and from the ChargePoint API
// WSDL types (using `encoding/xml` marshaling/unmarshaling rules).
//
// "API Guide" comments are quoted (with section numbers given before each
// quote) verbatim from the following reference:
//
//     ChargePoint Web Services API Programmer’s Guide
//     Document Part Number: 75-001102-01
//     Document Revision: 8
//     Revision Date: 2016-04-20
//
// In addition, the WSDL file that is modeled here, while unversioned, has the
// following MD5 hash:
//
//     79595c612a701f4bf7e3f9d1ede2f7a9
//
// Note that struct and field names have been chosen for clarity of semantic
// meaning, and not for strict equivalence with the corresponding WSDL types.
// Struct tags are used to map each field to the corresponding XML element name.
//
// In addition, the schema is defined by the element names in the XSD, not the
// XSD type names.  As such, no structs except those at the "top level" (i.e.,
// ones mapping to elements whose parents do not uniquely identify their name)
// should define an element name (via an `xml.Name` field), as it will always be
// overwritten by the containing element.
//
// TODO(james): As a result of the recursive way that SOAP Envelope messages are
// currently umarshaled, response types may not specify an XML namespace in
// their `xml.Name` field.  This limitation may be resolved by using an
// `interface{}` value for the body and using a single-pass parse, rather than
// the current recursive/two-pass approach.
package schema
