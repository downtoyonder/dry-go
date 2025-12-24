package gormdb

import (
	"context"
)

type Query struct {
	q   map[string]any
	not map[string]any
}

func (q *Query) Not(not map[string]any) *Query {
	q.not = not
	return q
}

func Q(q map[string]any) *Query {
	return &Query{q: q}
}

type (
	QueryOptFn  func(c *QueryOpt) *QueryOpt
	QueryOptFns []QueryOptFn
)

type QueryOpt struct {
	OmitNotFoundErrFn func(err error) error
	OrderBy           []string
	Preloads          []string
	TotalCount        int64
	PageNumber        int
	PageSize          int
	OmitNotFoundErr   bool
	Paginate          bool
}

func NewQueryOpt() *QueryOpt {
	return &QueryOpt{
		PageNumber: 1,  // Default to first page
		PageSize:   50, // Default page size
	}
}

func OmitNotFoundErr(errFn func(err error) error) QueryOptFn {
	if errFn == nil {
		panic("OmitNotFoundErr requires a non-nil error handling function")
	}

	return func(c *QueryOpt) *QueryOpt {
		c.OmitNotFoundErr = true
		c.OmitNotFoundErrFn = errFn
		return c
	}
}

// ListRes holds the result of a list query along with pagination information
type ListRes[T any] struct {
	Items     []*T  // The actual items retrieved
	Total     int64 // Total number of records matching the query (ignoring pagination)
	PageSize  int   // Size of each page
	PageCount int   // Total number of pages
	Page      int   // Current page number
}

// Pagination enables pagination with specified page number and size
func Pagination(pageNumber, pageSize int) QueryOptFn {
	return func(c *QueryOpt) *QueryOpt {
		c.Paginate = true
		if pageNumber > 0 {
			c.PageNumber = pageNumber
		}
		if pageSize > 0 {
			c.PageSize = pageSize
		}
		return c
	}
}

// OrderBy sets the order by columns
func OrderBy(orderBy ...string) QueryOptFn {
	return func(c *QueryOpt) *QueryOpt {
		c.OrderBy = orderBy
		return c
	}
}

// Preload sets the relations to preload
func Preload(preloads ...string) QueryOptFn {
	return func(c *QueryOpt) *QueryOpt {
		c.Preloads = preloads
		return c
	}
}

func (opts QueryOptFns) Build() *QueryOpt {
	c := NewQueryOpt()

	for _, fn := range opts {
		c = fn(c)
	}

	return c
}

func BuildOpt(opts ...QueryOptFn) *QueryOpt {
	c := NewQueryOpt()

	for _, fn := range opts {
		c = fn(c)
	}

	return c
}

type CRUD[T any] interface {
	// Create supports create one or multiple records
	// 创建完成后 ID，CreatedAt，UpdatedAt 会回填到 entities 中
	Create(ctx context.Context, entities ...*T) error
	// Get retrieve one record matches the conditions.
	Get(ctx context.Context, query *Query, opts ...QueryOptFn) (*T, error)
	// List retrieve all records matches the conditions.
	List(ctx context.Context, query *Query, opts ...QueryOptFn) (*ListRes[T], error)
	// Update set one or more records match the conditions according to updateParam
	Update(ctx context.Context, query *Query, uParam map[string]any) error
	// Delete supports delete one or multiple records
	Delete(ctx context.Context, query *Query, opts ...QueryOptFn) error

	// UpdateByFn updates an entity using a function that can contain business logic
	UpdateByFn(ctx context.Context, query *Query, updateFn func(*T) (bool, error)) error

	// Transaction executes operations within a database transaction
	Transaction(ctx context.Context, f func(ctx context.Context) error) error
}
