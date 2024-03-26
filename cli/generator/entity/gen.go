//go:generate go run .
//go:generate gofmt -w ../../../app/domain/entity

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"generator/transform"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	_ "embed"
)

const OutputBaseId = "../../../app/domain/entity"
const InputDir = "../../../docs/entity/**/"

type StructField struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Nullable bool   `yaml:"nullable"`
	Number   int    `yaml:"number"`
}

type StructInfo struct {
	Name    string                 `yaml:"name"`
	Package string                 `yaml:"package"`
	Fields  map[string]StructField `yaml:"structure"`
	Primary string                 `yaml:"primary"`
	Index   []string               `yaml:"index"`
}

type VoStructInfo struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Type    string `yaml:"type"`
}

//go:embed template.txt
var templateCode string

var voToPackages map[string]*VoStructInfo

func main() {
	log.Printf("Start entity/gen.go. output base directory: %s, input directory: %s", OutputBaseId, InputDir)

	yamlFiles, err := filepath.Glob(InputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return
	}

	log.Printf("Etity yaml file %d count.", len(yamlFiles))

	voToPackages = buildVoToPackages()

	for _, yamlFile := range yamlFiles {
		err := generateEntity(yamlFile, OutputBaseId)
		if err != nil {
			log.Fatalf("Error generating entity from YAML file %s: %v", yamlFile, err)
		}
	}

	log.Println("End entity/gen.go")
}

func buildVoToPackages() map[string]*VoStructInfo {
	voToPackages := map[string]*VoStructInfo{}
	inputDir := "../../../docs/vo/**/"

	yamlFiles, err := filepath.Glob(inputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return voToPackages
	}

	for _, yamlFile := range yamlFiles {
		structInfo, err := getVoStructInfo(yamlFile)
		if err != nil {
			log.Fatalf("Error generating vo from YAML file %s: %v", yamlFile, err)
		}

		voToPackages[fmt.Sprintf("%s.%s", structInfo.Package, structInfo.Name)] = structInfo
	}

	return voToPackages
}

func getVoStructInfo(yamlFile string) (*VoStructInfo, error) {
	yamlData, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file %s: %v", yamlFile, err)
	}

	var structInfo VoStructInfo
	if err := yaml.Unmarshal(yamlData, &structInfo); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML in file %s: %v", yamlFile, err)
	}

	return &structInfo, nil
}

func generateEntity(yamlFile string, outputBaseDir string) error {
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
		return fmt.Errorf("outputFilename file %s create error : %v", outputFileName, err)
	}

	if err := generateTemplate(structInfo, outputFile); err != nil {
		return fmt.Errorf("failed to generateTemplate: %v", err)
	}

	log.Printf("Created %s Entity in %s\n", structInfo.Name, outputFileName)

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
	tmpl, err := template.New("structTemplate").Funcs(template.FuncMap{
		"sortByNumber": sortByNumber,
	}).Parse(templateCode)
	if err != nil {
		return err
	}

	if err := tmpl.ExecuteTemplate(outputFile, "structTemplate", struct {
		Name       string
		SnakeName  string
		PluralName string
		Package    string
		Imports    []string
		Fields     map[string]StructField
	}{
		Name:       structInfo.Name,
		SnakeName:  transform.UpperCamelToSnake(structInfo.Name),
		PluralName: transform.SingularToPlural(structInfo.Name),
		Package:    structInfo.Package,
		Imports:    buildImport(structInfo.Fields),
		Fields:     structInfo.Fields,
	}); err != nil {
		return err
	}

	return nil
}

func sortByNumber(fields map[string]StructField) []struct {
	Name            string
	FieldInfo       StructField
	TypeWithPointer string
	Json            string
	Gorm            string
} {
	var sortedFields []struct {
		Name            string
		FieldInfo       StructField
		TypeWithPointer string
		Json            string
		Gorm            string
	}

	for name, fieldInfo := range fields {
		json := fieldInfo.Name
		gorm := fmt.Sprintf("column:%s", fieldInfo.Name)
		voStruct, voOk := voToPackages[fieldInfo.Type]
		if voOk {
			voPrimitiveType := ""
			switch voStruct.Type {
			case "int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64":
				voPrimitiveType = "int"
			case "uuid.UUID":
				voPrimitiveType = "string"
			default:
				voPrimitiveType = voStruct.Type
			}
			json += fmt.Sprintf(",%s", voPrimitiveType)

			gormType := voPrimitiveType
			switch gormType {
			case "string":
				gormType = "varchar(255)"
			}
			gorm += fmt.Sprintf(";type:%s", gormType)
		}

		sortedFields = append(sortedFields, struct {
			Name            string
			FieldInfo       StructField
			TypeWithPointer string
			Json            string
			Gorm            string
		}{
			Name:            name,
			FieldInfo:       fieldInfo,
			TypeWithPointer: getTypeWithPointer(fieldInfo),
			Json:            json,
			Gorm:            gorm,
		})
	}

	sort.SliceStable(sortedFields, func(i, j int) bool {
		return fields[sortedFields[i].Name].Number < fields[sortedFields[j].Name].Number
	})

	return sortedFields
}

func getTypeWithPointer(fieldInfo StructField) string {
	if fieldInfo.Nullable {
		return "*" + fieldInfo.Type
	}

	return fieldInfo.Type
}

func buildImport(structFileds map[string]StructField) []string {
	imports := []string{}

	for _, structField := range structFileds {
		types := strings.Split(structField.Type, ".")
		if len(types) <= 1 {
			continue
		}

		switch types[0] {
		case "time":
			imports = append(imports, "time")
		case "uuid":
			imports = append(imports, "github.com/google/uuid")
		default:
			imports = append(imports, fmt.Sprintf("app/domain/value/%s", types[0]))
		}
	}

	slices.Sort(imports)
	return slices.Compact(imports)
}
