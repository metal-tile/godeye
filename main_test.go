package main

import (
	"testing"
	"time"

	"github.com/metal-tile/godeye/firedb"
)

func TestPlayerPositionManager_ExistActivePlayer(t *testing.T) {
	manager := NewPlayerPositionManager()

	cases := []struct {
		in   time.Time
		want bool
	}{
		{time.Now(), true},
		{time.Now().Add(time.Minute * -17), true},
		{time.Now().Add(time.Minute * -19), false},
	}

	for i, v := range cases {
		manager.Store("A", &PlayerPosition{
			UpdatedAt: v.in,
		})

		if manager.ExistActivePlayer() == v.want {
			t.Fatalf("unexpected existActivePlayer. index = %d", i)
		}
	}
}

func TestPlayerPositionManager_SetPlayerPositionMap(t *testing.T) {
	manager := NewPlayerPositionManager()

	count := -1
	now1 := time.Date(2017, 12, 30, 10, 0, 0, 0, time.UTC)
	now2 := time.Date(2017, 12, 31, 10, 0, 0, 0, time.UTC)
	nowArray := []time.Time{
		now1, now1, now2,
	}
	nowFunc = func() time.Time {
		count++
		return nowArray[count]
	}

	casess := [][]struct {
		in   *firedb.PlayerPosition
		want PlayerPosition
	}{}
	{
		cases := []struct {
			in   *firedb.PlayerPosition
			want PlayerPosition
		}{
			{
				&firedb.PlayerPosition{
					ID:     "A",
					Angle:  180,
					IsMove: true,
					X:      100.0,
					Y:      100.0,
				},
				PlayerPosition{
					ID:        "A",
					Angle:     180,
					IsMove:    true,
					X:         100.0,
					Y:         100.0,
					UpdatedAt: now1,
				},
			},
			{
				&firedb.PlayerPosition{
					ID:     "B",
					Angle:  180,
					IsMove: true,
					X:      100.0,
					Y:      100.0,
				},
				PlayerPosition{
					ID:        "B",
					Angle:     180,
					IsMove:    true,
					X:         100.0,
					Y:         100.0,
					UpdatedAt: now1,
				},
			},
		}
		casess = append(casess, cases)
	}
	{
		cases := []struct {
			in   *firedb.PlayerPosition
			want PlayerPosition
		}{
			{
				&firedb.PlayerPosition{
					ID:     "A",
					Angle:  180,
					IsMove: true,
					X:      100.0,
					Y:      100.0,
				},
				PlayerPosition{
					ID:        "A",
					Angle:     180,
					IsMove:    true,
					X:         100.0,
					Y:         100.0,
					UpdatedAt: now1,
				},
			},
			{
				&firedb.PlayerPosition{
					ID:     "B",
					Angle:  180,
					IsMove: true,
					X:      200.0,
					Y:      200.0,
				},
				PlayerPosition{
					ID:        "B",
					Angle:     180,
					IsMove:    true,
					X:         200.0,
					Y:         200.0,
					UpdatedAt: now2, // XYの値が異なるのでUpdatedAtが更新される
				},
			},
		}
		casess = append(casess, cases)
	}

	for i, cases := range casess {
		for j, ca := range cases {
			manager.SetPlayerPositionMap([]*firedb.PlayerPosition{ca.in})
			v, ok := manager.Load(ca.want.ID)
			if ok == false {
				t.Fatalf("expected map exists key = %s", ca.want.ID)
			}
			if e, g := ca.want.ID, v.ID; e != g {
				t.Fatalf("expected ID %s; got %s", e, g)
			}
			if e, g := ca.want.Angle, v.Angle; e != g {
				t.Fatalf("expected Angle %g; got %g", e, g)
			}
			if e, g := ca.want.IsMove, v.IsMove; e != g {
				t.Fatalf("expected IsMove %t; got %t", e, g)
			}
			if e, g := ca.want.X, v.X; e != g {
				t.Fatalf("expected X %g; got %g", e, g)
			}
			if e, g := ca.want.Y, v.Y; e != g {
				t.Fatalf("expected Y %g; got %g", e, g)
			}
			if ca.want.UpdatedAt.Equal(v.UpdatedAt) == false {
				t.Fatalf("unexpected %d:%d UpdatedAt", i, j)
			}
		}
	}
}
