//go:generate go run .

package main

import (
	"fmt"
	"generator/transform"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"

	_ "embed"
)

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
	Foreign []string               `yaml:"foreign"`
}

type VoStructInfo struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Type    string `yaml:"type"`
}

//go:embed template.txt
var templateCode string

var voToTypes map[string]string

func main() {
	outputBaseDir := "../../../db/"
	inputDir := "../../../docs/entity/**/"

	log.Printf("Start ddl/gen.go. outputcd  base directory: %s, input directory: %s", outputBaseDir, inputDir)

	yamlFiles, err := filepath.Glob(inputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return
	}

	log.Printf("Entity yaml file %d count.", len(yamlFiles))

	voToTypes = buildVoToTypes()

	structInfos := map[string]StructInfo{}
	for _, yamlFile := range yamlFiles {
		structInfo, err := getStructInfo(yamlFile)
		if err != nil {
			log.Fatalf("Error structinfo: %v", err)
			return
		}

		structInfos[structInfo.Name] = *structInfo
	}

	for _, structInfo := range structInfos {
		err := generateDdl(structInfo.Name, structInfos, outputBaseDir)
		if err != nil {
			log.Fatalf("Error generating ddl %s: %v", structInfo.Name, err)
		}
	}

	log.Println("End ddl/gen.go")
}

func buildVoToTypes() map[string]string {
	voToTypes := map[string]string{}

	yamlFiles, err := filepath.Glob("../../../docs/vo/**/*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return voToTypes
	}

	for _, yamlFile := range yamlFiles {
		structInfo, err := getVoStructInfo(yamlFile)
		if err != nil {
			log.Fatalf("Error generating vo from YAML file %s: %v", yamlFile, err)
		}

		voToTypes[fmt.Sprintf("%s.%s", structInfo.Package, structInfo.Name)] = structInfo.Type
	}

	return voToTypes
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

func generateDdl(name string, structInfos map[string]StructInfo, outputBaseDir string) error {
	structInfo := structInfos[name]

	outputDir := filepath.Join(outputBaseDir, structInfo.Package)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("error createing output directory %s: %v", outputDir, err)
	}

	outputFileName := filepath.Join(outputDir, fmt.Sprintf("%s.gen.sql", transform.UpperCamelToSnake(structInfo.Name)))
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("outputFilename file %s create error : %v", outputFileName, err)
	}

	indexFields := []string{}
	primaryStrings := buildPrimaryStrings(&structInfo)
	indexFields = append(indexFields, primaryStrings)
	indexStrings := buildIndexStrings(&structInfo)
	indexFields = append(indexFields, indexStrings...)
	uniqueStrings := buildUniqueStrings(&structInfo)
	indexFields = append(indexFields, uniqueStrings...)
	foreignStrings := buildForeignStrings(&structInfo, structInfos)
	indexFields = append(indexFields, foreignStrings...)

	if err := generateTemplate(&structInfo, indexFields, outputFile); err != nil {
		return fmt.Errorf("faild to generateTemplate: %v", err)
	}

	log.Printf("Created %s Ddl in %s\n", structInfo.Name, outputFileName)

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

func buildPrimaryStrings(structInfo *StructInfo) string {
	primaryStrings := []string{}

	for _, column := range strings.Split(structInfo.Primary, ",") {
		for field := range structInfo.Fields {
			if column == field {
				primaryStrings = append(primaryStrings, fmt.Sprintf("`%s`", structInfo.Fields[field].Name))
			}
		}
	}

	return fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(primaryStrings, ", "))
}

func buildIndexStrings(structInfo *StructInfo) []string {
	indexStrings := []string{}

	for _, index := range structInfo.Index {
		names := []string{}
		parts := []string{}

		for _, column := range strings.Split(index, ",") {
			for field := range structInfo.Fields {
				if column == field {
					names = append(names, structInfo.Fields[field].Name)
					parts = append(parts, fmt.Sprintf("`%s`", structInfo.Fields[field].Name))
				}
			}
		}

		indexStrings = append(
			indexStrings,
			fmt.Sprintf("INDEX `idx_%s`(%s)", strings.Join(names, "_"), strings.Join(parts, ", ")),
		)
	}

	return indexStrings
}

