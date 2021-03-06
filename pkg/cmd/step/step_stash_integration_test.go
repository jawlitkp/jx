// +build integration

package step_test

import (
	"github.com/jenkins-x/jx/pkg/cmd/step"
	"github.com/jenkins-x/jx/pkg/cmd/testhelpers"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/gits"
	helm_test "github.com/jenkins-x/jx/pkg/helm/mocks"
	"github.com/jenkins-x/jx/pkg/tests"

	"github.com/stretchr/testify/assert"
)

func TestStepStash(t *testing.T) {
	originalJxHome, tempJxHome, err := testhelpers.CreateTestJxHomeDir()
	assert.NoError(t, err)
	defer func() {
		err := testhelpers.CleanupTestJxHomeDir(originalJxHome, tempJxHome)
		assert.NoError(t, err)
	}()
	originalKubeCfg, tempKubeCfg, err := testhelpers.CreateTestKubeConfigDir()
	assert.NoError(t, err)
	defer func() {
		err := testhelpers.CleanupTestKubeConfigDir(originalKubeCfg, tempKubeCfg)
		assert.NoError(t, err)
	}()

	tempDir, err := ioutil.TempDir("", "test-step-collect")
	assert.NoError(t, err)

	testData := "test_data/step_collect/junit.xml"

	o := &step.StepStashOptions{
		StepOptions: opts.StepOptions{
			CommonOptions: &opts.CommonOptions{},
		},
	}
	o.StorageLocation.Classifier = "tests"
	o.StorageLocation.BucketURL = "file://" + tempDir
	o.ToPath = "output"
	o.Pattern = []string{testData}
	o.ProjectGitURL = "https://github.com/jenkins-x/dummy-repo.git"
	o.ProjectBranch = "master"
	testhelpers.ConfigureTestOptions(o.CommonOptions, &gits.GitFake{}, helm_test.NewMockHelmer())

	err = o.Run()
	assert.NoError(t, err)

	generatedFile := filepath.Join(tempDir, o.ToPath, testData)
	assert.FileExists(t, generatedFile)

	tests.AssertTextFileContentsEqual(t, testData, generatedFile)
}
