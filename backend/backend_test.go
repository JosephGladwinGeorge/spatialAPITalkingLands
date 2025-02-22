package backend_test

import (
	"os"
	"spatialDB/backend"
	"testing"
)

func TestMain(m *testing.M) {
	app := backend.App{}
	app.Initialize()
	ensureTablesExist()
	code:= m.Run()

	clearPointsTable()
	clearPolygonsTable()

	os.Exit(code)
}