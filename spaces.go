package main

import (
	"log"
	"net/http"
	"strconv"

	b64 "encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	d "github.com/aws/aws-sdk-go/service/dynamodb"
	a "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)

var tableName = "Gemini"
var entityType = "space"
var idPrefix = "space_"
var maxLimit = 50

// Space name and coordinates for a layout
type Space struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Coordinates [][2]float32 `json:"coordinates,omitempty"`
	ClientID    string       `json:"client_id"`
	EntityType  string       `json:"entity_type" default:"space"`
	Area        float32      `json:"area"`
}

// SpaceRecord alias record for space
type SpaceRecord struct {
	PK          string       `json:"PK"`
	SK          string       `json:"SK"`
	Name        string       `json:"name"`
	Coordinates [][2]float32 `json:"coordinates,omitempty"`
	ClientID    string       `json:"client_id"`
	EntityType  string       `json:"entity_type" default:"space"`
	Area        float32      `json:"area"`
}

// ClientSpaceRecord client->space
type ClientSpaceRecord struct {
	PK         string  `json:"PK"`
	SK         string  `json:"SK"`
	Name       string  `json:"name"`
	EntityType string  `json:"entity_type" default:"space"`
	ClientID   string  `json:"client_id"`
	Area       float32 `json:"area"`
}

func toSpaceRecord(s Space) SpaceRecord {
	return SpaceRecord{
		PK:          s.ID,
		SK:          "A",
		Name:        s.Name,
		Coordinates: s.Coordinates,
		EntityType:  entityType,
		ClientID:    s.ClientID,
		Area:        s.Area,
	}
}

func fromSpaceRecord(r SpaceRecord) Space {
	return Space{
		ID:          r.PK,
		Name:        r.Name,
		Coordinates: r.Coordinates,
		EntityType:  r.EntityType,
		ClientID:    r.ClientID,
		Area:        r.Area,
	}
}

func toClientSpaceRecord(s Space) ClientSpaceRecord {
	return ClientSpaceRecord{
		PK:         s.ClientID,
		SK:         s.ID,
		Name:       s.Name,
		EntityType: entityType,
		ClientID:   s.ClientID,
		Area:       s.Area,
	}
}

func fromClientSpaceRecord(r ClientSpaceRecord) Space {
	return Space{
		ID:         r.SK,
		Name:       r.Name,
		EntityType: r.EntityType,
		ClientID:   r.PK,
		Area:       r.Area,
	}
}

func getSpaces(c *gin.Context) {
	ddb := DynamoClient()

	// FIXME read clientID from API Key
	clientID := "client_1ooV1qMUTNAJrmI3YJuRkrZvMGJ"

	limit := int64(10)
	if c.Query("limit") != "" {
		newLimit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		} else {
			limit = int64(newLimit)
		}
	}
	if limit > 50 {
		limit = 50
	}

	query := &d.QueryInput{
		TableName: aws.String("Gemini"),
		ExpressionAttributeValues: map[string]*d.AttributeValue{
			":pk": {S: aws.String(clientID)},
			":sk": {S: aws.String(idPrefix)},
		},
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		Limit:                  aws.Int64(limit),
		ReturnConsumedCapacity: aws.String(d.ReturnConsumedCapacityTotal),
	}

	if c.Query("cursor") != "" {
		cursor, err := b64.StdEncoding.DecodeString(c.Query("cursor"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		query.SetExclusiveStartKey(map[string]*d.AttributeValue{
			"PK": {S: aws.String(clientID)},
			"SK": {S: aws.String(string(cursor))},
		})
	}
	log.Print(query.ExclusiveStartKey)

	// run query: PK={id} and SK begins_with(SK, space_)
	res, err := ddb.Query(query)
	log.Print("ConsumedCapacity: ", res.ConsumedCapacity)
	log.Print("LastEvaluatedKey: ", res.LastEvaluatedKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// unmarshall to list of client-space records
	items := []ClientSpaceRecord{}
	if err = a.UnmarshalListOfMaps(res.Items, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// convert list of client-space records to space records
	spaces := []Space{}
	for _, item := range items {
		spaces = append(spaces, fromClientSpaceRecord(item))
	}

	lastKey := map[string]string{}
	if err = a.UnmarshalMap(res.LastEvaluatedKey, &lastKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	prevCursor := c.Query("cursor")
	nextCursor := b64.StdEncoding.EncodeToString([]byte(lastKey["SK"]))

	// respond
	c.JSON(http.StatusOK, gin.H{
		"prev": prevCursor,
		"next": nextCursor,
		"data": spaces,
	})
	return
}

func getSpace(c *gin.Context) {
	id := c.Param("id")
	ddb := DynamoClient()
	res, err := ddb.GetItem(&d.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*d.AttributeValue{
			"PK": {S: aws.String(id)},
			"SK": {S: aws.String("A")},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if res.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		return
	}
	item := SpaceRecord{}
	if err = a.UnmarshalMap(res.Item, &item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	space := fromSpaceRecord(item)
	c.JSON(200, space)
	return
}

func createSpace(c *gin.Context) {
	space := Space{}
	if err := c.BindJSON(&space); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	log.Print(space.Area)
	space.ID = GenerateID("space")
	space.EntityType = "space"
	saveSpace(space)
	c.JSON(http.StatusCreated, space)
}

func updateSpace(c *gin.Context) {
	spaceID := c.Param("id")
	item := Space{}
	err := c.BindJSON(&item)
	if err != nil {
		return
	}
	item.ID = spaceID
	c.JSON(http.StatusCreated, item)
}

func saveSpace(space Space) ([]map[string]*d.AttributeValue, error) {
	items := make([]map[string]*d.AttributeValue, 0)

	// client->space
	item, err := a.MarshalMap(toClientSpaceRecord(space))
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	// space alias
	item, err = a.MarshalMap(toSpaceRecord(space))
	if err != nil {
		return nil, err
	}
	items = append(items, item)

	ddb := DynamoClient()
	for _, item := range items {
		if _, err := ddb.PutItem(&d.PutItemInput{
			Item:      item,
			TableName: aws.String(tableName),
		}); err != nil {
			return nil, err
		}
	}

	log.Print("space added", items)

	return items, nil
}
