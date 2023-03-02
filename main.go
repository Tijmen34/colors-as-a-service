package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Colors []Color
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	gin.SetMode(gin.ReleaseMode)
	loadColors()
}

func loadColors() {
	colorsFile, err := os.Open("colors.json")
	if err != nil {
		log.Panic().Err(err).Msg("Error opening colors file.")
	}
	defer colorsFile.Close()

	byteValue, err := io.ReadAll(colorsFile)
	if err != nil {
		log.Panic().Err(err).Msg("Error reading colors file.")
	}

	var colors []Color
	err = json.Unmarshal(byteValue, &colors)
	if err != nil {
		log.Panic().Err(err).Msg("Error unmarshalling colors file")
	}
	Colors = colors
	log.Info().Msg("All colors loaded.")
}

func main() {
	r := gin.Default()
	r.Static("/assets", "./static/assets")
	r.LoadHTMLFiles("./static/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/color", getRandomColor)
	r.GET("/api/color/:colorName", getColor)
	r.GET("/api/color/:colorName/image", getColorImage)

	log.Info().Msg("Starting server on port 8080")
	r.Run("localhost:8080")
}

type Color struct {
	Name    string `json:"name"`
	HexCode string `json:"hex"`
	RGB     struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	}
}

func getColor(c *gin.Context) {
	colorName := c.Param("colorName")
	color := findColorByName(colorName)
	if color.Name == "" {
		c.JSON(404, gin.H{"status": 404, "message": "Color not found"})
	} else {
		c.JSON(200, color)
	}
}

func getRandomColor(c *gin.Context) {
	randomColor := Colors[rand.Intn(len(Colors))]
	c.JSON(200, randomColor)
}

func getColorImage(c *gin.Context) {
	colorName := c.Param("colorName")
	color := findColorByName(colorName)
	if color.Name == "" {
		c.JSON(404, gin.H{"status": 404, "message": "Color not found"})
	} else {
		imageName := createColorImage(color)
		c.File(imageName)
		os.Remove(imageName)
	}
}

func createColorImage(givenColor Color) string {
	width, height := 200, 200
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	red := color.RGBA{uint8(givenColor.RGB.R), uint8(givenColor.RGB.G), uint8(givenColor.RGB.B), 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			image.Set(x, y, red)
		}
	}

	imageName := fmt.Sprintf("%s.png", givenColor.Name)

	f, _ := os.Create(imageName)

	png.Encode(f, image)

	return imageName
}

func findColorByName(colorName string) Color {
	for _, color := range Colors {
		if color.Name == colorName {
			return color
		}
	}
	return Color{}
}
