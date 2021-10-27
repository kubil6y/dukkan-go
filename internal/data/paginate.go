package data

import (
	"math"
	"net/http"
	"strconv"

	"github.com/kubil6y/dukkan-go/internal/validator"
	"gorm.io/gorm"
)

type Paginate struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func ValidatePaginate(p *Paginate, v *validator.Validator) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Limit > 0, "limit", "must be greater than zero")
	v.Check(p.Page <= 10_000, "page", "must be a maximum of 10_000")
	v.Check(p.Limit <= 100, "limit", "must be a maximum of 100")
}

func NewPaginate(r *http.Request, v *validator.Validator, limitDefault, pageDefault int) *Paginate {
	return &Paginate{
		Limit: readInt(r, v, "limit", limitDefault),
		Page:  readInt(r, v, "page", pageDefault),
	}
}

// PaginatedResults is used when making db calls, example:
// err := m.DB.Scopes(p.PaginatedResults).Find(&users).Error
func (p Paginate) PaginatedResults(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit)
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetadata(p *Paginate, total int) Metadata {
	if total == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  p.Page,
		PageSize:     p.Limit,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(total) / float64(p.Limit))),
		TotalRecords: total,
	}
}

// NOTE code duplication for NewPaginate, cmd/api/helpers.go app.readInt
func readInt(r *http.Request, v *validator.Validator, key string, defaultValue int) int {
	qs := r.URL.Query()
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "invalid value")
		return defaultValue
	}
	return i
}
