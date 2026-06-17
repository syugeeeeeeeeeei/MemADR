package release

import "testing"

func TestValidateVersion(t *testing.T) {
	t.Parallel()

	valid := []string{"v0.1.0", "v1.2.3-rc.1", "v1.2.3+meta"}
	for _, item := range valid {
		if err := ValidateVersion(item); err != nil {
			t.Fatalf("ValidateVersion(%q) error: %v", item, err)
		}
	}

	invalid := []string{"0.1.0", "v1", "v1.2", "vx.y.z"}
	for _, item := range invalid {
		if err := ValidateVersion(item); err == nil {
			t.Fatalf("ValidateVersion(%q) succeeded", item)
		}
	}
}

func TestDevVersionUsesTaggedBase(t *testing.T) {
	t.Parallel()

	got := DevVersion("v1.2.3", "abcdef1234567890")
	if got != "v1.2.3-dev+abcdef1" {
		t.Fatalf("DevVersion = %q", got)
	}
}

func TestDevVersionFallsBackToZeroBase(t *testing.T) {
	t.Parallel()

	got := DevVersion("", "abcdef1234567890")
	if got != "v0.0.0-dev+abcdef1" {
		t.Fatalf("DevVersion = %q", got)
	}
}
