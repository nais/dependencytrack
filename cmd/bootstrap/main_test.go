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
