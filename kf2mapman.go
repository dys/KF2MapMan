package kf2mapman

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Unknwon/goconfig"
)

// MapExtension is the pattern used used to match kf2 maps
const MapExtension = ".kfm"

// MapCycleOption is the INI option containing the map rotation
const MapCycleOption = "GameMapCycles"

// MapCycleSection is the INI section where MapCycleOption is found
const MapCycleSection = "KFGame.KFGameInfo"

// MapCyclePrefix is a string to be removed from
// the start of MapCycleOption prior to splitting
const MapCyclePrefix = "(Maps=("

// MapCycleSuffix is a string to be removed from
// the end of MapCycleOption prior to splitting
const MapCycleSuffix = "))"

// MapSectionSuffix is joined to the map name
// with a space to make the map section header
const MapSectionSuffix = "KFMapSummary"

// MapSectionMapOption is the key name in the INI file
// which contains the map name
const MapSectionMapOption = "MapName"

// MapSectionScreenshotOption is the key name in the INI file
// which contains the screenshot path
const MapSectionScreenshotOption = "ScreenshotPathName"

// MapSectionDefaultScreenshot is the default path
// used for maps without an existing screenshot
const MapSectionDefaultScreenshot = "UI_MapPreview_TEX.UI_MapPreview_Placeholder"

// CreateSectionHeader returns a properly formatted
// map section header for the given name
func CreateSectionHeader(name string) string {
	return fmt.Sprintf("%s %s", name, MapSectionSuffix)
}

// GetMapSections returns a list of map names
// from the map sections of the given config
func GetMapSections(cfg *goconfig.ConfigFile) []string {
	sections := cfg.GetSectionList()
	names := []string{}
	for _, name := range sections {
		if strings.HasSuffix(name, MapSectionSuffix) {
			names = append(names,
				strings.TrimSuffix(name,
					fmt.Sprintf(" %s", MapSectionSuffix)))
		}
	}
	return names
}

// GetMapCycle returns a list of map names
// from the given goconfig.ConfigFile
func GetMapCycle(cfg *goconfig.ConfigFile) []string {
	mapcycle, _ := cfg.GetValue(MapCycleSection, MapCycleOption)
	mapcycle = strings.TrimPrefix(mapcycle, MapCyclePrefix)
	mapcycle = strings.TrimSuffix(mapcycle, MapCycleSuffix)
	mapcycle = strings.Replace(mapcycle, "\"", "", -1)
	return strings.Split(mapcycle, ",")
}

// CreateMapCycle returns a GameMapCycles string for the given names
func CreateMapCycle(names []string) string {
	return fmt.Sprintf("%s\"%s\"%s",
		MapCyclePrefix,
		strings.Join(names, "\",\""),
		MapCycleSuffix)
}

// MapInCycle returns true if a map is already in the map cycle for config
func MapInCycle(name string, cfg *goconfig.ConfigFile) bool {
	for _, m := range GetMapCycle(cfg) {
		if strings.ToLower(m) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

// AddMapToCycle appends a name to the config's rotation
func AddMapToCycle(name string, cfg *goconfig.ConfigFile) {
	if !MapInCycle(name, cfg) {
		cfg.SetValue(
			MapCycleSection,
			MapCycleOption,
			CreateMapCycle(append(GetMapCycle(cfg), name)))
	}
}

// AddMapSection adds a map section for name to the given config
func AddMapSection(name string, cfg *goconfig.ConfigFile) {
	cfg.SetValue(
		CreateSectionHeader(name),
		MapSectionMapOption,
		name)
	cfg.SetValue(
		CreateSectionHeader(name),
		MapSectionScreenshotOption,
		MapSectionDefaultScreenshot)
}

// AddMapsToConfig adds each map name to the given config
// as a map section, and updates the rotation
func AddMapsToConfig(names []string, cfg *goconfig.ConfigFile) {
	for _, name := range names {
		AddMapToCycle(name, cfg)
		AddMapSection(name, cfg)
	}
}

// FileIsMap returns true if the file matches MapExtension
func FileIsMap(name string) bool {
	return strings.HasSuffix(strings.ToLower(name), strings.ToLower(MapExtension))
}

// StripMapExtension returns name with the map extension removed
func StripMapExtension(name string) string {
	name = strings.TrimSuffix(name, strings.ToLower(MapExtension))
	name = strings.TrimSuffix(name, strings.ToUpper(MapExtension))
	return name
}

// GetMapsInDir returns a list of KF2 maps in dir
func GetMapsInDir(dir string) []string {
	files, _ := ioutil.ReadDir(dir)
	var maps []string
	for _, file := range files {
		if FileIsMap(file.Name()) {
			maps = append(maps,
				strings.TrimSuffix(StripMapExtension(file.Name()),
					MapExtension))
		}
	}
	return maps
}

func main() {
	// Get the user options
	mapDir := flag.String("mapdir", "",
		"The directory containing custom maps to add")
	configFile := flag.String("config", "",
		"The path to the PCServer-KFGame.ini file")
	flag.Parse()

	// Check the user options
	if *mapDir == "" {
		log.Fatal("-mapdir missing: ")
	}
	if *configFile == "" {
		log.Fatal("-config missing: ")
	}

	// Open the INI file
	ini, err := os.OpenFile(*configFile, os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("Failed to open INI file: ", err)
	}

	// Load the configuration
	cfg, err := goconfig.LoadFromReader(ini)
	if err != nil {
		log.Fatal("Failed to load INI file: ", err)
	}
	ini.Truncate(0)

	// Add maps
	AddMapsToConfig(GetMapsInDir(*mapDir), cfg)

	// Save the config
	err = goconfig.SaveConfigData(cfg, ini)
	if err != nil {
		log.Fatal("Failed to save Config: ", err)
	}
	ini.Close()
}
