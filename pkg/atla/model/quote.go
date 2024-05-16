package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/justverena/ATLA/pkg/atla/validator"
)

type Quote struct {
	ID        int       `json:"id"`
	Quote     string    `json:"quote"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type QuoteModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m QuoteModel) GetAll(quote string, filters Filters) ([]*Quote, Metadata, error) {
	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, quote, created_at, updated_at
		FROM quotes
		WHERE (quote ILIKE '%%' || $1 || '%%' OR $1 = '')
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
		`,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	// Organize our placeholder parameter values in a slice.
	args := []interface{}{quote, filters.limit(), filters.offset()}

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

	var quotes []*Quote
	for rows.Next() {
		var quote Quote
		err := rows.Scan(&totalRecords, &quote.ID, &quote.Quote, &quote.CreatedAt, &quote.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		// Add the Quote struct to the slice
		quotes = append(quotes, &quote)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination parameters
	// from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	// If everything went OK, then return the slice of quotes and metadata.
	return quotes, metadata, nil
}

func (m QuoteModel) Insert(quote *Quote) error {
	query := `
		INSERT INTO quotes (quote) 
		VALUES ($1) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{quote.Quote}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&quote.ID, &quote.CreatedAt, &quote.UpdatedAt)
}

func (m QuoteModel) Get(id int) (*Quote, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, quote, created_at, updated_at 
		FROM quotes
		WHERE id = $1
		`
	var quote Quote
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&quote.ID, &quote.Quote, &quote.CreatedAt, &quote.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		} else {
			return nil, fmt.Errorf("cannot retrieve quote with id: %v, %w", id, err)
		}
	}
	return &quote, nil
}

func (m QuoteModel) Update(quote *Quote) error {
	query := `
		UPDATE quotes
		SET quote = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at
		`
	args := []interface{}{quote.Quote, quote.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&quote.UpdatedAt)
}

func (m QuoteModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM quotes
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func ValidateQuote(v *validator.Validator, quote *Quote) {
	v.Check(quote.Quote != "", "quote", "must be provided")
}

func (m QuoteModel) GetQuotesByCharacterID(characterID int) ([]*Quote, error) {
	query := `
		SELECT q.id, q.quote, q.created_at, q.updated_at
		FROM quotes q
		JOIN characters_and_quotes cq ON q.id = cq.quote_id
		WHERE cq.character_id = $1
		ORDER BY q.id`

	rows, err := m.DB.Query(query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quotes := []*Quote{}
	for rows.Next() {
		var q Quote
		err := rows.Scan(&q.ID, &q.Quote, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, &q)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return quotes, nil
}
