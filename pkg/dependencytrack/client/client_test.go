package client

import (
	"bytes"
	"io"
	"testing"
)

func TestAPIClientDecodeBytes(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02, 0xff}
	var decoded []byte

	err := (&APIClient{}).decode(&decoded, payload, "application/octet-stream")

	if err != nil {
		t.Fatalf("decode returned an error: %v", err)
	}
	if !bytes.Equal(decoded, payload) {
		t.Fatalf("decoded payload = %v, want %v", decoded, payload)
	}
}

func TestAPIClientDecodeReader(t *testing.T) {
	payload := []byte("binary response")
	var decoded io.Reader

	err := (&APIClient{}).decode(&decoded, payload, "application/octet-stream")

	if err != nil {
		t.Fatalf("decode returned an error: %v", err)
	}
	result, err := io.ReadAll(decoded)
	if err != nil {
		t.Fatalf("reading decoded payload: %v", err)
	}
	if !bytes.Equal(result, payload) {
		t.Fatalf("decoded payload = %q, want %q", result, payload)
	}
}
