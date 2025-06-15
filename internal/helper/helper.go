package helper

const (
	//eror handle
	ErrGetUser            = "Error get user from database"
	ErrExistUser          = "Email or username already exists"
	ErrGeneratePassword   = "Error create password, please try again"
	ErrComparePassword    = "email or password is incorrect"
	ErrCreateToken        = "failed to create token"
	ErrCreateTokenRefresh = "failed to create refresh token"

	// routes
	ApiGroup = "MusicCatalog"
	SignIn   = "/signIn"
	LogIn    = "/logIn"

	// column users
	Id       = "id"
	Email    = "email"
	Username = "username"
)