func buildUniqueStrings(structInfo *StructInfo) []string {
	uniqueStrings := []string{}

	for _, unique := range structInfo.Unique {
		names := []string{}
		parts := []string{}

		for _, column := range strings.Split(unique, ",") {
			for field := range structInfo.Fields {
				if column == field {
					names = append(names, structInfo.Fields[field].Name)
					parts = append(parts, fmt.Sprintf("`%s`", structInfo.Fields[field].Name))
				}
			}
		}

		uniqueStrings = append(
			uniqueStrings,
			fmt.Sprintf("UNIQUE `uq_%s`(%s)", strings.Join(parts, "_"), strings.Join(parts, ", ")),
		)
	}

	return uniqueStrings
}

func buildForeignStrings(structInfo *StructInfo, structInfos map[string]StructInfo) []string {
	foreignStrings := []string{}

	for _, foreign := range structInfo.Foreign {
		fs := strings.Split(foreign, ",")
		if len(fs) > 2 {
			log.Fatalf("Error foreign key target %d", len(fs))
		}

		t := fs[0]
		f := fs[1]
		ff := strings.Split(f, ".")
		fff := structInfos[strings.Title(ff[0])]

		foreignStrings = append(
			foreignStrings,
			fmt.Sprintf(
				"CONSTRAINT `fk_%s` FOREIGN KEY (`%s`) REFERENCES `%s`(`%s`) ON DELETE CASCADE",
				transform.UpperCamelToSnake(structInfo.Name+t),
				transform.UpperCamelToSnake(t),
				transform.UpperCamelToSnake(fff.Name),
				strings.ToLower(ff[1]),
			),
		)
	}

	return foreignStrings
}

func generateTemplate(structInfo *StructInfo, indexFileds []string, outputFile *os.File) error {
	tmpl, err := template.New("structTemplate").Funcs(template.FuncMap{
		"sortByNumber": sortByNumber,
	}).Parse(templateCode)
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	if err := tmpl.ExecuteTemplate(outputFile, "structTemplate", struct {
		Table      string
		Fields     map[string]StructField
		PrimaryKey string
		IndexKey   string
		UniqueKey  string
	}{
		Table:    transform.UpperCamelToSnake(structInfo.Name),
		Fields:   structInfo.Fields,
		IndexKey: strings.Join(indexFileds, ",\n\t"),
	}); err != nil {
		return fmt.Errorf("template error: %v", err)
	}

	return nil
}

func sortByNumber(fields map[string]StructField) []struct {
	Name            string
	FieldInfo       StructField
	Column          string
	Type            string
	TypeWithPointer string
	Config          string
} {
	var sortedFields []struct {
		Name            string
		FieldInfo       StructField
		Column          string
		Type            string
		TypeWithPointer string
		Config          string
	}

	for name, fieldInfo := range fields {
		sortedFields = append(sortedFields, struct {
			Name            string
			FieldInfo       StructField
			Column          string
			Type            string
			TypeWithPointer string
			Config          string
		}{
			Name:            name,
			FieldInfo:       fieldInfo,
			Column:          fieldInfo.Name,
			Type:            getType(fieldInfo),
			TypeWithPointer: getTypeWithPointer(fieldInfo),
			Config:          getExtra(fieldInfo),
		})
	}

	sort.Slice(sortedFields, func(i, j int) bool {
		return fields[sortedFields[i].Name].Number < fields[sortedFields[j].Name].Number
	})

	return sortedFields
}

func getTypeWithPointer(fieldInfo StructField) string {
	if fieldInfo.Nullable {
		return " DEFAULT NULL"
	}

	return " NOT NULL"
}

func getType(field StructField) string {
	t := field.Type
	tt, ok := voToTypes[field.Type]
	if ok {
		t = tt
	}

	switch t {
	case "string", "uuid.UUID":
		return " VARCHAR(255)"
	case "int64":
		return " BIGINT"
	case "uint64":
		return " BIGINT UNSIGNED"
	case "int":
		return " INT"
	case "uint":
		return " INT UNSIGNED"
	case "bool":
		return " TINYINT"
	case "time.Time":
		return " DATETIME"
	}

	return ""
}

func getExtra(fieldInfo StructField) string {
	if fieldInfo.Name == "created_at" {
		return " DEFAULT CURRENT_TIMESTAMP"
	} else if fieldInfo.Name == "updated_at" {
		return " DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
	}

	return ""
}
