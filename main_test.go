package main

import (
	"testing"
	"time"
)

func TestPlayerPositionManager_existActivePlayer(t *testing.T) {
	pm := map[string]PlayerPosition{}
	pm["hoge"] = PlayerPosition{
		UpdatedAt: time.Now(),
	}

	ppm := PlayerPositionManager{}
	ppm.Map = pm

	if ppm.existActivePlayer() == false {
		t.Fatalf("unexpected existActivePlayer.")
	}
}
