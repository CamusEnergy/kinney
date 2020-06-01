package schema

import (
	"encoding/xml"
	"math/big"
)

// API Guide (§ 8.1): "Use this call to return a list of stations.  This will
// not return stations that you don't have access rights to.  For example, it
// will not return a public station unless you either own the station or have
// been granted rights by the station's owner."
//
// API Guide (§ 8.1.1): "Up to 500 stations will be returned by this method."
type GetStationsRequestParams struct {
	SearchQuery struct {
		// API Guide (§ 8.1.2): "A unique station identifier used in
		// ChargePoint.  This identifier never changes, even when the
		// station's head assembly is swapped.  Format:
		// CPNID:StationIdentifier."
		StationID           string `xml:"stationID,omitempty"`
		StationManufacturer string `xml:"stationManufacturer,omitempty"`
		StationModel        string `xml:"stationModel,omitempty"`
		// API Guide (§ 8.1.2): "Name of the station (wild card
		// characters are allowed).  It should be searched for by both
		// company name (the name of the organization that owns the
		// charging station) and station name.  Company name is
		// displayed on Line 1 of the charging station (if applicable)
		// and the station name is displayed on Line 2 of the charging
		// station (if applicable)."
		StationName string `xml:"stationName,omitempty"`

		// API Guide (§ 8.1.2): "Address around which you want to see
		// stations.  This can be street address or complete address
		// (street address, city, state, zip code, country)."
		Address    string `xml:"Address,omitempty"`
		City       string `xml:"City,omitempty"`
		State      string `xml:"State,omitempty"`
		Country    string `xml:"Country,omitempty"`
		PostalCode string `xml:"postalCode,omitempty"`

		Coordinate *Coordinate `xml:"Geo,omitempty"`
		// API Guide (§ 8.1.2): "Distance from the station's specified
		// lat/long (Geo) from which you want to retrieve station
		// information.  Default is 5"
		Proximity *big.Rat `xml:"Proximity,omitempty"`
		// API Guide (§ 8.1.2): "Default value for proximity unit is M.
		// Can have values: M (miles), N (Nautical miles), K
		// (Kilometer), F (Feet), I (Inches)."
		ProximityUnit string `xml:"proximityUnit,omitempty"`

		// WSDL: "Possible values 1 is Level 1, 2 is Level 2, 3 is Level
		// 1 2 ,4 is DC charger, 5 Level 1, Level2, DC charger"
		//
		// API Guide (§ 8.1.2): "Station level type where 1 is 'Level
		// 1', 2 is 'Level 2', 3 is 'Level 3', and 4 is 'DC Fast'.  If a
		// station has more than one level (for example, the station
		// provides both level 1 and level 2 charging), the response
		// will includ both level (1,2).  Note: This parameter is for
		// 'US Stations' and 'AU Stations' only (and is used instead of
		// 'Mode')."
		Level string `xml:"Level,omitempty"`
		// API Guide (§ 8.1.2): "Station mode type where 1 is 'Mode 1',
		// 3 is 'Mode 3', and 4 is 'DC Fast'.  If the station has more
		// than one mode (for example, the station provides both mode 1
		// with a domestic socket and mode 3 charging with an IEC 62196
		// Type 2 socket), the response will include both modes (1,3).
		// Note: This parameter is for "EU Stations" only (and is used
		// instead of "Level").
		Mode string `xml:"Mode,omitempty"`

		PricingSession *PricingSession `xml:"Pricing,omitempty"`

		// API Guide (§ 8.1.2): "Whether or not the station can be
		// reserved: '1' - the station can be reserved.  '0' - the
		// station cannot be reserved."
		Reservable uint8 `xml:"Reservable,omitempty"`

		// API Guide (§ 8.1.2): "Connector type.  For example: NEMA
		// 5-20R, J1772, ALFENL3, Shuko."
		Connector string `xml:"Connector,omitempty"`

		// API Guide (§ 8.1.2): "Nominal voltage (V)."
		Voltage string `xml:"Voltage,omitempty"`
		// API Guide (§ 8.1.2): "Current supported (A)."
		Current string `xml:"Current,omitempty"`
		// API Guide (§ 8.1.2): "Power supported (kW)."
		PowerKW string `xml:"Power,omitempty"`

		// API Guide (§ 8.1.2): "Array of serial numbers of stations
		// identified as a 'demo'.  Used only for client applications
		// that need to access stations identified as 'demo'.
		DemoStations *struct {
			SerialNumbers []string `xml:"serialNumber"`
		} `xml:"demoSerialNumber,omitempty"`

		// API Guide (§ 8.1.2): "The org identifier CPNID:CompanyID"
		OrganizationID   string `xml:"orgID,omitempty"`
		OrganizationName string `xml:"organizationName,omitempty"`
		StationGroupID   string `xml:"sgID,omitempty"`
		StationGroupName string `xml:"sgName,omitempty"`

		// API Guide (§ 8.1.2): "Start index for the stations that match
		// the query."
		StartRecord int32 `xml:"startRecord,omitempty"`
		// API Guide (§ 8.1.2): "Number of stations to return in the
		// response.  Maximum is 500, and if left blank, the method will
		// return up to 500 stations."
		NumRecords int32 `xml:"numStations,omitempty"`

		// Undocumented in the API Guide, but exist in the WSDL.
		SerialNumber          string      `xml:"serialNumber,omitempty"`
		StationActivationDate xsdDateTime `xml:"stationActivationDate,omitempty"`
	} `xml:"searchQuery"`
}

