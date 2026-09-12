package driver

import (
	"encoding/json"
	"testing"

	"github.com/signadot/hotrod/services/location"
)

func TestDispatchRequestMarshalUsesSnakeCase(t *testing.T) {
	request := DispatchRequest{
		PickupLocation:  &location.Location{ID: 1},
		DropoffLocation: &location.Location{ID: 2},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal dispatch request: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatalf("decode marshaled dispatch request: %v", err)
	}
	for _, field := range []string{"pickup_location", "dropoff_location"} {
		if _, ok := fields[field]; !ok {
			t.Errorf("marshaled dispatch request is missing %q", field)
		}
	}
	for _, field := range []string{"pickupLocation", "dropoffLocation"} {
		if _, ok := fields[field]; ok {
			t.Errorf("marshaled dispatch request contains legacy field %q", field)
		}
	}
}

func TestDispatchRequestUnmarshalFieldNames(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		wantPickupID   int64
		wantDropoffID  int64
		wantPickupNil  bool
		wantDropoffNil bool
	}{
		{
			name:          "snake case",
			payload:       `{"pickup_location":{"id":1},"dropoff_location":{"id":2}}`,
			wantPickupID:  1,
			wantDropoffID: 2,
		},
		{
			name:          "legacy camel case",
			payload:       `{"pickupLocation":{"id":3},"dropoffLocation":{"id":4}}`,
			wantPickupID:  3,
			wantDropoffID: 4,
		},
		{
			name:          "snake case takes precedence",
			payload:       `{"pickup_location":{"id":5},"pickupLocation":{"id":6},"dropoff_location":{"id":7},"dropoffLocation":{"id":8}}`,
			wantPickupID:  5,
			wantDropoffID: 7,
		},
		{
			name:           "explicit snake case null takes precedence",
			payload:        `{"pickup_location":null,"pickupLocation":{"id":9},"dropoff_location":null,"dropoffLocation":{"id":10}}`,
			wantPickupNil:  true,
			wantDropoffNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request DispatchRequest
			if err := json.Unmarshal([]byte(tt.payload), &request); err != nil {
				t.Fatalf("unmarshal dispatch request: %v", err)
			}

			if tt.wantPickupNil {
				if request.PickupLocation != nil {
					t.Fatalf("pickup location = %#v, want nil", request.PickupLocation)
				}
			} else if request.PickupLocation == nil || request.PickupLocation.ID != tt.wantPickupID {
				t.Errorf("pickup location = %#v, want ID %d", request.PickupLocation, tt.wantPickupID)
			}

			if tt.wantDropoffNil {
				if request.DropoffLocation != nil {
					t.Fatalf("dropoff location = %#v, want nil", request.DropoffLocation)
				}
			} else if request.DropoffLocation == nil || request.DropoffLocation.ID != tt.wantDropoffID {
				t.Errorf("dropoff location = %#v, want ID %d", request.DropoffLocation, tt.wantDropoffID)
			}
		})
	}
}
