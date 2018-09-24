package controllers

import(
	"net/http"
	"strings"
	"unicode/utf8"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-redis/redis"
	"API-Mutant/models"
)

var dnas = getSession().DB("dnas").C("dnas")

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var adn models.Dna

func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://127.0.0.1")

	if(err != nil){
		panic(err) 
	}

	return session
}

func IsMutant(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	
	err := decoder.Decode(&adn)

	if(err != nil){
		response(w, 400, "It's not a valid JSON")
		return
	}

	// DNA Sequence verification Len + characters
	valid := false
	adn.Dnax, valid = checkValidDnaSequence(w, adn)
	if(valid){
		//Because DNA makes us unique, we only need one row by human/mutant in our DB
		if(verifyIfExistOnDB(w, adn)){
			//DNA Verification Horizontal, Vertical & Diagonal
			verifyCombination(w, adn)
		}
	}
}

func checkValidDnaSequence(w http.ResponseWriter, adn models.Dna) (string, bool){
	validsChars := []string{"A", "C", "G", "T"}

	adn.Dnax = strings.Join(adn.Dna, "")

	if (len(adn.Dna) != 6 || utf8.RuneCountInString(adn.Dnax) != 36){
		response(w, 400, "It's not a valid DNA sequence")
		return adn.Dnax , false
		
	}

	a := []rune(adn.Dnax)
	
	for _, rd := range a {
		ok, _ := InArray(string(rd), validsChars)
		if (!ok){
			response(w, 400, "It's not a valid DNA sequence")
			return adn.Dnax , false
		}
	}

	return adn.Dnax , true
}

func verifyIfExistOnDB(w http.ResponseWriter, adn models.Dna) bool{
	var exist []models.Dna
	val, _ := client.Get(adn.Dnax).Result()

	if(val != "0" || val != "1"){
		query := bson.M{"dnax": adn.Dnax}

		iter := dnas.Find(query).Limit(1).Iter()
		err := iter.All(&exist)

		if(err != nil){
			panic(err)
			return false
		}

		if(len(exist) > 0){
			if(exist[0].IsMutant == true){
				response(w, 200, "The processed DNA belongs to a Mutant DB")
				client.Set(adn.Dnax, true, 0)
				return false
			}else{
				response(w, 403, "The processed DNA belongs to a human DB")
				client.Set(adn.Dnax, false, 0)
				return false
			}
		}
	}else if(val == "1"){
		response(w, 200, "The processed DNA belongs to a Mutant DB")
		return false
	}else{
		response(w, 403, "The processed DNA belongs to a human DB")
		return false
	}

	return true
}

func verifyCombination(w http.ResponseWriter, adn models.Dna) bool{
	nacth := []string{}
	counterCombH := 0
	counterCombV := 0
	counterCombD := 0

	//Horizontales
	for k, v := range adn.Dna{
		if(ValidCombination(v) == true){
			counterCombH++
		}

		nacth = convertColumsToRows(k,v,nacth)
	}
 	//Verticales
	for _, v := range nacth{
		if(ValidCombination(v) == true){
			counterCombV++
		}
	}
	//Diagonales 	
	nacth = convertDiagonalToRow(adn.Dna)

	for _, v := range nacth{
		if(ValidCombination(v) == true){
			counterCombD++
		}
	}
	
	if(counterCombH + counterCombV + counterCombD > 1){
		adn.IsMutant = true
		response(w, 200, "The processed DNA belongs to a Mutant")
		dnas.Insert(adn)
		return false
	}else{
		adn.IsMutant = false
		response(w, 403, "The processed DNA belongs to a human")
		dnas.Insert(adn)
		return false
	}

	return true
}


func ValidCombination(val string) bool{
	ch := 0
	
	if(strings.Index(val, "AAAA") > -1){
		ch++
	}else if (strings.Index(val, "TTTT") > -1){
		ch++
	}else if (strings.Index(val, "CCCC") > -1){
		ch++
	}else if (strings.Index(val, "GGGG") > -1){
		ch++
	}

	if(ch > 0){
		return true
	}else{
		return false
	}
}

func response(w http.ResponseWriter, http_code int, txt string){
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http_code)
	json.NewEncoder(w).Encode(txt)
}

func convertColumsToRows(row int, val string, array []string) []string{
	a := []rune(val)

	for i, r := range a {
		if(!isset(array,i)){
			array = append(array, string(r))
		}else{
			array[i] = array[i] + string(r)
		}
	}
	return array
}

func convertDiagonalToRow(array []string) []string{
	na := []string{"","","","",""}
	adnAux := [][]string{}

	for i := 0; i < len(array); i++{
		adnAux = append(adnAux,strings.Split(array[i], ""))
	}

	for i := 0; i < len(adnAux); i++{
		for j := 0; j < len(adnAux[i]); j++{
			if(i == j){
				na[0] = na[0] + adnAux[i][j]
			}else if((i + 1) == j){
				na[1] = na[1] + adnAux[i][j]
			}else if((i + 2) == j){
				na[2] = na[2] + adnAux[i][j]
			}else if((i - 1) == j){
				na[3] = na[3] + adnAux[i][j]
			}else if((i - 2 ) == j){
				na[4] = na[4] + adnAux[i][j]
			}
		}
	}

	return na
}

func isset(arr []string, index int) bool {
    return (len(arr) > index)
}

func InArray(val string, array []string) (ok bool, i int) {
    for i = range array {
        if ok = array[i] == val; ok {
            return
        }
    }
    return
}