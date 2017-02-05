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
ROTATION_STRING = '(Maps=("KF-Map1","KF-Map2"))'

def test_create_section():
    name = MAP_NAMES[0]
    expected = name + kf2mapman.MAP_SECTION_SUFFIX
    results = kf2mapman.create_section(MAP_NAMES[0])
    assert expected == results

def test_create_data_store():
    name = MAP_NAMES[0]
    results = kf2mapman.create_data_store(name)
    assert results[kf2mapman.MAP_SECTION_MAPNAME] == name
    assert results['ScreenshotPathName'] == kf2mapman.MAP_SECTION_SCREENSHOT

def test_create_map_rotation():
    results = kf2mapman.create_map_rotation(MAP_NAMES)
    assert results == ROTATION_STRING

def test_get_maps_in_dir():
    results = kf2mapman.get_maps_in_dir(TEST_DIR)
    assert results == NEW_MAP_NAMES
    assert 'NotAMap' not in results

def test_read_ini():
    results = kf2mapman.read_ini(PRE_EDIT_INI)
    assert isinstance(results, ConfigParser.ConfigParser)

def test_write_config():
    c = ConfigParser.ConfigParser()
    c.read(PRE_EDIT_INI)
    expected = [
        c.options(section)
        for section in c.sections()
    ]
    kf2mapman.write_config(c, EDITED_INI)
    c = ConfigParser.ConfigParser()
    c.read(EDITED_INI)
    os.remove(EDITED_INI)
    results = [
        c.options(section)
        for section in c.sections()
    ]
    assert expected == results

def test_get_maps_from_config():
    c = ConfigParser.ConfigParser()
    c.read(PRE_EDIT_INI)
    results = kf2mapman.get_maps_from_config(c)
    assert results == MAP_NAMES

def test_add_map_sections_to_config():
    c = ConfigParser.ConfigParser()
    results = kf2mapman.add_map_sections_to_config(MAP_NAMES, c)
    assert all([
        results.get(
            name + kf2mapman.MAP_SECTION_SUFFIX,
            kf2mapman.MAP_SECTION_MAPNAME
            ) == name
        for name in MAP_NAMES
    ])

def test_add_rotation_to_config():
    c = ConfigParser.ConfigParser()
    results = kf2mapman.add_rotation_to_config(MAP_NAMES, c)
    assert results.get(
        kf2mapman.MAP_ROTATION_SECTION,
        kf2mapman.MAP_ROTATION_OPTION
        ) == ROTATION_STRING

def test_get_rotation_from_config():
    c = ConfigParser.ConfigParser()
    c.read(PRE_EDIT_INI)
    results = kf2mapman.get_rotation_from_config(c)
    assert results == MAP_NAMES

#def test_edit_ini_file():
#    shutil.copy(PRE_EDIT_INI, EDITED_INI)
#    kf2mapman.add_maps_to_ini(NEW_MAP_NAMES, EDITED_INI)
#    expected = open(POST_EDIT_INI, 'r').read()
#    results = open(EDITED_INI, 'r').read()
#    assert results == expected
#    os.remove(EDITED_INI)
