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
	"errors"
	"strings"
)

type ImageName struct {
	RepositoryName string
	Tag            string
}

func SeparateImageName(imageName string) (ImageName, error) {
	separatedImageName := strings.Split(imageName, ":")
	switch len(separatedImageName) {
	case 1:
		return ImageName{
			RepositoryName: separatedImageName[0],
			Tag:            "latest",
		}, nil
	case 2:
		return ImageName{
			RepositoryName: separatedImageName[0],
			Tag:            separatedImageName[1],
		}, nil
	default:
		return ImageName{}, errors.New("image format is wrong")
	}
}
