package entity

type Todo struct {
	ID     uint
	Text   string
	Done   bool
	UserID uint
	User   User
}
