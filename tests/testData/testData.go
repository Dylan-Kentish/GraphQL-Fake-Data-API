package testData

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"golang.org/x/exp/maps"
)

const (
	numberOfUsers       int = 10
	numberOfUserAlbums  int = 10
	numberOfAlbumPhotos int = 10
)

func NewTestData() *data.Data {
	users := getUserData()
	albums := getAlbums(maps.Keys(users))
	photos := getPhotos(maps.Keys(albums))

	return data.NewData(users, albums, photos)
}

func getUserData() map[int]data.User {
	users := make(map[int]data.User, 0)
	for i := 0; i < numberOfUsers; i++ {
		iString := fmt.Sprint(i)
		users[i] = data.User{
			ID:       i,
			Name:     "User " + iString,
			Username: "User" + iString,
		}
	}

	return users
}

func getAlbums(userIDs []int) map[int]data.Album {
	albums := make(map[int]data.Album, len(userIDs)*numberOfUserAlbums)

	for _, userID := range userIDs {
		startIndex := userID * numberOfUserAlbums
		endIndex := (userID + 1) * numberOfUserAlbums
		for i := startIndex; i < endIndex; i++ {
			iString := fmt.Sprint(i)
			albums[i] = data.Album{
				ID:          i,
				UserID:      userID,
				Description: "Album " + iString,
			}
		}
	}

	return albums
}

func getPhotos(albumIDs []int) map[int]data.Photo {
	photos := make(map[int]data.Photo, len(albumIDs)*numberOfAlbumPhotos)

	for _, albumID := range albumIDs {
		startIndex := albumID * numberOfAlbumPhotos
		endIndex := (albumID + 1) * numberOfAlbumPhotos
		for i := startIndex; i < endIndex; i++ {
			iString := fmt.Sprint(i)
			photos[i] = data.Photo{
				ID:          i,
				AlbumID:     albumID,
				Description: "Photo " + iString,
			}
		}
	}

	return photos
}
