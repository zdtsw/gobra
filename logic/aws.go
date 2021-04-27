package main

// module for "aws"

import (
	"github.com/gin-gonic/gin"
	//	"context"
	//	"flag"
	"fmt"
	"log"
	// "os"
	// "time"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	// 	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var serviceList = [...]string{
	"ec2",
	"vpc",
	"cost",
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
func GetSession() *session.Session {
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		errorHandler(err)
	}
	return sess
}

//////////////////////////////////////EC2////////////////////////////////////////////////////////////////
func FilterOnTags(name string) *ec2.DescribeInstancesInput {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []*string{aws.String(name)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	return params
}

func GetEC2Instances(sess *session.Session) ([]map[string]string, error) {
	ec2Svc := ec2.New(sess)

	var filter *ec2.DescribeInstancesInput
	//only filter running/pending kingston ec2
	filter = FilterOnTags("kingston")

	result, err := ec2Svc.DescribeInstances(filter)
	var resultCollection []map[string]string

	if err != nil {
		fmt.Println("Error when retrieving information about EC2 instances:")
		errorHandler(err)
		return nil, err
	} else {
		instanceResult := make(map[string]string)
		for _, res := range result.Reservations {
			if res.RequesterId != nil {
				spaceOutputs(0, "EC2 Requester: "+*res.RequesterId)
			}

			for _, i := range res.Instances {
				spaceOutputs(1, "Instance ID and Owner")
				instanceResult["InstanceId"] = *i.InstanceId
				instanceResult["OwnerId"] = *res.OwnerId
				spaceOutputs(4, *i.InstanceId, *res.OwnerId)

				spaceOutputs(1, "Instance Image and InstanceType:")
				instanceResult["InstanceAMI"] = *i.ImageId
				instanceResult["InstanceType"] = *i.InstanceType
				spaceOutputs(4, *i.ImageId, *i.InstanceType)

				spaceOutputs(1, "LaunchTime:")
				instanceResult["LaunchTime"] = i.LaunchTime.Format("2006-01-02 15:04:05")
				spaceOutputs(4, i.LaunchTime.Format("2006-01-02 15:04:05"))

				if i.StateReason != nil && i.State != nil {
					if *i.State.Code != int64(16) { //any not in running status
						spaceOutputs(1, "Instance State and Reason:")
						instanceResult["State"] = *i.State.Name
						instanceResult["Reason"] = *i.StateReason.Message
						spaceOutputs(4, *i.State.Name, *i.StateReason.Message)
					} else { // only show IP for running one
						spaceOutputs(1, "PrivateIP:")
						instanceResult["PrivateIpAddress"] = *i.PrivateIpAddress
						spaceOutputs(4, *i.PrivateIpAddress)
					}
				}
			}
			/* to have it copy to a temp map,
			otherwise in the next loop once set new value to instanceResult,
			the previous resultCollection changes*/
			temp := make(map[string]string)
			for index, element := range instanceResult {
				temp[index] = element
			}
			resultCollection = append(resultCollection, temp)

		}
		return resultCollection, nil
	}
}

func EC2Handler(c *gin.Context) {
	sess := GetSession()
	ec2InfoList, err := GetEC2Instances(sess)
	if err != nil {
		errorHandler(err)
	}
	renderResponse(c, gin.H{
		"version":     render.VersionPage,
		"author":      render.ContactAuthor,
		"title":       "AWS Service EC2",
		"project":     "Kingston",
		"ec2InfoList": ec2InfoList,
	}, "aws/ec2.tmpl")
}

////////////////////////////////////////////////VPC/////////////////////////////////////////////////////////////////////////////////////
func GetVPC(sess *session.Session) error {
	//TODO
	fmt.Print("nice")
	return nil
}

func VPCHandler(c *gin.Context) {
	sess := GetSession()
	err := GetVPC(sess)
	if err != nil {
		errorHandler(err)
	}
}

////////////////////////////////////////////////////////////Main///////////////////////////////////////////////////////////////////////////////////
func AWSServiceHandler(c *gin.Context) {
	renderResponse(c, gin.H{
		"version":     render.VersionPage,
		"author":      render.ContactAuthor,
		"title":       "AWS Services",
		"serviceList": serviceList,
	}, "aws/summary.tmpl")
}

func Dispatcher(c *gin.Context) {
	log4Caller()
	log4Debug()
	log.Println("Load page in path: " + c.Request.URL.Path)

	switch svcName := c.Param("service"); svcName {
	case "ec2":
		EC2Handler(c)
	case "vpc":
		VPCHandler(c)
	// case "budget":
	// 	BudgetHandler()
	default:
		fmt.Println("Invalid service in Gobra AWS")
	}
}
