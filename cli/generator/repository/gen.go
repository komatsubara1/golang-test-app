//go:generate go run .
//go:generate gofmt -w ../../../app/domain/repository

package main

import (
	"fmt"
	"generator/transform"
	"github.com/ahmetalpbalkan/go-linq"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	_ "embed"
)

const OutputBaseDir = "../../../app/domain/repository"
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
	Unique  []string               `yaml:"unique"`
}

type MethodType struct {
	Script string
}

type VoStructInfo struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Type    string `yaml:"type"`
}

type ImportInfo struct {
	Alias   string
	Package string
}

//go:embed template.txt
var templateCode string

var voToPackages map[string]*VoStructInfo

func main() {
	log.Printf("Start reposiroty/gen.go. output base directory: %s, input directory: %s", OutputBaseDir, InputDir)

	yamlFiles, err := filepath.Glob(InputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return
	}

	log.Printf("Entity yaml file %d count.", len(yamlFiles))

	voToPackages = buildVoToPackages()

	for _, yamlFile := range yamlFiles {
		err := generateRepository(yamlFile, OutputBaseDir)
		if err != nil {
			log.Fatalf("Error generating repository from YAML file %s: %v", yamlFile, err)
		}
	}

	log.Println("End repository/gen.go")
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

func generateRepository(yamlFile string, outputBaseDir string) error {
	structInfo, err := getStructInfo(yamlFile)
	if err != nil {
		return err
	}

	outputDir := filepath.Join(outputBaseDir, structInfo.Package)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error createing output directory %s: %v", outputDir, err)
	}

	outputFileName := filepath.Join(outputDir, fmt.Sprintf("%s_repository.gen.go", transform.UpperCamelToSnake(structInfo.Name)))
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("outputFilename file %s create error : %v", outputFileName, err)
	}

	if err := generateTemplate(structInfo, outputFile); err != nil {
		return fmt.Errorf("failed to generateTemplate: %v", err)
	}

	log.Printf("Created %s Repository in %s\n", structInfo.Name, outputFileName)

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
	tmpl, err := template.New("repositoryTemplate").Parse(templateCode)
	if err != nil {
		return fmt.Errorf("error parsing repository template: %v", err)
	}

	generatedMethods, imports := generateMethods(structInfo)

	if err := tmpl.ExecuteTemplate(outputFile, "repositoryTemplate", struct {
		Name    string
		Package string
		Imports []ImportInfo
		Mock    string
		Methods map[string]MethodType
	}{
		Name:    structInfo.Name,
		Package: structInfo.Package,
		Imports: imports,
		Mock:    generateMock(structInfo),
		Methods: generatedMethods,
	}); err != nil {
		return fmt.Errorf("template error: %v", err)
	}

	return nil
}

func generateMock(structInfo *StructInfo) string {
	return fmt.Sprintf("//go:generate mockgen -source=./%s_repository.gen.go -destination=./%s_repository_mock.gen.go -package=%s", transform.UpperCamelToSnake(structInfo.Name), transform.UpperCamelToSnake(structInfo.Name), structInfo.Package)
}

func generateMethods(structInfo *StructInfo) (map[string]MethodType, []ImportInfo) {
	methods := make(map[string]MethodType)
	imports := []ImportInfo{}

	// FindByPrimary
	primaries := strings.Split(structInfo.Primary, ",")
	script, primaryImports := generateFindByKey(structInfo, primaries, false)
	methods["Find"] = MethodType{
		Script: script,
	}

	if len(primaryImports) > 1 {
		imports = append(imports, primaryImports...)
	}

	// FindByUnique
	for _, unique := range structInfo.Unique {
		uniques := strings.Split(unique, ",")
		script, uniqueImports := generateFindByKey(structInfo, uniques, false)
		methods[fmt.Sprintf("FindBy%s", strings.Join(uniques, "And"))] = MethodType{
			Script: script,
		}
		if len(uniqueImports) > 1 {
			for _, uimp := range uniqueImports {
				var exists = linq.From(imports).AnyWith(func(imp interface{}) bool {
					return imp.(ImportInfo).Alias == uimp.Alias
				})
				if exists {
					continue
				}

				imports = append(imports, uimp)
			}
		}
	}

	// FindByIndex
	for _, index := range structInfo.Index {
		indexes := strings.Split(index, ",")
		script, indexImports := generateFindByKey(structInfo, indexes, true)
		methods[fmt.Sprintf("FindBy%s", strings.Join(indexes, "And"))] = MethodType{
			Script: script,
		}
		if len(indexImports) > 1 {
			for _, iimp := range indexImports {
				var exists = linq.From(imports).AnyWith(func(imp interface{}) bool {
					return imp.(ImportInfo).Alias == iimp.Alias
				})
				if exists {
					continue
				}

				imports = append(imports, iimp)
			}
		}
	}

	// Save
	script, saveImport := generateSave(structInfo)
	methods["Save"] = MethodType{
		Script: script,
	}
	var exists = linq.From(imports).AnyWith(func(imp interface{}) bool {
		return imp.(ImportInfo).Alias == saveImport.Alias
	})
	if !exists {
		imports = append(imports, saveImport)
	}

	return methods, imports
}

func generateFindByKey(structInfo *StructInfo, keys []string, isPlural bool) (string, []ImportInfo) {
	imports := []ImportInfo{}

	arguments, argImports := generateArguments(keys, structInfo)
	imports = append(imports, argImports...)

	ret, retImport := generateEntity(structInfo, isPlural)
	imports = append(imports, retImport)

	script := fmt.Sprintf("FindBy%s(%s) (*%s, error)", strings.Join(keys, "And"), strings.Join(arguments, ", "), ret)

	return script, imports
}

func generateSave(structInfo *StructInfo) (string, ImportInfo) {
	script, importInfo := generateEntity(structInfo, false)
	return fmt.Sprintf("Save(c *gin.Context, entity %s) error", script), importInfo
}

func generateArguments(strs []string, structInfo *StructInfo) ([]string, []ImportInfo) {
	imports := []ImportInfo{}
	arguments := []string{"c *gin.Context"}

	for _, key := range strs {
		arg, imp := generateArgument(structInfo, key)

		arguments = append(arguments, arg)
		if imp != nil {
			imports = append(imports, *imp)
		}
	}

	return arguments, imports
}

func generateArgument(structInfo *StructInfo, key string) (string, *ImportInfo) {
	val, ok := structInfo.Fields[key]
	if !ok {
		log.Fatalf("Not Exists Column %s in Primary", key)
	}

	valType := val.Type
	voStruct, voOk := voToPackages[val.Type]
	imp := &ImportInfo{}
	if voOk {
		imp = &ImportInfo{
			Alias:   fmt.Sprintf("%s_value", voStruct.Package),
			Package: fmt.Sprintf("app/domain/value/%s", voStruct.Package),
		}
		valType = strings.Replace(val.Type, ".", "_value.", 1)
	} else if val.Type == "time.Time" {
		imp = &ImportInfo{
			Alias:   "",
			Package: "time",
		}
	} else {
		imp = nil
	}

	return fmt.Sprintf("%s %s", val.Name, valType), imp
}

func generateEntity(structInfo *StructInfo, isPlural bool) (string, ImportInfo) {
	name := structInfo.Name
	if isPlural {
		name = transform.SingularToPlural(name)
	}
	ret := fmt.Sprintf("%s_entity.%s", structInfo.Package, name)
	importInfo := ImportInfo{
		Alias:   fmt.Sprintf("%s_entity", structInfo.Package),
		Package: fmt.Sprintf("app/domain/entity/%s", structInfo.Package),
	}

	return ret, importInfo
}
