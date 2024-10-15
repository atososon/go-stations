package service

import (
	"context"
    "fmt"
	"database/sql"
    "strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
    // Prepare the insert statement
    stmt, err := s.db.PrepareContext(ctx, insert)
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    // Execute the insert statement
    res, err := stmt.ExecContext(ctx, subject, description)
    if err != nil {
        return nil, err
    }

    // Get the ID of the inserted row
    id, err := res.LastInsertId()
    if err != nil {
        return nil, err
    }

    // Prepare the confirm statement
    stmt, err = s.db.PrepareContext(ctx, confirm)
    if err != nil {
        return nil, err
    }

    // Query the inserted row
    row := stmt.QueryRowContext(ctx, id)

    // Scan the result into a TODO model
    todo := &model.TODO{}
    err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
    if err != nil {
        return nil, err
    }

    return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
    // Prepare the read statement
    stmt, err := s.db.PrepareContext(ctx, read)
    if prevID != 0 {
        stmt, err = s.db.PrepareContext(ctx, readWithID)
    }
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    // Execute the read statement
    var rows *sql.Rows
    if prevID != 0 {
        rows, err = stmt.QueryContext(ctx, prevID, size)
    } else {
        rows, err = stmt.QueryContext(ctx, size)
    }
    
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Scan the result into a slice of TODO models
    todos := []*model.TODO{}
    for rows.Next() {
        todo := &model.TODO{}
        err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
        if err != nil {
            return nil, err
        }
        todos = append(todos, todo)
    }

    return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
    // Prepare the update statement
    stmt, err := s.db.PrepareContext(ctx, update)
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    // Execute the update statement
    _, err = stmt.ExecContext(ctx, subject, description, id)
    if err != nil {
        return nil, err
    }

    // Prepare the confirm statement
    stmt, err = s.db.PrepareContext(ctx, confirm)
    if err != nil {
        return nil, &model.ErrNotFound{Message: err.Error()}
    }
    defer stmt.Close()

    // Execute the confirm statement
    _, err = stmt.ExecContext(ctx, id)
    if err != nil {
        return nil, &model.ErrNotFound{Message: err.Error()}
    }

    // Query the updated row
    row := stmt.QueryRowContext(ctx, id)

    // Scan the result into a TODO model
    todo := &model.TODO{}
    todo.ID = id
    err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
    if err != nil {
        return nil, &model.ErrNotFound{Message: err.Error()}
    }

    return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
    if len(ids) == 0 {
        return nil
    }
    // Prepare the delete statement with the correct number of placeholders
    placeholders := strings.Repeat(",?", len(ids)-1)
    stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(`DELETE FROM todos WHERE id IN (?%s)`, placeholders))
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Convert ids to []interface{}
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        args[i] = id
    }

    // Execute the delete statement
    res, err := stmt.ExecContext(ctx, args...)
    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return &model.ErrNotFound{Message: "no rows affected"}
    }

    return nil
}
