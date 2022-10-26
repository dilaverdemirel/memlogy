package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func EntryRepository() *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Create(entry Entry) (*Entry, error) {
	res, err := r.db.Exec("INSERT INTO entry(log_time, description) values(?,?)",
		entry.LogTime, entry.Description)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	entry.Id = id

	return &entry, nil
}

func (r *SQLiteRepository) All() ([]Entry, error) {
	rows, err := r.db.Query("SELECT * FROM entry")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.Id, &entry.LogTime, &entry.Description); err != nil {
			return nil, err
		}
		all = append(all, entry)
	}
	return all, nil
}

func (r *SQLiteRepository) GetById(id string) (*Entry, error) {
	row := r.db.QueryRow("SELECT * FROM entry WHERE id = ?", id)

	var entry Entry
	if err := row.Scan(&entry.Id, &entry.LogTime, &entry.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &entry, nil
}

func (r *SQLiteRepository) GetByDay(day time.Time) ([]Entry, error) {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())
	rows, err := r.db.Query("SELECT * FROM entry WHERE log_time >= ? and log_time <= ? order by id", startOfDay, endOfDay)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.Id, &entry.LogTime, &entry.Description); err != nil {
			return nil, err
		}
		all = append(all, entry)
	}
	return all, nil
}

func (r *SQLiteRepository) Update(id int64, updated Entry) (*Entry, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE entry SET log_time = ?, description = ? WHERE id = ?", updated.LogTime, updated.Description, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM entry WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
