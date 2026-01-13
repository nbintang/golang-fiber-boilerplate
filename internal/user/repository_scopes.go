package user

import "gorm.io/gorm"

type ScopeReturn func(db *gorm.DB) *gorm.DB

func WhereEmail(email string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}

func WhereID(id string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func SelectedFields(extra ...string) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		fields:= []string {
			"id",
			"name",
			"avatar_url",
			"password",
			"email",
			"is_email_verified",
		}
		fields = append(fields, extra...);
		return db.Select(fields)
	}
}
func Paginate(limit, offset int) ScopeReturn {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}
