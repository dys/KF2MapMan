package kf2mapman

import (
	"flag"
	"fmt"
	"io"
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

// LoadConfig returns a goconfig.ConfigFile from
// an io.Reader compatible source
func LoadConfig(in io.Reader) (*goconfig.ConfigFile, error) {
	cfg, err := goconfig.LoadFromReader(in)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// SaveConfig writes the given cfg to file
func SaveConfig(cfg *goconfig.ConfigFile, file string) {
	goconfig.SaveConfigFile(cfg, file)
}

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
	cycle := GetMapCycle(cfg)
	for _, m := range cycle {
		if m == name {
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
	if strings.HasSuffix(name, strings.ToUpper(MapExtension)) ||
		strings.HasSuffix(name, strings.ToLower(MapExtension)) {
		return true
	}
	return false
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
	maps := []string{}
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
		log.Fatal("-mapdir missing")
	}
	if *configFile == "" {
		log.Fatal("-config missing")
	}

	// Load the configuration file
	cfgdata, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := LoadConfig(cfgdata)
	if err != nil {
		log.Fatal(err)
	}

	AddMapsToConfig(GetMapsInDir(*mapDir), cfg)
	SaveConfig(cfg, *configFile)
}
