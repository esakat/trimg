/*
Copyright © 2020 esakat <esaka.tom@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package pkg

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	distreference "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/vbauerster/mpb"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

// convert image path
func ConvertImagePathForECR(imageName, region, accountId string) string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s", accountId, region, imageName)
}

// main func of transfer
func ImageTransfer(pullImageName, region, accountId string, wg *sync.WaitGroup, bar *mpb.Bar, resultMsg chan<- string) {

	defer wg.Done()

	// Step1. Pull Docker image from external registry.
	cl, err := client.NewEnvClient()
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}

	opts := types.ImagePullOptions{}
	ctx := context.Background()

	image, err := SeparateImageName(pullImageName)
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}

	resp, err := cl.ImagePull(ctx, pullImageName, opts)
	if err != nil {
		if err == distreference.ErrNameNotCanonical {
			imageNameLen := strings.Split(image.RepositoryName, "/")
			var pullingImagePrefix string
			switch len(imageNameLen) {
			case 1:
				pullingImagePrefix = "docker.io/library/"
			default:
				pullingImagePrefix = "docker.io/"
			}
			resp, err = cl.ImagePull(ctx, pullingImagePrefix+image.RepositoryName+":"+image.Tag, opts)
			if err != nil {
				resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
				return
			}
		} else {
			resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
			return
		}
	}
	defer resp.Close()

	jsonmessage.DisplayJSONMessagesStream(resp, ioutil.Discard, 0, false, nil)

	scanner := bufio.NewScanner(resp)
	for scanner.Scan() {
	}
	bar.Increment()

	// Step2. Create repository in ECR
	ecrSvc := ecr.New(session.New(&aws.Config{Region: aws.String(region)}))
	repositoryInfo := ecr.CreateRepositoryInput{
		RepositoryName: aws.String(image.RepositoryName),
	}
	_, err = ecrSvc.CreateRepository(&repositoryInfo)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() != ecr.ErrCodeRepositoryAlreadyExistsException {
				resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
				return
			}
		}
	}
	bar.Increment()

	// Step3. Get authorization for ECR
	loginAuth, err := ecrSvc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}

	decodedData, _ := base64.StdEncoding.DecodeString(*loginAuth.AuthorizationData[0].AuthorizationToken)
	decodedString := string(decodedData)

	// AuthorizationToken format is "user:password"
	authList := strings.Split(decodedString, ":")
	if len(authList) != 2 {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, "cannot get registry login token")
		return
	}
	username := authList[0]
	password := authList[1]
	serverAddress := *loginAuth.AuthorizationData[0].ProxyEndpoint

	auth := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: serverAddress,
	}
	_, err = cl.RegistryLogin(context.Background(), auth)
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}
	bar.Increment()

	// Step4. Tag image as ECR
	filtMap := map[string][]string{"reference": {image.RepositoryName + ":" + image.Tag}}
	filtBytes, _ := json.Marshal(filtMap)
	filt, err := filters.FromParam(string(filtBytes))
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}
	listOptions := types.ImageListOptions{
		All:     false,
		Filters: filt,
	}
	img, err := cl.ImageList(ctx, listOptions)
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}

	newImageTag := ConvertImagePathForECR(pullImageName, region, accountId)

	err = cl.ImageTag(ctx, img[0].ID, newImageTag)
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}
	bar.Increment()

	// Step5. Push image into ECR
	authJson := struct {
		Username string
		Password string
	}{
		Username: username,
		Password: password,
	}

	authBytes, _ := json.Marshal(authJson)
	authTokenBase64 := base64.StdEncoding.EncodeToString(authBytes)

	pushOpts := types.ImagePushOptions{
		RegistryAuth: authTokenBase64,
	}

	resp, err = cl.ImagePush(ctx, newImageTag, pushOpts)
	if err != nil {
		resultMsg <- fmt.Sprintf("%s failed to transfer. error message: %v", pullImageName, err)
		return
	}

	scanner = bufio.NewScanner(resp)
	for scanner.Scan() {
	}
	bar.Increment()
	resultMsg <- fmt.Sprintf("%s transfer to %s", pullImageName, newImageTag)

	// wait a few time, to display progress 100%
	time.Sleep(1 * time.Second)
}
