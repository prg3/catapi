package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"encoding/json"
	"github.com/go-redis/redis"
)

// Debug needs to be a global so we can use it everywhere
var debug bool
var redisclient *redis.Client

// Data structure for the cat images
type catImage struct {
	Url string `json:"url"`
	Id string `json:"id"`
	SourceUrl string `json:"source_url"`
}

type singleImage struct {
	Image catImage `json:"image"`
}

type manyImages struct {
	Images []catImage `json:"images"`
}


func GenSourceUrl (id string ) string {
	return "http://thecatapi.com/?id=" + id
}
func catHandler( apikey string ) (string, int) {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("x-api-key", apikey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	
	defer resp.Body.Close()
	
	cats := [] catImage{}
	cat := catImage{}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		if debug {
			log.Printf("Raw Response: %v\n", string(bodyBytes))
		}

		err= json.Unmarshal(bodyBytes, &cats)

		if err != nil {
			log.Fatal("Decoding error: ", err)
		}
		if debug {
			log.Printf("Received: %v\n", cats)
		}

		// Blatently disregard any other cats that may be on the array
		// because they shouldn't be there anyways
		cat = cats[0]		
	} 

	// Return the cat as a JSON response

	cat.SourceUrl = GenSourceUrl(cat.Id)
	catOutput := singleImage{}
	catOutput.Image = cat

	// Store the data in Redis if we have Redis
	if redisclient != nil {
		if debug {
			log.Println("In the Redis loop")
		}

		err := redisclient.Set(cat.Id, cat.Url, 0).Err()
		if err != nil {
			log.Fatal("Attempt to write to Redis failed: ", err)
		}
	}

	catJson, err := json.Marshal(catOutput)
	if err != nil {
		log.Fatal("Marshall failed: ", err)
	}

	return string(catJson), resp.StatusCode
}

func historyHandler() (string) {
	var cats []catImage
	var catOutput manyImages

	if redisclient == nil {
		return "{\"images\":[]}"
	}

	iter := redisclient.Scan(0, "", 0).Iterator()
	for iter.Next() {
		key := iter.Val()
		var cat catImage
		cat.Id = key
		cat.Url = redisclient.Get(key).Val()
		cat.SourceUrl = GenSourceUrl(key)
		if debug {
			log.Println(cat)
		}
		cats = append(cats, cat)
	}
	if err := iter.Err(); err != nil {
		log.Fatalln("Fatal error retrieving keys: ", err)
	}
	
	if debug {
		log.Println(cats)
	}

	catOutput.Images = cats

	if debug {
		log.Println(catOutput)
	}
	
	catJson, err := json.Marshal(catOutput)
	if err != nil {
		log.Fatal("Marshall failed: ", err)
	}

	return string(catJson)}

func main() {
	// Read environment variables
	debugstring := os.Getenv("DEBUG")
	if debugstring != "" {
		log.Println("Debug mode enabled, this could get noisy")
		debug = true
	}

	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		log.Println("APIKEY not defined in environment, quitting")
		return
	}

	// No redis isn't fatal, but we should warn about it being absent
	redisconnection := os.Getenv("REDIS")
	if redisconnection == "" {
		log.Println("REDIS not defined in environment, history will not be saved")
	} else {
		log.Println("Initializing Redis connection on server: ", redisconnection)
		// Current implementation does not use password or alternate DB
		redisclient = redis.NewClient(&redis.Options{
			Addr:     redisconnection,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err := redisclient.Ping().Result()
		if err != nil {
			log.Fatal("Redis connection failed: ", err)
		}
	}




	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Sorry, the path %s is invalid\n", r.URL.Path)
	})

	http.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		cat, httpCode := catHandler (apikey)
		if debug {
			log.Println("HTTP Response: ", httpCode)
			log.Println("Data: ", cat)
		}
		fmt.Fprintf(w, cat)
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		cats := historyHandler ()
		if debug {
			log.Println("Data: ", cats)
		}
		fmt.Fprintf(w, cats)
	})

	http.ListenAndServe(":80", nil)
}
