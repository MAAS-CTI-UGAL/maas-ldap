package config

import "testing"

func TestLoadListenAddressDefaultsWhenPortUnset(t *testing.T) {
	t.Setenv(envPort, "")

	if got := loadListenAddress(); got != defaultListenAddress {
		t.Fatalf("loadListenAddress() = %q, want %q", got, defaultListenAddress)
	}
}

func TestLoadListenAddressPrefixesBarePort(t *testing.T) {
	t.Setenv(envPort, "9091")

	if got := loadListenAddress(); got != ":9091" {
		t.Fatalf("loadListenAddress() = %q, want %q", got, ":9091")
	}
}

func TestLoadListenAddressPreservesHostPort(t *testing.T) {
	t.Setenv(envPort, "127.0.0.1:9091")

	if got := loadListenAddress(); got != "127.0.0.1:9091" {
		t.Fatalf("loadListenAddress() = %q, want %q", got, "127.0.0.1:9091")
	}
}

func TestLoadListenAddressPreservesColonPort(t *testing.T) {
	t.Setenv(envPort, ":9091")

	if got := loadListenAddress(); got != ":9091" {
		t.Fatalf("loadListenAddress() = %q, want %q", got, ":9091")
	}
}
