package api

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

type testData struct {
	Users  map[int]data.User
	Albums map[int]data.Album
	Photos map[int]data.Photo
}

// Public so that main.go can access it.
// This should be private once a 'real' data source has been added.
func NewTestData() *testData {
	users := getUserData()
	albums := getAlbums(maps.Keys(users))
	photos := getPhotos(maps.Keys(albums))
	return &testData{
		Users:  users,
		Albums: albums,
		Photos: photos,
	}
}

func (testData *testData) GetUsers() []data.User {
	return maps.Values(testData.Users)
}

func (testData *testData) GetAlbums() []data.Album {
	return maps.Values(testData.Albums)
}

func (testData *testData) GetPhotos() []data.Photo {
	return maps.Values(testData.Photos)
}

func (testData *testData) GetUser(id int) data.User {
	if user, ok := testData.Users[id]; ok {
		return user
	}
	return data.User{}
}

func (testData *testData) GetAlbum(id int) data.Album {
	if album, ok := testData.Albums[id]; ok {
		return album
	}
	return data.Album{}
}

func (testData *testData) GetPhoto(id int) data.Photo {
	if photo, ok := testData.Photos[id]; ok {
		return photo
	}
	return data.Photo{}
}

func (testData *testData) GetAlbumsByUserID(userID int) []data.Album {
	albums := make([]data.Album, 0)

	for _, album := range testData.Albums {
		if album.UserID == userID {
			albums = append(albums, album)
		}
	}

	return albums
}

func (testData *testData) GetPhotosByAlbumID(albumID int) []data.Photo {
	photos := make([]data.Photo, 0)

	for _, photo := range testData.Photos {
		if photo.AlbumID == albumID {
			photos = append(photos, photo)
		}
	}

	return photos
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
