package users

type exportUsers struct {
	Users []exportUser
}

type exportUser struct {
	Username string
	Password string
	Id       float64
	Uuid     string
}

type UserDetails struct {
	firstName string `faker:"first_name"`
	lastName  string `faker:"last_name"`
	password  string `faker:"password"`
	email     string `faker:"email"`
}
