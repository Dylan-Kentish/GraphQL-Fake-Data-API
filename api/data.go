package api

import (
	"fmt"

	"golang.org/x/exp/maps"
)

const (
	numberOfUsers       int = 10
	numberOfUserAlbums  int = 10
	numberOfAlbumPhotos int = 10
)

type data struct {
	Users  map[int]User
	Albums map[int]Album
	Photos map[int]Photo
}

type User struct {
	ID       int
	Name     string
	Username string
	Albums   []Album
}

type Album struct {
	ID          int
	UserID      int
	Description string
	Photos      []Photo
}

type Photo struct {
	ID          int
	AlbumID     int
	Description string
}

func NewData() *data {
	users := getUserData()
	albums := getAlbums(maps.Keys(users))
	return &data{
		Users:  users,
		Albums: albums,
		Photos: getPhotos(maps.Keys(albums)),
	}
}

func (data *data) getUser(id int) User {
	if user, ok := data.Users[id]; ok {
		return user
	}
	return User{}
}

func (data *data) getAlbum(id int) Album {
	if album, ok := data.Albums[id]; ok {
		return album
	}
	return Album{}
}

func (data *data) getAlbumsByUserID(userID int) []Album {
	albums := make([]Album, 0)

	for _, album := range data.Albums {
		if album.UserID == userID {
			albums = append(albums, album)
		}
	}

	return albums
}

func (data *data) getPhotosByAlbumID(albumID int) []Photo {
	photos := make([]Photo, 0)

	for _, photo := range data.Photos {
		if photo.AlbumID == albumID {
			photos = append(photos, photo)
		}
	}

	return photos
}

func getUserData() map[int]User {
	users := make(map[int]User, 0)
	for i := 0; i < numberOfUsers; i++ {
		iString := fmt.Sprint(i)
		users[i] = User{
			ID:       i,
			Name:     "User " + iString,
			Username: "User" + iString,
		}
	}

	return users
}

func getAlbums(userIDs []int) map[int]Album {
	albums := make(map[int]Album, len(userIDs)*numberOfUserAlbums)

	for _, userID := range userIDs {
		startIndex := userID * numberOfUserAlbums
		endIndex := (userID + 1) * numberOfUserAlbums
		for i := startIndex; i < endIndex; i++ {
			iString := fmt.Sprint(i)
			albums[i] = Album{
				ID:          i,
				UserID:      userID,
				Description: "Album " + iString,
			}
		}
	}

	return albums
}

func getPhotos(albumIDs []int) map[int]Photo {
	photos := make(map[int]Photo, len(albumIDs)*numberOfAlbumPhotos)

	for _, albumID := range albumIDs {
		startIndex := albumID * numberOfAlbumPhotos
		endIndex := (albumID + 1) * numberOfAlbumPhotos
		for i := startIndex; i < endIndex; i++ {
			iString := fmt.Sprint(i)
			photos[i] = Photo{
				ID:          i,
				AlbumID:     albumID,
				Description: "Photo " + iString,
			}
		}
	}

	return photos
}
