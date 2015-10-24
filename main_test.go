package main

import (
	"testing"
)

func Test_GetLightColorCode(t *testing.T) {
	lightColorCode := getLightColorCode("blue")
	if lightColorCode != "1" {
		t.Error("Expected", "1", "got", lightColorCode)
	}

	lightColorCode2 := getLightColorCode("blue_anime")
	if lightColorCode2 != "0" {
		t.Error("Expected", "0", "got", lightColorCode2)
	}

	lightColorCode3 := getLightColorCode("blue_anime")
	if lightColorCode3 != "0" {
		t.Error("Expected", "0", "got", lightColorCode3)
	}

	lightColorCode4 := getLightColorCode("blue_anime")
	if lightColorCode4 != "0" {
		t.Error("Expected", "0", "got", lightColorCode4)
	}
}
