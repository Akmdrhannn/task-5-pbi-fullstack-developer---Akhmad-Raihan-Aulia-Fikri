package app

type UserData struct {
	Username string `json:"username" valid:"required, stringlength(6|50)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required, stringlength(6|50)"`
}
