package kf2mapman

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/Unknwon/goconfig"
	"github.com/stretchr/testify/assert"
)

const TestData = "testdata/"
const PreEditIni = TestData + "pre-edit.ini"
const EditedIni = TestData + "to-edit.ini"
const IniFile = `
[MAP1 KFMapSummary]
MapName=MAP1
ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder

[KFGame.KFGameInfo]
GameMapCycles=(Maps=("MAP1","MAP2"))

[MAP2 KFMapSummary]
MapName=MAP2
ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
`

func Reader() io.Reader {
	return strings.NewReader(IniFile)
}

func Config() *goconfig.ConfigFile {
	cfg, err := LoadConfig(Reader())
	if err != nil {
		panic(err)
	}
	return cfg
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig(Reader())
	assert.Nil(t, err)
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

func TestCreateSectionHeader(t *testing.T) {
	expected := fmt.Sprintf("MAP1 %s", MapSectionSuffix)
	results := CreateSectionHeader("MAP1")
	assert.Equal(t, expected, results)
}

func TestGetMapSections(t *testing.T) {
	expected := []string{"MAP1", "MAP2"}
	results := GetMapSections(Config())
	assert.Equal(t, expected, results)
}

func TestGetMapCycle(t *testing.T) {
	expected := []string{"MAP1", "MAP2"}
	results := GetMapCycle(Config())
	assert.Equal(t, expected, results)
}

func TestCreateMapCycle(t *testing.T) {
	expected := "(Maps=(\"MAP1\",\"MAP2\"))"
	results := CreateMapCycle([]string{"MAP1", "MAP2"})
	assert.Equal(t, expected, results)
}

func TestMapInCycle(t *testing.T) {
	assert.True(t, MapInCycle("MAP1", Config()))
	assert.False(t, MapInCycle("MAP", Config()))
}

func TestAddMapToCycle(t *testing.T) {
	cfg := Config()
	expected := "(Maps=(\"MAP1\",\"MAP2\",\"MAP3\",\"MAP4\"))"
	AddMapToCycle("MAP3", cfg)
	AddMapToCycle("MAP4", cfg)
	results, _ := cfg.GetValue(
		MapCycleSection, MapCycleOption)
	assert.Equal(t, expected, results)
}

func TestDontAddDuplicateMaps(t *testing.T) {
	cfg := Config()
	expected := "(Maps=(\"MAP1\",\"MAP2\",\"MAP3\"))"
	AddMapToCycle("MAP3", cfg)
	AddMapToCycle("MAP3", cfg)
	results, _ := cfg.GetValue(
		MapCycleSection, MapCycleOption)
	assert.Equal(t, expected, results)
}

func TestAddMapSection(t *testing.T) {
	cfg := Config()
	name := "MAP3"
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
	cfg := Config()
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
