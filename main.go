package main

import (
    "fmt"
    "io/ioutil"
    "os"

    "github.com/alex-slynko/manifest_splitter/maputil"
    yaml "gopkg.in/yaml.v2"
)

func main() {
    if len(os.Args) != 3 {
        printUsage()
        os.Exit(1)
    }
    manifest := manifestReader(os.Args[1])
    smallManifest := manifestReader(os.Args[2])

    operations, _ := maputil.ExtractOperations(manifest, smallManifest)
    // fmt.Println(operations)
    o, err := yaml.Marshal(&operations)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(o))
}

func manifestReader(name string) map[string]interface{} {
    var manifest map[string]interface{}
    //var smallManifest map[string]interface{}
    file, _ := os.Open(name)
    b, _ := ioutil.ReadAll(file)

    yaml.Unmarshal(b, &manifest)
    return manifest
}

func printUsage() {
    fmt.Println("Usage:")
    fmt.Println("manifest_splitter original_manifest minimal_manifest")
}
