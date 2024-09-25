package typeconverter

import (
	"os/exec"

	"github.com/ravilmc/leo/reactssr/packages/utils"
	_ "github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

// Start starts the type converter
// It gets the name of structs in PropsStructsPath and generates a temporary file to run the type converter
func Start(structsFilePath, generatedTypesPath string) error {
	// Get struct names from file
	structNames, err := getStructNamesFromFile(structsFilePath)
	if err != nil {
		return err
	}
	// Create a folder for the temporary generator files
	cacheDir, err := utils.GetTypeConverterCacheDir()
	if err != nil {
		return err
	}
	// Create the generator file
	temporaryFilePath, err := createTemporaryFile(structsFilePath, generatedTypesPath, cacheDir, structNames)
	if err != nil {
		return err
	}

	// Run the file
	cmd := exec.Command("go", "run", temporaryFilePath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
