from os import listdir
from os.path import isfile, isdir, join
import re
import sys
import textwrap

DEFAULT_SCREENSHOT='UI_MapPreview_TEX.UI_MapPreview_Placeholder'
FILE_PATTERN = re.compile('.*.kfm$', re.IGNORECASE)
NAME_PATTERN = re.compile('^([^.]*).kfm$', re.IGNORECASE)
USAGE = 'Usage:\nkf2mapman.py <mapdir>'

# [KF-Default KFMapSummary]
# MapName=KF-Default
# ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
def data_store(name, screenshot=DEFAULT_SCREENSHOT):
    '''Returns a map data store for the given map name. Optionally provide a screenshot.'''
    return textwrap.dedent(
    """
    [{0} KFMapSummary]
    MapName={0}
    ScreenshotPathName={1}
    """.format(name, screenshot))

def get_map_names(directory):
    '''Return a list of map names for a given directory'''
    try:
        return maplist
    except UnboundLocalError:
        maplist = [
            re.search(NAME_PATTERN, name).group(1)
            for name in listdir(directory)
            if isfile(join(directory, name))
            and re.match(FILE_PATTERN, name)
        ]
        return maplist

def create_map_list(names):
    '''Returns a list of data stores for the given map names'''
    return ''.join(
        [data_store(name) for name in names]
    )

def create_map_rotation(names):
    return ','.join([
        '"%s"' % name
        for name in names
    ])

if __name__ == '__main__':
    if len(sys.argv) == 2 and isdir(sys.argv[1]):
        directory = sys.argv[1]
        print '\nrotation string:'
        print create_map_rotation(get_map_names(directory))
        print '\nmap entries:'
        print create_map_list(get_map_names(directory))
    else:
        print USAGE
        sys.exit()
