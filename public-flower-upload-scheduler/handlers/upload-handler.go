package handlers

import (
	"fmt"
	"log"
	"net/http"
	"public-flower-upload-scheduler/datastore"
	"public-flower-upload-scheduler/models"
	"public-flower-upload-scheduler/services"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func UploadFlowers(c echo.Context) error {
	var errLogs []string = []string{}

	db := datastore.NewPostgresql()
	flowers, err := services.FetchFlowerItems()
	if err != nil {
		c.Logger().Error(err)
	}

	var dataList []models.PublicBiddingFlower = []models.PublicBiddingFlower{}

	for i := 0; i < len(flowers); i++ {
		flower := flowers[i]

		saleDate, err := time.Parse("2006-01-02", flower.SaleDate)
		if err != nil {
			log.Fatalf(err.Error())
			msg := err.Error() + " => data:" + fmt.Sprint(flower)
			errLogs = append(errLogs, msg)
		}
		data := models.PublicBiddingFlower{
			FlowerType: flower.PumName,
			FlowerName: flower.GoodName,
			Grade:      flower.LvNm,
			Quantity:   strToInt(flower.TotQty, flower, &errLogs),
			MaxPrice:   strToInt(flower.MaxAmt, flower, &errLogs),
			MinPrice:   strToInt(flower.MinAmt, flower, &errLogs),
			AvgPrice:   strToInt(flower.AvgAmt, flower, &errLogs),
			TotalPrice: strToInt(flower.TotAmt, flower, &errLogs),
			BidDate:    saleDate,
		}
		dataList = append(dataList, data)
	}

	batchSize := len(dataList)

	db.CreateInBatches(dataList, batchSize)

	result := map[string]interface{}{
		"success": batchSize - len(errLogs),
		"failure": len(errLogs),
		"errors":  errLogs,
	}

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func strToInt(from string, flower models.FlowerItem, errLogs *[]string) int {
	target, err := strconv.Atoi(from)
	if err != nil {
		log.Fatalf(err.Error())
		msg := err.Error() + " => data:" + fmt.Sprint(flower)
		*errLogs = append(*errLogs, msg)
	}
	return target
}
