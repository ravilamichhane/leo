package forms

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ravilmc/leo/tygo"
)

func generate(name string, methodname string) {

	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Println("No mod file found:", err)
		return
	}
	defer file.Close()

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	basepackageName := ""

	// Read the first line
	if scanner.Scan() {
		firstLine := scanner.Text()
		withoutModule := strings.Replace(firstLine, "module", "", 1)
		basepackageName = strings.TrimSpace(withoutModule)
	}

	// Check for scanning error
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	if basepackageName == "" {
		log.Println("Invalid package")
		return
	}

	gen := tygo.New(&tygo.Config{
		Packages: []*tygo.PackageConfig{
			{
				Path:       basepackageName + "/app/controllers/" + strings.ToLower(name) + "controller",
				OutputPath: "resources/services/" + strings.ToLower(name) + "/" + methodname + "Form.tsx",
				Frontmatter: `
import { Fetcher, handleResponseError } from "./fetcher";
import {useForm} from "react-hook-form";
import { SafeParse } from "./safeparse";
import { generateFormData } from "./formdata";
import { useMutation } from "@tanstack/react-query";
import { Form } from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { DateInput } from "@/components/forms/date-input";
import { HTMLInput } from "@/components/forms/html-input";
import { ImageInput } from "@/components/forms/image-upload";
import { SelectInput } from "@/components/forms/select-input";
import { SwitchInput } from "@/components/forms/switch-input";
import { TextInput } from "@/components/forms/text-input";

				`,
				TypeMappings: map[string]string{
					"time.Time": "string",
				},
				ExcludeFiles: []string{
					strings.ToLower(name) + "controller.go",
				},
			},
		},
	})
	err = gen.GenerateForm(methodname)

	if err != nil {
		panic(err)
	}

	cmd := exec.Command("./node_modules/.bin/prettier", "resources/services/"+strings.ToLower(name)+"/"+methodname+"Form.tsx", "--write")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))

}
