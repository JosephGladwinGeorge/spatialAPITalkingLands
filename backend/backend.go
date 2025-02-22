package backend

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type App struct{
	DB *sqlx.DB
	Port string
	Router *gin.Engine
}

func (a *App) Initialize(){
	DB,err:=initDB()
	if err!=nil{
		log.Fatal(err.Error())
	}

	Router:=gin.Default()

	a.DB=DB
	a.Router=Router
	a.initializeRoutes()
}

func (a *App) initializeRoutes(){
	a.Router.POST("/points",a.newPoint)
	a.Router.GET("/points/:id",a.getPoint)
	a.Router.PUT("/points/:id",a.updatePoint)

	a.Router.POST("/polygons",a.newPolygon)
	a.Router.GET("/polygons/:id",a.getPolygon)
	a.Router.PUT("/polygons/:id",a.updatePolygon)
}

func (a *App) newPoint(c *gin.Context){
	var p Point

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := p.createPoint(a.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (a *App) getPoint(c *gin.Context){
	var p Point

	p.ID,_=uuid.Parse(c.Param("id"))

	if err:=p.fetchPoint(a.DB); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (a *App) updatePoint(c *gin.Context){
	var p Point

	if err:=c.ShouldBindBodyWithJSON(&p);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}

	if err:=p.updatePoint(a.DB); err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,p)
}

func (a *App) newPolygon(c *gin.Context){
	var pg Polygon

	if err := c.ShouldBindJSON(&pg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := pg.createPolygon(a.DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pg)
}

func (a *App) getPolygon(c *gin.Context){
	var pg Polygon

	pg.ID,_ = uuid.Parse(c.Param("id"))

	if err:=pg.fetchPolygon(a.DB); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pg)
}

func (a *App) updatePolygon(c *gin.Context){
	var pg Polygon

	if err:=c.ShouldBindBodyWithJSON(&pg);err!=nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}

	if err:=pg.updatePolygon(a.DB); err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,pg)
}


func (a *App) Run(){
	a.Router.Run(":" + a.Port)
}