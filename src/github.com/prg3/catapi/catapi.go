package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"encoding/json"
)

// Debug needs to be a global so we can use it everywhere
var debug bool


// Data structure for the cat images
// [{"breeds":[],"id":"a70","url":"https://cdn2.thecatapi.com/images/a70.jpg"}]
type catImage struct {
	Url string `json:"url"`
	Id string `json:"id"`
	SourceUrl string `json:"source_url"`
}

type singleImage struct {
	Image catImage `json:"image"`
}

// Target image type
// {
//     "url": "http://24.media.tumblr.com/tumblr_m3ay3e1zHp1qcxyrro1_1280.jpg",
//     "id": "c11",
//     "source_url": "http://thecatapi.com/?id=c11"
// }

func catHandler( r *http.Request, apikey string ) (string, int) {
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

	cat.SourceUrl = "http://thecatapi.com/?id=" + cat.Id
	catOutput := singleImage{}
	catOutput.Image = cat

	catJson, err := json.Marshal(catOutput)
	if err != nil {
		log.Fatal("Marshall failed: ", err)
	}

	return string(catJson), resp.StatusCode
}

func historyHandler(r *http.Request) string {
	return "History handler\n"
}

func main() {
	// Read environment variables
	apikey := os.Getenv("APIKEY")
	redisconnection := os.Getenv("REDIS")
	debugstring := os.Getenv("DEBUG")
	if apikey == "" {
		log.Println("APIKEY not defined in environment, quitting")
		return
	}

	// No redis isn't fatal, but we should warn about it being absent
	if redisconnection == "" {
		log.Println("REDIS not defined in environment, history will not be saved")
	}

	if debugstring != "" {
		log.Println("Debug mode enabled, this could get noisy")
		debug = true
	}


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Sorry, the path %s is invalid\n", r.URL.Path)
	})

	http.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
		cat, httpCode := catHandler (r, apikey)
		if debug {
			fmt.Println(cat)
			fmt.Println(httpCode)
		}
		fmt.Fprintf(w, cat)
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, historyHandler( r ))
	})

	http.ListenAndServe(":80", nil)
}
