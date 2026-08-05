package domain

import "testing"

func TestValidateContainer(t *testing.T) {
	t.Parallel()
	for _, container := range []Container{ContainerMKV, ContainerMP4} {
		if err := ValidateContainer(container); err != nil {
			t.Errorf("ValidateContainer(%q): %v", container, err)
		}
	}
	if err := ValidateContainer("avi"); err == nil {
		t.Fatal("unsupported container accepted")
	}
	NopReporter{}.Report(Event{Kind: EventComplete})
}
