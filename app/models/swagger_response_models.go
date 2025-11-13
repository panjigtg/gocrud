package models

// UsersListOK mendeskripsikan payload sukses untuk GET /users
type UsersListOK struct {
	Success bool        	`json:"success" example:"true"`
	Message string      	`json:"message" example:"succes"`
	Data    UserResponse 	`json:"data"`
}

// ErrorPayload standar untuk error (mirror helper.Response pada kasus error)
type ErrorPayload struct {
	Success bool   	`json:"success" example:"false"`
	Error   string 	`json:"error" example:"internal server error"`
}
