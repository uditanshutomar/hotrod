package driver

import (
	"github.com/signadot/hotrod/services/location"
)

type DispatchRequest struct {
	PickupLocation  *location.Location `json:"pickup_location"`
	DropoffLocation *location.Location `json:"dropoff_location"`
}

type Driver struct {
	DriverID    string
	Coordinates string
}
