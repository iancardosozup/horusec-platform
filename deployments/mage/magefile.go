// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build mage
// +build mage

// Horusec-Plaform mage file.
package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	// mage:import
	_ "github.com/ZupIT/horusec-devkit/pkg/utils/mageutils"
)

//UpdatePackageJson updates packge.json lock horusec version
func UpdatePackageJson(version string) error {
	if err := sh.RunV("npm", "install", "-g", "json"); err != nil {
		return err
	}
	return sh.RunV(fmt.Sprintf("json", "-I", "-f", "./manager/package.json", "-e", `'this.version="%s"'`, version))
}

func UpdateHorusecVersionInProject(actualVersion, releaseVersion string) error {
	return sh.RunV(fmt.Sprintf(
		`find ./ -type f -not -path "./.git/*" -not -path "./Makefile" -not -path "./manager/cypress/*" -not -name "*.sum" -not -name "*.mod"| xargs  sed -t -i "s/%s/%s/g"`,
		actualVersion, releaseVersion))

}

func DockerPushPlatformGoProjects(tag string) error {
	for _, image := range getImages() {
		if err := sh.RunV("docker", "push", fmt.Sprintf("%s:%s", image, tag)); err != nil {
			return fmt.Errorf("failed to push %s with error %s", fmt.Sprintf("%s:%s", image, tag), err.Error())
		}
	}
	return nil
}

func DockerSignPlatformGoProjects(tag string) error {
	mg.Deps(hasAllNecessaryEnvs, isCosignInstalled)
	err := os.MkdirAll("./tmp", 0700)
	if err != nil {
		return err
	}
	err = os.WriteFile("./tmp/cosign.key", []byte(os.Getenv("COSIGN_KEY")), 0700)
	if err != nil {
		return err
	}
	for _, image := range getImages() {
		if err := sh.RunV("cosign", "sign", "-key", "./tmp/cosign.key", fmt.Sprintf("%s:%s", image, tag)); err != nil {
			return fmt.Errorf("failed to sign %s with error %s", fmt.Sprintf("%s:%s", image, tag), err.Error())
		}
	}
	return nil
}

const (
	ImageMessages      = "horuszup/horusec-messages"
	ImageWebhook       = "horuszup/horusec-webhook"
	ImageAuth          = "horuszup/horusec-auth"
	ImageAnalytic      = "horuszup/horusec-analytic"
	ImageVulnerability = "horuszup/horusec-vulnerability"
	ImageMigrations    = "horuszup/horusec-migrations"
	ImageCore          = "horuszup/horusec-core"
	ImageApi           = "horuszup/horusec-api"
)

func getImages() []string {
	return []string{
		ImageMessages,
		ImageWebhook,
		ImageAuth,
		ImageAnalytic,
		ImageVulnerability,
		ImageMigrations,
		ImageCore,
		ImageApi,
	}
}
func isCosignInstalled() error {
	return sh.RunV("cosign", "version")
}

func hasAllNecessaryEnvs() error {
	var result []string

	for k, v := range getConsingEnvs() {
		if v == "" {
			result = append(result, k)
		}
	}

	if len(result) != 0 {
		return fmt.Errorf("missing some env var: %v", result)
	}
	if err := os.Setenv("COSIGN_PASSWORD", os.Getenv("COSIGN_PWD")); err != nil {
		return err
	}

	return nil
}

func getConsingEnvs() map[string]string {
	return map[string]string{
		"COSIGN_PWD": os.Getenv("COSIGN_PWD"),
		"COSIGN_KEY": os.Getenv("COSIGN_KEY"),
	}

}
