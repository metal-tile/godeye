package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/sinmetal/slog"
	"google.golang.org/api/iterator"
)

// ProjectID is GCP Project ID
const ProjectID = "metal-tile-dev1"

// PlayerPositionManager is PlayerのPositionを管理するもの
type PlayerPositionManager struct {
	Map map[string]PlayerPosition
}

// PlayerPosition is Firestoreの/world-{worldname}-player-positionのデータを入れるstruct
type PlayerPosition struct {
	ID        string
	X         float64
	Y         float64
	UpdatedAt time.Time
}

func (ppm *PlayerPositionManager) existActivePlayer() bool {
	for _, v := range ppm.Map {
		if v.UpdatedAt.After(time.Now().Add(time.Minute * -18)) {
			return true
		}
	}
	return false
}

func main() {
	crs := 0
	for {
		t := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-t.C:
				log := slog.Start(time.Now())
				go func(log *slog.Log) {
					ppm := PlayerPositionManager{
						Map: map[string]PlayerPosition{},
					}

					p, err := getPlayerPositions(log, ProjectID)
					if err != nil {
						log.Errorf("failed getPlayerPositions from Firestore. err = %s", err.Error())
					}
					ppm.Map = p
					b := ppm.existActivePlayer()
					nrs := 0
					if b {
						nrs = 1
					}
					if crs != nrs {
						log.Infof("Update land-node Replicas %d -> %d", crs, nrs)
						err := updateReplicas("default", "land-node", int32(nrs))
						if err != nil {
							log.Errorf("error:%+v\n", err)
							return
						}
						crs = nrs
					}
				}(&log)
				log.Flush()
			}
		}
	}
}

// getPlayerPositions is FirestoreからPlayerPositionを取得する
// TODO Firestore周りは後で別パッケージに分けて、Mockを用意したほうがいいかも
func getPlayerPositions(log *slog.Log, projectID string) (map[string]PlayerPosition, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ppm := map[string]PlayerPosition{}
	// FIXME world name
	iter := client.Collection("world-default-player-position").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var pp PlayerPosition
		err = doc.DataTo(&pp)
		if err != nil {
			return nil, err
		}
		pp.ID = doc.Ref.ID
		ppm[pp.ID] = pp
		// FIXME UpdatedAt を更新処理がないので、必ず人がいるとみなされる
	}

	return ppm, nil
}
