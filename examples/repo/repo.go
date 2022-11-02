package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type TestEntity struct {
	Content string
}

type repo struct {
	conn *pgx.Conn
}

func connString(user, password, address, db string, port int) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		address,
		port,
		db,
		user,
		password)
}

func New(user, password, address, db string, port int) (*repo, error) {
	conn, err := pgx.Connect(
		context.Background(),
		connString(user, password, address, db, port))
	if err != nil {
		return nil, err
	}

	return &repo{conn: conn}, nil
}

func (r *repo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *repo) InsertContent(ctx context.Context, testEntities []TestEntity) error {
	for i := range testEntities {
		_, err := r.conn.Exec(context.Background(), "INSERT INTO test_table(content) VALUES($1)", testEntities[i].Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repo) SelectContent(ctx context.Context) ([]TestEntity, error) {
	rows, err := r.conn.Query(ctx, "SELECT content FROM test_table")
	if err != nil {
		return nil, err
	}

	var testEntities []TestEntity
	for rows.Next() {
		var content string
		err := rows.Scan(&content)
		if err != nil {
			return nil, err
		}
		testEntities = append(testEntities, TestEntity{Content: content})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()

	return testEntities, nil
}
