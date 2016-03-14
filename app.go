package main

import (
	"github.com/drone/routes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Profile struct {
	Country       string	`json:"country"`
	Email         string	`json:"email"`
	Favorite_Color string	`json:"favorite_color"`
	Favorite_Sport string	`json:"favorite_sport"`
	Food          struct {
			      Drink_Alcohol string	`json:"drink_alcohol"`
			      Type         string	`json:"type"`
		      }	`json:"food"`
	Is_Smoking string	`json:"is_smoking"`
	Movie     struct {
			      Movies  []string	`json:"movies"`
			      Tv_Shows []string	`json:"tv_shows"`
		      }	`json:"movie"`
	Music struct {
			      Spotify_User_ID string	`json:"spotify_user_id"`
		      }	`json:"music"`
	Profession string	`json:"profession"`
	Travel     struct {
			      Flight struct {
					     Seat string	`json:"seat"`
				     }	`json:"flight"`
		      }	`json:"travel"`
	Zip string	`json:"zip"`
}

var profiles_map = make(map[string]Profile)

func main() {

	mux := routes.New()

	mux.Get("/profile/:email", GetProfile)
	mux.Post("/profile", PostProfile)
	mux.Del("/profile/:email", DeleteProfile)
	mux.Put("/profile/:email", PutProfile)

	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func GetProfile(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()                        //Get all the request parameters from the URL
	email := params.Get(":email")                    //Get the emailid from the URL
	log.Println(profiles_map[email])                 //Get the profile object corresponding to email key and print it
	mapB, _ := json.Marshal(profiles_map[email])     //Get the profile object corresponding to email key and marshal it to JSON
	w.Write([]byte(mapB))				 //Write the JSON to response
}

func DeleteProfile(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	email := params.Get(":email")
	delete(profiles_map, email)
	w.WriteHeader(http.StatusNoContent)
}

func PostProfile(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)		//Read the http body
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(string(body))			//convert byte to string and print the body

	var profile Profile				//define an object for profile
	err = json.Unmarshal(body, &profile)		//unmarshall json into object
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(profile)				//Print the object
	var email = profile.Email			//use the email field from profile as key for the map
	profiles_map[email] = profile			//set the whole profile object as value for the map
	rw.WriteHeader(http.StatusCreated)
}

func PutProfile(rw http.ResponseWriter, req *http.Request){
	params := req.URL.Query()                        //Get all the request parameters from the URL
	email := params.Get(":email")                    //Get the emailid from the URL
	log.Println(profiles_map[email])                 //Get the profile object corresponding to email key and print it

	currentProfileJson, _ := json.Marshal(profiles_map[email]) //Get the struct instance for the email from profile_map and marshal it to json

	var currentProfileMap map[string]*json.RawMessage
	var err = json.Unmarshal(currentProfileJson, &currentProfileMap) //Unmarshall the json into a map

	body, err := ioutil.ReadAll(req.Body)		//Read the http body (json format)
	log.Println(string(body))

	if err != nil {
		log.Println(err.Error())
	}
	var newProfilemap map[string]*json.RawMessage
	err = json.Unmarshal(body, &newProfilemap)	//Unmarshall the json into a map
	if err != nil {
		log.Println(err.Error())
	}

	for key, value := range newProfilemap {		//for each key, value in the new profile (from put body), update the corresponding key in already existing profile
		log.Println(string(key))
		log.Println(value)
		currentProfileMap[key] = value
	}

	jsonStr, _ := json.Marshal(currentProfileMap)	//Marshal the map into json
	log.Println(string(jsonStr))

	var newProfile Profile
	err = json.Unmarshal(jsonStr, &newProfile)	//Unmarshal the json into struct instance
     	profiles_map[email] = newProfile
	rw.WriteHeader(http.StatusNoContent)
}

