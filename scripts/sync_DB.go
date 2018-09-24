package main

import(
	"gopkg.in/mgo.v2"
	"github.com/go-redis/redis"
	"API-Mutant/models"
)

var dnas = getSession().DB("dnas").C("dnas")

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://127.0.0.1")

	if(err != nil){
		panic(err) 
	}

	return session
}

func main(){
	var results []models.Dna

	err := dnas.Find(nil).All(&results)

	if(err != nil){
		panic(err)
		return
	}

	for _, v := range results{
		client.Set(v.Dnax, v.IsMutant, 0)
	}
}

