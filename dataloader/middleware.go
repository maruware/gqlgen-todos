package dataloader

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/maruware/gqlgen-todos/entity"
	"github.com/maruware/gqlgen-todos/graph/model"
	"gorm.io/gorm"
)

type loadersKeyType string

const loadersKey loadersKeyType = "dataloaders"

type Loaders struct {
	UserByID *UserLoader
}

func DataLoaderMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			UserByID: NewUserLoader(UserLoaderConfig{
				MaxBatch: 100,
				Wait:     1 * time.Millisecond,
				Fetch: func(keys []string) ([]*model.User, []error) {
					ids := make([]int, len(keys))
					for i, key := range keys {
						id, err := strconv.Atoi(key)
						if err != nil {
							return nil, []error{err}
						}
						ids[i] = id
					}

					var records []entity.User
					if err := db.Debug().Find(&records, ids).Error; err != nil {
						return nil, []error{err}
					}

					userByID := map[string]*model.User{}
					for _, record := range records {
						user := model.NewUserFromEntity(&record)
						userByID[user.ID] = user
					}

					users := make([]*model.User, len(ids))
					for i, key := range keys {
						users[i] = userByID[key]
					}

					return users, nil
				},
			}),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
