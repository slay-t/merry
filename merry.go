package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/connect"
	"os"
	"time"
)

type IotButtonEvent struct {
	DeviceInfo struct {
		DeviceID      string  `json:"deviceId"`
		Type          string  `json:"type"`
		RemainingLife float64 `json:"remainingLife"`
		Attributes    struct {
			ProjectRegion      string `json:"projectRegion"`
			ProjectName        string `json:"projectName"`
			PlacementName      string `json:"placementName"`
			DeviceTemplateName string `json:"deviceTemplateName"`
		} `json:"attributes"`
	} `json:"deviceInfo"`
	DeviceEvent struct {
		ButtonClicked struct {
			ClickType    string    `json:"clickType"`
			ReportedTime time.Time `json:"reportedTime"`
		} `json:"buttonClicked"`
	} `json:"deviceEvent"`
	PlacementInfo struct {
		ProjectName   string `json:"projectName"`
		PlacementName string `json:"placementName"`
		Attributes    struct {
			DestinationPhoneNumber string `json:"DestinationPhoneNumber"`
			ContactFlowId          string `json:"ContactFlowId"`
		} `json:"attributes"`
		Devices struct {
			SampleRequest string `json:"Sample-Request"`
		} `json:"devices"`
	} `json:"placementInfo"`
}

type MyResponse struct {
	Message   string `json:"states"`
	ClickType string `json:"clickType"`
}

func merry(event IotButtonEvent) (MyResponse, error) {

	outboundInout := connect.StartOutboundVoiceContactInput{}
	outboundInout.
		SetContactFlowId(event.PlacementInfo.Attributes.ContactFlowId).
		SetDestinationPhoneNumber(event.PlacementInfo.Attributes.DestinationPhoneNumber).
		SetSourcePhoneNumber(os.Getenv("SourcePhoneNumber")).
		SetInstanceId(os.Getenv("InstanceId"))

	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("Region"))})

	if err != nil {
		return MyResponse{Message: "faild", ClickType: event.DeviceEvent.ButtonClicked.ClickType}, err
	}

	svc := connect.New(sess)

	outboundOutput, err := svc.StartOutboundVoiceContact(&outboundInout)

	if err != nil {
		return MyResponse{Message: fmt.Sprintf("faild ContactId:%s", outboundOutput.ContactId), ClickType: event.DeviceEvent.ButtonClicked.ClickType}, err
	}

	return MyResponse{Message: "success", ClickType: event.DeviceEvent.ButtonClicked.ClickType}, nil
}

func main() {
	lambda.Start(merry)
}
