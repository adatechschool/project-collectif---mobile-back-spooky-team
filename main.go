package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Spot struct {
	ID          string  `json:"id"`
	ImageName   string  `json:"imageName"`
	Description string  `json:"description"`
	City        string  `json:"city"`
	Longitude   float64 `json:"longitude"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Name        string  `json:"name"`
}

type Allspots struct {
	Allspots []Spot `json:"spots"`
}

// var spots = allspots{
// 	{
// 		ID:          "1",
// 		Name:       "Introduction to Golang",
// 		Description: "Come join us for a chance to learn how golang works and get to spotually try it out",
// 	},
// }

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

// func createspot(w http.ResponseWriter, r *http.Request) {

// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the spot Name and description only in order to update")
// 	}

// 	json.Unmarshal(reqBody, &newspot)
// 	spots = append(Spots, newspot)
// 	w.WriteHeader(http.StatusCreated)

// 	json.NewEncoder(w).Encode(newspot)
// }

func getOnespot(w http.ResponseWriter, r *http.Request) {

	parseJson := parsingJson()
	spotID := mux.Vars(r)["id"]

	for _, singlespot := range parseJson {
		if singlespot.ID == spotID {
			json.NewEncoder(w).Encode(singlespot)
		}
	}
}

func getAllspots(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	json.NewEncoder(w).Encode(parseJson)
}

// func updatespot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]
// 	var updatedspot Spot

// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the spot Name and description only in order to update")
// 	}
// 	json.Unmarshal(reqBody, &updatedspot)

// 	for i, singlespot := range Spots {
// 		if singlespot.ID == spotID {
// 			singlespot.Name = updatedspot.Name
// 			singlespot.Description = updatedspot.Description
// 			Spots = append(Spots[:i], singlespot)
// 			json.NewEncoder(w).Encode(singlespot)
// 		}
// 	}
// }

// func deletespot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]

// 	for i, singlespot := range spots {
// 		if singlespot.ID == spotID {
// 			spots = append(spots[:i], spots[i+1:]...)
// 			fmt.Fprintf(w, "The spot with ID %v has been deleted successfully", spotID)
// 		}
// 	}
// }

func parsingJson() []Spot {

	// Open our jsonFile
	jsonFile, err := os.Open("spots.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened spots.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var spots Allspots

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &spots)

	return spots.Allspots
}

func main() {
	parseSpot := parsingJson()
	fmt.Println(parseSpot)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	// router.HandleFunc("/spot", createspot).Methods("POST")
	router.HandleFunc("/spots", getAllspots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOnespot).Methods("GET")
	// router.HandleFunc("/spots/{id}", updatespot).Methods("PATCH")
	// router.HandleFunc("/spots/{id}", deletespot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
