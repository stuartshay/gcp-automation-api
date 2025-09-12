package handlers_test

import (
	"testing"

	"github.com/stuartshay/gcp-automation-api/pkg/utils"
)

func TestValidateProjectID(t *testing.T) {
	testCases := []struct {
		projectID string
		valid     bool
	}{
		{"my-project-123", true},
		{"valid-project", true},
		{"project123", true},
		{"My-Project", false},        // uppercase
		{"project_123", false},       // underscore
		{"123project", false},        // starts with number
		{"a", false},                 // too short
		{"project-", false},          // ends with hyphen
		{"very-long-project-name-that-exceeds-limit", false}, // too long
	}

	for _, tc := range testCases {
		err := utils.ValidateProjectID(tc.projectID)
		if tc.valid && err != nil {
			t.Errorf("Expected %s to be valid, got error: %v", tc.projectID, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Expected %s to be invalid, but got no error", tc.projectID)
		}
	}
}

func TestValidateBucketName(t *testing.T) {
	testCases := []struct {
		bucketName string
		valid      bool
	}{
		{"my-bucket-123", true},
		{"valid-bucket", true},
		{"bucket123", true},
		{"My-Bucket", false},         // uppercase
		{"bucket_123", false},        // underscore
		{"a", false},                 // too short
		{"bucket-", false},           // ends with hyphen
		{"very-long-bucket-name-that-exceeds-the-maximum-allowed-limit-for-sure", false}, // too long
	}

	for _, tc := range testCases {
		err := utils.ValidateBucketName(tc.bucketName)
		if tc.valid && err != nil {
			t.Errorf("Expected %s to be valid, got error: %v", tc.bucketName, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("Expected %s to be invalid, but got no error", tc.bucketName)
		}
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	
	if !utils.ContainsString(slice, "apple") {
		t.Error("Expected to find 'apple' in slice")
	}
	
	if utils.ContainsString(slice, "orange") {
		t.Error("Did not expect to find 'orange' in slice")
	}
}

func TestRemoveString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	result := utils.RemoveString(slice, "banana")
	
	if len(result) != 2 {
		t.Errorf("Expected length 2, got %d", len(result))
	}
	
	if utils.ContainsString(result, "banana") {
		t.Error("Expected 'banana' to be removed from slice")
	}
}