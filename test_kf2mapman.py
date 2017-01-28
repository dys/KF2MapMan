import sys
import textwrap
import kf2mapman

def test_data_store():
    expected = textwrap.dedent("""
        [KF-Default KFMapSummary]
        MapName=KF-Default
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.data_store('KF-Default')
    assert results == expected

def test_map_list():
    expected = textwrap.dedent("""
        [KF-Default KFMapSummary]
        MapName=KF-Default
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder

        [KF-Default2 KFMapSummary]
        MapName=KF-Default2
        ScreenshotPathName=UI_MapPreview_TEX.UI_MapPreview_Placeholder
    """)
    results = kf2mapman.map_list(['KF-Default', 'KF-Default2'])
    assert results == expected
