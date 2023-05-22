package auth

type UserAuthInfo struct {
	User  User
	Token string
}

type User struct {
	ID       uint
	Email    string
	Password string
	Name     string
	Role     int
	IsVerify bool
}
