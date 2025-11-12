package auth

// CreateInput represents input for user registration
type CreateInput struct {
	Username string
	Password string
	Email    string
}

// LoginInput represents input for user login
type LoginInput struct {
	Username string
	Password string
}

// LoginView represents output for user login
type LoginView struct {
	AccessToken  string
	RefreshToken string
}

// RegisterView represents output for user registration
type RegisterView struct {
	Username string
	Message  string
}

