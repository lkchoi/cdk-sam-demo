package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var addrs = []string{"123 Sesame St", "800 Airport Blvd", "221B Baker Street", "12 Grimmauld Place", "1938 Sulivan Ln", "South Park basin", "425 Grove Street", "1313 Webfoot Walk", "Skypad Apartments", "Apartment 5A", "124 Conch Street", "17 Cherry Tree Lane", "Saffron Hill", "No. 7 Saville Row", "Cemetery Ridge", "112 ½ Beacon Street", "Oxenthorpe Rd", "740 Evergreen Terrace", "31 Spooner Street", "Southfork Ranch", "510 Glenview", "Los Santos", "Vice City", "56B Whitehavens Mansions", "32 Windsor Gardens", "186 Fleet Street"}

type Space struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Coordinates [][2]float32 `json:"coordinates,omitempty"`
	ClientID    string       `json:"client_id"`
}

func getRandomSpace() Space {
	space := Space{}

	space.ID = generateID("space")
	space.Name = randomElement(addrs)
	space.Coordinates = randomPoints(randomInt(4, 20))
	space.ClientID = generateID("client")

	return space
}

func randomPoints(n int) [][2]float32 {
	points := make([][2]float32, n)
	for i := 0; i < n; i++ {
		points[i] = randomPoint(100)
	}
	return points
}

func randomPoint(max float32) [2]float32 {
	return [2]float32{randomFloat(max), randomFloat(max)}
}

func getSpaces(c *gin.Context) {
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
	spaces := make([]Space, limit)

	for i := 0; i < limit; i++ {
		spaces[i] = getRandomSpace()
		spaces[i].Coordinates = make([][2]float32, 0)
	}

	c.JSON(200, spaces)
}

func getSpace(c *gin.Context) {
	spaceID := c.Param("id")
	randomSpace := getRandomSpace()
	randomSpace.ID = spaceID
	c.JSON(200, randomSpace)
}

func createSpace(c *gin.Context) {
	newSpace := Space{}
	err := c.BindJSON(&newSpace)

	if err != nil {
		return
	}

	newSpace.ID = generateID("space")
	c.JSON(http.StatusAccepted, newSpace)
}
