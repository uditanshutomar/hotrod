package driver

import (
	"encoding/json"
	"fmt"

	"github.com/signadot/hotrod/services/location"
)

type DispatchRequest struct {
	PickupLocation  *location.Location `json:"pickup_location"`
	DropoffLocation *location.Location `json:"dropoff_location"`
}

// UnmarshalJSON accepts the current snake_case field names and the legacy
// camelCase names used by older producers. This keeps mixed-version deployments
// working while Marshal continues to emit only snake_case fields.
func (request *DispatchRequest) UnmarshalJSON(data []byte) error {
	var fields struct {
		PickupLocation        json.RawMessage `json:"pickup_location"`
		DropoffLocation       json.RawMessage `json:"dropoff_location"`
		LegacyPickupLocation  json.RawMessage `json:"pickupLocation"`
		LegacyDropoffLocation json.RawMessage `json:"dropoffLocation"`
	}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	pickup := fields.PickupLocation
	if pickup == nil {
		pickup = fields.LegacyPickupLocation
	}
	dropoff := fields.DropoffLocation
	if dropoff == nil {
		dropoff = fields.LegacyDropoffLocation
	}

	var decoded DispatchRequest
	if pickup != nil {
		if err := json.Unmarshal(pickup, &decoded.PickupLocation); err != nil {
			return fmt.Errorf("decode pickup location: %w", err)
		}
	}
	if dropoff != nil {
		if err := json.Unmarshal(dropoff, &decoded.DropoffLocation); err != nil {
			return fmt.Errorf("decode dropoff location: %w", err)
		}
	}

	*request = decoded
	return nil
}

type Driver struct {
	DriverID    string
	Coordinates string
}
