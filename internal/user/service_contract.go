package user

import "context"

type UserService interface {
	FindAllUsers(ctx context.Context, page, limit, offset int) ([]UserResponseDTO, int64, error)
	FindUserByID(ctx context.Context, id string) (*UserResponseDTO, error)
	UpdateProfile(ctx context.Context, id string, dto UserUpdateDTO) error
}
