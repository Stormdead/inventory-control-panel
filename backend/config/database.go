package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Stormdead/inventory-control-panel/backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	log.Println("Conexión a MySQL exitosa")

	//Auto-Migration: Crea las tablas automáticamente
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Movement{},
	)
	if err != nil {
		log.Fatal("Error en la migración:", err)
	}

	log.Println("Tablas creadas/actualizadas correctamente")
}
