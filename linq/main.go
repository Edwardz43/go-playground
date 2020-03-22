package main

import (
	"fmt"
	"log"

	//. "github.com/ahmetb/go-linq"
	"github.com/go-redis/redis"
)

type result struct {
	bet int
}

func main() {
	rr := getInstance()

	t := make([]string, 7)

	for i := 0; i < 7; i++ {
		t[i] = fmt.Sprintf("Com1:20200321%2d", i+17)
	}

	//.Format("2006010215")

	cmd := redis.NewStringIntMapCmd("HMGet", t)
	rr.Client.Process(cmd)
	//redis.Client.

	//r := redis.Client.HMGet("com", t[0], t[1], t[2], t[3], t[4], t[5], t[6])
	cmd.Result()
	// if r.Err() != nil {
	// 	log.Fatalln(r.Err())
	// }

	var res int64 = 0

	//tmp, _ := r.Result()
	tmp, _ := cmd.Result()

	for _, v := range tmp {
		//`if i, e := strconv.Atoi(v.(string)); e == nil {
		res += v
		//}
	}
	log.Println(res)
	//log.Println(From(res).SumInts())
	// From(cars).Where(func(c interface{}) bool {
	// 	return c.(Car).year >= 2015
	// }).Select(func(c interface{}) interface{} {
	// 	return c.(Car).owner
	// }).ToSlice(&owners)

	// for _, owner := range owners {
	// 	fmt.Printf("owner:%s\n", owner)
	// }
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
