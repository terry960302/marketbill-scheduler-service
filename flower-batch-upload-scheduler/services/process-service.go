package services

import (
	"errors"
	"log"
	"public-flower-upload-scheduler/models"
	"time"

	"gorm.io/gorm"
)

// 공공데이터를 가공하여 새로운 꽃, 꽃품목 데이터를 업로드
func ProcessFlowerRawData(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) (*models.FlowerBatchProcessLogs, error) {
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
				log := models.FlowerBatchProcessLogs{
					NewFlowerCount:      -1,
					NewFlowerTypeCount:  -1,
					AffectedFlowerCount: -1,
					Status:              "FAILURE",
					ErrLogs:             err.Error(),
				}
				db.Create(log)
			}
		}
	}()

	newFlowerTypesLen, err := UploadNewFlowerTypes(db, publicFlowers)
	if err != nil {
		return nil, err
	}
	newFlowersLen, err := UploadNewFlowers(db, publicFlowers)
	if err != nil {
		return nil, err
	}
	biddingFlowersLen, err := UploadNewFlowerBiddingDate(db, publicFlowers)
	if err != nil {
		return nil, err
	}

	log := models.FlowerBatchProcessLogs{
		NewFlowerTypeCount:  newFlowerTypesLen,
		NewFlowerCount:      newFlowersLen,
		AffectedFlowerCount: biddingFlowersLen,
		Status:              "SUCCESS",
		ErrLogs:             "",
	}

	tx := db.Create(&log)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &log, nil
}

// 공공데이터를 기반으로 디비에 없는 새로운 꽃을 추가합니다.(꽃 품목 우선 생성 실행)
func UploadNewFlowers(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) (int, error) {
	var newFlowers []models.Flowers = []models.Flowers{}
	var dbFlowers []models.Flowers
	var dbFlowerTypes []models.FlowerTypes
	flowerTable := db.Table("flowers")
	flowerTypeTable := db.Table("flower_types")

	// 디비에서 꽃과 종류 데이터 가져오기(디비 부하를 줄이기 위해 서버에서 처리)
	flowerTx := flowerTable.Find(&dbFlowers)
	if flowerTx.Error != nil {
		return -1, flowerTx.Error
	}
	flowerTypeTx := flowerTypeTable.Find(&dbFlowerTypes)
	if flowerTypeTx.Error != nil {
		return -1, flowerTypeTx.Error
	}

	// 디비 데이터 contains처리를 위한
	dbFlowersMap := convertFlowersToMap(dbFlowers)
	dbFlowerTypesMap := convertFlowerTypesToMap(dbFlowerTypes)

	// 공공데이터에서 새로운 꽃인지 탐색
	for _, pFlower := range publicFlowers {
		flowerRelationMap, flowerNameExists := dbFlowersMap[pFlower.FlowerName]
		flowerTypeID, isFlowerTypeExists := dbFlowerTypesMap[pFlower.FlowerType]

		if !isFlowerTypeExists {
			continue
		}
		_, flowerWithTypeExists := flowerRelationMap[flowerTypeID]

		// 새로 추가하는 케이스
		// - 같은 이름의 꽃이 존재하나, 다른 꽃품목인 경우
		// - 디비에 없는 명칭의 꽃인 경우
		if (flowerNameExists && !flowerWithTypeExists) || (!flowerNameExists) {
			newFlowers = append(newFlowers, models.Flowers{
				Name:         pFlower.FlowerName,
				FlowerTypeID: flowerTypeID,
			})
		}
	}

	// 꽃 배치 업로드
	createTx := flowerTable.CreateInBatches(newFlowers, len(newFlowers))
	if createTx.Error != nil {
		return -1, createTx.Error
	}

	log.Printf("New Flowers Length : %d", len(newFlowers))
	return len(newFlowers), nil

}

