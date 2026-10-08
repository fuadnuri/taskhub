package service 


type TokenService interface{
	generateToken()string
	validateToken() string 
}