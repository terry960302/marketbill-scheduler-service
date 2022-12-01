package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"public-flower-upload-scheduler/handlers"
	"public-flower-upload-scheduler/models"
	"public-flower-upload-scheduler/services"
	"public-flower-upload-scheduler/test"
	"public-flower-upload-scheduler/utils"

	"github.com/aws/aws-lambda-go/events"
)

func init() {
	profile := os.Getenv("PROFILE")
	log.Print("PROFILE : ", profile)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	r := models.NewLambdaResponse()
	log.Printf("[LOG] method : %s", request.HTTPMethod)
	switch request.HTTPMethod {
	case "GET":
		return handlers.HealthCheck(r, request)
	case "POST":
		return handlers.HandleUpload(r, request)
	default:
		return r.Error(http.StatusBadRequest, "Wrong http method")
	}
}

// func main() {
// 	lambda.Start(HandleRequest)
// }

// test
func main() {
	envSetter := test.NewEnvSetter("local")
	envSetter.SetEnv()

	url := services.GenerateUrl("2022-11-30")
	fmt.Println(url)

	resp, _ := http.Get(url)

	var apiResp models.PublicApiResponse
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &apiResp)

	flowerMap := map[string]int{}

	for _, item := range apiResp.Response.Items {
		fName := fmt.Sprintf("%s(%s)", item.GoodName, item.PumName)
		count, exists := flowerMap[fName]
		if exists {
			flowerMap[fName] = count + 1
		} else {
			flowerMap[fName] = 1
		}
	}

	utils.PrettyPrint(flowerMap)

	// db, err := datastore.NewPostgresql()
	// if err != nil {
	// 	e := fmt.Sprintf("[NewPostgresql] %s", err.Error())
	// 	panic(e)
	// }

	// db.AutoMigrate(&models.FlowerBatchUploadLogs{}, &models.FlowerBatchProcessLogs{})

	// rawFlowers, err := services.FetchRawFlowerItems()
	// if err != nil {
	// 	panic(err)
	// }

	// publicFlowers, uploadLog, err := services.UploadFlowers(db, rawFlowers)
	// if err != nil {
	// 	panic(err)
	// }

	// utils.PrettyPrint(uploadLog)

	// processLog, err := services.ProcessFlowerRawData(db, *publicFlowers)
	// if err != nil {
	// 	panic(err)
	// }
	// utils.PrettyPrint(processLog)

}
