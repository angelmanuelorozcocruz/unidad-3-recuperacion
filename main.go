package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// variable globales
var DB *sql.DB

// Estructura Model Automovil
type automovil struct {
	Id          int    `json:"id"`
	Precio      string `json:"precio"`
	Descripcion string `json:"descripcion"`
	Marca       string `json:"marca"`
	Modelo      string `json:"modelo"`
	Creado      string `json:"creado"`
}

func main() {
	r := gin.Default()

	DB, err := sql.Open("mysql", "root:Angel1997.@/automotores")
	if err != nil {
		panic(err.Error())
	}
	defer DB.Close()

	//Mainpage to start
	r.GET("/homepage", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "WELCOME TO THE REGISTER OF CARS with APIREST",
		})
	})

	//Insert inside of automoviles
	r.POST("/automoviles", func(c *gin.Context) {
		automovil := automovil{}
		err := c.ShouldBindJSON(&automovil)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		stmt, err := DB.Query("INSERT INTO automoviles (`precio`,`descripcion`,`marca`,`modelo`,`creado`) VALUES (?,?,?,?,?)", automovil.Precio, automovil.Descripcion, automovil.Marca, automovil.Modelo, automovil.Creado)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		c.JSON(200, automovil)
	})

	//Select all automoviles
	r.GET("/automoviles", func(c *gin.Context) {
		rows, err := DB.Query("SELECT * FROM automoviles")
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var automoviles []automovil
		for rows.Next() {
			var automovil automovil
			rows.Scan(&automovil.Id, &automovil.Precio, &automovil.Descripcion, &automovil.Marca, &automovil.Modelo, &automovil.Creado)
			automoviles = append(automoviles, automovil)
		}

		c.JSON(200, automoviles)
	})

	//Select one only automovil
	r.GET("/automoviles/:id", func(c *gin.Context) {
		id := c.Param("id")

		var automovil automovil // crear estructura donde se guardara el json de usuario
		err := DB.QueryRow("SELECT * FROM automoviles WHERE id=?", id).Scan(&automovil.Id, &automovil.Precio, &automovil.Descripcion, &automovil.Marca, &automovil.Modelo, &automovil.Creado)
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"msg": "user not found"})
			return
		}
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, automovil)

	})

	//Update one only automovil
	r.PUT("/automoviles/:id", func(c *gin.Context) {
		var automovil automovil
		id := c.Param("id")

		stmt, err := DB.Prepare("UPDATE automoviles SET precio=?, descripcion=?, marca=?, modelo=?, creado=? WHERE id=?")
		if err != nil {
			panic(err)
		}

		// execute
		res, err := stmt.Exec(automovil.Precio, automovil.Descripcion, automovil.Marca, automovil.Modelo, automovil.Creado, id)
		if err != nil {
			fmt.Println(err.Error())
		}

		a, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println(a)

		c.JSON(200, automovil)
	})

	//Delete one only automovil
	r.DELETE("/automoviles/:id", func(c *gin.Context) {
		id := c.Param("id")
		DB.QueryRow("DELETE FROM automoviles WHERE id=?", id)

		c.JSON(200, id)
	})

	r.Run(":8069") // Listening for windows "localhost:8069"
}
