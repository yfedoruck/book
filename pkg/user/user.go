package user

type Data struct {
	Id       int
	Name     string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}

func New() *Data {
	return new(Data)
}