type GetStationsRequest struct {
	XMLName xml.Name `xml:"urn:dictionary:com.chargepoint.webservices getStations"`
	GetStationsRequestParams
}

type PricingSession struct {
	StartTime xsdTime `xml:"startTime"`

	// WSDL: "Expected duration of charging session in minutes."
	//
	// API Guide (§ 8.1.2): "Estimated duration of session in hours"
	ExpectedDurationMinutes int32 `xml:"Duration"`
	// API Guide (§ 8.1.2): "Estimated energy needed for a charging session in kWh."
	ExpectedDurationKWh float64 `xml:"energyRequired"`

	// This field is part of the XSD type used by "getStations", but is not
	// documented for that request.  It is only documented in § 5.1
	// ("getPublicStations"), § 5.2 ("getPublicStationStatus"), and § 8.2
	// ("getStatusStatus").  It does however, appear in the sample
	// "getStations" request in § 8.1.4.
	//
	// API Guide (§ 5.1.2): "If a session is active, present amount of power
	// in kW being delivered to the vehicle."
	PowerDrawKW float64 `xml:"vehiclePower"`
}

type GetStationsResponse struct {
	XMLName xml.Name `xml:"getStationsResponse"`

	commonResponseParameters

	Stations []struct {
		StationID           string `xml:"stationID,omitempty"`
		StationManufacturer string `xml:"stationManufacturer,omitempty"`
		StationModel        string `xml:"stationModel,omitempty"`
		StationMACAddress   string `xml:"stationMacAddr,omitempty"`
		StationSerialNumber string `xml:"stationSerialNum,omitempty"`

		StationGroupID   string `xml:"sgID,omitempty"`
		StationGroupName string `xml:"sgName,omitempty"`
		OrganizationID   string `xml:"orgID"`
		OrganizationName string `xml:"organizationName"`

		// API Guide (§ 8.1.3): "Complete address (street address, city,
		// state, zip code, country)."
		Address    string `xml:"Address,omitempty"`
		City       string `xml:"City,omitempty"`
		State      string `xml:"State,omitempty"`
		Country    string `xml:"Country,omitempty"`
		PostalCode string `xml:"postalCode,omitempty"`

		NumPorts int32  `xml:"numPorts,omitempty"`
		Ports    []Port `xml:"Port,omitempty"`

		// API Guide (§ 8.1.3): "The ISO 4217 code for the currency used
		// on the station.  For eample, US Dollar = USD, Canadian Dollar
		// = CAD, Euro = EUR."
		CurrencyCode string `xml:"currencyCode,omitempty"`

		// `maxOccurs` for this element in the WSDL is 2.
		PricingSpecification []PricingSpecification `xml:"Pricing,omitempty"`

		DriverSupportPhoneNumber string `xml:"mainPhone,omitempty"`

		// Undocumented in the API Guide, but exist in the WSDL.
		StationActivationDate xsdDateTime `xml:"stationActivationDate,omitempty"`
		DriverName            string      `xml:"driverName,omitempty"`
		DriverAddress         string      `xml:"driverAddress,omitempty"`
		DriverEmail           string      `xml:"driverEmail,omitempty"`
		DriverPhoneNumber     string      `xml:"driverPhoneNumber,omitempty"`
		LastModifiedDate      xsdDateTime `xml:"lastModifiedDate,omitempty"`
		ModTimeStamp          xsdDateTime `xml:"modTimeStamp,omitempty"`
		TimezoneOffset        string      `xml:"timezoneOffset,omitempty"`
	} `xml:"stationData,omitempty"`

	// API Guide (§ 8.1.3): "Indicates that the number of stations that
	// match this query is greater than the maximum number of stations that
	// can be returned in one call (currently 500), and therefore the list
	// was truncated."
	//
	// This field has `type="xsd:int"` in the WSDL, but is semantically
	// boolean (the value should either be "0" or "1"), and can safely be
	// parsed as such.
	Truncated bool `xml:"moreFlag,omitempty"`
}

