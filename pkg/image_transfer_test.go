package pkg

import "testing"

func TestConvertImagePathForECR(t *testing.T) {
	expected := "123456789012.dkr.ecr.ap-northeast1.amazonaws.com/nginx:latest"
	actual := ConvertImagePathForECR("nginx:latest", "ap-northeast1", "123456789012")

	if expected != actual {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}
