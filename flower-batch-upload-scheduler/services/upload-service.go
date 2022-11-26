package services

import (
	"fmt"
	"log"
	"public-flower-upload-scheduler/models"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func UploadRawFlowerData(db *gorm.DB, flowers []models.FlowerItem) (*models.FlowerBatchUploadLogs, error) {

	var errLogs []string = []string{}
	var dataList []models.PublicBiddingFlowers = []models.PublicBiddingFlowers{}

	for i := 0; i < len(flowers); i++ {
		flower := flowers[i]

		saleDate, err := time.Parse("2006-01-02", flower.SaleDate)
		if err != nil {
			log.Fatalf(err.Error())
			msg := err.Error() + " => data:" + fmt.Sprint(flower)
			errLogs = append(errLogs, msg)
		}
		data := models.PublicBiddingFlowers{
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

	tx := db.CreateInBatches(dataList, batchSize)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &models.FlowerBatchUploadLogs{
		Success: batchSize - len(errLogs),
		Failure: len(errLogs),
		Total:   batchSize,
		ErrLogs: strings.Join(errLogs, " || "),
	}, nil

}

func strToInt(from string, flower models.FlowerItem, errLogs *[]string) int {
	target, err := strconv.Atoi(from)
	if err != nil {
		log.Print(err.Error())
		msg := err.Error() + " => data:" + fmt.Sprint(flower)
		*errLogs = append(*errLogs, msg)
	}
	return target
}
