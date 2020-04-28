//this file is part of notepad++
//Copyright (C)2003 Don HO <donho@altern.org>
//
//This program is free software; you can redistribute it and/or
//modify it under the terms of the GNU General Public License
//as published by the Free Software Foundation; either
//version 2 of the License, or (at your option) any later version.
//
//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with this program; if not, write to the Free Software
//Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.

#include "PluginDefinition.h"
#include "menuCmdID.h"
#include "bmdCipher.hpp"

//
// The plugin data that Notepad++ needs
//
FuncItem funcItem[] = {
	// name func cmdID init2Check shortKey
	{ TEXT("Hello Notepad++"),		hello,		0,	false,	NULL },
	{ TEXT("Hello (with dialog)"),	helloDlg,	0,	false,	NULL },
	{ TEXT("bmd/ozd/ozg decode"),	handleDec,	0,	false,	NULL },
};

int funcItemLen() {
	return sizeof(funcItem) / sizeof(funcItem[0]);
}

//
// The data of Notepad++ that you can use in your plugin commands
//
NppData nppData;

//
// Initialize your plugin data here
// It will be called while plugin loading   
void pluginInit(HANDLE /*hModule*/){}

//
// Here you can do the clean up, save the parameters (if any) for the next session
//
void pluginCleanUp(){}

//
// Initialization of your plugin commands
// You should fill your plugins commands here
void commandMenuInit(){}

//
// Here you can do the clean up (especially for the shortcut)
//
void commandMenuCleanUp(){}

//----------------------------------------------//
//-- STEP 4. DEFINE YOUR ASSOCIATED FUNCTIONS --//
//----------------------------------------------//
void hello()
{
    // Open a new document
    ::SendMessageW(nppData._nppHandle, NPPM_MENUCOMMAND, 0, IDM_FILE_NEW);

    // Get the current scintilla
    int which = -1;
    ::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
    if (which == -1)
        return;
    HWND curScintilla = (which == 0)?nppData._scintillaMainHandle:nppData._scintillaSecondHandle;

    // Say hello now :
    // Scintilla control has no Unicode mode, so we use (char *) here
    ::SendMessage(curScintilla, SCI_SETTEXT, 0, (LPARAM)"Hello, Notepad++!");
}

void helloDlg()
{
    ::MessageBox(NULL, TEXT("Hello, Notepad++!"), TEXT("Notepad++ Plugin Template"), MB_OK);
}

void handleDec() {
	// Use current scintilla
	bmdenc(nullptr, nullptr, 0);
}
