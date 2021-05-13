package mtab

import (
	"strings"
	"testing"
)

func TestParseOptsEmpty(t *testing.T) {
	opts := ""
	optsMap := parseOpts(opts)
	if len(optsMap) != 0 {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("len(optsMap) = %d, want 0", len(optsMap))
	}
}

func TestParseOptsEqualsOnly(t *testing.T) {
	opts := "="
	optsMap := parseOpts(opts)
	if len(optsMap) != 0 {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("len(optsMap) = %d, want 0", len(optsMap))
	}
}

func TestParseOptsRoutine(t *testing.T) {
	opts := "a=b,c=d,e"
	optsMap := parseOpts(opts)
	if optsMap["a"] != "b" {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("optsMap[\"a\"] = %s, want \"b\"", optsMap["a"])
	} else if optsMap["c"] != "d" {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("optsMap[\"c\"] = %s, want \"d\"", optsMap["c"])
	} else if optsMap["e"] != "true" {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("optsMap[\"e\"] = %s, want \"true\"", optsMap["e"])
	} else if len(optsMap) != 3 {
		t.Logf("optsMap = %v", optsMap)
		t.Fatalf("len(optsMap) = %d, want 3", len(optsMap))
	}
}

func TestParseMTabLineEmpty(t *testing.T) {
	line := ""
	entry, err := ParseMTabLine(line)
	if err == nil {
		t.Fatalf("entry = %v, want MTabEntry{}", entry)
	}
}

func TestParseMTabLineOverfull(t *testing.T) {
	line := "a,b,c,d,e,f"
	entry, err := ParseMTabLine(line)
	if err == nil {
		t.Fatalf("entry = %v, want MTabEntry{}", entry)
	}
}

func TestParseMTabLineNoNums(t *testing.T) {
	line := "a b c d"
	entry, err := ParseMTabLine(line)
	if err != nil {
		t.Fatalf("err = %s, want err = nil", err)
	} else if entry.DumpFreq != 0 {
		t.Fatalf("entry.DumpFreq = %d, want 0", entry.DumpFreq)
	}
}

func TestParseMTabLineInvalidNums(t *testing.T) {
	line := "a b c d e f"
	_, err := ParseMTabLine(line)
	if err == nil {
		t.Fatalf("err = nil, want err != nil")
	}
}

func TestParseMTabLineRoutine(t *testing.T) {
	line := "a b c d,e=fgh,j=klm 1 2"
	entry, err := ParseMTabLine(line)
	if err != nil {
		t.Fatalf("err = %s, want err == nil", err)
	}

	if entry.DevicePath != "a" {
		t.Fatalf("entry.DevicePath = \"%s\", want \"a\"", entry.DevicePath)
	} else if entry.MountPath != "b" {
		t.Fatalf("entry.MountPath = \"%s\", want \"b\"", entry.MountPath)
	} else if entry.FileSystem != "c" {
		t.Fatalf("entry.FileSystem = \"%s\", want \"c\"", entry.FileSystem)
	} else if len(entry.Options) != 3 {
		t.Fatalf("len(entry.Options) = \"%d\", want 3", len(entry.Options))
	} else if val, ok := entry.Options["d"]; val != "true" || !ok {
		t.Fatalf(
			"entry.Options[\"d\"] = \"%s\", want \"true\""+
				"ok = %v, want true",
			entry.Options["d"],
			ok,
		)
	} else if val, ok := entry.Options["e"]; val != "fgh" || !ok {
		t.Fatalf(
			"entry.Options[\"e\"] = \"%s\", want \"fgh\""+
				"ok = %v, want true",
			entry.Options["e"],
			ok,
		)
	} else if val, ok := entry.Options["j"]; val != "klm" || !ok {
		t.Fatalf(
			"entry.Options[\"j\"] = \"%s\", want \"klm\""+
				"ok = %v, want true",
			entry.Options["j"],
			ok,
		)
	}
}

func TestReadMTabEOF(t *testing.T) {
	sr := strings.NewReader("")
	tabs := ReadMTab(sr)
	if len(tabs) != 0 {
		t.Fatalf("len(tabs) = %d, want 0", len(tabs))
	}
}

func TestReadMTabRoutine(t *testing.T) {
	sr := strings.NewReader(
		"a b c d,e=fgh,j=klm 1 2\n" + "a b c d\n",
	)
	tabs := ReadMTab(sr)
	if len(tabs) != 2 {
		t.Fatalf("len(tabs) = %d, want 2", len(tabs))
	}
}
