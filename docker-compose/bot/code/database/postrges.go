package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"github.com/wcharczuk/go-chart"
	"golang.org/x/net/context"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var sslmode = os.Getenv("SSLMODE")
var dbname = os.Getenv("DBNAME")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)


func CollectData(userID int, chatID int64, message string, answer string)  error {

    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return err
    }
    defer db.Close()

    query := `INSERT INTO user_stats (user_id, chat_id, message, answer) VALUES ($1, $2, $3, $4)`

    _, err = db.Exec(query, userID, chatID, message, answer)
    if err != nil {
        log.Println("error inserting data:", err)
        return err
    }

    return nil
}

func CreateTable() error {
    db, err := sql.Open("postgres", dbInfo)
    if err != nil {
        return err
    }
    defer db.Close()

    query := `CREATE TABLE IF NOT EXISTS user_stats (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        chat_id BIGINT NOT NULL,
        message TEXT,
        answer TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

    _, err = db.Exec(query)
    if err != nil {
        log.Println("error creating table:", err)
        return err
    }

    return nil
}


func GetStatistics(chat_id int64) (string, string, error){

    conn, err := pgx.Connect(context.Background(), dbInfo)
    if err != nil {
        return "","", err
    }

    sqlQueryALL := `SELECT COUNT(*) AS valRequests FROM public.user_stats WHERE chat_id = $1`
    sqlQueryValFirstDate := `SELECT MIN(created_at) AS valRequestsFirst, COUNT(*) AS valRequests FROM user_stats WHERE chat_id = $1`
    sqlQueryValLastDate := `SELECT MAX(created_at) AS valRequestsLast, COUNT(*) AS valRequests FROM user_stats WHERE chat_id = $1`
    sqlQueryMessage := `SELECT message FROM public.user_stats WHERE chat_id = $1`
    tx, _ := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())

    var dateFirst, dateLast time.Time
    var countFirst, countLast, countAll int
    var cities []string
    var city string

    rowFirst := tx.QueryRow(context.Background(), sqlQueryValFirstDate, chat_id)
    err = rowFirst.Scan(&dateFirst, &countFirst)
    if err != nil {
        log.Println("sqlQueryValFirstDate", err)
        return "","", err
    }

   rowLast := tx.QueryRow(context.Background(), sqlQueryValLastDate, chat_id)
    err = rowLast.Scan(&dateLast, &countLast)
    if err != nil {
        log.Println("sqlQueryValLastDate", err)
        return "","", err
    }
  
    row := tx.QueryRow(context.Background(), sqlQueryALL, chat_id)
    err = row.Scan(&countAll)
    if err != nil {
        log.Println("sqlQueryALL", err)
        return "","", err
    }
    
   //statistics := fmt.Sprintf("In total, you have made %d requests. \nThe first request was made on %v, you made %d requests that day. \nThe last request was made on %v, you made %d requests that day.",countAll,dateFirst.Format(time.ANSIC),countFirst,dateLast.Format(time.ANSIC),countLast)
    statistics := fmt.Sprintf("Всего вы сделали %d запросов. \nПервый запрос был сделан %v, в тот день вы сделали %d запросов. \nПоследний запрос был %v, вы сделали %d запросов.",countAll,format_time(dateFirst),countFirst,format_time(dateLast),countLast)
    //Далее построение piecharm
    rows, err := tx.Query(context.Background(),sqlQueryMessage,chat_id)
    if err != nil {
        log.Println("sqlQueryMessage", err)
        return statistics,"", err
    }
    defer rows.Close()

    cities = make([]string, 0)

    for rows.Next() {
        err := rows.Scan(&city)
        if err != nil {
            log.Println("sqlQueryMessage", err)
            return statistics,"", err
        }
        cities = append(cities, city)
    }

    if err := rows.Err(); err != nil {
        log.Println("sqlQueryMessage", err)
        return statistics,"", err
    }

    counts := make(map[string]int)
    for _, city := range cities {
        counts[city]++
    }

    data := make([]chart.Value, 0, len(counts))
    for city, count := range counts {
        data = append(data, chart.Value{Label: city, Value: float64(count)})
    }

    pie := chart.PieChart{
        Title: "Список городов",
        Values: data,
    }

    uuid := uuid.NewV4()
    filename := fmt.Sprintf("%s.png", uuid)

    f, _ := os.Create(filename)
    defer f.Close()
    pie.Render(chart.PNG, f)

    return statistics, filename, nil
}

func format_time(t time.Time) string {
    return t.Format("2006.01.02-15.04.05")
  }

