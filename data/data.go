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

type IData interface {
	GetUsers() []User
	GetAlbums() []Album
	GetPhotos() []Photo

	GetUser(id int) User
	GetAlbum(id int) Album
	GetPhoto(id int) Photo

	GetAlbumsByUserID(userID int) []Album
	GetPhotosByAlbumID(albumID int) []Photo
}
