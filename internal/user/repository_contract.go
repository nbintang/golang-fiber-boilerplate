package user

import "context"

type UserRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]User,int64, error)
	FindByID(ctx context.Context, id string) (*User, error) 
	FindByIDWithRole(ctx context.Context, id string) (*User, error) 
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, dto *User) error
	Update(ctx context.Context, id string, dto *User) error
}