// 공공데이터를 기반으로 디비에 없는 새로운 꽃 품목을 추가합니다.
func UploadNewFlowerTypes(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) (int, error) {
	var flowerTypes []models.FlowerTypes
	var newFlowerTypes []models.FlowerTypes = []models.FlowerTypes{}

	table := db.Table("flower_types")
	flowerTypeTx := table.Find(&flowerTypes)
	if flowerTypeTx.Error != nil {
		return -1, flowerTypeTx.Error
	}

	// 공공데이터 꽃 품목 중복 삭제
	publicFlowerTypeMap := map[string]int{}
	nonDuplicatePublicFlowerTypes := []string{}
	for _, pFlower := range publicFlowers {
		_, exists := publicFlowerTypeMap[pFlower.FlowerType]
		if !exists {
			publicFlowerTypeMap[pFlower.FlowerType] = 1
			nonDuplicatePublicFlowerTypes = append(nonDuplicatePublicFlowerTypes, pFlower.FlowerType)
		}
	}

	// 기존 디비 데이터 안에 있는지 contains 처리를 위한 map
	dbFlowerTypeMap := convertFlowerTypesToMap(flowerTypes)

	// 꽃 품목 중 '새로운 품목' 배열 생성
	for _, pFlowerType := range nonDuplicatePublicFlowerTypes {
		_, exists := dbFlowerTypeMap[pFlowerType]
		if !exists {
			newFlowerTypes = append(newFlowerTypes, models.FlowerTypes{
				Name: pFlowerType,
			})
		}
	}

	// 꽃 품목 배치 업로드
	createTx := table.CreateInBatches(newFlowerTypes, len(newFlowerTypes))
	if createTx.Error != nil {
		return -1, createTx.Error
	}

	log.Printf("New Flower Types Length : %d", len(newFlowerTypes))
	return len(newFlowerTypes), nil
}

// 공공데이터를 기반으로 기존 꽃데이터에 꽃 경매일자를 추가합니다.
func UploadNewFlowerBiddingDate(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) (int, error) {
	var dbFlowers []models.Flowers
	var dbFlowerTypes []models.FlowerTypes
	flowerTable := db.Table("flowers")
	flowerTypeTable := db.Table("flower_types")
	biddingFlowerTable := db.Table("bidding_flowers")

	// 디비에서 꽃과 종류 데이터 가져오기(디비 부하를 줄이기 위해 서버에서 처리)
	flowerTx := flowerTable.Find(&dbFlowers)
	if flowerTx.Error != nil {
		return -1, flowerTx.Error
	}
	flowerTypeTx := flowerTypeTable.Find(&dbFlowerTypes)
	if flowerTypeTx.Error != nil {
		return -1, flowerTypeTx.Error
	}

	// 디비 데이터 contains처리를 위한
	dbFlowersMap := convertFlowersToMap(dbFlowers)
	dbFlowerTypesMap := convertFlowerTypesToMap(dbFlowerTypes)

	// 공공데이터에 해당하는 꽃 ID 추출
	curDate := time.Now()
	newBiddingFlowers := []models.BiddingFlowers{}
	for _, pFlower := range publicFlowers {
		flowerRelationMap, isFlowerNameExists := dbFlowersMap[pFlower.FlowerName]
		flowerTypeID, isFlowerTypeExists := dbFlowerTypesMap[pFlower.FlowerType]

		if !isFlowerNameExists || !isFlowerTypeExists {
			continue
		}
		flowerID, isFlowerExists := flowerRelationMap[flowerTypeID]

		if !isFlowerExists {
			continue
		}
		newBiddingFlowers = append(newBiddingFlowers, models.BiddingFlowers{
			FlowerID:    flowerID,
			BiddingDate: curDate,
		})
	}

	// 꽃 경매일자 배치 업로드
	createTx := biddingFlowerTable.CreateInBatches(newBiddingFlowers, len(newBiddingFlowers))
	if createTx.Error != nil {
		return -1, createTx.Error
	}

	log.Printf("New BiddingFlowers Length : %d", len(newBiddingFlowers))
	return len(newBiddingFlowers), nil
}

// 꽃 탐색용
// - 꽃은 이름은 같지만 품목이 다른 경우가 존재(value가 map인 이유)
func convertFlowersToMap(flowers []models.Flowers) map[string]map[uint]uint {
	dbFlowerMap := map[string]map[uint]uint{}
	for _, f := range flowers {
		_, exists := dbFlowerMap[f.Name]
		if exists {
			dbFlowerMap[f.Name][f.FlowerTypeID] = f.ID
		} else {
			flowerRelationMap := map[uint]uint{} // {flowerTypeID : flowerID}
			flowerRelationMap[f.FlowerTypeID] = f.ID
			dbFlowerMap[f.Name] = flowerRelationMap
		}
	}
	return dbFlowerMap
}

// 꽃 품목 탐색용
func convertFlowerTypesToMap(flowerTypes []models.FlowerTypes) map[string]uint {
	dbFlowerTypesMap := map[string]uint{}
	for _, t := range flowerTypes {
		dbFlowerTypesMap[t.Name] = t.ID
	}
	return dbFlowerTypesMap
}
