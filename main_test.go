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
	if lightColorCode2 != "10" {
		t.Error("Expected", "10", "got", lightColorCode2)
	}

	lightColorCode3 := getLightColorCode("yellow")
	if lightColorCode3 != "2" {
		t.Error("Expected", "2", "got", lightColorCode3)
	}

	lightColorCode4 := getLightColorCode("yellow_anime")
	if lightColorCode4 != "20" {
		t.Error("Expected", "20", "got", lightColorCode4)
	}
}
