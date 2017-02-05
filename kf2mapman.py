from os import listdir
from os.path import isfile, isdir, join
import ConfigParser
import re
import sys
import textwrap

USAGE = 'Usage:\nkf2mapman.py <mapdir>'
FILE_PATTERN = '.*.kfm$'
NAME_PATTERN = '^([^.]*).kfm$'
MAP_SECTION_SUFFIX = ' KFMapSummary'
MAP_SECTION_MAPNAME = 'MapName'
MAP_SECTION_SCREENSHOT = 'UI_MapPreview_TEX.UI_MapPreview_Placeholder'
MAP_ROTATION_SECTION = 'KFGame.KFGameInfo'
MAP_ROTATION_OPTION = 'GameMapCycles'
MAP_ROTATION_PREFIX = '(Maps=('
MAP_ROTATION_SUFFIX = '))'

# [KF-Default KFMapSummary]
# MapName=KF-Default
# ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder

def create_section(name):
    '''Returns a given map name as a section header string'''
    return name + MAP_SECTION_SUFFIX

def create_data_store(name):
    '''Returns a map data store entry for the given map name.'''
    return {
        'MapName': name,
        'ScreenshotPathName': MAP_SECTION_SCREENSHOT
    }

def create_map_rotation(names):
    '''Create the map rotation string for MapCycle'''
    return '{0}{1}{2}'.format(
        MAP_ROTATION_PREFIX,
        ','.join([
            '"%s"' % name
            for name in names
        ]),
        MAP_ROTATION_SUFFIX
        )

def get_maps_in_dir(directory):
    '''Return a list of map names (minus extension) for a given directory.'''
    return [
        re.search(re.compile(
            NAME_PATTERN, re.IGNORECASE), name) .group(1)
        for name in listdir(directory)
        if isfile(join(directory, name))
        and re.match(re.compile(
            FILE_PATTERN, re.IGNORECASE), name)
    ]

def read_ini(ini):
    '''Returns the given ini file as a ConfigParser object'''
    with open(ini, 'r') as configfile:
        config = ConfigParser.ConfigParser()
        config.read(configfile)
        return config

def write_config(config, ini):
    '''Writes the ConfigParser object to the ini path'''
    with open(ini, 'w') as configfile:
        config.write(configfile)

def get_maps_from_config(config):
    '''Returns a list of maps from the given ConfigParser object'''
    return [
        config.get(name, 'MapName')
        for name in config.sections()
        if MAP_SECTION_SUFFIX in name
    ]

def add_map_sections_to_config(names, config):
    '''Add the maps to a ConfigParser object'''
    for name in names:
        try:
            config.add_section(create_section(name))
        except ConfigParser.DuplicateSectionError:
            pass
        for k,v in create_data_store(name).iteritems():
            config.set(create_section(name), k, v)
    return config

def add_rotation_to_config(names, config):
    '''Adds the given map names to the config's rotation'''
    print config.sections()
    print [
            o
            for s in config.sections()
            for o in config.options(s)
            ]
    try:
        config.add_section(MAP_ROTATION_SECTION)
    except ConfigParser.DuplicateSectionError:
        pass
    config.set(
        MAP_ROTATION_SECTION,
        MAP_ROTATION_OPTION,
        create_map_rotation([
            name for name in names
            if name not in get_rotation_from_config(config)
        ])
    )
    return config

def get_rotation_from_config(config):
    '''Returns a list of map names in the config's rotation'''
    return config.get(MAP_ROTATION_SECTION, MAP_ROTATION_OPTION) \
        .replace(MAP_ROTATION_PREFIX, '') \
        .replace(MAP_ROTATION_SUFFIX, '') \
        .replace('"', '') \
        .split(',')

if __name__ == '__main__':
    if len(sys.argv) == 2 and isdir(sys.argv[1]):
        directory = sys.argv[1]
        print '\nrotation string:'
        print create_map_rotation(get_maps_in_dir(directory))
        print '\nmap entries:'
        print create_map_list(get_maps_in_dir(directory))
    else:
        print USAGE
        sys.exit()
