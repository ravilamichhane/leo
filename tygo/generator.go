package tygo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Generator for one or more input packages, responsible for linking
// them together if necessary.
type Tygo struct {
	conf *Config

	packageGenerators map[string]*PackageGenerator
}

// Responsible for generating the code for an input package
type PackageGenerator struct {
	conf    *PackageConfig
	pkg     *packages.Package
	GoFiles []string
}

func New(config *Config) *Tygo {
	return &Tygo{
		conf:              config,
		packageGenerators: make(map[string]*PackageGenerator),
	}
}

func (g *Tygo) SetTypeMapping(goType string, tsType string) {
	for _, p := range g.conf.Packages {
		p.TypeMappings[goType] = tsType
	}
}

func (g *Tygo) Generate() error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles,
	}, g.conf.PackageNames()...)
	if err != nil {
		return err
	}

	for i, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return fmt.Errorf("%+v", pkg.Errors)
		}

		if len(pkg.GoFiles) == 0 {
			return fmt.Errorf("no input go files for package index %d", i)
		}

		pkgConfig := g.conf.PackageConfig(pkg.ID)

		pkgGen := &PackageGenerator{
			conf:    pkgConfig,
			GoFiles: pkg.GoFiles,
			pkg:     pkg,
		}
		g.packageGenerators[pkg.PkgPath] = pkgGen
		code, err := pkgGen.Generate()
		if err != nil {
			return err
		}

		outPath := pkgGen.conf.ResolvedOutputPath(filepath.Dir(pkg.GoFiles[0]))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = ioutil.WriteFile(outPath, []byte(code), 0664)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (g *Tygo) GenerateForm() error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedFiles,
	}, g.conf.PackageNames()...)
	if err != nil {
		return err
	}

	for i, pkg := range pkgs {
		if len(pkg.Errors) > 0 {
			return fmt.Errorf("%+v", pkg.Errors)
		}

		if len(pkg.GoFiles) == 0 {
			return fmt.Errorf("no input go files for package index %d", i)
		}

		pkgConfig := g.conf.PackageConfig(pkg.ID)

		pkgGen := &PackageGenerator{
			conf:    pkgConfig,
			GoFiles: pkg.GoFiles,
			pkg:     pkg,
		}
		g.packageGenerators[pkg.PkgPath] = pkgGen
		code, err := pkgGen.GenerateForm()
		if err != nil {
			return err
		}

		outPath := pkgGen.conf.ResolvedOutputPath(filepath.Dir(pkg.GoFiles[0]))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		path := strings.Split(cwd, "app/routes")
		if len(path) < 2 {
			log.Println("Invalid path")
		}

		basePath := path[0]

		err = os.WriteFile(outPath, []byte(code), 0664)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		cmd := exec.Command(basePath+"./node_modules/.bin/prettier", pkgGen.conf.OutputPath, "--write")
		stdout, err := cmd.Output()
		// Print the output
		fmt.Println(string(stdout))
		if err != nil {

			return nil
		}

	}
	return nil
}
