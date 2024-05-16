package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/justverena/ATLA/pkg/atla/validator"
)

type Character struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Status    string `json:"status"`
	Nation    string `json:"nation"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CharacterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m CharacterModel) GetAll(name string, age_from int, age_to int, filters Filters) ([]*Character, Metadata, error) {

	query := fmt.Sprintf(
		`
		SELECT count(*) OVER(), id, name, age, gender, status, nation, created_at, updated_at
		FROM characters
		WHERE (LOWER(name) = LOWER($1) OR $1 = '')
		AND (age >= $2 OR $2 = 0)
		AND (age <= $3 OR $3 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
		`,
		filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, age_from, age_to, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	totalRecords := 0

	var characters []*Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&totalRecords, &character.ID,
			&character.Name,
			&character.Age,
			&character.Gender,
			&character.Status,
			&character.Nation,
			&character.CreatedAt,
			&character.UpdatedAt)
		if err != nil {
			return nil, Metadata{}, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return characters, metadata, nil
}

func (m CharacterModel) Insert(character *Character) error {
	query := `
		INSERT INTO characters (name, age, gender, status, nation) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Status, character.Nation}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.CreatedAt, &character.UpdatedAt)
}

func (m CharacterModel) Get(id int) (*Character, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, name, age, gender, status, nation, created_at, updated_at 
		FROM characters
		WHERE id = $1
		`
	var character Character
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&character.ID, &character.Name, &character.Age, &character.Gender, &character.Status, &character.Nation, &character.CreatedAt, &character.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive character with id: %v, %w", id, err)
	}
	return &character, nil
}

func (m CharacterModel) Update(character *Character) error {
	query := `
		UPDATE characters
		SET name = $1, age = $2, gender = $3, status = $4, nation = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6 and updated_at = $7
		RETURNING updated_at
		`
	args := []interface{}{character.Name, character.Age, character.Gender, character.Status, character.Nation, character.ID, character.UpdatedAt}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&character.UpdatedAt)
}

func (m CharacterModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM characters
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

func (m *CharacterModel) GetByEpisode(episodeID int) ([]*Character, error) {
	query := `
        SELECT c.id, c.name, c.age, c.gender, c.status, c.nation, c.created_at, c.updated_at
        FROM characters c
        JOIN characters_and_episodes ce ON c.id = ce.character_id
        WHERE ce.episode_id = $1
        ORDER BY c.id
    `

	rows, err := m.DB.Query(query, episodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*Character
	for rows.Next() {
		var character Character
		err := rows.Scan(&character.ID, &character.Name, &character.Age, &character.Gender, &character.Status, &character.Nation, &character.CreatedAt, &character.UpdatedAt)
		if err != nil {
			return nil, err
		}
		characters = append(characters, &character)
	}

	return characters, nil
}

func (m *CharacterModel) GetByQuote(quoteID int) (*Character, error) {
	query := `
        SELECT c.id, c.name, c.age, c.gender, c.status, c.nation, c.created_at, c.updated_at
        FROM characters c
        JOIN characters_and_quotes cq ON c.id = cq.character_id
        WHERE cq.quote_id = $1
        ORDER BY c.id
    `

	var character Character
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, quoteID)
	err := row.Scan(&character.ID, &character.Name, &character.Age, &character.Gender, &character.Status, &character.Nation, &character.CreatedAt, &character.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("cannot retrive character with id: %v, %w", quoteID, err)
	}
	return &character, nil
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	v.Check(character.Name != "", "name", "must be provided")
	v.Check(character.Age <= 10000, "age", "must not be more than 10000 bytes long")
	v.Check(character.Gender != "", "gender", "must be provided")
	v.Check(character.Status != "", "status", "must be provided")
	v.Check(character.Nation != "", "nation", "must be provided")

}
