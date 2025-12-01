package be

import (
	"testing"
)

func Equal[e comparable](t testing.TB, actual, expected e) {
	if actual != expected {
		t.Fatalf("not equal: %+v != %+v", actual, expected)
	}
}

func True(t testing.TB, actual bool) {
	if !actual {
		t.Fatalf("not true: %+v", actual)
	}
}

func NoError(t testing.TB, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %+v", err)
	}
}

func AnError(t testing.TB, err error) {
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
