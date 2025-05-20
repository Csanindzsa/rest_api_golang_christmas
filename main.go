package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	_ "rest_api_golang_christmas/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Partnerek struct {
	Id         int    `json:"id"`
	PartnerNev string `json:"partnernev"`
	Email      string `json:"email"`
}

type Szamlak struct {
	Id         int    `json:"id"`
	Szamlaszam string `json:"szamlaszam"`
	Tetelszam  int    `json:"tetelszam"`
	Megjegyzes string `json:"megjegyzes"`
	PartnerId  int    `json:"partnerid"`
}

var db *sql.DB

// @title Számlák és Partnerek kezelése
// @version 1.0
// @description Számlák és Partnerek kezelése.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:10090
// @BasePath /

func createTables() error {
	partnerstable := `
    CREATE TABLE IF NOT EXISTS Partnerek (
        Id INT PRIMARY KEY AUTO_INCREMENT,
        PartnerNev VARCHAR(50),
        Email VARCHAR(50)
    );`

	receiptstable := `
    CREATE TABLE IF NOT EXISTS Szamlak (
        Id INT PRIMARY KEY AUTO_INCREMENT,
        Szamlaszam VARCHAR(50),
        Tetelszam INT,
        Megjegyzes VARCHAR(50),
		PartnerId INT,
        FOREIGN KEY (PartnerId) REFERENCES Partnerek(Id)
    );`

	_, err := db.Exec(partnerstable)
	if err != nil {
		return fmt.Errorf("hiba a partnerek tábla létrehozásában: %v", err)
	}

	_, err = db.Exec(receiptstable)
	if err != nil {
		return fmt.Errorf("hiba a szamlak tábla létrehozásában: %v", err)
	}

	return nil
}

// @Summary Partnerek lekérdezése
// @Description Partnerek lekérdezése
// @Tags partnerek
// @Accept  json
// @Produce  json
// @Success 200 {array} Partnerek
// @Router /partnerek [get]
func getPartners(c *gin.Context) {
	rows, err := db.Query("SELECT Id, PartnerNev, Email FROM Partnerek")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	var partners []Partnerek
	for rows.Next() {
		var partner Partnerek
		if err := rows.Scan(&partner.Id, &partner.PartnerNev, &partner.Email); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		partners = append(partners, partner)
	}

	c.IndentedJSON(http.StatusOK, partners)

}

// @Summary Partner hozzáadása
// @Description Partner hozzáadása
// @Tags partnerek
// @Accept  json
// @Produce  json
// @Param partnerek body Partnerek true "Partnerek to add"
// @Success 201 {object} Partnerek
// @Router /partnerek [post]
func postPartner(c *gin.Context) {
	var newPartner Partnerek
	if err := c.BindJSON(&newPartner); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás lekérdezés"})
		return
	}

	_, err := db.Exec("INSERT INTO Partnerek (PartnerNev, Email) VALUES (?, ?)", newPartner.PartnerNev, newPartner.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPartner)
}

// @Summary Egy partner lekérdezése
// @Description Egy partner lekérdezése
// @Tags partnerek
// @Accept  json
// @Produce  json
// @Param id path int true "Partner ID"
// @Success 200 {object} Partnerek
// @Router /partnerek/{id} [get]
func getPartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	var partner Partnerek
	err = db.QueryRow("SELECT Id, PartnerNev, Email FROM Partnerek WHERE Id = ?", id).Scan(&partner.Id, &partner.PartnerNev, &partner.Email)
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partner nem találva"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, partner)

}

// @Summary Partner információ frissítése
// @Description Partner információ frissítése
// @Tags partnerek
// @Accept  json
// @Produce  json
// @Param id path int true "Partner ID"
// @Param partnerek body Partnerek true "Partner to update"
// @Success 200 {object} Partnerek
// @Router /partnerek/{id} [put]
func putPartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	var updatedPartner Partnerek
	if err := c.BindJSON(&updatedPartner); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás lekérdezés"})
		return
	}

	_, err = db.Exec("UPDATE Partnerek SET PartnerNev = ?, Email = ? WHERE Id = ?", updatedPartner.PartnerNev, updatedPartner.Email, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedPartner)

}

// @Summary Partner törlése
// @Description Partner törlése
// @Tags partnerek
// @Accept  json
// @Produce  json
// @Param id path int true "Partner ID"
// @Success 200 {object} gin.H
// @Router /partnerek/{id} [delete]
func deletePartner(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Tranzakció hiba: " + err.Error()})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// First delete linked invoices
	_, err = tx.Exec("DELETE FROM Szamlak WHERE PartnerId = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Kapcsolódó számlák törlése sikertelen: " + err.Error()})
		return
	}

	// Then delete the partner
	result, err := tx.Exec("DELETE FROM Partnerek WHERE Id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Partner törlése sikertelen: " + err.Error()})
		return
	}

	// Check if any row was affected (if partner exists)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if rowsAffected == 0 {
		tx.Rollback()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Partner nem található"})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Tranzakció lezárása sikertelen: " + err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Partner és kapcsolódó számlák törölve"})
}

