package mysql
import (
	"database/sql"
	// "time"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"encoding/json"
	"fmt"
)
// Config is mysql config.
type Config struct {
	DSN         string
	Active      int            // pool
	Idle        int            // pool
	// IdleTimeout time.Time    // connect max life time.
  }

func getConf() (*Config,error){
	file, err := os.Open("./conf/database.json")
	if err!=nil{
		fmt.Println(err)
	}
  	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		return &Config{},err
	}
	return &conf,nil
}
// NewMysql initialize mysql connection .
func NewMysql() (db *sql.DB) {
	// TODO add query exec and transaction timeout .
	c,err:=getConf()
	if err!=nil{
		fmt.Println(err)
	}
	db, err = Open(c)
	if err != nil {
		panic(err)
	}
	return db
}
func Open(c *Config) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.DSN)
	if err != nil {
		return nil, err
	}
	// db.SetMaxOpenConns(c.Active)
	// db.SetMaxIdleConns(c.Idle)
	// db.SetConnMaxLifetime(time.Hour)
	return db, nil
}