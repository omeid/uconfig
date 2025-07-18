package f

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalerStringSlice(t *testing.T) {
	expect := TextUnmarshalerStringSlice{"a", "b", "c"}
	value := TextUnmarshalerStringSlice{}

	err := value.UnmarshalText([]byte("a.b.c"))
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(expect, value); diff != "" {
		t.Error(diff)
	}
}
