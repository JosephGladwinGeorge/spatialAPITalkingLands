package backend

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Polygon struct {
	ID      uuid.UUID `db:"id" json:"id"`
	Name    string    `db:"name" json:"name"`
	GeoJSON string    `db:"geojson" json:"geojson"`
}

func (pg *Polygon) createPolygon(db *sqlx.DB) error {
	pg.ID = uuid.New()

	query := `INSERT INTO polygons_data (id, name, geojson) VALUES (:id, :name, :geojson)`
	_, err := db.NamedExec(query, pg)
	return err
}

func (pg *Polygon) fetchPolygon(db *sqlx.DB) error {
	query := `SELECT id, name, geojson FROM polygons_data WHERE id = $1`
	return db.Get(pg, query, pg.ID)
}

func (pg *Polygon) updatePolygon(db *sqlx.DB) error {
	query := `UPDATE polygons_data SET name = :name, geojson = :geojson WHERE id = :id`
	_, err := db.NamedExec(query, pg)
	return err
}