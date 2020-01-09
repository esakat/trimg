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
package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/esakat/trimg/pkg"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"os"
	"sync"
)

var (
	filename string
	dryRun   bool
)

// transferCmd represents the transfer command
var transferCmd = &cobra.Command{
	Use:   "transfer <imagename>",
	Short: "pull image from external registry, and then create ECR repository, finally push it into ECR",
	Long: `transfer subcommand execute multiple process a batch
firstly it pull docker images from external registry, e.g. DockerHub,
and then create ECR repository, retag images for ECR, finally push into it into ECR

you can specify image path by manual and from kubernetes manifest

Specify image paths by manual:
  trimg transfer nginx:latest redis golang:1.13.5

Get image paths from kubernetes manifest:
  trimg transfer -f kubernetes-manifest.yml

`,
	Run: func(cmd *cobra.Command, args []string) {

		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			fmt.Printf("you should do `export AWS_DEFAULT_REGION=...`")
			os.Exit(1)
		}

		if accountId == "" {
			svc := sts.New(session.New(&aws.Config{Region: aws.String(region)}))
			t, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
			accountId = *t.Account
		}

		// run image transfer
		if filename == "" {
			if len(args) == 0 {
				fmt.Printf("You should set image paths")
				os.Exit(1)
			}

			args = removeDuplicateImage(args)

			if dryRun {
				fmt.Println("following images will be transfer")
				for _, imagePath := range args {
					newImagePath := pkg.ConvertImagePathForECR(imagePath, region, accountId)
					fmt.Printf("%s -> %s\n", imagePath, newImagePath)
				}
			} else {
				var wg sync.WaitGroup
				p := mpb.New(mpb.WithWaitGroup(&wg))
				steps, numBars := 5, len(args)
				wg.Add(numBars)

				resultMsg := make(chan string, len(args))

				for _, imagePath := range args {
					name := fmt.Sprintf("[%s]", imagePath)
					bar := p.AddBar(int64(steps),
						mpb.PrependDecorators(
							decor.Name(name),
						),
						mpb.AppendDecorators(
							decor.Percentage(decor.WCSyncSpace),
						),
					)
					go pkg.ImageTransfer(imagePath, region, accountId, &wg, bar, resultMsg)
				}
				// wait all task finish
				wg.Wait()

				// output result
				for i, _ := range args {
					msg := <-resultMsg
					fmt.Printf("%d: %s\n", i+1, msg)
				}
			}
		} else {
			// parse yaml file
			yamls, err := pkg.ParseMultiDocYaml(filename)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			// get image paths to transfer
			var imagePaths []string
			for _, y := range yamls {
				images, err := pkg.GetUsingImages(y)
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
				imagePaths = append(imagePaths, images...)
			}

			imagePaths = removeDuplicateImage(imagePaths)

			// if dryRun "true", just output target image paths
			if dryRun {
				fmt.Println("following images will be transfer")
				for _, imagePath := range imagePaths {
					newImagePath := pkg.ConvertImagePathForECR(imagePath, region, accountId)
					fmt.Printf("%s -> %s\n", imagePath, newImagePath)
				}
			} else {
				var wg sync.WaitGroup
				p := mpb.New(mpb.WithWaitGroup(&wg))
				steps, numBars := 5, len(imagePaths)
				wg.Add(numBars)

				resultMsg := make(chan string, len(imagePaths))

				for _, imagePath := range imagePaths {
					name := fmt.Sprintf("[%s]", imagePath)
					bar := p.AddBar(int64(steps),
						mpb.PrependDecorators(
							decor.Name(name),
						),
						mpb.AppendDecorators(
							decor.Percentage(decor.WCSyncSpace),
						),
					)
					go pkg.ImageTransfer(imagePath, region, accountId, &wg, bar, resultMsg)
				}
				// wait all task finish
				wg.Wait()

				// output result
				for i, _ := range imagePaths {
					msg := <-resultMsg
					fmt.Printf("%d: %s\n", i+1, msg)
				}
			}
		}
	},
}

func removeDuplicateImage(images []string) []string {
	results := make([]string, 0, len(images))
	encountered := map[string]bool{}
	for i := 0; i < len(images); i++ {
		if !encountered[images[i]] {
			encountered[images[i]] = true
			results = append(results, images[i])
		}
	}
	return results
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.PersistentFlags().StringVar(&accountId, "account-id", "", "target of pushing images, default: your IAM AccountId")
	transferCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "specify kubernetes manifest filepath")
	transferCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "only print the object that would be replaced, without transfer it.")
}
