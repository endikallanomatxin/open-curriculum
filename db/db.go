package db

import (
    "fmt"
	"os"
	"database/sql"
	
    _ "github.com/lib/pq"
)

var db *sql.DB

var connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))


func InitializeDB() {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    createTables := `
    CREATE TABLE IF NOT EXISTS groups (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        group_id INT,
        FOREIGN KEY (group_id) REFERENCES groups(id)
    );

    CREATE TABLE IF NOT EXISTS units (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        group_id INT,
        FOREIGN KEY (group_id) REFERENCES groups(id)
    );

    CREATE TABLE IF NOT EXISTS dependencies (
        unit_id INT,
        depends_on_id INT,
        PRIMARY KEY (unit_id, depends_on_id),
        FOREIGN KEY (unit_id) REFERENCES units(id),
        FOREIGN KEY (depends_on_id) REFERENCES units(id)
    );`

    // Ejecutar los comandos SQL
    _, err = db.Exec(createTables)
    if err != nil {
        panic(err)
    }

    fmt.Println("Tablas creadas con éxito")
}

type Unit struct {
    ID          int
    Name        string
    Description string
}

func GetUnits() []Unit {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, name, description FROM units")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    units := []Unit{}
    for rows.Next() {
        var u Unit
        err := rows.Scan(&u.ID, &u.Name, &u.Description)
        if err != nil {
            panic(err)
        }
        units = append(units, u)
    }

    return units
}

func CreateUnit(u Unit) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    _, err = db.Exec("INSERT INTO units (name, description) VALUES ($1, $2)", u.Name, u.Description)
    if err != nil {
        panic(err)
    }
}

func GetUnit(id int) Unit {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    var u Unit
    err = db.QueryRow("SELECT id, name, description FROM units WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Description)
    if err != nil {
        panic(err)
    }

    return u
}

func DeleteUnit(id int) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM units WHERE id = $1", id)
    if err != nil {
        panic(err)
    }
}
