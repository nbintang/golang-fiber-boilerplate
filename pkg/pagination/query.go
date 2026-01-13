package pagination
import "errors"

type Query struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (q Query) Normalize(defaultLimit, maxLimit int) Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = defaultLimit
	}

	if maxLimit > 0 && q.Limit > maxLimit {
		q.Limit = maxLimit
	}
	return q
}


func (q Query) RequirePaginate() error {
	if q.Page <= 0 {
		return errors.New("page must be greater than 0")
	}
	if q.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}
	return nil
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}
