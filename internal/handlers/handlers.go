package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mugglemath/go-dew/internal/discord"
	"github.com/mugglemath/go-dew/internal/weather"
)

var outdoorDewpoint string = ""

func HandleOutdoorDewpoint(c *gin.Context) {
	response, err := weather.NwsAPIResponse(os.Getenv("OFFICE"), os.Getenv("GRID_X"), os.Getenv("GRID_Y"), os.Getenv("NWS_USER_AGENT"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	parsed, err := weather.ParseOutdoorDewpoint(response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	outdoorDewpoint = fmt.Sprintf("%.2f", parsed)
	c.JSON(http.StatusOK, parsed)
}

func HandleSensorData(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indoorTemperature := data["indoorTemperature"].(float64)
	indoorHumidity := data["indoorHumidity"].(float64)
	indoorDewpoint := data["indoorDewpoint"].(float64)
	outdoorDewpoint := data["outdoorDewpoint"].(float64)
	dewpointDelta := data["dewpointDelta"].(float64)
	windowsShouldBe := data["windowsShouldBe"].(string)
	humidityAlert := data["humidityAlert"].(bool)
	isoTimestamp := time.Now().Format(time.RFC3339)

	message := fmt.Sprintf("%s\n"+
		"Indoor Temperature = %.2f C\n"+
		"Indoor Humidity = %.2f %%\n"+
		"Indoor Dewpoint = %.2f C\n"+
		"Outdoor Dewpoint = %.2f C\n"+
		"Dewpoint Delta = %.2f C\n"+
		"Windows Should Be = %s\n"+
		"Humidity Alert = %t",
		isoTimestamp, indoorTemperature, indoorHumidity, indoorDewpoint,
		outdoorDewpoint, dewpointDelta, windowsShouldBe, humidityAlert)

	discord.SendDiscordMessage(message, os.Getenv("DISCORD_SENSOR_FEED_WEBHOOK_URL"))
	fmt.Printf("Received data from C++ app: %v\n", data)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "POST request received"})
}

func HandleWindowAlert(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indoorDewpoint := data["indoorDewpoint"].(float64)
	outdoorDewpoint := data["outdoorDewpoint"].(float64)
	dewpointDelta := data["dewpointDelta"].(float64)
	windowsShouldBe := data["windowsShouldBe"].(string)
	isoTimestamp := time.Now().Format(time.RFC3339)

	message := fmt.Sprintf("%s\n@everyone\n"+
		"Indoor Dewpoint = %.2f C\n"+
		"Outdoor Dewpoint = %.2f C\n"+
		"Dewpoint Delta = %.2f C\n"+
		"Windows Should Be = %s\n",
		isoTimestamp, indoorDewpoint, outdoorDewpoint, dewpointDelta, windowsShouldBe)

	discord.SendDiscordMessage(message, os.Getenv("DISCORD_WINDOW_ALERT_WEBHOOK_URL"))
	fmt.Printf("Received data from C++ app: %v\n", data)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "POST request received"})
}

func HandleHumidityAlert(c *gin.Context) {
	var data map[string]interface{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	indoorHumidity := data["indoorHumidity"].(float64)
	humidityAlert := data["humidityAlert"].(bool)
	isoTimestamp := time.Now().Format(time.RFC3339)

	message := fmt.Sprintf("%s\n@everyone\nIndoor Humidity = %.2f %%\nHumidity Alert = %t",
		isoTimestamp, indoorHumidity, humidityAlert)

	discord.SendDiscordMessage(message, os.Getenv("DISCORD_HUMIDITY_ALERT_WEBHOOK_URL"))
	fmt.Printf("Received data from C++ app: %v\n", data)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "POST request received"})
}
