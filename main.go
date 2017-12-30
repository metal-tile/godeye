package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/metal-tile/godeye/firedb"
	"github.com/sinmetal/slog"
)

// ProjectID is GCP Project ID
const ProjectID = "metal-tile-dev1"

// PlayerPositionManager is PlayerのPositionを管理するもの
type PlayerPositionManager struct {
	Map *sync.Map
}

// NewPlayerPositionManager is PlayerPositionManagerを作成
func NewPlayerPositionManager() PlayerPositionManager {
	return PlayerPositionManager{
		Map: &sync.Map{},
	}
}

// Load is Load from PlayerPosition
func (manager *PlayerPositionManager) Load(key string) (playerPosition *PlayerPosition, ok bool) {
	v, ok := manager.Map.Load(key)
	if ok {
		return v.(*PlayerPosition), ok
	}
	return nil, ok
}

// Store is Store to PlayerPosition
func (manager *PlayerPositionManager) Store(key string, playerPosition *PlayerPosition) {
	manager.Map.Store(key, playerPosition)
}

// SetPlayerPositionMap is Update PlayerPositionMap
// Firestoreから取得したPlayerPositionを保持する
// 状態が変わっているものはUpdatedAtを更新して保存する
// TODO 減ってるプレイヤーがある場合の考慮が必要か？
func (manager *PlayerPositionManager) SetPlayerPositionMap(pps []*firedb.PlayerPosition) {
	for _, v := range pps {
		cpp, ok := manager.Load(v.ID)
		if ok {
			if cpp.X == v.X && cpp.Y == v.Y {
				continue
			}
			cpp.Angle = v.Angle
			cpp.IsMove = v.IsMove
			cpp.X = v.X
			cpp.Y = v.Y
			cpp.UpdatedAt = Now()
			manager.Store(v.ID, cpp)
		} else {
			manager.Store(v.ID, &PlayerPosition{
				ID:        v.ID,
				Angle:     v.Angle,
				IsMove:    v.IsMove,
				X:         v.X,
				Y:         v.Y,
				UpdatedAt: Now(),
			})
		}
	}

}

// ExistActivePlayer is ActiveなPlayerが存在するかのチェック
func (manager *PlayerPositionManager) ExistActivePlayer() bool {
	exist := false
	manager.Map.Range(func(key interface{}, value interface{}) bool {
		if value.(*PlayerPosition).UpdatedAt.After(Now().Add(time.Minute * -18)) {
			exist = true
			return false
		}
		return true
	})

	return exist
}

// PlayerPosition is Firestoreの/world-{worldname}-player-positionのデータを入れるstruct
type PlayerPosition struct {
	ID        string    `firestore:"-" json:"id"`
	Angle     float64   `json:"angle"`
	IsMove    bool      `json:"isMove"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func main() {
	ppm := NewPlayerPositionManager()
	ch := make(chan error)
	go func() {
		ch <- watchActivePlayer(ppm)
	}()

	err := <-ch
	fmt.Println(err.Error())
}

func watchActivePlayer(manager PlayerPositionManager) error {
	existActivePlayer := false
	playerStore := firedb.NewPlayerStore()
	for {
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-t.C:
				log := slog.Start(time.Now())
				ctx := context.Background()
				pps, err := playerStore.GetPlayerPositions(ctx)
				if err != nil {
					log.Errorf("playerStore.GetPlayerPositions. %s", err.Error())
					log.Flush()
					continue
				}

				manager.SetPlayerPositionMap(pps)
				neap := manager.ExistActivePlayer()
				if existActivePlayer != neap {
					var replicas int32
					if neap {
						replicas = 1
					}

					log.Infof("try update land replica size = %d", replicas)
					err := updateReplicas("default", "land", replicas)
					if err != nil {
						log.Errorf("failed update land replica size. %+v", err)
						log.Flush()
						continue
					}
					existActivePlayer = neap
				}

				// debug log
				j, err := json.Marshal(pps)
				if err != nil {
					log.Errorf("json.Marshal. %s", err.Error())
					log.Flush()
				}
				log.Infof(string(j))
				log.Flush()
			}
		}
	}
}
