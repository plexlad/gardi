package main

import (
	"net/http"
	"time"

	//"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/plexlad/gardi/server/lib"
)

const (
	CollectionInstances = "instances"
	CollectionSchemas   = "schemas"
)

type NewSchemaRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NewInstanceRequest struct {
	Name string `json:"name"`
	Description string `json:"description"`
	SchemaID string `json:"schema_id"`
}

func main() {
	db := NewJsonDB("./data")

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.POST("/:user/schemas/new", func(c echo.Context) error {
		user := c.Param("user")

		var req NewSchemaRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
		}

		instanceID := uuid.New().String()
		schema := lib.Schema{
			ID:          instanceID,
			Version:     1,
			UserVersion: 1,
			Name:        req.Name,
			Description: req.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := db.Set(CollectionInstances, user, instanceID, schema)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid schema",
			})
		}

		return c.JSON(http.StatusOK, schema)
	})

	router.POST("/:user/instances/new", func(c echo.Context) error {
		user := c.Param("user")

		var req NewInstanceRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
		}

		var schema lib.Schema
		err := db.Get(CollectionSchemas, user, req.SchemaID, &schema)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid schema",
			})
		}

		instanceID := uuid.New().String()
		instance := lib.Instance{
			ID: instanceID,
			SchemaID: req.SchemaID,
			Name: req.Name,
			Description: req.Description,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		return c.JSON(http.StatusOK, instance)
		
	})

	router.POST("/:user/instances/save", func(c echo.Context) error {
		user := c.Param("user")

		var req lib.Instance

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
		}

		db.Set(CollectionInstances, user, req.ID, req)

		return c.String(http.StatusOK, "instance saved")
	})

	router.POST("/:user/schemas/save", func(c echo.Context) error {
		user := c.Param("user")

		var req lib.Schema

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
		}

		db.Set(CollectionSchemas, user, req.ID, req)

		return c.String(http.StatusOK, "schema saved")
	})

	router.GET("/:user/schemas", func(c echo.Context) error {
		user := c.Param("user")

		schemaIDs, err := db.List(CollectionSchemas, user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		schemas := []lib.Schema{}
		for _, id := range schemaIDs {
			var schema lib.Schema
			if err := db.Get(CollectionSchemas, user, id, &schema); err != nil {
				continue
			}
			schemas = append(schemas, schema)
		}

		return c.JSON(http.StatusOK, schemas)
	})

	router.GET("/:user/instances", func(c echo.Context) error {
		user := c.Param("user")

		instanceIDs, err := db.List("instances", user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		instances := []lib.Instance{}
		for _, id := range instanceIDs {
			var instance lib.Instance
			if err := db.Get(CollectionInstances, user, id, &instance); err != nil {
				continue
			}
			instances = append(instances, instance)
		}

		return c.JSON(http.StatusOK, instances)
	})

	// run server with error checking
	router.Logger.Fatal(router.Start(":5497"))
}
