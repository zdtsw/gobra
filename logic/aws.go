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
	//"sort"

	"github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	// 	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var serviceList = [...]string{
	"ec2",
	"ami",
	"iam",
	"vpc",
	"cost",
}

const awsProject string = "WenProjectName"

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
// TODO: impletment as interface
// type FilterOnTags interface {
//     filtering() ***
// }
func FilterOnTags(proj string) *ec2.DescribeInstancesInput {
	filters := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []*string{aws.String(proj)},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}
	return filters
}

func GetEC2Instances(sess *session.Session, filter *ec2.DescribeInstancesInput) ([]map[string]string, error) {
	ec2Svc := ec2.New(sess)

	result, err := ec2Svc.DescribeInstances(filter)
	var resultCollection []map[string]string

	if err != nil {
		fmt.Println("Error when retrieving information about EC2 instances:")
		errorHandler(err)
		return nil, err
	} else {

		for _, res := range result.Reservations {
			if res.RequesterId != nil {
				spaceOutputs(0, "EC2 Requester: "+*res.RequesterId)
			}

			for _, i := range res.Instances {
				instanceResult := make(map[string]string)
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
				resultCollection = append(resultCollection, instanceResult)
			}
		}
		return resultCollection, nil
	}
}


// EC2Handler godoc
// @Summary ec2
// @Description show information of ec2
// @Tags aws
// @Accept json
// @Produce html
// @Success 200 {string} string
// @Router /api/v1/aws/svc/ec2 [get]
func EC2Handler(c *gin.Context) {
	sess := GetSession()

	var filter *ec2.DescribeInstancesInput
	//only filter running/pending awsProject ec2
	filter = FilterOnTags(awsProject)

	ec2InfoList, err := GetEC2Instances(sess, filter)
	if err != nil {
		errorHandler(err)
	}
	renderResponse(c, gin.H{
		"version":     render.VersionPage,
		"author":      render.ContactAuthor,
		"title":       "AWS Service EC2",
		"project":     "WenProject",
		"ec2InfoList": ec2InfoList,
	}, "aws/ec2.tmpl")
}

////////////////////////////////////////////////AMI///////////////////////////////////
func FilterOnType(proj, env, build string) *ec2.DescribeImagesInput {
	filters := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Project"),
				Values: []*string{aws.String(proj)},
			},
			{
				Name:   aws.String("tag:Environment"),
				Values: []*string{aws.String(env)},
			},
			{
				Name:   aws.String("tag:Type"),
				Values: []*string{aws.String(build)},
			},
		},
		Owners: []*string{aws.String("1355687902")}, // account number
	}
	return filters
}

func sortAMI(images []*ec2.Image) []*ec2.Image {
	// TODO: implement sorting by time
	sortedImages := images
	// sort.Sort(imageSort(sortedImages))
	return sortedImages
}

func GetAMI(sess *session.Session, filters *ec2.DescribeImagesInput) ([]map[string]string, error) {
	ec2Svc := ec2.New(sess)
	resultAMI, err := ec2Svc.DescribeImages(filters)
	if err != nil {
		fmt.Println("Error when retrieving information about AMI:")
		errorHandler(err)
	}
	if len(resultAMI.Images) == 0 {
		fmt.Println("Could not find AMI matching filters:")
		return nil, nil
	}
	sortedImages := sortAMI(resultAMI.Images)
	var resultCollection []map[string]string

	spaceOutputs(1, "Sorted AMI:")
	for _, i := range sortedImages {
		ImagesInfo := make(map[string]string)
		spaceOutputs(4, *i.ImageId, *i.Name, *i.CreationDate)
		ImagesInfo["ImageId"] = *i.ImageId
		ImagesInfo["Name"] = *i.Name
		ImagesInfo["CreationDate"] = *i.CreationDate
		resultCollection = append(resultCollection, ImagesInfo)
	}
	return resultCollection, nil
}

// AMIHandler godoc
// @Summary ami
// @Description show information of ec2
// @Tags aws
// @Accept json
// @Produce html
// @Success 200 {string} string
// @Router /api/v1/aws/svc/ami [get]
func AMIHandler(c *gin.Context) {
	sess := GetSession()
	var filter *ec2.DescribeImagesInput
	filter = FilterOnType(awsProject, "product", "nightly")
	AMIList, err := GetAMI(sess, filter)
	if err != nil {
		errorHandler(err)
	}
	renderResponse(c, gin.H{
		"version": render.VersionPage,
		"author":  render.ContactAuthor,
		"title":   "AWS Service AMI",
		"project": "WenProject",
		"AMIList": AMIList,
	}, "aws/ami.tmpl")
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

//////////////////////////////////////////////IAM//////////////////////////////////////
func GetIAM(sess *session.Session) error {
	//TODO
	fmt.Print("nice")
	return nil
}

func IAMHandler(c *gin.Context) {
	sess := GetSession()
	err := GetIAM(sess)
	if err != nil {
		errorHandler(err)
	}
}

/////////////////////////////////////////////Budget//////////////////////////////////
func GetCost(sess *session.Session) error {
	//TODO
	fmt.Print("nice")
	return nil
}

func CostHandler(c *gin.Context) {
	sess := GetSession()
	err := GetCost(sess)
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
	case "ami":
		AMIHandler(c)
	case "vpc":
		VPCHandler(c)
	case "iam":
		IAMHandler(c)
	case "cost":
		CostHandler(c)
	default:
		fmt.Println("Invalid service in Gobra AWS")
	}
}
