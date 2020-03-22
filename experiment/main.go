package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"time"

	"github.com/go-redis/redis"
)

var (
	errCount          uint64 = 0
	normalCount       uint64 = 0
	totalResponseTime int64  = 0
)

func main() {

	forever := make(chan struct{})
	redis := getInstance()

	log.Println("start work")
	for i := 0; i < 50; i++ {
		work(redis)
	}

	<-forever
}

func work(s *Service) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			generate(s)
		}
	}
}

type user struct {
	ID  int   `json:"id"`
	Bet int64 `json:"bet"`
}

func generate(s *Service) {

	//r := time.Now().Second()

	randCompany := []string{"Com1", "Com2"}

	randSystem := []string{"Sys1", "Sys2", "Sys3"}

	randWebsite := []string{"Web1", "Web2", "Web3", "Web4"}

	//ra := rand.

	ran := rand.Int63n(1000)

	randUser := user{
		ID:  rand.Intn(10),
		Bet: ran,
	}

	date := time.Now().Format("2006010215")

	//com
	comKey := fmt.Sprintf("%s:%s", randCompany[rand.Intn(2)], date)
	d := make(map[string]interface{})

	d[comKey] = randUser.Bet

	exsit, _ := s.Client.HExists("com", comKey).Result()

	if exsit {
		s.Client.HIncrBy("com", comKey, randUser.Bet)
	} else {
		s.Client.HMSet("com", d)
	}

	//sys
	sysKey := fmt.Sprintf("%s:%s:%s", randCompany[rand.Intn(2)], randSystem[rand.Intn(3)], date)

	d = make(map[string]interface{})

	d[sysKey] = randUser.Bet

	exsit, _ = s.Client.HExists("sys", sysKey).Result()

	if exsit {
		s.Client.HIncrBy("sys", sysKey, randUser.Bet)
	} else {
		s.Client.HMSet("sys", d)
	}

	//web
	webKey := fmt.Sprintf("%s:%s:%s:%s", randCompany[rand.Intn(2)], randSystem[rand.Intn(3)], randWebsite[rand.Intn(4)], date)

	d = make(map[string]interface{})

	d[webKey] = randUser.Bet

	exsit, _ = s.Client.HExists("web", webKey).Result()

	if exsit {
		s.Client.HIncrBy("web", webKey, randUser.Bet)
	} else {
		s.Client.HMSet("web", d)
	}

	//user
	userKey := fmt.Sprintf("%s:%s:%s:%d:%s", randCompany[rand.Intn(2)], randSystem[rand.Intn(3)], randWebsite[rand.Intn(4)], randUser.ID, date)
	userData, _ := json.Marshal(randUser)
	r := s.Client.LPush(userKey, userData)
	if r.Err() != nil {
		log.Fatal(r.Err())
	}
}

// Service ...
type Service struct {
	Client *redis.Client
}

type redisConfig struct {
	Addr     string `json:"addr"`
	PASSWORD string `json:"password"`
	INDEX    int    `json:"index"`
}

// getInstance return a redis service.
func getInstance() *Service {

	return &Service{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       1,  // use default DB
		}),
	}
}
