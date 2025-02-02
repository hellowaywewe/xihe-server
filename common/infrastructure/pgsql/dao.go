package pgsql

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	errRowExists   = errors.New("row exists")
	errRowNotFound = errors.New("row not found")
)

type SortByColumn struct {
	Column string
	Ascend bool
}

func (s SortByColumn) order() string {
	v := " ASC"
	if !s.Ascend {
		v = " DESC"
	}
	return s.Column + v
}

type Pagination struct {
	PageNum      int
	CountPerPage int
}

func (p Pagination) pagination() (limit, offset int) {
	limit = p.CountPerPage

	if limit > 0 && p.PageNum > 0 {
		offset = (p.PageNum - 1) * limit
	}

	return
}

type dbTable struct {
	name string
}

func NewDBTable(name string) dbTable {
	return dbTable{name: name}
}

func (t dbTable) DB() *gorm.DB {
	return db
}

func (t dbTable) Create(result interface{}) error {
	return db.Table(t.name).
		Create(result).
		Error
}

func (t dbTable) Updates(filter, result interface{}) error {
	return db.Table(t.name).
		Where(filter).
		Updates(result).
		Error
}

func (t dbTable) GetRecords(
	filter, result interface{}, p Pagination,
	sort []SortByColumn,
) (err error) {
	query := db.Table(t.name).Where(filter)

	var orders []string
	for _, v := range sort {
		orders = append(orders, v.order())
	}

	if len(orders) >= 0 {
		query.Order(strings.Join(orders, ","))
	}

	if limit, offset := p.pagination(); limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	err = query.Find(result).Error

	return
}

func (t dbTable) Count(filter interface{}) (int, error) {
	var total int64
	err := db.Table(t.name).Where(filter).Count(&total).Error

	return int(total), err
}

func (t dbTable) Filter(filter, result interface{}) error {
	return db.Table(t.name).
		Find(result, filter).Error
}

func (t dbTable) First(filter, result interface{}) error {
	return db.Table(t.name).
		First(result, filter).
		Error
}

func (t dbTable) GetOrderOneRecord(filter, order, result interface{}) error {
	err := db.Table(t.name).Where(filter).Order(order).Limit(1).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errRowNotFound
	}

	return nil
}

func (t dbTable) GetRecord(filter, result interface{}) error {
	err := db.Table(t.name).Where(filter).First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errRowNotFound
	}

	return err
}

func (t dbTable) UpdateRecord(filter, update interface{}) (err error) {
	query := db.Table(t.name).Where(filter).Updates(update)
	if err = query.Error; err != nil {
		return
	}

	if query.RowsAffected == 0 {
		err = errRowNotFound
	}

	return
}

func (t dbTable) IsRowNotFound(err error) bool {
	return errors.Is(err, errRowNotFound)
}

func (t dbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}
