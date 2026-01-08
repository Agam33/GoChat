package request

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInRequst struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetRefreshToken struct {
	RefreshToken string `json:"refreshToken"`
}
