package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Spot struct {
	ID          *string  `json:"idspot"`
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	City        *string  `json:"city"`
	Country     *string  `json:"country"`
	Longitude   *float64 `json:"longitude"`
	Latitude    *float64 `json:"latitude"`
	ImageURL    *string  `json:"image_url"`
}

type Shortspot struct {
	ID       *string `json:"idspot"`
	Name     *string `json:"name"`
	ImageURL *string `json:"image_URL"`
}


func createspot(w http.ResponseWriter, r *http.Request) {
	// on déclare une variable de type Shark qui recevra les infos de notre nouveau requin
	var newspot Spot

	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le
	// Unmarshal afin de le stocker dans notre variable newshark
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the spot informations in order to update")
	}

	json.Unmarshal(reqBody, &newspot)

	stmt, err := db.Prepare("INSERT INTO spot (name, description, city, country, longitude, latitude, image_url) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(newspot.Name, newspot.Description, newspot.City, newspot.Country, newspot.Longitude, newspot.Latitude, newspot.ImageURL)

	fmt.Fprintf(w, "new spot created!")

}

func getOnespot(w http.ResponseWriter, r *http.Request) {
	var spot Spot
	spotID := mux.Vars(r)["id"]

	stmt, err := db.Query("SELECT * FROM spot WHERE idspot = ?", spotID)

	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	for stmt.Next() {
		err := stmt.Scan(&spot.ID, &spot.Name, &spot.Description, &spot.City, &spot.Country, &spot.Longitude, &spot.Latitude, &spot.ImageURL)
		if err != nil {
			panic(err.Error())
		}
	}

	json.NewEncoder(w).Encode(spot)
}

func getAllspots(w http.ResponseWriter, r *http.Request) {
	var spots []Spot

	stmt, err := db.Query("SELECT * FROM spot")

	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	for stmt.Next() {
		var spot Spot
		err := stmt.Scan(&spot.ID, &spot.Name, &spot.Description, &spot.City, &spot.Country, &spot.Longitude, &spot.Latitude, &spot.ImageURL)
		if err != nil {
			panic(err.Error())
		}
		spots = append(spots, spot)
	}

	json.NewEncoder(w).Encode(spots)
}

func getList(w http.ResponseWriter, r *http.Request) {
	var spots []Shortspot

	stmt, err := db.Query("SELECT idspot, name, image_url FROM spot")

	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	for stmt.Next() {
		var spot Shortspot
		err := stmt.Scan(&spot.ID, &spot.Name, &spot.ImageURL)
		if err != nil {
			panic(err.Error())
		}
		spots = append(spots, spot)
	}

	json.NewEncoder(w).Encode(spots)
}

// func updatespot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]
// 	var updatedspot Spot
// 	parseJson := parsingJson()

// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Fprintf(w, "Kindly enter data with the spot Name and description only in order to update")
// 	}
// 	json.Unmarshal(reqBody, &updatedspot)
// }

// func deletespot(w http.ResponseWriter, r *http.Request) {
// 	spotID := mux.Vars(r)["id"]
// 	
// 	fmt.Fprintf(w, "The spots with ID %v has been deleted successfully", spotID)
// 	}


// // on crée un objet db et un objet erreur
var db *sql.DB
var err error

func main() {
	// on ouvre la connection à la bdd et on utilise defer pour lui demander de rester ouverte
	// jusqu'à ce que qu'on ait fini.
	// attention "Root5003" est le mot de passe de mon serveur SQl, à remplacer par le votre
	// dans l'idéal ce serait une variable d'env
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/spooky_spot")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Opened DB !")
	}
	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/spot", createspot).Methods("POST")
	router.HandleFunc("/spots", getAllspots).Methods("GET")
	router.HandleFunc("/spot/{id}", getOnespot).Methods("GET")
	router.HandleFunc("/list", getList).Methods("GET")
	// router.HandleFunc("/spots/{id}", updatespot).Methods("PATCH")
	// router.HandleFunc("/spots/{id}", deletespot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
