package auth

type AuthService interface{}

type authService struct {
	authRepo AuthRepository
}

func NewAuthService(authRepo AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}
