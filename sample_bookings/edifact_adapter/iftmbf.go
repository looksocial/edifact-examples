package edifact_adapter

import (
	"github.com/looksocial/edifact/examples/bookings/models"
	"github.com/looksocial/edifact/internal/model"
)

// IFTMBFAdapter converts IFTMBF EDIFACT messages to Booking models.
type IFTMBFAdapter struct{}

func NewIFTMBFAdapter() *IFTMBFAdapter {
	return &IFTMBFAdapter{}
}

func (a *IFTMBFAdapter) CanHandle(messageType string) bool {
	return messageType == "IFTMBF"
}

func (a *IFTMBFAdapter) Handle(message *model.Message) (interface{}, error) {
	booking := &models.Booking{}
	for _, segment := range message.Segments {
		switch segment.Tag {
		case "BGM":
			if len(segment.Elements) > 1 {
				booking.ReferenceNumber = segment.Elements[1].Value
			}
		case "NAD":
			if len(segment.Elements) > 0 {
				party := segment.Elements[0].Value
				switch party {
				case "CA": // Carrier/Booking party
					if len(segment.Elements) > 2 {
						booking.BookingParty = segment.Elements[2].Value
					}
				case "CN": // Consignee
					if len(segment.Elements) > 2 {
						booking.Consignee = segment.Elements[2].Value
					}
				}
			}
		case "TDT":
			if len(segment.Elements) > 1 && segment.Elements[1].IsComposite {
				booking.VesselName = segment.Elements[1].Components[0]
			}
			if len(segment.Elements) > 2 {
				booking.VoyageNumber = segment.Elements[2].Value
			}
		case "LOC":
			if len(segment.Elements) > 0 {
				locType := segment.Elements[0].Value
				if locType == "9" && len(segment.Elements) > 1 && segment.Elements[1].IsComposite {
					booking.PortOfLoading = segment.Elements[1].Components[0]
				}
				if locType == "11" && len(segment.Elements) > 1 && segment.Elements[1].IsComposite {
					booking.PortOfDischarge = segment.Elements[1].Components[0]
				}
			}
		case "EQD":
			if len(segment.Elements) > 1 {
				container := segment.Elements[1].Value
				booking.Containers = append(booking.Containers, container)
			}
		}
	}
	return booking, nil
}