type PricingSpecification struct {
	// API Guide (§ 8.1.3): "Pricing Type (Session, Hourly, or kWh)"
	Type string `xml:"Type,omitempty"`

	// API Guide (§ 8.1.3): "The start time of a pricing session."
	StartTime xsdTime `xml:"startTime,omitempty"`
	// API Guide (§ 8.1.3): "The end time of a pricing session."
	EndTime xsdTime `xml:"endTime,omitempty"`

	// API Guide (§ 8.1.3): "Maximum time allowed for a session."
	MaxSessionTime string `xml:"sessionTime,omitempty"`

	// API Guide (§ 8.1.3): "The minimum price charged for a session."
	MinPrice float64 `xml:"minPrice,omitempty"`
	// API Guide (§ 8.1.3): "The maximum price charged for a session."
	MaxPrice float64 `xml:"maxPrice,omitempty"`

	// API Guide (§ 8.1.3): "The hourly price if this mode of pricing is
	// enabled"
	UnitPricePerHour float64 `xml:"unitPricePerHour,omitempty"`
	// API Guide (§ 8.1.3): "The session price if this mode of pricing is
	// enabled"
	UnitPricePerSession float64 `xml:"unitPricePerSession,omitempty"`
	// API Guide (§ 8.1.3): "The kWh price if this mode of pricing is
	// enabled"
	UnitPricePerKWh float64 `xml:"unitPricePerKWh,omitempty"`

	// This field is documented in the API Guide, but does not exist in the
	// WSDL.
	//
	// API Guide (§ 8.1.3): "The hourly price for the first portion of the
	// pricing specification if pricing varies by length of time"
	UnitPriceForFirst float64 `xml:"unitPriceForFirst,omitempty"`

	// This field is in the WSDL, but not in the API Guide.  It is most
	// likely the duration of that "first unit" (corresponding to
	// `UnitPriceForFirst`).
	initialUnitPriceDuration int32 `xml:"initialUnitPriceDuration,omitempty"`

	// API Guide (§ 8.1.3): "The hourly price for the second portion of the
	// pricing specification if pricing varies by length of time"
	UnitPricePerHourThereafter float64 `xml:"unitPricePerHourThereafter,omitempty"`
}

type Port struct {
	// API Guide (§ 8.1.3): "Identifier of the port.  This ID is 1 based."
	PortNumber string `xml:"portNumber,omitempty"`

	StationName string      `xml:"stationName,omitempty"`
	Coordinate  *Coordinate `xml:"Geo,omitempty"`
	Reservable  uint8       `xml:"Reservable,omitempty"`
	Level       string      `xml:"Level,omitempty"`
	Mode        string      `xml:"Mode,omitempty"`
	Connector   string      `xml:"Connector,omitempty"`
	Voltage     string      `xml:"Voltage,omitempty"`
	Current     string      `xml:"Current,omitempty"`
	PowerKW     string      `xml:"Power,omitempty"`

	// Undocumented in the API Guide, but exist in the WSDL
	Description   string      `xml:"Desription,omitempty"`
	Status        string      `xml:"Status,omitempty"`
	Timestamp     xsdDateTime `xml:"timeStamp,omitempty"`
	EstimatedCost float64     `xml:"estimatedCost,omitempty"`
}
