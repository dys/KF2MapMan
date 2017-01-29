import os
import textwrap
import kf2mapman

def test_data_store():
    expected = textwrap.dedent("""
        [KF-Default KFMapSummary]
        MapName=KF-Default
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.data_store('KF-Default')
    print results
    assert results == expected

def test_create_map_list():
    expected = textwrap.dedent("""
        [KF-Default KFMapSummary]
        MapName=KF-Default
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder

        [KF-Default2 KFMapSummary]
        MapName=KF-Default2
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.create_map_list(['KF-Default', 'KF-Default2'])
    print results
    assert results == expected

def test_get_map_names():
    testmaps_dir = os.path.join(
        os.path.dirname(os.path.realpath(__file__)),
        'testmaps'
    )
    expected = ['KF-Default', 'KF-Default2']
    results = kf2mapman.get_map_names(testmaps_dir)
    print results
    assert expected == results

def test_create_map_rotation():
    names = ['KF-Default', 'KF-Default2']
    expected = '"KF-Default","KF-Default2"'
    results = kf2mapman.create_map_rotation(names)
    print results
    assert expected == results
