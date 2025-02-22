package backend

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initDB()(*sqlx.DB,error){
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "admin"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "admin123"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "spatialdb"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;`)
	if err != nil {
		return nil, fmt.Errorf("failed to enable PostGIS extension: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS points_data (
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	    name TEXT NOT NULL,
	    latitude DOUBLE PRECISION NOT NULL,
	    longitude DOUBLE PRECISION NOT NULL,
	    geom GEOMETRY(Point, 4326) GENERATED ALWAYS AS (
	        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
	    ) STORED
	);

	CREATE TABLE IF NOT EXISTS polygons_data (
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	    name TEXT NOT NULL,
	    geojson JSONB NOT NULL,
	    geom GEOMETRY(Polygon, 4326) GENERATED ALWAYS AS (
	        ST_GeomFromGeoJSON(geojson)
	    ) STORED
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	fmt.Println("Connected to PostGIS database!")
	return db,nil
}