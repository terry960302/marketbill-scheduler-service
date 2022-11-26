package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"public-flower-upload-scheduler/models"
	"strconv"
	"time"
)

const (
	API_HOST    = "https://flower.at.or.kr/api/returnData.api?kind=f001"
	FLOWER_TYPE = 1
	COUNT       = 999999
	DATE_FORMAT = "2006-01-02"
	DATA_TYPE   = "json"
)

func FetchFlowerItems() ([]models.FlowerItem, error) {
	url := GenerateUrl()
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var apiResp models.PublicApiResponse
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &apiResp); err != nil {
		return nil, err
	}

	return apiResp.Response.Items, nil
}

func GenerateUrl(baseDate ...string) string {
	apiKey := os.Getenv("API_KEY")

	now := time.Now()
	var _baseDate string = ""
	if len(baseDate) == 0 {
		_baseDate = now.Format(DATE_FORMAT)
	} else {
		_baseDate = baseDate[0]
	}

	flowerType := strconv.Itoa(FLOWER_TYPE)
	count := strconv.Itoa(COUNT)

	url := API_HOST +
		"&serviceKey=" + apiKey +
		"&baseDate=" + _baseDate +
		"&flowerGubn=" + flowerType +
		"&dataType=" + DATA_TYPE +
		"&countPerPage=" + count

	return url
}
