//---------------------------------------------------------------------------
// # Developed by:	S@nek[BoR];
// # Created:		10.10.2010 18-00;
// # Last Update:	10.10.2010 19-05;
// # Status:		Completed;
// # Description:	Color class for C++ program;
//---------------------------------------------------------------------------

//--------------------------------------------------------------------------------------------------------------------------------
#ifndef TCOLOR_H
#define TCOLOR_H
//--------------------------------------------------------------------------------------------------------------------------------

#include <windows.h>
//--------------------------------------------------------------------------------------------------------------------------------

namespace TColor
{
enum T
{
    AliceBlue				= RGB(0xF0, 0xF8, 0xFF),
    AntiqueWhite			= RGB(0xFA, 0xEB, 0xD7),
    Aqua					= RGB(0x00, 0xFF, 0xFF),
    Aquamarine				= RGB(0x7F, 0xFF, 0xD4),
    Azure					= RGB(0xF0, 0xFF, 0xFF),
    Beige					= RGB(0xF5, 0xF5, 0xDC),
    Bisque					= RGB(0xFF, 0xE4, 0xC4),
    Black					= RGB(0x00, 0x00, 0x00),
    BlanchedAlmond			= RGB(0xFF, 0xEB, 0xCD),
    Blue					= RGB(0x00, 0x00, 0xFF),
    BlueViolet				= RGB(0x8A, 0x2B, 0xE2),
    Brown					= RGB(0xA5, 0x2A, 0x2A),
    BurlyWood				= RGB(0xDE, 0xB8, 0x87),
    CadetBlue				= RGB(0x5F, 0x9E, 0xA0),
    Chartreuse				= RGB(0x7F, 0xFF, 0x00),
    Chocolate				= RGB(0xD2, 0x69, 0x1E),
    Coral					= RGB(0xFF, 0x7F, 0x50),
    CornflowerBlue			= RGB(0x64, 0x95, 0xED),
    Cornsilk				= RGB(0xFF, 0xF8, 0xDC),
    Crimson					= RGB(0xDC, 0x14, 0x3C),
    Cyan					= RGB(0x00, 0xFF, 0xFF),
    DarkBlue				= RGB(0x00, 0x00, 0x8B),
    DarkCyan				= RGB(0x00, 0x8B, 0x8B),
    DarkGoldenrod			= RGB(0xB8, 0x86, 0x0B),
    DarkGray				= RGB(0xA9, 0xA9, 0xA9),
    DarkGreen				= RGB(0x00, 0x64, 0x00),
    DarkKhaki				= RGB(0xBD, 0xB7, 0x6B),
    DarkMagenta				= RGB(0x8B, 0x00, 0x8B),
    DarkOliveGreen			= RGB(0x55, 0x6B, 0x2F),
    DarkOrange				= RGB(0xFF, 0x8C, 0x00),
    DarkOrchid				= RGB(0x99, 0x32, 0xCC),
    DarkRed					= RGB(0x8B, 0x00, 0x00),
    DarkSalmon				= RGB(0xE9, 0x96, 0x7A),
    DarkSeaGreen			= RGB(0x8F, 0xBC, 0x8B),
    DarkSlateBlue			= RGB(0x48, 0x3D, 0x8B),
    DarkSlateGray			= RGB(0x2F, 0x4F, 0x4F),
    DarkTurquoise			= RGB(0x00, 0xCE, 0xD1),
    DarkViolet				= RGB(0x94, 0x00, 0xD3),
    DeepPink				= RGB(0xFF, 0x14, 0x93),
    DeepSkyBlue				= RGB(0x00, 0xBF, 0xFF),
    DimGray					= RGB(0x69, 0x69, 0x69),
    DodgerBlue				= RGB(0x1E, 0x90, 0xFF),
    Firebrick				= RGB(0xB2, 0x22, 0x22),
    FloralWhite				= RGB(0xFF, 0xFA, 0xF0),
    ForestGreen				= RGB(0x22, 0x8B, 0x22),
    Fuchsia					= RGB(0xFF, 0x00, 0xFF),
    Gainsboro				= RGB(0xDC, 0xDC, 0xDC),
    GhostWhite				= RGB(0xF8, 0xF8, 0xFF),
    Gold					= RGB(0xFF, 0xD7, 0x00),
    Goldenrod				= RGB(0xDA, 0xA5, 0x20),
    Gray					= RGB(0x80, 0x80, 0x80),
    Green					= RGB(0x00, 0x80, 0x00),
    GreenYellow				= RGB(0xAD, 0xFF, 0x2F),
    Honeydew				= RGB(0xF0, 0xFF, 0xF0),
    HotPink					= RGB(0xFF, 0x69, 0xB4),
    IndianRed				= RGB(0xCD, 0x5C, 0x5C),
    Indigo					= RGB(0x4B, 0x00, 0x82),
    Ivory					= RGB(0xFF, 0xFF, 0xF0),
    Khaki					= RGB(0xF0, 0xE6, 0x8C),
    Lavender				= RGB(0xE6, 0xE6, 0xFA),
    LavenderBlush			= RGB(0xFF, 0xF0, 0xF5),
    LawnGreen				= RGB(0x7C, 0xFC, 0x00),
    LemonChiffon			= RGB(0xFF, 0xFA, 0xCD),
    LightBlue				= RGB(0xAD, 0xD8, 0xE6),
    LightCoral				= RGB(0xF0, 0x80, 0x80),
    LightCyan				= RGB(0xE0, 0xFF, 0xFF),
    LightGoldenrodYellow	= RGB(0xFA, 0xFA, 0xD2),
    LightGray				= RGB(0xD3, 0xD3, 0xD3),
    LightGreen				= RGB(0x90, 0xEE, 0x90),
    LightPink				= RGB(0xFF, 0xB6, 0xC1),
    LightSalmon				= RGB(0xFF, 0xA0, 0x7A),
    LightSeaGreen			= RGB(0x20, 0xB2, 0xAA),
    LightSkyBlue			= RGB(0x87, 0xCE, 0xFA),
    LightSlateGray			= RGB(0x77, 0x88, 0x99),
    LightSteelBlue			= RGB(0xB0, 0xC4, 0xDE),
    LightYellow				= RGB(0xFF, 0xFF, 0xE0),
    Lime					= RGB(0x00, 0xFF, 0x00),
    LimeGreen				= RGB(0x32, 0xCD, 0x32),
    Linen					= RGB(0xFA, 0xF0, 0xE6),
    Magenta					= RGB(0xFF, 0x00, 0xFF),
    Maroon					= RGB(0x80, 0x00, 0x00),
    MediumAquamarine		= RGB(0x66, 0xCD, 0xAA),
    MediumBlue				= RGB(0x00, 0x00, 0xCD),
    MediumOrchid			= RGB(0xBA, 0x55, 0xD3),
    MediumPurple			= RGB(0x93, 0x70, 0xDB),
    MediumSeaGreen			= RGB(0x3C, 0xB3, 0x71),
    MediumSlateBlue			= RGB(0x7B, 0x68, 0xEE),
    MediumSpringGreen		= RGB(0x00, 0xFA, 0x9A),
    MediumTurquoise			= RGB(0x48, 0xD1, 0xCC),
    MediumVioletRed			= RGB(0xC7, 0x15, 0x85),
    MidnightBlue			= RGB(0x19, 0x19, 0x70),
    MintCream				= RGB(0xF5, 0xFF, 0xFA),
    MistyRose				= RGB(0xFF, 0xE4, 0xE1),
    Moccasin				= RGB(0xFF, 0xE4, 0xB5),
    NavajoWhite				= RGB(0xFF, 0xDE, 0xAD),
    Navy					= RGB(0x00, 0x00, 0x80),
    OldLace					= RGB(0xFD, 0xF5, 0xE6),
    Olive					= RGB(0x80, 0x80, 0x00),
    OliveDrab				= RGB(0x6B, 0x8E, 0x23),
    Orange					= RGB(0xFF, 0xA5, 0x00),
    OrangeRed				= RGB(0xFF, 0x45, 0x00),
    Orchid					= RGB(0xDA, 0x70, 0xD6),
    PaleGoldenrod			= RGB(0xEE, 0xE8, 0xAA),
    PaleGreen				= RGB(0x98, 0xFB, 0x98),
    PaleTurquoise			= RGB(0xAF, 0xEE, 0xEE),
    PaleVioletRed			= RGB(0xDB, 0x70, 0x93),
    PapayaWhip				= RGB(0xFF, 0xEF, 0xD5),
    PeachPuff				= RGB(0xFF, 0xDA, 0xB9),
    Peru					= RGB(0xCD, 0x85, 0x3F),
    Pink					= RGB(0xFF, 0xC0, 0xCB),
    Plum					= RGB(0xDD, 0xA0, 0xDD),
    PowderBlue				= RGB(0xB0, 0xE0, 0xE6),
    Purple					= RGB(0x80, 0x00, 0x80),
    Red						= RGB(0xFF, 0x00, 0x00),
    RosyBrown				= RGB(0xBC, 0x8F, 0x8F),
    RoyalBlue				= RGB(0x41, 0x69, 0xE1),
    SaddleBrown				= RGB(0x8B, 0x45, 0x13),
    Salmon					= RGB(0xFA, 0x80, 0x72),
    SandyBrown				= RGB(0xF4, 0xA4, 0x60),
    SeaGreen				= RGB(0x2E, 0x8B, 0x57),
    SeaShell				= RGB(0xFF, 0xF5, 0xEE),
    Sienna					= RGB(0xA0, 0x52, 0x2D),
    Silver					= RGB(0xC0, 0xC0, 0xC0),
    SkyBlue					= RGB(0x87, 0xCE, 0xEB),
    SlateBlue				= RGB(0x6A, 0x5A, 0xCD),
    SlateGray				= RGB(0x70, 0x80, 0x90),
    Snow					= RGB(0xFF, 0xFA, 0xFA),
    SpringGreen				= RGB(0x00, 0xFF, 0x7F),
    SteelBlue				= RGB(0x46, 0x82, 0xB4),
    Tan						= RGB(0xD2, 0xB4, 0x8C),
    Teal					= RGB(0x00, 0x80, 0x80),
    Thistle					= RGB(0xD8, 0xBF, 0xD8),
    Tomato					= RGB(0xFF, 0x63, 0x47),
    Transparent				= RGB(0xFF, 0xFF, 0xFF),
    Turquoise				= RGB(0x40, 0xE0, 0xD0),
    Violet					= RGB(0xEE, 0x82, 0xEE),
    Wheat					= RGB(0xF5, 0xDE, 0xB3),
    White					= RGB(0xFF, 0xFF, 0xFF),
    WhiteSmoke				= RGB(0xF5, 0xF5, 0xF5),
    Yellow					= RGB(0xFF, 0xFF, 0x00),
};
};
//--------------------------------------------------------------------------------------------------------------------------------



//--------------------------------------------------------------------------------------------------------------------------------
#endif
//--------------------------------------------------------------------------------------------------------------------------------


