package structs

type UserRole string

const (
	ADMIN  UserRole = "admin"
	SELLER UserRole = "seller"
	BUYER  UserRole = "buyer"
)

type UserRet struct {
	ID       int
	Username string
	Name     string
	Role     UserRole
}

type UserCred struct {
	Username string
	Name     string
	Role     UserRole
	Password string
}
