package models

// Booking represents a booking record in your database.
type Booking struct {
	ReferenceNumber string   `json:"reference_number"`
	BookingParty    string   `json:"booking_party"`
	Consignee       string   `json:"consignee"`
	VesselName      string   `json:"vessel_name"`
	VoyageNumber    string   `json:"voyage_number"`
	PortOfLoading   string   `json:"port_of_loading"`
	PortOfDischarge string   `json:"port_of_discharge"`
	Containers      []string `json:"containers"`
}
