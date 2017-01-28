import sys
import textwrap

DEFAULT_SCREENSHOT='UI_MapPreview_TEX.UI_MapPreview_Placeholder'

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

def map_list(names):
    '''Returns a list of data stores for the given map names'''
    return ''.join(
        [data_store(name) for name in names]
    )

if __name__ == __name__:
    with open(sys.argv[1], 'r') as f:
        maplist = [m for m in f.read().split()]
    print map_list(maplist)
