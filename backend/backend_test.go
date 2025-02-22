package backend_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"spatialDB/backend"
	"testing"

	"github.com/google/uuid"
)

var app backend.App

const tablePointCreationQuery = `CREATE TABLE IF NOT EXISTS points_data (
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	    name TEXT NOT NULL,
	    latitude DOUBLE PRECISION NOT NULL,
	    longitude DOUBLE PRECISION NOT NULL,
	    geom GEOMETRY(Point, 4326) GENERATED ALWAYS AS (
	        ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)
	    ) STORED
	);`

const tablePolygonCreationQuery = `CREATE TABLE IF NOT EXISTS polygons_data (
	    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	    name TEXT NOT NULL,
	    geojson JSONB NOT NULL,
	    geom GEOMETRY(Polygon, 4326) GENERATED ALWAYS AS (
	        ST_GeomFromGeoJSON(geojson)
	    ) STORED
	);`

func TestMain(m *testing.M) {
	app = backend.App{}
	app.Initialize()
	ensureTablesExist()
	code:= m.Run()

	clearPointsTable()
	clearPolygonsTable()

	os.Exit(code)
}

func ensureTablesExist(){
	if _,err:=app.DB.Exec(tablePointCreationQuery);err!=nil{
		log.Fatal(err)
	}
	if _,err:=app.DB.Exec(tablePolygonCreationQuery);err!=nil{
		log.Fatal(err)
	}
}

func TestCreatePoint(t *testing.T){
	clearPointsTable()

	payload := []byte(`{"name":"TestPoint","latitude":12.34,"longitude":56.78}`)

	req,_:=http.NewRequest("POST","/points",bytes.NewBuffer(payload))

	response:=executeRequest(req)

	checkResponseCode(t,http.StatusOK,response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TestPoint"{
		t.Errorf("Expected point name to be 'TestPoint' but found '%s'",m["name"])
	}
	if m["latitude"] != 12.34{
		t.Errorf("Expected point latitude to be 12.34 but found '%d'",m["latitude"])
	}
	if m["longitude"] != 56.78{
		t.Errorf("Expected point longitude to be 56.78 but found '%d'",m["longitude"])
	}
	

}

func TestGetPoint(t *testing.T) {
	clearPointsTable()
	id:=addPoint()
	fmt.Println(id,id.String())

	req,_:=http.NewRequest("GET",fmt.Sprintf("/points/%s", id.String()),nil)

	response:=executeRequest(req)

	checkResponseCode(t,http.StatusOK,response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TestPoint"{
		t.Errorf("Expected point name to be TestPoint found '%s'",m["name"])
	}
	if m["latitude"] != 12.34{
		t.Errorf("Expected point latitude to be 12.34 found '%d'",m["latitude"])
	}
	if m["longitude"] != 56.78{
		t.Errorf("Expected point longitude to be 56.78 found '%d'",m["longitude"])
	}

}

func TestUpdatePoint(t *testing.T){
	clearPointsTable()
	id:=addPoint()

	payload := []byte(`{"name":"UpdatePoint","latitude":24.34,"longitude":65.78}`)

	req,_:=http.NewRequest("PUT",fmt.Sprintf("/points/%s",id.String()),bytes.NewBuffer(payload))

	response:=executeRequest(req)

	checkResponseCode(t,http.StatusOK,response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&m)

	if m["name"] != "UpdatePoint"{
		t.Errorf("Expected point name to be UpdatePoint found '%s'",m["name"])
	}
	if m["latitude"] != 24.34 {
		t.Errorf("Expected point latitude to be 24.34 found '%d'",m["latitude"])
	}
	if m["longitude"] != 65.78 {
		t.Errorf("Expected point longitude to be 65.78 found '%d'",m["longitude"]) 
	}
}

func addPoint() uuid.UUID{
	id:=uuid.New()

	res,err:=app.DB.Exec("INSERT INTO points_data (id, name, latitude, longitude) VALUES ($1,'TestPoint', 12.34, 56.78)",id)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(res)

	return id

}

func TestCreatePolygon(t *testing.T) {
	clearPolygonsTable()
	payload := []byte(`{"name":"TestPolygon", "geojson": "{\"type\": \"Polygon\", \"coordinates\": [[[30, 10], [40, 40], [20, 40], [10, 20], [30, 10]]]}"}`)
	req, _ := http.NewRequest("POST", "/polygons", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TestPolygon" {
		t.Errorf("Expected name to be 'TestPolygon'. Got '%v'", m["name"])
	}
		
}
func TestGetPolygon(t *testing.T) {
	clearPolygonsTable()
	id:=addPolygon()

	req,_:=http.NewRequest("GET", fmt.Sprintf("/polygons/%s", id),nil)

	response:=executeRequest(req)

	checkResponseCode(t,http.StatusOK,response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "TestPolygon" {
		t.Errorf("Expected name to be 'TestPolygon'. Got '%s'", m["name"])
	}
	if m["geojson"] != "{\"type\": \"Polygon\", \"coordinates\": [[[30, 10], [40, 40], [20, 40], [10, 20], [30, 10]]]}" {
		t.Errorf("Expected geojson to be {\"type\": \"Polygon\", \"coordinates\": [[[30, 10], [40, 40], [20, 40], [10, 20], [30, 10]]]} \n found : %v",m["geojson"])
	}
}

func TestUpdatePolygon(t *testing.T) {
	clearPolygonsTable()
	id:=addPolygon()

	payload:=[]byte(`{"name":"UpdatePolygon", "geojson": "{\"type\": \"Polygon\", \"coordinates\": [[[40, 10], [50, 40], [10, 40], [30, 20], [30, 10]]]}"}`)

	req,_:=http.NewRequest("PUT",fmt.Sprintf("/polygons/%s",id),bytes.NewBuffer(payload))

	response := executeRequest(req)

	checkResponseCode(t,http.StatusOK,response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(),&m)

	if m["name"] != "UpdatePolygon"{
		t.Errorf("Expected name to be UpdatePolygon found '%s'",m["name"])
	}
	if m["geojson"] != "{\"type\": \"Polygon\", \"coordinates\": [[[40, 10], [50, 40], [10, 40], [30, 20], [30, 10]]]}"{
		t.Errorf("Expected geojson to be {\"type\": \"Polygon\", \"coordinates\": [[[40, 10], [50, 40], [10, 40], [30, 20], [30, 10]]]} \n found '%v'",m["geojson"])
	}
}
func addPolygon() uuid.UUID{
	id:=uuid.New()
	app.DB.Exec("INSERT INTO polygons_data (id, name, geojson) VALUES ($1, 'TestPolygon', '{\"type\": \"Polygon\", \"coordinates\": [[[30, 10], [40, 40], [20, 40], [10, 20], [30, 10]]]}')",id)

	return id
}

func clearPointsTable(){
	app.DB.Exec("DELETE FROM points_data")
}

func clearPolygonsTable(){
	app.DB.Exec("DELETE FROM polygons_data")
}


func executeRequest(req *http.Request) *httptest.ResponseRecorder{
	rr:=httptest.NewRecorder()
	app.Router.ServeHTTP(rr,req)

	return rr
}
func checkResponseCode(t *testing.T, expected,actual int){
	if expected != actual{
		t.Errorf("expected response code %d, got %d", expected, actual)
	}
}