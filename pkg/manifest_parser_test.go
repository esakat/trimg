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
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestParseMultiDocYaml(t *testing.T) {
	expected1 := map[interface{}]interface{}{
		"test":   1,
		"sample": 2,
	}

	expected2 := map[interface{}]interface{}{
		"hoge": "foo",
		"mo":   true,
		"hogehoge": map[interface{}]interface{}{
			"hoge": 2,
			"fuga": 3,
		},
	}

	expected3 := map[interface{}]interface{}{
		"hoge": "hoge",
		"mo":   false,
		"hogehoge": map[interface{}]interface{}{
			"sample": "mooooo",
			"fuga":   3,
		},
	}

	actual, err := ParseMultiDocYaml("../testfiles/input/parse_test.yml")
	if err != nil {
		t.Fatalf("failed to parse")
	}

	if len(actual) != 3 {
		t.Fatalf("failed to parse")
	}

	if actual[0]["test"] != expected1["test"] {
		t.Fatalf("failed to parse")
	}
	if actual[0]["sample"] != expected1["sample"] {
		t.Fatalf("failed to parse")
	}

	if actual[1]["hoge"] != expected2["hoge"] {
		t.Fatalf("failed to parse")
	}
	if actual[1]["mo"] != expected2["mo"] {
		t.Fatalf("failed to parse")
	}
	if actual[1]["hogehoge"].(map[interface{}]interface{})["hoge"] != expected2["hogehoge"].(map[interface{}]interface{})["hoge"] {
		t.Fatalf("failed to parse")
	}
	if actual[1]["hogehoge"].(map[interface{}]interface{})["fuga"] != expected2["hogehoge"].(map[interface{}]interface{})["fuga"] {
		t.Fatalf("failed to parse")
	}

	if actual[2]["hoge"] != expected3["hoge"] {
		t.Fatalf("failed to parse")
	}
	if actual[2]["mo"] != expected3["mo"] {
		t.Fatalf("failed to parse")
	}
	if actual[2]["hogehoge"].(map[interface{}]interface{})["sample"] != expected3["hogehoge"].(map[interface{}]interface{})["sample"] {
		t.Fatalf("failed to parse")
	}
	if actual[2]["hogehoge"].(map[interface{}]interface{})["fuga"] != expected3["hogehoge"].(map[interface{}]interface{})["fuga"] {
		t.Fatalf("failed to parse")
	}
}

func TestGetUsingImagesDeployment(t *testing.T) {
	deployment, _ := ParseMultiDocYaml("../testfiles/input/deployment.yml")
	actualDeployment, err := GetUsingImages(deployment[0])
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to get deployment images")
	}
	expectedDeployment := []string{"nginx:latest", "nginx:latest", "initImage:latest"}
	if len(expectedDeployment) != len(actualDeployment) {
		t.Log(len(actualDeployment))
		t.Fatalf("failed to get deployment images")
	}

	sort.Strings(expectedDeployment)
	sort.Strings(actualDeployment)

	for i := range expectedDeployment {
		if expectedDeployment[i] != actualDeployment[i] {
			t.Fatalf("failed to get deployment images")
		}
	}
}

func TestGetUsingImagesPod(t *testing.T) {
	pod, _ := ParseMultiDocYaml("../testfiles/input/pod.yml")
	actualPod, err := GetUsingImages(pod[0])
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to get pod images")
	}
	expectedPod := []string{"nginx", "initPod:v2.0.0"}
	if len(expectedPod) != len(actualPod) {
		t.Log(len(actualPod))
		t.Fatalf("failed to get pod images")
	}

	sort.Strings(expectedPod)
	sort.Strings(actualPod)

	for i := range expectedPod {
		if expectedPod[i] != actualPod[i] {
			t.Fatalf("failed to get pod images")
		}
	}
}

func TestGetUsingImagesJob(t *testing.T) {
	job, _ := ParseMultiDocYaml("../testfiles/input/job.yml")
	actualJob, err := GetUsingImages(job[0])
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to get job images")
	}
	expectedJob := []string{"perl"}
	if len(expectedJob) != len(actualJob) {
		t.Log(len(actualJob))
		t.Fatalf("failed to get job images")
	}

	sort.Strings(expectedJob)
	sort.Strings(actualJob)

	for i := range expectedJob {
		if expectedJob[i] != actualJob[i] {
			t.Fatalf("failed to get job images")
		}
	}
}

func TestGetUsingImagesCronJob(t *testing.T) {
	cronJob, _ := ParseMultiDocYaml("../testfiles/input/cronjob.yml")
	actualCronJob, err := GetUsingImages(cronJob[0])
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to get cronJob images")
	}
	expectedCronJob := []string{"redis", "busybox"}
	if len(expectedCronJob) != len(actualCronJob) {
		t.Log(len(actualCronJob))
		t.Fatalf("failed to get cronJob images")
	}

	sort.Strings(expectedCronJob)
	sort.Strings(actualCronJob)

	for i := range expectedCronJob {
		if expectedCronJob[i] != actualCronJob[i] {
			t.Fatalf("failed to get cronJob images")
		}
	}
}

