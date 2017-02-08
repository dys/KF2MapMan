package kf2mapman

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Unknwon/goconfig"
	"github.com/stretchr/testify/assert"
)

var TestData = "testdata/"
var PreEditIni = TestData + "pre-edit.ini"
var EditedIni = TestData + "to-edit.ini"
var PostEditIni = TestData + "post-edit.ini"

func TestCreateSectionHeader(t *testing.T) {
	expected := fmt.Sprintf("MAP1 %s", MapSectionSuffix)
	results := CreateSectionHeader("MAP1")
	assert.Equal(t, expected, results)
}

func TestGetMapSections(t *testing.T) {
	cfg, _ := goconfig.LoadConfigFile(PreEditIni)
	expected := []string{"MAP1", "MAP2"}
	results := GetMapSections(cfg)
	assert.Equal(t, expected, results)
}

func TestGetMapCycle(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue(MapCycleSection, MapCycleOption,
		"(Maps=(\"MAP1\",\"MAP2\"))")
	expected := []string{"MAP1", "MAP2"}
	results := GetMapCycle(cfg)
	assert.Equal(t, expected, results)
}

func TestCreateMapCycle(t *testing.T) {
	expected := "(Maps=(\"MAP1\",\"MAP2\"))"
	results := CreateMapCycle([]string{"MAP1", "MAP2"})
	assert.Equal(t, expected, results)
}

func TestMapInCycle(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue(MapCycleSection, MapCycleOption,
		"(Maps=(\"MAP1\",\"MAP2\"))")
	assert.True(t, MapInCycle("MAP1", cfg))
	assert.False(t, MapInCycle("MAP", cfg))
}

func TestAddMapToCycle(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue(MapCycleSection, MapCycleOption,
		"(Maps=(\"MAP1\",\"MAP2\"))")
	expected := "(Maps=(\"MAP1\",\"MAP2\",\"MAP3\",\"MAP4\"))"
	AddMapToCycle("MAP3", cfg)
	AddMapToCycle("MAP4", cfg)
	results, _ := cfg.GetValue(
		MapCycleSection, MapCycleOption)
	assert.Equal(t, expected, results)
}

func TestDontAddDuplicateMaps(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue(MapCycleSection, MapCycleOption,
		"(Maps=(\"MAP1\",\"MAP2\"))")
	expected := "(Maps=(\"MAP1\",\"MAP2\",\"MAP3\"))"
	AddMapToCycle("MAP3", cfg)
	AddMapToCycle("MAP3", cfg)
	results, _ := cfg.GetValue(
		MapCycleSection, MapCycleOption)
	assert.Equal(t, expected, results)
}

func TestAddMapSection(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	name := "MAP1"
	section := fmt.Sprintf("%s %s", name, MapSectionSuffix)
	AddMapSection(name, cfg)
	resultsSection, _ := cfg.GetSection(section)
	assert.NotNil(t, resultsSection)
	resultsMapname, _ := cfg.GetValue(
		section, MapSectionMapOption)
	assert.Equal(t, name, resultsMapname)
	resultsScreenshot, _ := cfg.GetValue(
		section, MapSectionScreenshotOption)
	assert.Equal(t, MapSectionDefaultScreenshot, resultsScreenshot)
}

func TestAddMapsToConfig(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue(MapCycleSection, MapCycleOption,
		"(Maps=(\"MAP1\",\"MAP2\"))")
	names := []string{"MAP3", "MAP4"}
	AddMapsToConfig(names, cfg)
	for _, name := range names {
		section := fmt.Sprintf("%s %s", name, MapSectionSuffix)
		resultsSection, _ := cfg.GetSection(section)
		if assert.NotNil(t, resultsSection) {
			resultsMapname, _ := cfg.GetValue(
				section, MapSectionMapOption)
			assert.Equal(t, name, resultsMapname)
			resultsScreenshot, _ := cfg.GetValue(
				section, MapSectionScreenshotOption)
			assert.Equal(t, MapSectionDefaultScreenshot, resultsScreenshot)
			resultsCycle, _ := cfg.GetValue(
				MapCycleSection, MapCycleOption)
			assert.True(t, strings.Contains(resultsCycle, name))
		}
	}
}

func TestFileIsMap(t *testing.T) {
	assert.True(t, FileIsMap("MAP.kfm"))
	assert.True(t, FileIsMap("MAP.KFM"))
	assert.False(t, FileIsMap("MAP.TXT"))
	assert.False(t, FileIsMap("KFM.TXT"))
}

func TestStripMapExtension(t *testing.T) {
	assert.Equal(t, StripMapExtension("MAP.kfm"), "MAP")
	assert.Equal(t, StripMapExtension("MAP.KFM"), "MAP")
}

func TestGetMapsInDir(t *testing.T) {
	expected := []string{"MAP3", "MAP4"}
	dir, _ := os.Getwd()
	results := GetMapsInDir(path.Join(dir, TestData))
	assert.Equal(t, expected, results)
}

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig(PreEditIni)
	names := []string{"MAP1", "MAP2"}
	for _, name := range names {
		section := fmt.Sprintf("%s %s", name, MapSectionSuffix)
		resultsSection, _ := cfg.GetSection(section)
		if assert.NotNil(t, resultsSection) {
			resultsMapname, _ := cfg.GetValue(
				section, MapSectionMapOption)
			assert.Equal(t, name, resultsMapname)
			resultsScreenshot, _ := cfg.GetValue(
				section, MapSectionScreenshotOption)
			assert.Equal(t, MapSectionDefaultScreenshot, resultsScreenshot)
			resultsCycle, _ := cfg.GetValue(
				MapCycleSection, MapCycleOption)
			assert.True(t, strings.Contains(resultsCycle, name))
		}
	}
}

func TestSaveConfig(t *testing.T) {
	cfg, _ := goconfig.LoadFromData([]byte(""))
	cfg.SetValue("MAP", "MAP", "MAP")
	SaveConfig(cfg, EditedIni)
	file, err := ioutil.ReadFile(EditedIni)
	if err != nil {
		panic(err)
	}
	contents := string(file)
	assert.True(t, strings.Contains(contents, "MAP"))
}

func TestMain(t *testing.T) {
	CopyFile(PreEditIni, EditedIni)
	os.Args = []string{"cmd",
		fmt.Sprintf("-mapdir=%s", TestData),
		fmt.Sprintf("-config=%s", EditedIni)}
	main()
	cfg, _ := goconfig.LoadConfigFile(EditedIni)
	names := []string{"MAP1", "MAP2", "MAP3", "MAP4"}
	for _, name := range names {
		section := fmt.Sprintf("%s %s", name, MapSectionSuffix)
		resultsSection, _ := cfg.GetSection(section)
		if assert.NotNil(t, resultsSection) {
			resultsMapname, _ := cfg.GetValue(
				section, MapSectionMapOption)
			assert.Equal(t, name, resultsMapname)
			resultsScreenshot, _ := cfg.GetValue(
				section, MapSectionScreenshotOption)
			assert.Equal(t, MapSectionDefaultScreenshot, resultsScreenshot)
			resultsCycle, _ := cfg.GetValue(
				MapCycleSection, MapCycleOption)
			assert.True(t, strings.Contains(resultsCycle, name))
		}
	}
}
