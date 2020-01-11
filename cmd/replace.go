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
	"gopkg.in/yaml.v2"
	"os"

	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace <filepath>",
	Short: "replace kubernetes manifest `image path` to ECR path",
	Long: `replace subcommand replace kubernetes manifest
get the value of the image from the manifest file and replace it to the path of the ECR will be sent by the transfer command
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("you can only specify one filepath")
			os.Exit(1)
		}

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

		// parse yaml file
		yamls, err := pkg.ParseMultiDocYaml(args[0])
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		result := ""
		for i, y := range yamls {
			var d []byte
			replacedManifest, err := pkg.ReplaceUsingImages(y, region, accountId)
			if err != nil {
				d, err = yaml.Marshal(y)
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
			} else {
				d, err = yaml.Marshal(replacedManifest)
				if err != nil {
					fmt.Printf("%v\n", err)
					os.Exit(1)
				}
			}
			if i != 0 {
				result += "---\n" + string(d)
			} else {
				result += string(d)
			}

		}

		fmt.Printf("%s", result)

	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)
	replaceCmd.PersistentFlags().StringVar(&accountId, "account-id", "", "target of pushing images, default: your IAM AccountId")
}
