package user

import "time"

// buat ORM
// untuk structnya, harus sesuai dengan nama DB, namun di singular kan
// pastikan penamaannya benar (singular dan pluralnya) untuk memudahkan koneksi databasenya

type User struct {
	ID             int
	Name           string
	Occupation     string
	Email          string
	PasswordHash   string
	AvatarFileName string
	Role           string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
