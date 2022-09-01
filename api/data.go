package api

import (
	"fmt"

	"golang.org/x/exp/maps"
)

const (
	numberOfUsers      int = 10
	numberOfUserAlbums int = 10
)

type data struct {
	Users  map[int]User
	Albums map[int]Album
}

type User struct {
	ID       string
	Name     string
	Username string
}

type Album struct {
	ID          string
	UserID      string
	Description string
}

func generateData() data {
	users := getUserData()
	return data{
		Users:  users,
		Albums: getAlbums(maps.Keys(users)),
	}
}

func getUserData() map[int]User {
	users := make(map[int]User, 0)
	for i := 0; i < numberOfUsers; i++ {
		iString := fmt.Sprint(i)
		users[i] = User{
			ID:       iString,
			Name:     "User " + iString,
			Username: "User" + iString,
		}
	}

	return users
}

func getAlbums(userIDs []int) map[int]Album {
	albums := make(map[int]Album, len(userIDs)*numberOfUserAlbums)

	for _, userID := range userIDs {
		userIDString := fmt.Sprint(userID)
		startIndex := userID * numberOfUserAlbums
		endIndex := (userID + 1) * numberOfUserAlbums
		for i := startIndex; i < endIndex; i++ {
			iString := fmt.Sprint(i)
			albums[i] = Album{
				ID:          iString,
				UserID:      userIDString,
				Description: "An Album",
			}
		}
	}

	return albums
}
