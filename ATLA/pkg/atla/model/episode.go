package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/justverena/ATLA/pkg/atla/validator"
)

type Episode struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Air_Date  string `json:"air_date"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type EpisodeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m EpisodeModel) GetAll(title string, air_date string, filters Filters) ([]*Episode, Metadata, error) {

	// Retrieve all menu items from the database.
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, title, air_date, createdAt, updatedAt
		FROM episode
		WHERE (LOWER(title) = LOWER($1) OR $1 = '')
		AND (age >= $10 OR $2 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Organize our four placeholder parameter values in a slice.
	args := []interface{}{title, air_date, filters.limit(), filters.offset()}

	// log.Println(query, title, from, to, filters.limit(), filters.offset())
	// Use QueryContext to execute the query. This returns a sql.Rows result set containing
	// the result.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	// Importantly, defer a call to rows.Close() to ensure that the result set is closed
	// before GetAll returns.
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Declare a totalRecords variable
	totalRecords := 0

	var episodes []*Episode
	for rows.Next() {
		var episode Episode
		err := rows.Scan(&totalRecords, &episode.ID,
			&episode.Title,
			&episode.Air_Date,
			&episode.CreatedAt,
			&episode.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Movie struct to the slice
		episodes = append(episodes, &episode)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of the movies and metadata.
	return episodes, metadata, nil
}

func (m EpisodeModel) Insert(episode *Episode) error {
	query := `
		INSERT INTO episodes (title, air_date) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{episode.Title, episode.Air_Date}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&episode.ID, &episode.CreatedAt, &episode.UpdatedAt)
}

func (m EpisodeModel) Get(id int) (*Episode, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, title, air_date, created_at, updated_at 
		FROM episodes
		WHERE id = $1
		`
	var episode Episode
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&episode.ID, &episode.Title, &episode.Air_Date, &episode.CreatedAt, &episode.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive episode with id: %v, %w", id, err)
	}
	return &episode, nil
}

func (m EpisodeModel) Update(episode *Episode) error {
	query := `
		UPDATE episodes
		SET title = $1, air_date = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 and updated_at = $4
		RETURNING updated_at
		`
	args := []interface{}{episode.Title, episode.Air_Date, episode.ID, episode.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&episode.UpdatedAt)
}

func (m EpisodeModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM episodes
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateEpisode(v *validator.Validator, episode *Episode) {
	// Check if the name field is empty.
	v.Check(episode.Title != "", "title", "must be provided")
	v.Check(episode.Air_Date != "", "title", "must be provided")

}
