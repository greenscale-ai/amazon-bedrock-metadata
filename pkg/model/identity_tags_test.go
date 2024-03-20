package model

import (
	"testing"
)

func TestParseIamEntity(t *testing.T) {
	identityTags := NewIdentityTagsBuilder(nil)
	iamEntityType, entityName := identityTags.parseIamEntity("arn:aws:iam::123456789012:role/example-role")
	if iamEntityType != "role" {
		t.Errorf("got %q, wanted %q", iamEntityType, "role")
	}

	if entityName != "example-role" {
		t.Errorf("got %q, wanted %q", entityName, "example-role")
	}
}
