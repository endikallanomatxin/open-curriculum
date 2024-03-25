package db

import (
    "fmt"
	"os"
	"database/sql"
	
    _ "github.com/lib/pq"
)

func InitializeDB() {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    createTables := `
    CREATE TABLE IF NOT EXISTS grupos (
        id SERIAL PRIMARY KEY,
        nombre VARCHAR(255) NOT NULL,
        descripcion TEXT,
        FOREIGN KEY (grupo_id) REFERENCES grupos(id)
    );

    CREATE TABLE IF NOT EXISTS unidades (
        id SERIAL PRIMARY KEY,
        nombre VARCHAR(255) NOT NULL,
        descripcion TEXT,
        grupo_id INT,
        FOREIGN KEY (grupo_id) REFERENCES grupos(id)
    );

    CREATE TABLE IF NOT EXISTS dependencias (
        unidad_id INT,
        depende_de_unidad_id INT,
        PRIMARY KEY (unidad_id, depende_de_unidad_id),
        FOREIGN KEY (unidad_id) REFERENCES unidades(id),
        FOREIGN KEY (depende_de_unidad_id) REFERENCES unidades(id)
    );`

    // Ejecutar los comandos SQL
    _, err = db.Exec(createTables)
    if err != nil {
        panic(err)
    }

    fmt.Println("Tablas creadas con éxito")
}