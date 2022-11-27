package services

import (
	"log"
	"public-flower-upload-scheduler/models"

	"gorm.io/gorm"
)

// 공공데이터를 가공하여 새로운 꽃, 꽃품목 데이터를 업로드
func ProcessFlowerRawData(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) error {

	if err := UploadNewFlowerTypes(db, publicFlowers); err != nil {
		return err
	}
	if err := UploadNewFlowers(db, publicFlowers); err != nil {
		return err
	}

	return nil
}

// 공공데이터를 기반으로 디비에 없는 새로운 꽃을 추가합니다.(꽃 품목 우선 생성 실행)
func UploadNewFlowers(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) error {
	var newFlowers []models.Flowers = []models.Flowers{}
	var dbFlowers []models.Flowers
	var dbFlowerTypes []models.FlowerTypes
	flowerTable := db.Table("flowers")
	flowerTypeTable := db.Table("flower_types")

	// 디비에서 꽃과 종류 데이터 가져오기(디비 부하를 줄이기 위해 서버에서 처리)
	flowerTx := flowerTable.Find(&dbFlowers)
	if flowerTx.Error != nil {
		return flowerTx.Error
	}
	flowerTypeTx := flowerTypeTable.Find(&dbFlowerTypes)
	if flowerTypeTx.Error != nil {
		return flowerTypeTx.Error
	}

	// 디비 데이터 contains처리를 위한
	dbFlowersMap := convertFlowersToMap(dbFlowers)
	dbFlowerTypesMap := convertFlowerTypesToMap(dbFlowerTypes)

	// 공공데이터에서 새로운 꽃인지 탐색
	for _, pFlower := range publicFlowers {
		_, flowerNameExists := dbFlowersMap[pFlower.FlowerName]
		flowerTypeID, flowerTypeExists := dbFlowerTypesMap[pFlower.FlowerType] // 꽃품목을 먼저 업데이트하기 때문에 새로운 품목이 있을리는 없음(단지 ID 추출용)

		if !flowerNameExists && flowerTypeExists {
			newFlowers = append(newFlowers, models.Flowers{
				Name:         pFlower.FlowerName,
				FlowerTypeID: flowerTypeID,
			})
		}
	}

	// 꽃 배치 업로드
	createTx := flowerTable.CreateInBatches(newFlowers, len(newFlowers))
	if createTx.Error != nil {
		return createTx.Error
	}

	log.Printf("New Flowers Length : %d", len(newFlowers))
	return nil

}

// 공공데이터를 기반으로 디비에 없는 새로운 꽃 품목을 추가합니다.
func UploadNewFlowerTypes(db *gorm.DB, publicFlowers []models.PublicBiddingFlowers) error {
	var flowerTypes []models.FlowerTypes
	var newFlowerTypes []models.FlowerTypes = []models.FlowerTypes{}

	table := db.Table("flower_types")
	flowerTypeTx := table.Find(&flowerTypes)
	if flowerTypeTx.Error != nil {
		return flowerTypeTx.Error
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
		return createTx.Error
	}

	log.Printf("New Flower Types Length : %d", len(newFlowerTypes))
	return nil
}

func convertFlowersToMap(flowers []models.Flowers) map[string]uint {
	dbFlowerMap := map[string]uint{}
	for _, f := range flowers {
		dbFlowerMap[f.Name] = f.ID
	}
	return dbFlowerMap
}
func convertFlowerTypesToMap(flowerTypes []models.FlowerTypes) map[string]uint {
	dbFlowerTypesMap := map[string]uint{}
	for _, t := range flowerTypes {
		dbFlowerTypesMap[t.Name] = t.ID
	}
	return dbFlowerTypesMap
}
