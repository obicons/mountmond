package mtab

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type MTabEntry struct {
	DevicePath string
	MountPath  string
	FileSystem string
	Options    map[string]string
	DumpFreq   uint
	FSPassNo   uint
}

const (
	minFields = 4
	maxFields = 6
)

/*
 * Reads formatted lines (see man fstab) from reader into stream.
 */
func ReadMTab(reader io.Reader, stream chan<- MTabEntry) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		tab, err := ParseMTabLine(scanner.Text())
		if err != nil {
			log.Println(err)
		} else {
			stream <- tab
		}
	}
	close(stream)
}

/*
 * We'll parse the line, or know the reason why.
 */
func ParseMTabLine(line string) (MTabEntry, error) {
	fields := strings.Split(line, " ")
	if len(fields) > maxFields {
		return MTabEntry{}, fmt.Errorf("Too many fields in %s", line)
	} else if len(fields) < minFields {
		return MTabEntry{}, fmt.Errorf("Too few fields in %s", line)
	}

	dumpFreq := uint64(0)
	if len(fields) > minFields {
		var err error
		dumpFreq, err = strconv.ParseUint(fields[minFields], 10, 32)
		if err != nil {
			return MTabEntry{}, err
		}
	}

	passNo := uint64(0)
	if len(fields) == maxFields {
		var err error
		passNo, err = strconv.ParseUint(fields[maxFields-1], 10, 32)
		if err != nil {
			return MTabEntry{}, err
		}
	}

	return MTabEntry{
		DevicePath: fields[0],
		MountPath:  fields[1],
		FileSystem: fields[2],
		Options:    parseOpts(fields[3]),
		DumpFreq:   uint(dumpFreq),
		FSPassNo:   uint(passNo),
	}, nil
}

/*
 * Transforms a string formatted as opt1=val1,opt2,opt3=val3... into a map.
 * If a value is not provided for an option, "true" is provided.
 */
func parseOpts(opts string) map[string]string {
	out := make(map[string]string)
	fields := strings.Split(opts, ",")
	for _, field := range fields {
		// skip blank fields
		if len(field) == 0 || field == "=" {
			continue
		}

		splitField := strings.Split(field, "=")
		if len(splitField) == 1 {
			out[splitField[0]] = "true"
		} else {
			out[splitField[0]] = splitField[1]
		}
	}
	return out
}
