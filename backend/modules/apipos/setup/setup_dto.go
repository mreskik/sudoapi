package setup

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DataUser struct {
	Id       int    `bun:"id"`
	Username string `bun:"username"`
	Password string `bun:"password"`
}
