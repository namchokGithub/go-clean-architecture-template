package port

// Repository interfaces are defined here as features are added.
// Example pattern:
//
//	type ConfigurationRepository interface {
//	    FindAll(ctx context.Context) ([]domain.Configuration, error)
//	    FindByID(ctx context.Context, id int64) (domain.Configuration, error)
//	    Create(ctx context.Context, entity domain.Configuration) (domain.Configuration, error)
//	    Update(ctx context.Context, entity domain.Configuration) error
//	    Delete(ctx context.Context, id int64) error
//	}
