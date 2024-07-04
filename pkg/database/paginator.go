package database

import (
	"context"
	"math"
	"strconv"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type (
	PaginatorInterface interface {
		Raw(query string, vars []interface{}, countQuery string, countVars []interface{}) *Paginator
		UpdatePageInfo(ctx context.Context)
		Find(ctx context.Context) *gorm.DB
		FindWithOrderBy(ctx context.Context, orderBy string) *gorm.DB
	}
)

// Paginator structure containing pagination information and result records.
// Can be sent to the client directly.
type Paginator struct {
	DB *gorm.DB `json:"-"`

	Records interface{} `json:"records"`

	rawQuery          string
	rawQueryVars      []interface{}
	rawCountQuery     string
	rawCountQueryVars []interface{}

	MaxPage     int64 `json:"max_page"`
	Total       int64 `json:"total"`
	PageSize    int   `json:"page_size"`
	CurrentPage int   `json:"current_page"`

	loadedPageInfo bool
}

func paginateScope(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// NewPaginator create a new Paginator.
//
// Given DB transaction can contain clauses already, such as WHERE, if you want to
// filter results.
//
//	articles := []model.Article{}
//	tx := database.Conn().Where("title LIKE ?", "%"+sqlutil.EscapeLike(search)+"%")
//	paginator := database.NewPaginator(tx, page, pageSize, &articles)
//	result := paginator.Find()
//	if response.HandleDatabaseError(result) {
//	    response.JSON(http.StatusOK, paginator)
//	}
func NewPaginator(db *gorm.DB, page, pageSize int, dest interface{}) PaginatorInterface {
	return &Paginator{
		DB:          db,
		CurrentPage: page,
		PageSize:    pageSize,
		Records:     dest,
	}
}

// Raw set a raw SQL query and count query.
// The Paginator will execute the raw queries instead of automatically creating them.
// The raw query should not contain the "LIMIT" and "OFFSET" clauses, they will be added automatically.
// The count query should return a single number (`COUNT(*)` for example).
func (p *Paginator) Raw(query string, vars []interface{}, countQuery string, countVars []interface{}) *Paginator {
	p.rawQuery = query
	p.rawQueryVars = vars
	p.rawCountQuery = countQuery
	p.rawCountQueryVars = vars
	return p
}

// UpdatePageInfo executes count request to calculate the `Total` and `MaxPage`.
func (p *Paginator) UpdatePageInfo(ctx context.Context) {
	count := int64(0)
	db := p.DB.WithContext(ctx).Session(&gorm.Session{})
	prevPreloads := db.Statement.Preloads
	if len(prevPreloads) > 0 {
		db.Statement.Preloads = map[string][]interface{}{}
		defer func() {
			db.Statement.Preloads = prevPreloads
		}()
	}
	var err error
	if p.rawCountQuery != "" {
		err = db.Raw(p.rawCountQuery, p.rawCountQueryVars...).Scan(&count).Error
	} else {
		err = db.Model(p.Records).Count(&count).Error
	}
	if err != nil {
		panic(err)
	}
	p.Total = count
	p.MaxPage = int64(math.Ceil(float64(count) / float64(p.PageSize)))
	if p.MaxPage == 0 {
		p.MaxPage = 1
	}
	p.loadedPageInfo = true
}

// Find requests page information (total records and max page) and
// executes the transaction. The Paginate struct is updated automatically, as
// well as the destination slice given in NewPaginator().
func (p *Paginator) Find(ctx context.Context) *gorm.DB {
	if !p.loadedPageInfo {
		p.UpdatePageInfo(ctx)
	}
	if p.rawQuery != "" {
		return p.rawStatement(ctx).Scan(p.Records)
	}
	return p.DB.Scopes(paginateScope(p.CurrentPage, p.PageSize)).Find(p.Records)
}

func (p *Paginator) FindWithOrderBy(ctx context.Context, orderBy string) *gorm.DB {
	if !p.loadedPageInfo {
		p.UpdatePageInfo(ctx)
	}
	if p.rawQuery != "" {
		return p.rawStatementCustomOrderBy(ctx, orderBy).Scan(p.Records)
	}
	return p.DB.Scopes(paginateScope(p.CurrentPage, p.PageSize)).Find(p.Records)
}

func (p *Paginator) rawStatement(ctx context.Context) *gorm.DB {
	offset := (p.CurrentPage - 1) * p.PageSize
	db := p.DB.WithContext(ctx).Raw(p.rawQuery, p.rawQueryVars...)

	db.Statement.SQL.WriteString(" ")

	if db.Dialector.Name() == "sqlserver" {
		if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil {
			db.Statement.SQL.WriteString("ORDER BY ")
			db.Statement.WriteQuoted(db.Statement.Schema.PrioritizedPrimaryField.DBName)
			db.Statement.SQL.WriteByte(' ')
		} else {
			db.Statement.SQL.WriteString("ORDER BY (SELECT NULL) ")
		}

		if p.CurrentPage > 0 {
			db.Statement.SQL.WriteString("OFFSET ")
			db.Statement.SQL.WriteString(strconv.Itoa(offset))
			db.Statement.SQL.WriteString(" ROWS")
		}

		if p.PageSize > 0 {
			if p.CurrentPage == 0 {
				db.Statement.SQL.WriteString("OFFSET 0 ROW")
			}
			db.Statement.SQL.WriteString(" FETCH NEXT ")
			db.Statement.SQL.WriteString(strconv.Itoa(p.PageSize))
			db.Statement.SQL.WriteString(" ROWS ONLY")
		}
	} else {
		clause.Limit{Limit: p.PageSize, Offset: offset}.Build(db.Statement)
	}

	return db
}

func (p *Paginator) rawStatementCustomOrderBy(ctx context.Context, orderBy string) *gorm.DB {
	offset := (p.CurrentPage - 1) * p.PageSize
	db := p.DB.WithContext(ctx).Raw(p.rawQuery, p.rawQueryVars...)

	db.Statement.SQL.WriteString(" ")

	if db.Statement.Schema != nil && db.Statement.Schema.PrioritizedPrimaryField != nil {
		db.Statement.SQL.WriteString("ORDER BY ")
		db.Statement.WriteQuoted(db.Statement.Schema.PrioritizedPrimaryField.DBName)
		db.Statement.SQL.WriteByte(' ')
	} else {
		db.Statement.SQL.WriteString(orderBy)
	}

	if p.CurrentPage > 0 {
		db.Statement.SQL.WriteString(" OFFSET ")
		db.Statement.SQL.WriteString(strconv.Itoa(offset))
		db.Statement.SQL.WriteString(" ROWS")
	}

	if p.PageSize > 0 {
		if p.CurrentPage == 0 {
			db.Statement.SQL.WriteString(" OFFSET 0 ROW")
		}
		db.Statement.SQL.WriteString(" FETCH NEXT ")
		db.Statement.SQL.WriteString(strconv.Itoa(p.PageSize))
		db.Statement.SQL.WriteString(" ROWS ONLY")
	}

	return db
}
