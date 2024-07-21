package main

import (
	"encoding/base64"
	"io"
	"log"
	"os"

	"github.com/dghubble/oauth1"

	"github.com/automatis-io/go-twitter/twitter"
)

func main() {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	// A single pixel PNG image
	base64Image := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQAAAAA3bvkkAAAACklEQVR4AWNgAAAAAgABc3UBGAAAAABJRU5ErkJggg=="
	// Decode to []byte
	image, _ := base64.StdEncoding.DecodeString(base64Image)

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// v1 API - Upload media
	clientV1 := twitter.NewClient(httpClient)
	media, respUpload, err := clientV1.Media.Upload(
		image,
		twitter.UploadParams{
			MediaType:        twitter.UploadMediaTypeTweetImage,
			AdditionalOwners: []string{"3805104374"}, // corresponds to @furni test account
		},
	)
	defer respUpload.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Media ID:", media.MediaID)

	mediaStatus, respStatus, err := clientV1.Media.Status(media.MediaID)
	defer respStatus.Body.Close()
	if err != nil {
		respBody, _ := io.ReadAll(respStatus.Body)
		log.Printf("Error getting media status: %v", respBody)
		log.Fatal(err)
	}

	log.Println("Media Status:", mediaStatus.ProcessingInfo.State)
}
