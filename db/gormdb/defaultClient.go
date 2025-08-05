package gormdb

import (
	"context"

	"gorm.io/gorm"
)

var _ CRUD[struct{}] = (*crud[struct{}])(nil)

type crud[T any] struct {
	*gorm.DB
}

func NewCRUD[T any](db *gorm.DB) CRUD[T] {
	return &crud[T]{db}
}

func (r *crud[T]) Create(ctx context.Context, entities ...*T) error {
	if err := r.DB.WithContext(ctx).Create(entities).Error; err != nil {
		return err
	}

	return nil
}

func (r *crud[T]) Get(ctx context.Context, query *Query, opts ...QueryOptFn) (*T, error) {
	result := new(T)
	o := BuildOpt(opts...)

	db := r.DB.WithContext(ctx).Where(query.q).Not(query.not)

	// Apply preloads if specified
	for _, preload := range o.Preloads {
		db = db.Preload(preload)
	}

	if err := db.First(result).Error; err != nil && o.OmitNotFoundErr {
		return nil, o.OmitNotFoundErrFn(err)
	}

	return result, nil
}

func (r *crud[T]) List(ctx context.Context, query *Query, opts ...QueryOptFn) (*ListRes[T], error) {
	results := make([]*T, 0)
	o := BuildOpt(opts...)

	db := r.DB.WithContext(ctx).Where(query.q).Not(query.not)

	// Apply sorting if specified
	for _, orderBy := range o.OrderBy {
		db = db.Order(orderBy)
	}

	// Count total records if pagination is enabled
	if o.Paginate {
		if err := db.Model(new(T)).Count(&o.TotalCount).Error; err != nil {
			return nil, err
		}

		// Apply pagination
		offset := (o.PageNumber - 1) * o.PageSize
		db = db.Offset(offset).Limit(o.PageSize)
	}

	// Apply preloads if specified
	for _, preload := range o.Preloads {
		db = db.Preload(preload)
	}

	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	pageCount := int(o.TotalCount / int64(o.PageSize))
	if o.TotalCount%int64(o.PageSize) > 0 {
		pageCount++
	}
	return &ListRes[T]{Items: results, Total: o.TotalCount, PageSize: o.PageSize, PageCount: pageCount, Page: o.PageNumber}, nil
}

func (r *crud[T]) Update(ctx context.Context, query *Query, uParam map[string]any) error {
	updatedEntity := new(T)
	// 创建完成后 ID，CreatedAt，UpdatedAt 会回填到 updatedEntity 中吗？待确认
	if err := r.DB.WithContext(ctx).Model(updatedEntity).Where(query.q).Not(query.not).Updates(uParam).Error; err != nil {
		return err
	}

	return nil
}

func (r *crud[T]) Delete(ctx context.Context, query *Query, opts ...QueryOptFn) error {
	o := BuildOpt(opts...)

	var t T

	if err := r.DB.WithContext(ctx).Where(query.q).Not(query.not).Delete(&t).Error; err != nil && o.OmitNotFoundErr {
		return o.OmitNotFoundErrFn(err)
	}

	return nil
}

// 如此，repo 层就没有业务逻辑代码了，updateFn 虽然参数只有 *T，
// 不过在业务层可以临时闭包函数的形式捕获业务层变量，以更新 *T
// 这种方式称做 updateFn pattern
func (r *crud[T]) UpdateByFn(ctx context.Context, query *Query, updateFn func(*T) (bool, error)) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		updatedEntity := new(T)

		if err := tx.WithContext(ctx).Where(query.q).Not(query.not).First(updatedEntity).Error; err != nil {
			return err
		}

		updated, err := updateFn(updatedEntity)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		if err := tx.WithContext(ctx).Save(updatedEntity).Error; err != nil {
			return err
		}

		return nil
	})
}

// Implementation of transaction for CRUD operations
func (r *crud[T]) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx)
	})
}
