package main

import (
	"fmt"
	"math/rand"
	"rest-fiber/internal/enums" 
	"rest-fiber/internal/user"
	"rest-fiber/pkg/password"

	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Options struct {
	Count int
}

func InitSeeds(db *gorm.DB, opt Options) error {
	domain := "example.com"

	if err := resetTables(db); err != nil {
		return err
	}
	hashed, err := password.Hash("Password123")
	if err != nil {
		return err
	}
	for i := 0; i < opt.Count; i++ {
		name := faker.Name()

		userSeed := user.User{
			Name:            name,
			Email:           realisticEmail(name, i, domain),
			AvatarURL:       realisticAvatar(name),
			Password:        hashed,
			Role:            evenOddRole(i),
			IsEmailVerified: rand.Intn(2) == 1,
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&userSeed).Error; err != nil {
			fmt.Printf("Error when create user: %s\n", userSeed.Name)
			return err
		}
	}
	return nil
}

func evenOddRole(i int) user.Role {
	if i%2 == 0 {
		return user.Role(enums.Member)
	}
	return user.Role(enums.Admin)
}


func realisticEmail(fullName string, i int, domain string) string {
	// domain tanpa "@"
	domain = strings.TrimPrefix(domain, "@")
	user := slug.Make(fullName)
	return fmt.Sprintf("%s.%d@%s", user, i+1, domain)
}

func realisticAvatar(fullName string) string {
	// gampang & rapi, avatar dari initials
	// contoh: https://ui-avatars.com/api/?name=John+Doe&background=random
	name := strings.ReplaceAll(fullName, " ", "+")
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", name)
}


func realisticPostBody() string {
	// faker.Paragraph() biasanya lebih enak dari Sentence()
	// kalau faker versi kamu tidak punya Paragraph, ganti jadi faker.Sentence() beberapa kali
	return faker.Paragraph()
}

func resetTables(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE TABLE posts RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		if err := tx.Exec("TRUNCATE TABLE categories RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		if err := tx.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
		return nil
	})
}
