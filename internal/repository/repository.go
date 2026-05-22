package repository

import (
	"gorm.io/gorm"
)

// Dependencies holds all infrastructure deps for repository constructors.
type Dependencies struct {
	DB *gorm.DB
}

// Repository is the aggregate that holds all feature repositories.
// Add a field here for each new feature repository.
type Repository struct {
	// Example: Health *health.Repository
}

// New creates the Repository aggregate. Register feature repositories here.
func New(d Dependencies) *Repository {
	return &Repository{}
}
