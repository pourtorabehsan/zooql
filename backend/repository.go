package main

import (
	"database/sql"
	"sync"
	"time"
)

type znodeRepository struct {
	mu sync.RWMutex
	db *sql.DB
}

func (z *znodeRepository) create(path, data string) error {
	z.mu.Lock()
	defer z.mu.Unlock()

	insert := "INSERT INTO znodes (path, data) VALUES (?, ?) ON CONFLICT(path) DO UPDATE SET data=?"
	_, err := z.db.Exec(insert, path, data, data)
	return err
}

func (z *znodeRepository) delete(path string) error {
	z.mu.Lock()
	defer z.mu.Unlock()

	delete := "DELETE FROM znodes WHERE path=?"
	_, err := z.db.Exec(delete, path)
	return err
}

func (z *znodeRepository) update(path, data string) error {
	z.mu.Lock()
	defer z.mu.Unlock()

	update := "UPDATE znodes SET data=? WHERE path=?"
	_, err := z.db.Exec(update, data, path)
	return err
}

type queryResult struct {
	Columns []string   `json:"columns"`
	Rows    [][]string `json:"rows"`
	Elapsed string     `json:"elapsed"`
}

func (z *znodeRepository) query(q string) (*queryResult, error) {
	z.mu.RLock()
	defer z.mu.RUnlock()

	now := time.Now()
	rows, err := z.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	resultRows := make([][]string, 0)

	for rows.Next() {
		resultRows = append(resultRows, rowToArray(rows, cols))
	}

	result := &queryResult{
		Columns: cols,
		Rows:    resultRows,
		Elapsed: time.Since(now).String(),
	}

	return result, nil
}

func rowToArray(rows *sql.Rows, columns []string) []string {
	buff := make([]interface{}, len(columns))
	data := make([]string, len(columns))
	for i := 0; i < len(buff); i++ {
		buff[i] = &data[i]
	}
	rows.Scan(buff...)
	return data
}
