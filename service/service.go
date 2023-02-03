package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/tg"
)

func SaveUsers(ctx context.Context, sender *message.Sender, client *tg.Client, channelName string) {
	queryB, err := sender.Resolve(fmt.Sprintf("https://t.me/%s", channelName)).AsInputChannel(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	participants := query.GetParticipants(client, queryB)

	size, err := participants.Count(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Paticipants: ", size)
	iter := query.GetParticipants(client, queryB).BatchSize(1000).Iter()

	n := 0
	for iter.Next(ctx) {
		n++
		user, _ := iter.Value().User()
		fmt.Println(n, user.ID, user.Username, user.Phone)

		if user.Photo == nil {
			continue
		}
		photo, ok := user.Photo.AsNotEmpty()
		if !ok {
			continue
		}
		fileName := fmt.Sprintf("%d_%s", user.ID, user.Username)
		savePhoto(ctx, client, photo, user.ID, fileName, user.AccessHash, channelName)
	}

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func savePhoto(
	ctx context.Context,
	client *tg.Client,
	photo *tg.UserProfilePhoto,
	userID int64,
	fileName string,
	accessHash int64,
	groupName string) {

	fileLocation := &tg.InputPeerPhotoFileLocation{
		Big: true,
		Peer: &tg.InputPeerUser{
			UserID:     userID,
			AccessHash: accessHash,
		},
		PhotoID: photo.PhotoID,
	}

	fileRquest := tg.UploadGetFileRequest{
		Location: fileLocation,
		Offset:   0,
		Limit:    1024 * 1024,
	}

	final, err := client.UploadGetFile(ctx, &fileRquest)
	if err != nil {
		fmt.Println("Error upload", err)
	}
	switch result := final.(type) {
	case *tg.UploadFile:
		saveFile(result.Bytes, groupName, fileName)
	default:
	}
}

func saveFile(imgByte []byte, folder string, fileName string) {

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}
	path := fmt.Sprintf("%s/%s", "exportImgs", folder)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	out, err := os.Create(fmt.Sprintf("%s/%s%s", path, fileName, ".jpeg"))
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 100

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.Println(err)
	}

}
