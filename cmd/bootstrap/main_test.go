package main

import "testing"

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
