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
	"github.com/pkg/errors"
	"testing"
)

func TestSeparateImageName(t *testing.T) {
	patterns := []struct {
		imageName string
		expected  ImageName
		err       error
	}{
		{"nginx", ImageName{"nginx", "latest"}, nil},
		{"docker.io/nginx", ImageName{"docker.io/nginx", "latest"}, nil},
		{"docker.io/library/nginx", ImageName{"docker.io/library/nginx", "latest"}, nil},
		{"esaka/cowsay", ImageName{"esaka/cowsay", "latest"}, nil},
		{"kubernetesui/dashboard:v2.0.0-beta6", ImageName{"kubernetesui/dashboard", "v2.0.0-beta6"}, nil},
		{"kubernetesui/metrics-scraper:v1.0.2", ImageName{"kubernetesui/metrics-scraper", "v1.0.2"}, nil},

		// Failed Case
		{"wrong:wrong:worng", ImageName{}, errors.New("image format is wrong")},

		// TODO: incorrect image name or tag are checked when docker pull. we should introduce regex check.
		{"gafas gafas fas", ImageName{"gafas gafas fas", "latest"}, nil},
		{"ほげ", ImageName{"ほげ", "latest"}, nil},
	}

	for idx, pattern := range patterns {
		actual, err := SeparateImageName(pattern.imageName)
		if err != nil {
			if pattern.err.Error() != err.Error() {
				t.Errorf("pattern %d: want %v, actual %v", idx, pattern.err.Error(), err.Error())
			}
		} else if pattern.expected != actual {
			t.Errorf("pattern %d: want %v, actual %v", idx, pattern.expected, actual)
		}
	}
}
