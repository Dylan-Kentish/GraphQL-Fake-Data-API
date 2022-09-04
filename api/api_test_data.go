package api

import (
	"fmt"

	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/data"
	"github.com/Dylan-Kentish/GraphQLFakeDataAPI/utils"
	"golang.org/x/exp/maps"
)

const (
	numberOfUsers       int = 10
	numberOfUserAlbums  int = 10
	numberOfAlbumPhotos int = 10
)

type testData struct {
	users  map[int]data.User
	albums map[int]data.Album
	photos map[int]data.Photo
}

// Public so that main.go can access it.
// This should be private once a 'real' data source has been added.
func NewTestData() data.IData {
	users := getUserData()
	albums := getAlbums(maps.Keys(users))
	photos := getPhotos(maps.Keys(albums))
	return &testData{
		users:  users,
		albums: albums,
		photos: photos,
	}
}

func (testData *testData) GetUsers() []data.User {
	return utils.OrderedValues(testData.users)
}

func (testData *testData) GetAlbums() []data.Album {
	return utils.OrderedValues(testData.albums)
}

func (testData *testData) GetPhotos() []data.Photo {
	return utils.OrderedValues(testData.photos)
}

func (testData *testData) GetUser(id int) data.User {
	if user, ok := testData.users[id]; ok {
		return user
	}
	return data.User{}
}

func (testData *testData) GetAlbum(id int) data.Album {
	if album, ok := testData.albums[id]; ok {
		return album
	}
	return data.Album{}
}

func (testData *testData) GetPhoto(id int) data.Photo {
	if photo, ok := testData.photos[id]; ok {
		return photo
	}
	return data.Photo{}
}

func (testData *testData) GetAlbumsByUserID(userID int) []data.Album {
	return utils.ValuesWhere(testData.albums, func(album data.Album) bool {
		return album.UserID == userID
	})
}

func (testData *testData) GetPhotosByAlbumID(albumID int) []data.Photo {
	return utils.ValuesWhere(testData.photos, func(photo data.Photo) bool {
		return photo.AlbumID == albumID
	})
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
