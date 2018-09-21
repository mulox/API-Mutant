package controllers

import(
	"net/http"
	"encoding/json"
	"API-Mutant/models"
)

// var dnas = getSession().DB("dnas").C("dnas")

// func getSession() *mgo.Session{
// 	session, err := mgo.Dial("mongodb://127.0.0.1")

// 	if(err != nil){
// 		panic(err) 
// 	}

// 	return session
// }//Global connection 


func GetStats(w http.ResponseWriter, r *http.Request) {
	var adnlist []models.Dna
	// var adn models.Dna

	err := dnas.Find(nil).All(&adnlist)

	if(err != nil){
		response(w, 400, "It's not a valid JSON")
	}

	var result models.Stats
	result.Mutant = 0
	result.Human = 0
	result.Ratio = 0

	for _, v := range adnlist{

		if(v.IsMutant == true){
			result.Mutant++
		}else{
			result.Human++
		}
	}

	if(result.Human > 0){
		result.Ratio = float32(result.Mutant)/float32(result.Human)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(result)
}

// func response(w http.ResponseWriter, http_code int, txt string){
// 	w.Header().Set("Content-type", "application/json")

// 	w.WriteHeader(http_code)
// 	json.NewEncoder(w).Encode(txt)
// }