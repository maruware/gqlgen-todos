package loader

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/maruware/gqlgen-todos/entity"
	"github.com/maruware/gqlgen-todos/graph/model"
	"gorm.io/gorm"
)

type loadersKeyType string

const loadersKey loadersKeyType = "dataloaders"

type Loaders struct {
	UserByID *dataloader.Loader
}

func newLoaders(db *gorm.DB) *Loaders {
	return &Loaders{
		UserByID: dataloader.NewBatchedLoader(
			newUserLoaderFunc(db),
			dataloader.WithWait(1*time.Millisecond)),
	}
}

func newUserLoaderFunc(db *gorm.DB) dataloader.BatchFunc {
	return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
		ids := []int{}
		for _, key := range keys {
			id, err := strconv.Atoi(key.String())
			if err != nil {
				continue
			}
			ids = append(ids, id)
		}

		var records []entity.User
		if err := db.Debug().Find(&records, ids).Error; err != nil {
			return []*dataloader.Result{}
		}

		userByID := map[string]*model.User{}
		for _, record := range records {
			user := model.NewUserFromEntity(&record)
			userByID[user.ID] = user
		}

		results := make([]*dataloader.Result, len(keys))
		for i, key := range keys {
			k := key.String()
			results[i] = &dataloader.Result{Data: nil, Error: nil}
			if user, ok := userByID[k]; ok {
				results[i].Data = user
			} else {
				results[i].Error = fmt.Errorf("user[key=%s] not found", k)
			}
		}

		return results
	}
}

func LoadUser(ctx context.Context, id string) (*model.User, error) {
	loader := ctx.Value(loadersKey).(*Loaders)
	thunk := loader.UserByID.Load(ctx, dataloader.StringKey(id))
	data, err := thunk()
	if err != nil {
		return nil, err
	}
	return data.(*model.User), nil
}

func DataLoaderMiddleware(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, newLoaders(db))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
