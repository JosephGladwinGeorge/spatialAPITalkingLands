package backend

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Point struct {
	ID        uuid.UUID   `db:"id" json:"id"`         
	Name      string  `db:"name" json:"name"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}

func (p *Point) createPoint(db *sqlx.DB) error {
	p.ID = uuid.New()
	query := `INSERT INTO points_data (id, name, latitude, longitude) VALUES (:id, :name, :latitude, :longitude) RETURNING id`
	rows, err := db.NamedQuery(query, p)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&p.ID)
	}
	return err
}

func (p *Point) fetchPoint(db *sqlx.DB) error {
	query := `SELECT id, name, latitude, longitude FROM points_data WHERE id = $1`
	return db.Get(p, query, p.ID)
}

func (p *Point) updatePoint(db *sqlx.DB) error {
	query := `UPDATE points_data SET name = :name, latitude = :latitude, longitude = :longitude WHERE id = :id`
	_, err := db.NamedExec(query, p)
	return err
}


