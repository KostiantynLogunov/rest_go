package user

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"` //bson for mongo
	Email        string `json:"email" bson:"email"`
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
}

type CreateUserDto struct {
	Username string `json:"username"` //CreateUserDto has NOTHING with db that why we don't have bson here
	Password string `json:"password"`
	Email    string `json:"email"`
}
