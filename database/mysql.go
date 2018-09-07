package database
import (
	"database/sql"
	// "time"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"encoding/json"
)
// Config is mysql config.
type MysqlConfig struct {
	DSN         string
	Database	string
	Charset		string
	Active      int            // pool
	Idle        int            // pool
	// IdleTimeout time.Time    // connect max life time.
  }

func getMysqlConf() (*MysqlConfig, error){
	file, err := os.Open("./conf/mysql.json")
	if err!=nil{
		return nil,err
	}
  	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := MysqlConfig{}
	err = decoder.Decode(&conf)
	if err != nil {
		return nil,err
	}
	return &conf,nil
}
// NewMysql initialize mysql connection .
func NewMysql(database ...string) (db *sql.DB ,err error) {
	// TODO add query exec and transaction timeout .
	c,err:=getMysqlConf()
	if err!=nil{
		return nil,err
	}
	if database != nil{
		c.Database=database[0]
	}
	db, err = Open(c)
	if err != nil {
		panic(err)
	}
	return db,nil
}
func Open(c *MysqlConfig) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", c.DSN+"/"+c.Database+"?charset="+c.Charset)
	if err != nil {
		return nil, err
	}
	// db.SetMaxOpenConns(c.Active)
	// db.SetMaxIdleConns(c.Idle)
	// db.SetConnMaxLifetime(time.Hour)
	return db, nil
}