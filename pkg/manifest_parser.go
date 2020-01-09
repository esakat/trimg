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
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

func ParseMultiDocYaml(filepath string) ([]map[interface{}]interface{}, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(f)

	var yamls []map[interface{}]interface{}

	var tmp map[interface{}]interface{}
	for (dec.Decode(&tmp)) == nil {
		yamls = append(yamls, tmp)
		tmp = nil
	}

	return yamls, nil
}

func GetUsingImages(manifest map[interface{}]interface{}) ([]string, error) {

	kind, ok := manifest["kind"]
	if !ok {
		return nil, errors.New("invalid format manifest")
	}

	var images []string

	switch kind {
	case "Deployment", "ReplicaSet", "StatefulSet", "Job":
		y, err := DigYaml(manifest, "spec", "template", "spec")
		if err != nil {
			return nil, err
		}
		podSpec, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("invalid format manifest")
		}
		for psKey, psValue := range podSpec {
			if psKey == "containers" || psKey == "initContainers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for _, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							images = append(images, cValue.(string))
						}
					}
				}
			}
		}

	case "Pod":
		for key, value := range manifest {
			if key == "spec" {
				podSpec, ok := value.(map[interface{}]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for psKey, psValue := range podSpec {
					if psKey == "containers" || psKey == "initContainers" {
						tmp, ok := psValue.([]interface{})
						if !ok {
							return nil, errors.New("invalid format manifest")
						}
						for _, tmpContainer := range tmp {
							container, ok := tmpContainer.(map[interface{}]interface{})
							if !ok {
								return nil, errors.New("invalid format manifest")
							}
							for cKey, cValue := range container {
								if cKey == "image" {
									images = append(images, cValue.(string))
								}
							}
						}
					}
				}
			}
		}
	case "CronJob":
		y, err := DigYaml(manifest, "spec", "jobTemplate", "spec", "template", "spec")
		if err != nil {
			return nil, err
		}
		podSpec, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("invalid format manifest")
		}
		for psKey, psValue := range podSpec {
			if psKey == "containers" || psKey == "initContainers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for _, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							images = append(images, cValue.(string))
						}
					}
				}
			}
		}
	default:
		// don't have images resorces
		return nil, nil
	}

	return images, nil
}

func ReplaceUsingImages(manifest map[interface{}]interface{}, region, accountId string) (map[interface{}]interface{}, error) {

	kind, ok := manifest["kind"]
	if !ok {
		return nil, errors.New("invalid format manifest")
	}

	switch kind {
	case "Deployment", "ReplicaSet", "StatefulSet", "Job":
		y, err := DigYaml(manifest, "spec", "template", "spec")
		if err != nil {
			return nil, err
		}
		podSpec, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("invalid format manifest")
		}
		for psKey, psValue := range podSpec {
			if psKey == "containers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for i, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							manifest["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["containers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
								ConvertImagePathForECR(cValue.(string), region, accountId)
						}
					}
				}
			} else if psKey == "initContainers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for i, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							manifest["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["initContainers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
								ConvertImagePathForECR(cValue.(string), region, accountId)
						}
					}
				}
			}
		}

	case "Pod":
		for key, value := range manifest {
			if key == "spec" {
				podSpec, ok := value.(map[interface{}]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for psKey, psValue := range podSpec {
					if psKey == "containers" {
						tmp, ok := psValue.([]interface{})
						if !ok {
							return nil, errors.New("invalid format manifest")
						}
						for i, tmpContainer := range tmp {
							container, ok := tmpContainer.(map[interface{}]interface{})
							if !ok {
								return nil, errors.New("invalid format manifest")
							}
							for cKey, cValue := range container {
								if cKey == "image" {
									manifest["spec"].(map[interface{}]interface{})["containers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
										ConvertImagePathForECR(cValue.(string), region, accountId)
								}
							}
						}
					} else if psKey == "initContainers" {
						tmp, ok := psValue.([]interface{})
						if !ok {
							return nil, errors.New("invalid format manifest")
						}
						for i, tmpContainer := range tmp {
							container, ok := tmpContainer.(map[interface{}]interface{})
							if !ok {
								return nil, errors.New("invalid format manifest")
							}
							for cKey, cValue := range container {
								if cKey == "image" {
									manifest["spec"].(map[interface{}]interface{})["initContainers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
										ConvertImagePathForECR(cValue.(string), region, accountId)
								}
							}
						}
					}
				}
			}
		}
	case "CronJob":
		y, err := DigYaml(manifest, "spec", "jobTemplate", "spec", "template", "spec")
		if err != nil {
			return nil, err
		}
		podSpec, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("invalid format manifest")
		}
		for psKey, psValue := range podSpec {
			if psKey == "containers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for i, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							manifest["spec"].(map[interface{}]interface{})["jobTemplate"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["containers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
								ConvertImagePathForECR(cValue.(string), region, accountId)
						}
					}
				}
			} else if psKey == "initContainers" {
				tmp, ok := psValue.([]interface{})
				if !ok {
					return nil, errors.New("invalid format manifest")
				}
				for i, tmpContainer := range tmp {
					container, ok := tmpContainer.(map[interface{}]interface{})
					if !ok {
						return nil, errors.New("invalid format manifest")
					}
					for cKey, cValue := range container {
						if cKey == "image" {
							manifest["spec"].(map[interface{}]interface{})["jobTemplate"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["template"].(map[interface{}]interface{})["spec"].(map[interface{}]interface{})["initContainers"].([]interface{})[i].(map[interface{}]interface{})["image"] =
								ConvertImagePathForECR(cValue.(string), region, accountId)
						}
					}
				}
			}
		}

	default:
		// don't have images resorces
	}

	return manifest, nil
}

func DigYaml(y interface{}, keys ...interface{}) (interface{}, error) {

	if len(keys) == 0 {
		return nil, errors.New("keys is required")
	}

	head := keys[0]
	t := reflect.TypeOf(head).Kind()

	if len(keys) == 1 {
		if t != reflect.String {
			return nil, errors.New("key should be string")
		}
		item, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("it's not yaml format")
		}
		for key, value := range item {
			if key == head.(string) {
				return value, nil
			}
		}
		errorMessage := fmt.Sprintf("key: %v not exists", head)
		return nil, errors.New(errorMessage)
	}

	// recursive
	switch t {
	case reflect.Int:
		// Array
		item, ok := y.([]interface{})
		if !ok {
			return nil, errors.New("it's not array")
		}
		return DigYaml(item[head.(int)], keys[1:]...)
	case reflect.String:
		// key
		item, ok := y.(map[interface{}]interface{})
		if !ok {
			return nil, errors.New("it's not yaml format")
		}
		for key, value := range item {
			if key == head.(string) {
				return DigYaml(value, keys[1:]...)
			}
		}

		errorMessage := fmt.Sprintf("key: %v not exists", head)
		return nil, errors.New(errorMessage)
	default:
		return nil, errors.New("keys should be string or int")
	}
}