func TestGetUsingImagesReplicaset(t *testing.T) {
	replicaset, _ := ParseMultiDocYaml("../testfiles/input/replicaset.yml")
	actualReplicaset, err := GetUsingImages(replicaset[0])
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to get replicaset images")
	}
	expectedReplicaset := []string{"gcr.io/google_samples/gb-frontend:v3"}
	if len(expectedReplicaset) != len(actualReplicaset) {
		t.Log(len(actualReplicaset))
		t.Fatalf("failed to get replicaset images")
	}

	sort.Strings(expectedReplicaset)
	sort.Strings(actualReplicaset)

	for i := range expectedReplicaset {
		if expectedReplicaset[i] != actualReplicaset[i] {
			t.Fatalf("failed to get replicaset images")
		}
	}
}

func TestGetUsingImagesStatefulset(t *testing.T) {
	statefulset, _ := ParseMultiDocYaml("../testfiles/input/statefulset.yml")
	actualStatefulset, err := GetUsingImages(statefulset[0])
	if err != nil {
		t.Fatalf("failed to get statefulset images")
	}
	expectedStatefulset := []string{"k8s.gcr.io/nginx-slim:0.8"}
	if len(expectedStatefulset) != len(actualStatefulset) {
		t.Fatalf("failed to get statefulset images")
	}

	sort.Strings(expectedStatefulset)
	sort.Strings(actualStatefulset)

	for i := range expectedStatefulset {
		if expectedStatefulset[i] != actualStatefulset[i] {
			t.Fatalf("failed to get statefulset images")
		}
	}
}

func TestReplaceUsingImagesDeployment(t *testing.T) {
	deployment, _ := ParseMultiDocYaml("../testfiles/input/deployment.yml")
	actualManifest, err := ReplaceUsingImages(deployment[0], "ap-northeast-1", "111222333444")
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to replace manifest")
	}

	f, _ := os.Open("../testfiles/expected/replaceDeployment.yml")
	defer f.Close()

	dec := yaml.NewDecoder(f)

	var expectedManifest map[interface{}]interface{}
	dec.Decode(&expectedManifest)

	if !reflect.DeepEqual(actualManifest, expectedManifest) {
		t.Fatalf("expected: %v, got: %v", expectedManifest, actualManifest)
	}
}

func TestReplaceUsingImagesPod(t *testing.T) {
	pod, _ := ParseMultiDocYaml("../testfiles/input/pod.yml")
	actualManifest, err := ReplaceUsingImages(pod[0], "ap-northeast-1", "333222333444")
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to replace manifest")
	}

	f, _ := os.Open("../testfiles/expected/replacePod.yml")
	defer f.Close()

	dec := yaml.NewDecoder(f)

	var expectedManifest map[interface{}]interface{}
	dec.Decode(&expectedManifest)

	if !reflect.DeepEqual(actualManifest, expectedManifest) {
		t.Fatalf("expected: %v, got: %v", expectedManifest, actualManifest)
	}
}

func TestReplaceUsingImagesCronJob(t *testing.T) {
	cronjob, _ := ParseMultiDocYaml("../testfiles/input/cronjob.yml")
	actualManifest, err := ReplaceUsingImages(cronjob[0], "us-west-1", "999222333444")
	if err != nil {
		t.Fatalf(err.Error())
		t.Fatalf("failed to replace manifest")
	}

	f, _ := os.Open("../testfiles/expected/replaceCronjob.yml")
	defer f.Close()

	dec := yaml.NewDecoder(f)

	var expectedManifest map[interface{}]interface{}
	dec.Decode(&expectedManifest)

	if !reflect.DeepEqual(actualManifest, expectedManifest) {
		t.Fatalf("expected: %v, got: %v", expectedManifest, actualManifest)
	}
}

func TestDigYaml(t *testing.T) {
	testdata := map[interface{}]interface{}{
		"hoge": "foo",
		"mo":   true,
		"hogehoge": map[interface{}]interface{}{
			"hoge": 2,
			"fuga": 3,
		},
		"sample": []interface{}{
			map[interface{}]interface{}{
				"test": 1,
			},
			map[interface{}]interface{}{
				"foo": 2,
			},
			map[interface{}]interface{}{
				"mac": "testmsg",
			},
		},
	}

	actual1, err := DigYaml(testdata, "hogehoge", "hoge")
	if err != nil {
		t.Fatalf("failed to dig yaml: %v", err)
	}
	expected1 := 2

	if actual1 != expected1 {
		t.Fatalf("expected: %v, got: %v", expected1, actual1)
	}

	actual2, err := DigYaml(testdata, "sample", 2, "mac")
	if err != nil {
		t.Fatalf("failed to dig yaml: %v", err)
	}
	expected2 := "testmsg"

	if actual2 != expected2 {
		t.Fatalf("expected: %v, got: %v", expected2, actual2)
	}
}
