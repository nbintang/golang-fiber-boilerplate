package user

type UserUpdateDTO struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	AvatarURL string `json:"avatar_url" validate:"required,url"`
}
