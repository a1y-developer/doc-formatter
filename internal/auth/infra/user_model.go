package infra

type UserModel struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email      string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	IsVerified bool   `gorm:"default:false"`
}
