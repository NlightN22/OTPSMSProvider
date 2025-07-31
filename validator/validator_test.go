package validator

import "testing"

func TestE164Regexp(t *testing.T) {
	valid := []string{
		"+1234567890",
		"+19876543210",
	}
	invalid := []string{
		"123456",
		"+0123",
		"+123abc",
		"",
	}
	for _, s := range valid {
		if !e164Regexp.MatchString(s) {
			t.Errorf("e164Regexp should match %q", s)
		}
	}
	for _, s := range invalid {
		if e164Regexp.MatchString(s) {
			t.Errorf("e164Regexp should not match %q", s)
		}
	}
}
