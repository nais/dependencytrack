package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/nais/dependencytrack/pkg/dependencytrack"
)

func TestBootstrapTrivyIgnoreUnfixed(t *testing.T) {
	for _, tt := range []struct {
		name    string
		stored  string
		desired string
	}{
		{name: "disable ignoring unfixed", stored: "true", desired: "false"},
		{name: "enable ignoring unfixed", stored: "false", desired: "true"},
		{name: "already false", stored: "false", desired: "false"},
		{name: "already true", stored: "true", desired: "true"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			originalConfig, originalArgs := *cfg, os.Args
			t.Cleanup(func() {
				*cfg, os.Args = originalConfig, originalArgs
			})
			os.Args = []string{"bootstrap"}
			t.Chdir(t.TempDir())
			flag.VisitAll(func(f *flag.Flag) {
				if !strings.HasPrefix(f.Name, "test.") {
					t.Setenv(strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")), f.DefValue)
				}
			})
			if err := os.Unsetenv("TRIVY_IGNORE_UNFIXED"); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(".env", []byte("TRIVY_IGNORE_UNFIXED="+tt.desired+"\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile("users.yaml", []byte("users: []\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			t.Setenv("USERS_FILE", "users.yaml")
			t.Setenv("TRIVY_API_TOKEN", "test-token")

			updates := make(chan string, 2)
			token := "eyJhbGciOiJub25lIn0." + base64.RawURLEncoding.EncodeToString(
				fmt.Appendf(nil, `{"exp":%d}`, time.Now().Add(time.Hour).Unix()),
			) + "."
			mux := http.NewServeMux()
			mux.HandleFunc("POST /api/v1/user/login", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				_, _ = w.Write([]byte(token))
			})
			mux.HandleFunc("GET /api/v1/team", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`[{"name":"Administrators","uuid":"test-team","apiKeys":[{"key":"test-key"}]}]`))
			})
			mux.HandleFunc("GET /api/v1/configProperty", func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `[{"groupName":"integrations","propertyName":"trivy.ignore.unfixed","propertyValue":%q,"propertyType":"BOOLEAN"}]`, tt.stored)
			})
			mux.HandleFunc("POST /api/v1/configProperty", func(w http.ResponseWriter, r *http.Request) {
				var property dependencytrack.ConfigProperty
				if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if property.PropertyValue == nil {
					t.Error("missing property value")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				updates <- *property.PropertyValue
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{}`))
			})
			server := httptest.NewServer(mux)
			defer server.Close()
			t.Setenv("BASE_URL", server.URL+"/api")

			main()

			select {
			case got := <-updates:
				if tt.stored == tt.desired {
					t.Errorf("unexpected update to %q when value is already set", got)
				} else if got != tt.desired {
					t.Errorf("trivy.ignore.unfixed = %q, want %q", got, tt.desired)
				}
			default:
				if tt.stored != tt.desired {
					t.Errorf("trivy.ignore.unfixed was not updated from %q to %q", tt.stored, tt.desired)
				}
			}
		})
	}
}

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

func TestGithubAdvisoryMirrorCadence(t *testing.T) {
	c := &Config{GithubAdvisoryMirrorCadence: "1"}
	if got := c.taskSchedulerCadences()["ghsa.mirror.cadence"]; got != "1" {
		t.Fatalf("ghsa.mirror.cadence = %q, want 1", got)
	}
}

func TestVulnerabilityCadenceFlags(t *testing.T) {
	original := *cfg
	t.Cleanup(func() { *cfg = original })

	tests := []struct {
		name     string
		mirror   string
		analysis string
	}{
		{name: "custom intervals", mirror: "2", analysis: "12"},
		{name: "empty leaves unchanged", mirror: "", analysis: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
			for _, name := range []string{"osv-mirror-cadence", "portfolio-vulnerability-analysis-cadence"} {
				f := flag.Lookup(name)
				if f == nil {
					t.Fatalf("flag %q is not registered", name)
				}
				flags.Var(f.Value, f.Name, f.Usage)
			}
			if err := flags.Parse([]string{
				"--osv-mirror-cadence=" + tt.mirror,
				"--portfolio-vulnerability-analysis-cadence=" + tt.analysis,
			}); err != nil {
				t.Fatal(err)
			}
			cadences := cfg.taskSchedulerCadences()
			for property, want := range map[string]string{
				"osv.mirror.cadence":                       tt.mirror,
				"portfolio.vulnerability.analysis.cadence": tt.analysis,
			} {
				if got := cadences[property]; got != want {
					t.Errorf("%s = %q, want %q", property, got, want)
				}
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
