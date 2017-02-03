from os import listdir
from os.path import isfile, isdir, join
import ConfigParser
import re
import sys
import textwrap

USAGE = 'Usage:\nkf2mapman.py <mapdir>'
DEFAULT_SCREENSHOT='UI_MapPreview_TEX.UI_MapPreview_Placeholder'
FILE_PATTERN = '.*.kfm$'
NAME_PATTERN = '^([^.]*).kfm$'

# [KF-Default KFMapSummary]
# MapName=KF-Default
# ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
def create_data_store(name):
    '''Returns a map data store entry for the given map name.'''
    return textwrap.dedent(
    """
    [{0} KFMapSummary]
    MapName={0}
    ScreenshotPathName={1}
    """.format(name, DEFAULT_SCREENSHOT))

def create_map_list(names):
    '''Returns a list of data stores for the given map names'''
    return ''.join(
        [create_data_store(name) for name in names]
    )

def create_map_rotation(names):
    '''Create the map rotation string for MapCycle'''
    return ','.join([
        '"%s"' % name
        for name in names
    ])


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

def get_current_maps(config):
    '''Returns a list of maps from the given ConfigParser object'''
    return [
        config.get(name, 'MapName')
        for name in config.sections()
        if 'KFMapSummary' in name
    ]


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
