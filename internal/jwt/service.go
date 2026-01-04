package jwt

type JwtService interface {
	GenerateAccessToken(userId uint64) (string, error)
	GenerateRefreshToken(userId uint64) (string, error)
	ValidateRefreshToken(token string) (*Jwtuser, error)
	ValidateAccessToken(token string) (*Jwtuser, error)
}

type jwtService struct {
	jwtConfig *JwtConfig
}

func NewJwtService(jwtConfig *JwtConfig) JwtService {
	return &jwtService{
		jwtConfig: jwtConfig,
	}
}

func (s *jwtService) GenerateAccessToken(userId uint64) (string, error) {
	return "", nil
}

func (s *jwtService) GenerateRefreshToken(userId uint64) (string, error) {
	return "", nil
}

func (s *jwtService) ValidateRefreshToken(token string) (*Jwtuser, error) {
	claims, err := ValidateJWT(token, s.jwtConfig.RefreshSecret)
	if err != nil {
		return nil, err
	}

	if !CheckClaims(claims, "user_id", "exp") {
		return nil, err
	}

	return &Jwtuser{
		UserId: claims["userId"].(int64),
		Claims: &claims,
	}, nil
}

func (s *jwtService) ValidateAccessToken(token string) (*Jwtuser, error) {
	claims, err := ValidateJWT(token, s.jwtConfig.AccessSecret)
	if err != nil {
		return nil, err
	}

	if !CheckClaims(claims, "user_id", "exp") {
		return nil, err
	}
	return &Jwtuser{
		UserId: claims["userId"].(int64),
		Claims: &claims,
	}, nil
}
