package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var companies = []string{"Leffler, Bauch and D'Amore", "Rolfson, Gleason and Hammes", "Feeney - Bergnaum", "O'Reilly Group", "Ziemann - Johnston", "Hammes - Bradtke", "Kub, Brekke and Pfeffer", "Monahan, Roberts and Gutkowski", "Reynolds and Sons", "Dickens - Dickinson", "Bednar, Schmitt and Satterfield", "Ward - Grant", "Mraz - Kiehn", "Rau and Sons", "Medhurst, Labadie and Ebert", "Nikolaus Group", "Murray, Conroy and Armstrong", "Cartwright LLC", "Hane and Sons", "VonRueden LLC", "Huel - Weber", "Hamill - Hermiston", "Adams - Osinski", "Donnelly LLC", "Lindgren and Sons", "Schmeler - Dare", "Upton LLC", "Sawayn, Cruickshank and White", "Emard Inc", "Halvorson, Mertz and Kemmer"}

type Client struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func getRandomClient() Client {
	client := Client{}
	client.ID = GenerateID("client")
	client.Name = RandomElement(companies)
	return client
}

func getClients(c *gin.Context) {
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
	clients := make([]Client, limit)

	for i := 0; i < limit; i++ {
		clients[i] = getRandomClient()
	}

	c.JSON(200, clients)
}

func getClient(c *gin.Context) {
	clientID := c.Param("id")
	randomClient := getRandomClient()
	randomClient.ID = clientID
	c.JSON(200, randomClient)
}

func createClient(c *gin.Context) {
	newClient := Client{}
	err := c.BindJSON(&newClient)

	if err != nil {
		return
	}

	newClient.ID = GenerateID("client")
	c.JSON(http.StatusAccepted, newClient)
}
