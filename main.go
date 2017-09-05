package main

import (
	"context"
	"fmt"
	"log"

	"gocloud-exp/firestore"
)

func main() {
	// Get Firestore project ID from environment.
	projectID := "sinmetal-firestore"

	// Get a Firestore client authenticated with
	// Application Default Credentials.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Error getting client: %v", err)
	}

	// Close client when done.
	defer client.Close()

	// Get the 'NYC' doc from the 'cities' collection
	docsnap, err := client.Collection("cities").Doc("SF").Get(ctx)
	if err != nil {
		log.Fatalf("Error getting document: %v", err)
	}

	// Print document data.
	fmt.Printf("Got document: %v\n", docsnap.Data())
}
