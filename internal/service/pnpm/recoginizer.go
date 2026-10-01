package pnpm

import (
	"encoding/json/v2"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/goccy/go-yaml"
)

type PnpmWorkspaceRoot struct {
	Packages []*WorkspacePackage
}

type WorkspacePackage struct {
	Name    string
	Version string
	Scripts map[string]string
}

func RecognizeWorkspace() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Could not get the current Directory!")
	}

	log.Println("Current Directory:", cwd)

	worksPaceRootPath, err := findWorkspaceRoot(cwd)
	if err != nil {
		log.Fatalln("Filed parsing workspace root:", err)
	}

	log.Println("Workspace Root:", worksPaceRootPath)

	_, err = constructWorkspace(worksPaceRootPath)
	if err != nil {
		log.Fatalln("Filed parsing workspace root:", err)
	}
}

func findWorkspaceRoot(cwd string) (string, error) {
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return "", err
	}

	rootCnt := 0
	rootPath := ""
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}

		if rootCnt > 1 {
			return "", errors.New("Multiple workspace roots found")
		}

		if entry.Name() == "pnpm-workspace.yaml" {
			rootPath = filepath.Join(cwd, "pnpm-workspace.yaml")
			rootCnt++
		}
	}

	if rootPath == "" {
		return "", errors.New("No workspace root found")
	}

	return rootPath, nil
}

type PnpmWorkspacePackages struct {
	Packages []string `yaml:"packages"`
}

type PackageJson struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Scipts  map[string]string `json:"scripts"`
}

func constructWorkspace(rootPath string) (PnpmWorkspaceRoot, error) {
	pnpmWorkspaceBytes, err := os.ReadFile(rootPath)
	if err != nil {
		return PnpmWorkspaceRoot{}, err
	}

	var pnpmWorkspacePackages PnpmWorkspacePackages

	if err := yaml.Unmarshal(pnpmWorkspaceBytes, &pnpmWorkspacePackages); err != nil {
		return PnpmWorkspaceRoot{}, err
	}

	log.Println("pnpmWorkspacePackages:", pnpmWorkspacePackages)

	packages := make([]string, 0, len(pnpmWorkspacePackages.Packages))
	for _, glob := range pnpmWorkspacePackages.Packages {
		files, err := doublestar.FilepathGlob(
			glob,
			doublestar.WithFilesOnly(),
			doublestar.WithFailOnIOErrors(),
		)
		if err != nil {
			return PnpmWorkspaceRoot{}, err
		}
		packages = append(packages, filterforPackages(files)...)
	}

	log.Println("packages:", packages)

	workspacePackages := make([]*WorkspacePackage, 0, len(packages))
	for _, packagePath := range packages {
		packageJson, err := readPackageJson(packagePath)
		if err != nil {
			return PnpmWorkspaceRoot{}, err
		}
		log.Println("packageJson:", packageJson)
		workspacePackages = append(workspacePackages, &WorkspacePackage{
			Name:    packageJson.Name,
			Version: packageJson.Version,
			Scripts: packageJson.Scipts,
		})
	}

	return PnpmWorkspaceRoot{
		Packages: workspacePackages,
	}, nil

}

func filterforPackages(files []string) []string {
	packages := make([]string, 0, len(files))
	for _, file := range files {
		if strings.Contains(file, "node_modules") {
			continue
		}

		if strings.HasSuffix(file, "package.json") {
			packages = append(packages, file)
		}
	}
	return packages
}

func readPackageJson(filePath string) (PackageJson, error) {
	packageJsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return PackageJson{}, err
	}

	var packageJson PackageJson
	if err := json.Unmarshal(packageJsonBytes, &packageJson); err != nil {
		return PackageJson{}, err
	}

	return packageJson, nil
}
