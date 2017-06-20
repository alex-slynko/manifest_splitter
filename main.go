package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 3 {
		os.Exit(1)
	}
	manifest := manifestReader("manifest.yml")
	smallManifest := manifestReader("small_manifest.yml")

	fmt.Printf("%#v", manifest)
	fmt.Printf("%#v", smallManifest)

}

func manifestReader(name string) map[string]interface{} {
	var manifest map[string]interface{}
	//var smallManifest map[string]interface{}
	file, _ := os.Open(name)
	b, _ := ioutil.ReadAll(file)

	yaml.Unmarshal(b, &manifest)
	return manifest
}
