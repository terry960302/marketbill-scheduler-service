package services

import (
	"errors"
	"fmt"
	"log"
	"public-flower-upload-scheduler/models"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UploadFlowers(db *gorm.DB, rawFlowers []models.FlowerItem) (*[]models.PublicBiddingFlowers, *models.FlowerBatchUploadLogs, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	defer func() {
		var err error = nil
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}

			if err != nil {
				log := models.FlowerBatchUploadLogs{
					Success: -1,
					Failure: -1,
					Total:   -1,
					ErrLogs: err.Error(),
				}
				logger.Error(err.Error())
				db.Create(log)
			}
		}
	}()

	publicFlowers, log, err := UploadRawFlowerData(db, rawFlowers)
	if err != nil {
		return nil, nil, err
	}

	tx := db.Create(log)
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	logger.Info("completed")

	return publicFlowers, log, nil
}

func UploadRawFlowerData(db *gorm.DB, flowers []models.FlowerItem) (*[]models.PublicBiddingFlowers, *models.FlowerBatchUploadLogs, error) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

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
		return nil, nil, tx.Error
	}

	logger.Info("completed")

	return &dataList, &models.FlowerBatchUploadLogs{
		Success: batchSize - len(errLogs),
		Failure: len(errLogs),
		Total:   batchSize,
		ErrLogs: strings.Join(errLogs, " || "),
	}, nil

}

func strToInt(from string, flower models.FlowerItem, errLogs *[]string) int {
	target, err := strconv.Atoi(from)
	if err != nil {
		msg := err.Error() + " => data:" + fmt.Sprint(flower)
		*errLogs = append(*errLogs, msg)
	}
	return target
}
