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

type List struct {
	List []Shortspot `json:"Shortspot"`
}

type Shortspot struct {
	ID        string `json:"id"`
	ImageName string `json:"imageName"`
	Name      string `json:"name"`
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

func createspot(w http.ResponseWriter, r *http.Request) {
	// on déclare une variable de type Shark qui recevra les infos de notre nouveau requin 
	var newspot Spot
	// on récup de Json parsé 
	parseJson := parsingJson()
	// on récup le corps de la requête, en affichant une erreur si c'est mal formaté et on le 
	// Unmarshal afin de le stocker dans notre variable newshark 
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the spot informations in order to update")
	}

	json.Unmarshal(reqBody, &newspot)

	// on ajoute newShark au tableau de requins dans notre objet parseJson 
	// mais attention à ce stade là, rien n'est encore écrit dans notre fichier Json
	parseJson.Allspots = append(parseJson.Allspots, newspot)

	fmt.Println(parseJson)

	// afin de pouvoir l'écrire dans le Json, on Marshal notre parseJson 
	modifJson, err := json.Marshal(parseJson)
	if err != nil {
		fmt.Println(err)
	}
	// on écrit notre modifJson (parseJson "marshalisé") dans le fichier sharks.json grace à ioutil 
	err = ioutil.WriteFile("spots.json", modifJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func getOnespot(w http.ResponseWriter, r *http.Request) {

	parseJson := parsingJson()
	spotID := mux.Vars(r)["id"]

	for _, singlespot := range parseJson.Allspots {
		if singlespot.ID == spotID {
			json.NewEncoder(w).Encode(singlespot)
		}
	}
}

func getAllspots(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	json.NewEncoder(w).Encode(parseJson)
}

func getList(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
	leng := len(parseJson.Allspots)
	vlist := make([]Shortspot, leng)
	for i, singlespot := range parseJson.Allspots {

		vlist[i].Name = singlespot.Name
		vlist[i].ID = singlespot.ID
		vlist[i].ImageName = singlespot.ImageName

	}
	json.NewEncoder(w).Encode(vlist)
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

 func deletespot(w http.ResponseWriter, r *http.Request) {
	parseJson := parsingJson()
 	spotID := mux.Vars(r)["id"]
	  // une boucle for pour chercher le spot concerné
	for i, singlespot := range parseJson.Allspots {
		if singlespot.ID == spotID {
		// on supprime le spot concerné en décalant les valeurs du tableau vers la gauche 
		// à partir de l'ID trouvé
			parseJson.Allspots = append(parseJson.Allspots[:i], parseJson.Allspots[i+1:]...)
			fmt.Fprintf(w, "The spots with ID %v has been deleted successfully", spotID)
		}
	}

	// afin de pouvoir l'écrire dans le Json, on Marshal notre parseJson 
	modifJson, err := json.Marshal(parseJson)
	if err != nil {
		fmt.Println(err)
	}
	// on écrit notre modifJson (parseJson "marshalisé") dans le fichier spots.json grace à ioutil 
	err = ioutil.WriteFile("spots.json", modifJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newspot)

}

func parsingJson() Allspots {

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

	return spots
}

func main() {
	parseSpot := parsingJson()
	fmt.Println(parseSpot)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/spot", createspot).Methods("POST")
	router.HandleFunc("/spots", getAllspots).Methods("GET")
	router.HandleFunc("/spots/{id}", getOnespot).Methods("GET")
	router.HandleFunc("/list", getList).Methods("GET")
	router.HandleFunc("/spots/{id}", deletespot).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
