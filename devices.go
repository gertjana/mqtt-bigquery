package main

import (
	"fmt"
)

const (
	PressureDeviceType = 1
	TrackingDeviceType = 2
)

var devices = map[string]int{
	"sensor01":  PressureDeviceType,
	"tracker01": TrackingDeviceType,
}

type PressureDevice struct {
	Timestamp string
	Celcius float64
	Mbar    float64
}

type TrackingDevice struct {
	Timestamp string
	Latitude  float64
	Longitude float64
}

func deviceFromFields(device string, timestamp string, fields map[string]interface{}) interface{} {
	dt := devices[device]
	switch dt {
	case PressureDeviceType:
		return PressureDevice{
			timestamp, 
			fields["celcius"].(float64), 
			fields["mbar"].(float64),
		}
	case TrackingDeviceType:
		return TrackingDevice{
			timestamp, 
			fields["latitude"].(float64), 
			fields["longitude"].(float64),
		}
	default:
		panic(fmt.Sprintf("Unknown DeviceType: %s", dt))
	}
}
