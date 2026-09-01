package main

import (
	"strings"
	"testing"
)

func TestResolveBomValidationMode(t *testing.T) {
	tests := []struct {
		name        string
		disabled    string
		wantMode    string
		wantResolve bool
	}{
		{
			name:        "unset does not resolve",
			disabled:    "",
			wantMode:    "",
			wantResolve: false,
		},
		{
			name:        "true disables validation",
			disabled:    "true",
			wantMode:    "DISABLED",
			wantResolve: true,
		},
		{
			name:        "false enables validation",
			disabled:    "false",
			wantMode:    "ENABLED",
			wantResolve: true,
		},
		{
			name:        "case insensitive true disables validation",
			disabled:    "TRUE",
			wantMode:    "DISABLED",
			wantResolve: true,
		},
		{
			name:        "unexpected value does not resolve",
			disabled:    "maybe",
			wantMode:    "",
			wantResolve: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMode, gotResolve := resolveBomValidationMode(tt.disabled)
			if gotMode != tt.wantMode {
				t.Fatalf("resolveBomValidationMode(%q) mode = %q, want %q", tt.disabled, gotMode, tt.wantMode)
			}
			if gotResolve != tt.wantResolve {
				t.Fatalf("resolveBomValidationMode(%q) resolved = %v, want %v", tt.disabled, gotResolve, tt.wantResolve)
			}
		})
	}
}

func TestResolveCadence(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		wantHours   string
		wantResolve bool
	}{
		{name: "unset does not resolve", value: "", wantHours: "", wantResolve: false},
		{name: "positive integer resolves", value: "3", wantHours: "3", wantResolve: true},
		{name: "surrounding whitespace trimmed", value: "  24 ", wantHours: "24", wantResolve: true},
		{name: "zero does not resolve", value: "0", wantHours: "", wantResolve: false},
		{name: "negative does not resolve", value: "-1", wantHours: "", wantResolve: false},
		{name: "non-numeric does not resolve", value: "hourly", wantHours: "", wantResolve: false},
		{name: "fractional does not resolve", value: "1.5", wantHours: "", wantResolve: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHours, gotResolve := resolveCadence(tt.value)
			if gotHours != tt.wantHours {
				t.Fatalf("resolveCadence(%q) hours = %q, want %q", tt.value, gotHours, tt.wantHours)
			}
			if gotResolve != tt.wantResolve {
				t.Fatalf("resolveCadence(%q) resolved = %v, want %v", tt.value, gotResolve, tt.wantResolve)
			}
		})
	}
}

func TestParseEcosystemList(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "empty", raw: "", want: ""},
		{name: "trims and sorts", raw: "npm; Go ;Maven", want: "Go;Maven;npm"},
		{name: "de-duplicates", raw: "Go;Go;PyPI;PyPI", want: "Go;PyPI"},
		{name: "drops blank entries", raw: "Go;;;Maven;", want: "Go;Maven"},
		{name: "keeps names with dots", raw: "crates.io;Go", want: "Go;crates.io"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.Join(parseEcosystemList(tt.raw), ";")
			if got != tt.want {
				t.Fatalf("parseEcosystemList(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestSameEcosystemSet(t *testing.T) {
	ptr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		stored  *string
		desired string
		want    bool
	}{
		{name: "nil stored", stored: nil, desired: "Go", want: false},
		{name: "same set different order", stored: ptr("Maven;Go"), desired: "Go;Maven", want: true},
		{name: "same set different spacing", stored: ptr(" Go ; Maven "), desired: "Go;Maven", want: true},
		{name: "subset is not equal", stored: ptr("Go"), desired: "Go;Maven", want: false},
		{name: "superset is not equal", stored: ptr("Go;Maven;npm"), desired: "Go;Maven", want: false},
		{name: "disjoint", stored: ptr("PyPI"), desired: "Go", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sameEcosystemSet(tt.stored, tt.desired); got != tt.want {
				t.Fatalf("sameEcosystemSet(%v, %q) = %v, want %v", tt.stored, tt.desired, got, tt.want)
			}
		})
	}
}
