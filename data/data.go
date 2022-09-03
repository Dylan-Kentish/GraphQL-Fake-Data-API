package data

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

type Data struct {
	Users  map[int]User
	Albums map[int]Album
	Photos map[int]Photo
}

func NewData(
	users map[int]User,
	albums map[int]Album,
	photos map[int]Photo) *Data {
	return &Data{
		Users:  users,
		Albums: albums,
		Photos: photos,
	}
}

func (data *Data) GetUser(id int) User {
	if user, ok := data.Users[id]; ok {
		return user
	}
	return User{}
}

func (data *Data) GetAlbum(id int) Album {
	if album, ok := data.Albums[id]; ok {
		return album
	}
	return Album{}
}

func (data *Data) GetPhoto(id int) Photo {
	if photo, ok := data.Photos[id]; ok {
		return photo
	}
	return Photo{}
}

func (data *Data) GetAlbumsByUserID(userID int) []Album {
	albums := make([]Album, 0)

	for _, album := range data.Albums {
		if album.UserID == userID {
			albums = append(albums, album)
		}
	}

	return albums
}

func (data *Data) GetPhotosByAlbumID(albumID int) []Photo {
	photos := make([]Photo, 0)

	for _, photo := range data.Photos {
		if photo.AlbumID == albumID {
			photos = append(photos, photo)
		}
	}

	return photos
}
