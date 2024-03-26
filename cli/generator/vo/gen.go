//go:generate go run .
//go:generate gofmt -w ../../../app/domain/value/

package main

import (
	"fmt"
	"generator/transform"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"

	_ "embed"
)

const OutputBaseDir = "../../../app/domain/value"
const InputDir = "../../../docs/vo/**/"

//go:embed template.txt
var templateCode string

type StructInfo struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
	Type    string `yaml:"type"`
}

func main() {
	log.Printf("Start vo/gen.go. output base directory: %s, input directory: %s", OutputBaseDir, InputDir)

	yamlFiles, err := filepath.Glob(InputDir + "*.yaml")
	if err != nil {
		log.Fatalf("Error finding YAML files: %v", err)
		return
	}

	log.Printf("Vo yaml file %d count.", len(yamlFiles))

	for _, yamlFile := range yamlFiles {
		err := generateVo(yamlFile, OutputBaseDir)
		if err != nil {
			log.Fatalf("Error generating vo from YAML file %s: %v", yamlFile, err)
		}
	}

	log.Printf("End vo/gen.go")
}

func generateVo(yamlFile string, outputBaseDir string) error {
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

	log.Printf("Created %s VO in %s\n", structInfo.Name, outputFileName)

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

	scanScript := generateScanScript(structInfo)
	unmarshalScript := generateUnmarshalScript(structInfo)

	if err := tmpl.ExecuteTemplate(outputFile, "structTemplate", struct {
		Name            string
		Package         string
		Imports         []string
		Type            string
		ScanScript      string
		UnmarshalScript string
	}{
		Name:            structInfo.Name,
		Package:         structInfo.Package,
		Imports:         buildImports(structInfo.Type),
		Type:            structInfo.Type,
		ScanScript:      scanScript,
		UnmarshalScript: unmarshalScript,
	}); err != nil {
		return err
	}

	return nil
}

func generateScanScript(structInfo *StructInfo) string {
	additionalCase := ""
	switch structInfo.Type {
	case "uuid.UUID":
		additionalCase = fmt.Sprintf(
			`case string:
			uuid, err := uuid.Parse(vt)
			if err != nil {
				return err
			}
			*v = New%s(uuid)
		case []uint8:
			uuid, err := uuid.ParseBytes(vt)
			if err != nil {
				return err
			}
			*v = New%s(uuid)`,
			structInfo.Name,
			structInfo.Name,
		)
		break
	case "uint64":
		additionalCase = fmt.Sprintf(
			`case int64: *v = New%s(uint64(vt))`,
			structInfo.Name,
		)
	}

	scanScript := fmt.Sprintf(
		`func (v *%s) Scan(value interface{}) error {
	switch vt := value.(type) {
	case %s:
		*v = New%s(vt)
	%s
	default:
		return fmt.Errorf("invalid type. type=%s", vt)
	}
	return nil
}`,
		structInfo.Name,
		structInfo.Type,
		structInfo.Name,
		additionalCase,
		"%s",
	)

	return scanScript
}

func generateUnmarshalScript(structInfo *StructInfo) string {
	// Numberをjsonパースするとfloat64にされるので変換用のswithを追加する
	additional := ""
	switch structInfo.Type {
	case "uint64", "int64":
		additional = fmt.Sprintf(`	switch t.(type) {
	case float64:
		v.Scan(%s(t.(float64)))
		return nil
}`, structInfo.Type)
	}

	script := fmt.Sprintf(
		`func (v *%s) UnmarshalJSON(data []byte) error {
	var t any
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	%s

	v.Scan(t)
	return nil
}`,
		structInfo.Name,
		additional,
	)

	return script
}

func buildImports(structType string) []string {
	imports := []string{
		"app/domain/value",
		"encoding/json",
		"fmt",
	}

	switch structType {
	case "time.Time":
		imports = append(imports, "time")
	case "uuid.UUID":
		imports = append(imports, "github.com/google/uuid")
	}

	return imports
}
