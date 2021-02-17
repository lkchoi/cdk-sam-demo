package main

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var breeds = []string{"Afghan Hound", "Beagle", "Bernese Mountain Dog", "Bloodhound", "Dalmatian", "Jack Russell Terrier", "Norwegian Elkhound"}
var names = []string{"Bailey", "Bella", "Max", "Lucy", "Charlie", "Molly", "Buddy", "Daisy", "Rocky", "Maggie", "Jake", "Sophie", "Jack", "Sadie", "Toby", "Chloe", "Cody", "Bailey", "Buster", "Lola", "Duke", "Zoe", "Cooper", "Abby", "Riley", "Ginger", "Harley", "Roxy", "Bear", "Gracie", "Tucker", "Coco", "Murphy", "Sasha", "Lucky", "Lily", "Oliver", "Angel", "Sam", "Princess", "Oscar", "Emma", "Teddy", "Annie", "Winston", "Rosie"}

type Pet struct {
	ID          string    `json:"id"`
	Breed       string    `json:"breed"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

func getRandomPet() Pet {
	pet := Pet{}

	pet.ID = generateID("pet")
	pet.Breed = randomElement(breeds)
	pet.Name = randomElement(names)

	pet.DateOfBirth = randomDate()

	return pet
}

func randomDate() time.Time {
	now := time.Now()
	start := now.AddDate(-15, 0, 0)
	delta := now.Unix() - start.Unix()

	sec := rand.Int63n(delta) + start.Unix()
	return time.Unix(sec, 0)
}

func getPets(c *gin.Context) {
	limit := 10
	if c.Query("limit") != "" {
		newLimit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		} else {
			limit = newLimit
		}
	}
	if limit > 50 {
		limit = 50
	}
	pets := make([]Pet, limit)

	for i := 0; i < limit; i++ {
		pets[i] = getRandomPet()
	}

	c.JSON(200, pets)
}

func getPet(c *gin.Context) {
	petID := c.Param("id")
	randomPet := getRandomPet()
	randomPet.ID = petID
	c.JSON(200, randomPet)
}

func createPet(c *gin.Context) {
	newPet := Pet{}
	err := c.BindJSON(&newPet)

	if err != nil {
		return
	}

	newPet.ID = generateID("pet")
	c.JSON(http.StatusAccepted, newPet)
}
