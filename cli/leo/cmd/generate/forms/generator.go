package forms

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ravilmc/leo/tygo"
)

func generate() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := strings.Split(cwd, "app/routes")
	if len(path) < 2 {
		log.Println("Invalid path")
		return
	}

	basePath := path[0]

	file, err := os.Open(basePath + "go.mod")
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
				Path:       basepackageName + "/app/routes" + path[1],
				OutputPath: "app/form.tsx",
				Frontmatter: `
import { Fetcher } from "@/lib/fetcher";
import {useForm} from "react-hook-form";
import { SafeParse } from "@lib/utils";
import { handleResponseError } from "@/lib/handle-error";

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
			},
		},
	})
	err = gen.GenerateForm()

	if err != nil {
		panic(err)
	}

}
