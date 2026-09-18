package cli

import (
	"reflect"
	"testing"
)

func TestNormalizeArgs(t *testing.T) {
	for _, args := range [][]string{{"launch", "--device", "USB", "--preset", "Très léger", "--yes"}, {"--device", "USB", "launch", "--preset", "Très léger", "--yes"}} {
		want := []string{"--device", "USB", "--preset", "Très léger", "--yes", "launch"}
		if got := normalizeArgs(args); !reflect.DeepEqual(got, want) {
			t.Fatal(got)
		}
	}
	got := normalizeArgs([]string{"preview", "--args", "launch", "--preset=Qualité"})
	want := []string{"--args", "launch", "--preset=Qualité", "preview"}
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
}

func TestNormalizePreservesTerminator(t *testing.T) {
	got := normalizeArgs([]string{"preview", "--", "--yes"})
	want := []string{"--", "preview", "--yes"}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("option after terminator was activated", got)
	}
}
