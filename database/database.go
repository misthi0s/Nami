package NamiDatabase

import (
	"database/sql"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

func CreateInitialDatabase() {
	os.Mkdir("./db", 0755)
	os.Create("./db/Nami.db")

	database, _ := sql.Open("sqlite", "./db/Nami.db")
	database.Exec("CREATE TABLE `sessions` (`id` INTEGER PRIMARY KEY, `uuid` VARCHAR(38), `checkin_time` VARCHAR(40), `hostname` VARCHAR(40), `username` VARCHAR(40), `implant` VARCHAR(40))")
	defer database.Close()
}

func AddToDatabase(uuid string, hostname string, username string, implant string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	database, _ := sql.Open("sqlite", "./db/Nami.db")
	statement, _ := database.Prepare("INSERT INTO sessions (uuid, hostname, username, implant, checkin_time) VALUES (?,?,?,?,?)")
	statement.Exec(uuid, hostname, username, implant, string(currentTime))
	defer database.Close()
}

func QuerySessions() *sql.Rows {
	database, _ := sql.Open("sqlite", "./db/Nami.db")
	results, _ := database.Query("SELECT id, uuid, hostname, username, implant, checkin_time FROM `sessions`")
	defer database.Close()
	return results
}

func QuerySessionById(id int) *sql.Row {
	database, _ := sql.Open("sqlite", "./db/Nami.db")
	statement, _ := database.Prepare("SELECT uuid FROM `sessions` where id = ?")
	results := statement.QueryRow(id)
	return results
}

func UpdateCheckInTime(uuid string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	database, _ := sql.Open("sqlite", "./db/Nami.db")
	statement, _ := database.Prepare("UPDATE sessions SET checkin_time = ? WHERE uuid = ?")
	statement.Exec(string(currentTime), uuid)
	defer database.Close()
}

func CleanUpDatabase() {
	// Need to save UUIDs to delete temporarily in slice to prevent database lock
	deleteUUIDs := []string{}
	currentTime := time.Now().Local()
	database, _ := sql.Open("sqlite", "./db/Nami.db")
	results, _ := database.Query("SELECT uuid,checkin_time FROM `sessions`")
	var uuid string
	var checkin_time string
	for results.Next() {
		results.Scan(&uuid, &checkin_time)
		checkin_formmatted, _ := time.ParseInLocation("2006-01-02 15:04:05", checkin_time, time.Local)
		diff := currentTime.Sub(checkin_formmatted)
		if diff.Hours() > 24 {
			deleteUUIDs = append(deleteUUIDs, uuid)
		}
	}
	statement, _ := database.Prepare("DELETE FROM `sessions` where uuid = ?")
	for _, uuid := range deleteUUIDs {
		statement.Exec(uuid)
	}
	defer database.Close()
}

func main() {
	CreateInitialDatabase()
}
