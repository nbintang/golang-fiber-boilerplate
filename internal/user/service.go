package user

import (
	"context"
	"encoding/json"
	"errors" 
	"rest-fiber/internal/infra/infraapp"
	"rest-fiber/internal/infra/rediscache"
	"rest-fiber/pkg/slice" 
)

type userServiceImpl struct {
	userRepo     UserRepository
	logger       *infraapp.AppLogger
	redisService rediscache.Service
}

func NewUserService(userRepo UserRepository, logger *infraapp.AppLogger, redisService rediscache.Service) UserService {
	return &userServiceImpl{userRepo, logger, redisService}
}

func (s *userServiceImpl) FindAllUsers(ctx context.Context, page, limit, offset int) ([]UserResponseDTO, int64, error) {
	version := s.getListVersion(ctx, usersVersionKey)
	cacheKey := usersListKey(version, limit, offset)

	if cached, err := s.redisService.Get(ctx, cacheKey); err == nil && cached != "" {
		var payload usersListCache
		if err := json.Unmarshal([]byte(cached), &payload); err == nil {
			return payload.Items, payload.Total, nil
		}
	}

	users, total, err := s.userRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	userResponses := slice.Map[User, UserResponseDTO](users, func(u User) UserResponseDTO {
		return UserResponseDTO{
			ID:        u.ID,
			Name:      u.Name,
			AvatarURL: u.AvatarURL,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	})

	payload := usersListCache{Total: total, Items: userResponses}
	if b, err := json.Marshal(payload); err == nil {
		_ = s.redisService.Set(ctx, cacheKey, string(b), usersListTTL)
	}

	return userResponses, total, nil
}

func (s *userServiceImpl) FindUserByID(ctx context.Context, id string) (*UserResponseDTO, error) {
	key := userByIDKey(id)

	if cached, err := s.redisService.Get(ctx, key); err == nil && cached != "" {
		var dto UserResponseDTO
		if err := json.Unmarshal([]byte(cached), &dto); err == nil {
			return &dto, nil
		}
	}
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User Not Found")
	}

	dto := &UserResponseDTO{
		ID:        user.ID,
		AvatarURL: user.AvatarURL,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	if b, err := json.Marshal(dto); err == nil {
		_ = s.redisService.Set(ctx, key, string(b), userByIDTTL)
	}

	return dto, nil
}

func (s *userServiceImpl) UpdateProfile(ctx context.Context, id string, dto UserUpdateDTO) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	user.Name = dto.Name
	user.AvatarURL = dto.AvatarURL
	return s.userRepo.Update(ctx, id, user)
}
