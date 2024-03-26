//go:generate go run .
//go:generate gofmt -w ../../../app/domain/enum

package main

import (
	"fmt"
	"generator/transform"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	_ "embed"
)

const OutputBaseDir = "../../../app/domain/enum"
const InputDir = "../../../docs/enum/"

type StructInfo struct {
	Name    string   `yaml:"name"`
	Package string   `yaml:"package"`
	Fields  []string `yaml:"structure"`
}

//go:embed template.txt
var templateCode string

func main() {
	log.Printf("Start enum/gen.go. output base directory: %s, input directory: %s", OutputBaseDir, InputDir)

	yamlFiles, err := filepath.Glob(InputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
	}

	log.Printf("Enum yaml file %d count.", len(yamlFiles))

	for _, yamlFile := range yamlFiles {
		err := generateEnum(yamlFile, OutputBaseDir)
		if err != nil {
			log.Fatalf("Error generating enum from YAML file %s: %v", yamlFile, err)
		}
	}

	log.Println("End enum/gen.go")
}

func generateEnum(yamlFile string, outputBaseDir string) error {
	structInfo, err := getStructInfo(yamlFile)
	if err != nil {
		return err
	}

	outputDir := filepath.Join(outputBaseDir, structInfo.Package)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error createing output directory %s: %v", outputDir, err)
	}

	outputFileName := filepath.Join(outputDir, fmt.Sprintf("%s.gen.go", transform.UpperCamelToSnake(structInfo.Name)))
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("outputFileName file %s create error: %v", outputFileName, err)
	}

	if err := generateTemplate(structInfo, outputFile); err != nil {
		return fmt.Errorf("failed to generateTemplate: %v", err)
	}

	log.Printf("Created %s Enum in %s\n", structInfo.Name, outputFileName)

	return nil
}

func getStructInfo(yamlFile string) (*StructInfo, error) {
	yamlData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file %s: %v", yamlFile, err)
	}

	var structInfo StructInfo
	if err := yaml.Unmarshal(yamlData, &structInfo); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML in file %s: %v", yamlFile, err)
	}

	return &structInfo, nil
}

func generateTemplate(structInfo *StructInfo, outputFile *os.File) error {
	tmpl, err := template.New("structTemplate").Parse(templateCode)
	if err != nil {
		return err
	}

	constList := make([]string, len(structInfo.Fields))
	_ = copy(constList, structInfo.Fields)
	constList[0] = fmt.Sprintf("%s %s = iota", constList[0], structInfo.Name)
	constBlock := strings.Join(constList, "\n\t")

	script := fmt.Sprintf(
		`const (
	%s
)

func (e %s) ToString() string {
	names := [...]string{%s}
	if e < %s || e > %s {
		return "Unknown"
	}

	return names[e]
}`,
		constBlock,
		structInfo.Name,
		`"`+strings.Join(structInfo.Fields, `", "`)+`"`,
		structInfo.Fields[0],
		structInfo.Fields[len(structInfo.Fields)-1],
	)

	if err := tmpl.ExecuteTemplate(outputFile, "structTemplate", struct {
		Name    string
		Package string
		Script  string
	}{
		Name:    structInfo.Name,
		Package: structInfo.Package,
		Script:  script,
	}); err != nil {
		return err
	}

	return nil
}
