package user

import "gorm.io/gorm"

type Repository interface {
	//kumpulan method kosong
	Save(user User) (User, error)

	//untuk cari email login
	FindByEmail(email string) (User, error)
}

//huruf kecil, karena sifatnya private
type repository struct {
	//define untuk akses ke DB
	db *gorm.DB
}

//untuk create new object
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}
