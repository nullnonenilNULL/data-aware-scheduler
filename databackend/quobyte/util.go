package quobyte

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func parseXattrSegments(xattrString string) []*segment {
	segments := []*segment{}
	segs := strings.Split(xattrString, "segment")

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '_'
	}

	for _, seg := range segs {
		fields := strings.FieldsFunc(seg, f)
		if fields[0] == "posix_attrs" {
			continue
		}

		startOffset, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Printf("Error during parse: %v", err)
		}

		length, err := strconv.Atoi(fields[3])
		if err != nil {
			log.Printf("Error during parse: %v", err)
		}

		s := &segment{
			startOffset: startOffset,
			length:      length,
		}

		st := &stripe{}

		for i := 5; i < len(fields); i++ {
			switch fields[i] {
			case "version":
				v, err := strconv.Atoi(fields[i+1])
				if err != nil {
					log.Printf("Error during parse: %v", err)
				}
				st.version = v
				i++
			case "device_id":
				id, err := strconv.ParseUint(fields[i+1], 10, 64)
				if err != nil {
					log.Printf("Error during parse: %v", err)
				}
				st.deviceIDs = append(st.deviceIDs, id)
				i++
			}
		}
		s.stripe = st

		segments = append(segments, s)
	}

	return segments
}

func validateAPIURL(apiURL string) error {
	url, err := url.Parse(apiURL)
	if err != nil {
		return err
	}

	if url.Scheme == "" {
		return fmt.Errorf("Scheme is no set in URL: %s", apiURL)
	}

	if url.Host == "" {
		return fmt.Errorf("Scheme is no set in URL: %s", apiURL)
	}

	return nil
}

func getAllFilesInsideDir(dir string) []string {
	resFiles := []string{}

	err := filepath.Walk(dir, func(searchPath string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			resFiles = append(resFiles, searchPath)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return resFiles
}
