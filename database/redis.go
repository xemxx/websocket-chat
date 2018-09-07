package database

import (
	"github.com/go-redis/redis"
	"os"
	"encoding/json"
)

type RedisConf struct{
	Addr 		string
    Password 	string
    DB  		int
}


func getRedisConf() (*RedisConf,error){
	file, err := os.Open("./conf/redis.json")
	if err!=nil{
		return nil,err
	}
  	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := &RedisConf{}
	err = decoder.Decode(conf)
	if err != nil {
		return nil,err
	}
	return conf,nil
}

func NewRedis() (*redis.Client,error){
	conf,_ :=getRedisConf()
	client := redis.NewClient(&redis.Options{
        Addr: conf.Addr,
        Password: conf.Password,
        DB: conf.DB,  	       
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil,err
	}
	return client,nil
}