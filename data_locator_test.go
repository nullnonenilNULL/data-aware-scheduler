package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestParseXattrOutput(t *testing.T) {
	s, err := ioutil.ReadFile("test/xattr_example.txt")
	if err != nil {
		log.Fatalln(err)
	}

	segments := parseXattrSegments(string(s))

	expectedSegments := 4
	if len(segments) != expectedSegments {
		log.Printf("Expected: %d got %d", expectedSegments, len(segments))
		t.FailNow()
	}
}

func TestConvertSegmentsToDevices(t *testing.T) {
	expectedResults := map[uint64]*device{
		3: &device{
			host:       "",
			deviceType: "",
			dataSize:   32212254720,
			id:         3,
		},
		4: &device{
			host:       "",
			deviceType: "",
			dataSize:   32212254720,
			id:         4,
		},
		5: &device{
			host:       "",
			deviceType: "",
			dataSize:   32212254720,
			id:         5,
		},
		10: &device{
			host:       "",
			deviceType: "",
			dataSize:   32212254720,
			id:         10,
		},
	}

	s, err := ioutil.ReadFile("test/xattr_example.txt")
	if err != nil {
		log.Fatalln(err)
	}

	devices := convertSegmentsToDevices(parseXattrSegments(string(s)))
	if len(devices) != len(expectedResults) {
		log.Printf("Expected: %d got %d", len(expectedResults), len(devices))
		t.FailNow()
	}

	for k, dev := range devices {
		expectedDevice, ok := expectedResults[k]
		if !ok {
			log.Printf("Expected: Device ID %d but not found in %v", k, devices)
			t.FailNow()
		}

		if dev.dataSize != expectedDevice.dataSize {
			log.Printf("Expected DataSize for Device %d: %d got %d", k, expectedDevice.dataSize, dev.dataSize)
			t.FailNow()
		}
	}
}