// @Summary Számlák lekérdezése
// @Description Számlák lekérdezése
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Success 200 {array} Szamlak
// @Router /szamlak [get]
func getReceipts(c *gin.Context) {
	rows, err := db.Query("SELECT Id, Szamlaszam, Tetelszam, Megjegyzes, PartnerId FROM Szamlak")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	var receipts []Szamlak
	for rows.Next() {
		var receipt Szamlak
		if err := rows.Scan(&receipt.Id, &receipt.Szamlaszam, &receipt.Tetelszam, &receipt.Megjegyzes, &receipt.PartnerId); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		receipts = append(receipts, receipt)
	}

	c.IndentedJSON(http.StatusOK, receipts)

}

// @Summary Számla hozzáadása
// @Description Számla hozzáadása
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Param szamlak body Szamlak true "Receipt to add"
// @Success 201 {object} Szamlak
// @Router /szamlak [post]
func postReceipt(c *gin.Context) {
	var newReceipt Szamlak
	if err := c.BindJSON(&newReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás lekérés"})
		return
	}

	_, err := db.Exec("INSERT INTO Szamlak (Szamlaszam, Tetelszam, Megjegyzes, PartnerId) VALUES (?, ?, ?, ?)", newReceipt.Szamlaszam, newReceipt.Tetelszam, newReceipt.Megjegyzes, newReceipt.PartnerId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newReceipt)
}

// @Summary Egy számla lekérdezése
// @Description Egy számla lekérdezése
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Param id path int true "Receipt ID"
// @Success 200 {object} Szamlak
// @Router /szamlak/{id} [get]
func getReceipt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	var receipt Szamlak
	err = db.QueryRow("SELECT Id, Szamlaszam, Tetelszam, Megjegyzes, PartnerId FROM Szamlak WHERE Id = ?", id).Scan(&receipt.Id, &receipt.Szamlaszam, &receipt.Tetelszam, &receipt.Megjegyzes, &receipt.PartnerId)
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Számla nem található"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, receipt)

}

// @Summary Számla frissítése
// @Description Számla frissítése
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Param id path int true "Receipt ID"
// @Param szamlak body Szamlak true "Receipt to update"
// @Success 200 {object} Szamlak
// @Router /szamlak/{id} [put]
func putReceipt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	var updatedReceipt Szamlak
	if err := c.BindJSON(&updatedReceipt); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás lekérés"})
		return
	}

	_, err = db.Exec("UPDATE Szamlak SET Szamlaszam = ?, Tetelszam = ?, Megjegyzes = ?, PartnerId = ?  WHERE Id = ?", updatedReceipt.Szamlaszam, updatedReceipt.Tetelszam, updatedReceipt.Megjegyzes, updatedReceipt.PartnerId, id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedReceipt)

}

// @Summary Számla törlése
// @Description Számla törlése
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Param id path int true "Receipt ID"
// @Success 200 {object} gin.H
// @Router /szamlak/{id} [delete]
func deleteReceipt(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Hibás ID"})
		return
	}

	_, err = db.Exec("DELETE FROM Szamlak WHERE Id = ?", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Számla törölve"})

}

// @Summary Számla keresése név alapján
// @Description Számla keresése név alapján
// @Tags szamlak
// @Accept  json
// @Produce  json
// @Param szamlaszam query string true "Receipt number"
// @Success 200 {array} Szamlak
// @Router /szamlak/search [get]
func searchReceipt(c *gin.Context) {
	receiptnumber := c.Query("szamlaszam")
	if receiptnumber == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Számlaszám paraméter szükséges"})
		return
	}

	rows, err := db.Query("SELECT Id, Szamlaszam, Tetelszam, Megjegyzes, PartnerId FROM Szamlak WHERE Szamlaszam LIKE ?", "%"+receiptnumber+"%")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	var receipts []Szamlak
	for rows.Next() {
		var receipt Szamlak
		if err := rows.Scan(&receipt.Id, &receipt.Szamlaszam, &receipt.Tetelszam, &receipt.Megjegyzes, &receipt.PartnerId); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		receipts = append(receipts, receipt)
	}

	c.IndentedJSON(http.StatusOK, receipts)

}

func main() {
	var err error
	// Connect to MySQL server without specifying a database
	connString := "root:@tcp(localhost:3306)/"
	db, err = sql.Open("mysql", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}

	// Create the database with UTF-8 collation if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS nyilvantartas CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatal("Error creating database: ", err.Error())
	}

	// Select the database for use
	_, err = db.Exec("USE nyilvantartas")
	if err != nil {
		log.Fatal("Error selecting database: ", err.Error())
	}

	// Create tables
	err = createTables()
	if err != nil {
		log.Fatal("Error creating tables: ", err.Error())
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/partnerek", getPartners)
	router.POST("/partnerek", postPartner)
	router.GET("/partnerek/:id", getPartner)
	router.PUT("/partnerek/:id", putPartner)
	router.DELETE("/partnerek/:id", deletePartner)
	router.GET("/szamlak", getReceipts)
	router.POST("/szamlak", postReceipt)
	router.GET("/szamlak/:id", getReceipt)
	router.PUT("/szamlak/:id", putReceipt)
	router.DELETE("/szamlak/:id", deleteReceipt)
	router.GET("/szamlak/search", searchReceipt)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":10090")
}
