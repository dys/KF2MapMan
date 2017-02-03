import kf2mapman

import ConfigParser
import os
import shutil
import textwrap

TEST_DIR = os.path.join(
    os.path.dirname(os.path.realpath(__file__)),
    'testmaps')
PRE_EDIT_INI = os.path.join(TEST_DIR, 'pre-edit.ini')
EDITED_INI = os.path.join(TEST_DIR, 'to-edit.ini')
POST_EDIT_INI = os.path.join(TEST_DIR, 'post-edit.ini')
MAP_NAMES = ['KF-Map1', 'KF-Map2']
NEW_MAP_NAMES = ['KF-Newmap1', 'KF-Newmap2']

def test_create_data_store():
    expected = textwrap.dedent("""
        [KF-Newmap1 KFMapSummary]
        MapName=KF-Newmap1
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.create_data_store('KF-Newmap1')
    assert results == expected

def test_create_map_list():
    expected = textwrap.dedent("""
        [KF-Newmap1 KFMapSummary]
        MapName=KF-Newmap1
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder

        [KF-Newmap2 KFMapSummary]
        MapName=KF-Newmap2
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.create_map_list(NEW_MAP_NAMES)
    assert results == expected

def test_create_map_rotation():
    expected = '"KF-Newmap1","KF-Newmap2"'
    results = kf2mapman.create_map_rotation(NEW_MAP_NAMES)
    assert results == expected


def test_get_maps_in_dir():
    results = kf2mapman.get_maps_in_dir(TEST_DIR)
    assert results == NEW_MAP_NAMES

def test_read_ini():
    #config = ConfigParser.ConfigParser()
    results = kf2mapman.read_ini(EDITED_INI)
    assert isinstance(results, ConfigParser.ConfigParser)

def test_get_maps_from_config():
    config = ConfigParser.ConfigParser()
    config.read(PRE_EDIT_INI)
    results = kf2mapman.get_current_maps(config)
    assert results == MAP_NAMES


#def test_edit_ini_file():
#    shutil.copy(PRE_EDIT_INI, EDITED_INI)
#    kf2mapman.add_maps_to_ini(NEW_MAP_NAMES, EDITED_INI)
#    expected = open(POST_EDIT_INI, 'r').read()
#    results = open(EDITED_INI, 'r').read()
#    assert results == expected
#    os.remove(EDITED_INI)
