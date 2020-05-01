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
#include "random_cipher.hpp"

#define VERSION "1.0"

#ifdef _UNICODE
#define tstring std::wstring
#else
#define tstring std::string
#endif

//
// The plugin data that Notepad++ needs
//
FuncItem funcItem[] = {
	// name func cmdID init2Check shortKey
	{ TEXT("Hello Notepad++"),		hello,		0,	false,	NULL },
	{ TEXT("Hello (with dialog)"),	helloDlg,	0,	false,	NULL },
	{ TEXT("---"),					nullptr,	0,	false,	NULL },
	{ TEXT("ozg -> gfx"),			ozg2gfx,	0,	false,	NULL },
	{ TEXT("gfx -> ozg"),			gfx2ozg,	0,	false,	NULL },
	{ TEXT("---"),					nullptr,	0,	false,	NULL },
	{ TEXT("about"),				about,		0,	false,	NULL },
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

void ozg2gfx() {
	unsigned char* bufPlain = nullptr;
	unsigned char* bufCipher = nullptr;
	tstring text;
	[&] {
		__try {
			// validate ozg extension
			TCHAR fileName[MAX_PATH] = { 0 };
			::SendMessage(nppData._nppHandle, NPPM_GETFILENAME, MAX_PATH, (LPARAM)fileName);
			TCHAR* ext = wcsrchr(fileName, TEXT('.'));
			if (ext == nullptr || wcscmp(ext+1, TEXT("ozg")) != 0) {
				text += TEXT("[");
				text += fileName;
				text += TEXT("]");
				text += TEXT(" is not a [*.ozg]");
				::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
				return;
			}

			// Get the current scintilla
			int which = -1;
			::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
			if (which == -1)
				return;
			HWND curScintilla = (which == 0) ? nppData._scintillaMainHandle : nppData._scintillaSecondHandle;

			// get current text
			int lenCipher = ::SendMessage(curScintilla, SCI_GETLENGTH, 0, 0);
			bufCipher = new unsigned char[lenCipher + 1];
			::SendMessage(curScintilla, SCI_GETTEXT, lenCipher + 1, (LPARAM)bufCipher);

			// decrypt
			size_t lenPlain = RandomCipher::decrypt(nullptr, bufCipher, lenCipher);
			// validate cipher length
			if (lenPlain == 0) {
				text += TEXT("convert failed [");
				text += fileName;
				text += TEXT("] is not a valid [*.ozg]");
				::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
				return;
			}
			bufPlain = new unsigned char[lenPlain];
			RandomCipher::decrypt(bufPlain, bufCipher, lenCipher);

#if 0
			// marshal

			// Scintilla control has no Unicode mode, so we use (char *) here
			::SendMessage(curScintilla, SCI_SETTEXT, 0, (LPARAM)bufPlain);

			TCHAR path[MAX_PATH] = { 0 };
			::SendMessage(curScintilla, NPPM_GETCURRENTDIRECTORY, MAX_PATH, (LPARAM)path);
#endif
			::SendMessage(curScintilla, SCI_REPLACETARGET, lenPlain, (LPARAM)bufPlain);
		}
		__finally {
			if (bufPlain != nullptr) {
				delete[]bufPlain;
			}
			if (bufCipher != nullptr) {
				delete[]bufCipher;
			}
		}
	}();
}

void gfx2ozg() {
	unsigned char* bufPlain = nullptr;
	unsigned char* bufCipher = nullptr;
	tstring text;
	[&] {
		__try {
			// get file name
			TCHAR fileName[MAX_PATH] = { 0 };
			::SendMessage(nppData._nppHandle, NPPM_GETFILENAME, MAX_PATH, (LPARAM)fileName);

			// Get the current scintilla
			int which = -1;
			::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
			if (which == -1)
				return;
			HWND curScintilla = (which == 0) ? nppData._scintillaMainHandle : nppData._scintillaSecondHandle;

			// get current text
			int lenPlain = ::SendMessage(curScintilla, SCI_GETLENGTH, 0, 0);
			if (lenPlain == 0) {
				text += TEXT("convert failed [");
				text += fileName;
				text += TEXT("] is empty");
				::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
				return;
			}
			bufPlain = new unsigned char[lenPlain + 1];
			::SendMessage(curScintilla, SCI_GETTEXT, lenPlain + 1, (LPARAM)bufPlain);

			// encrypt
			size_t lenCipher = RandomCipher::encrypt(nullptr, bufPlain, lenPlain);
			bufCipher = new unsigned char[lenCipher];
			RandomCipher::encrypt(bufCipher, bufPlain, lenPlain);

#if 0
			// marshal

			// Scintilla control has no Unicode mode, so we use (char *) here
			::SendMessage(curScintilla, SCI_SETTEXT, 0, (LPARAM)bufPlain);

			TCHAR path[MAX_PATH] = { 0 };
			::SendMessage(curScintilla, NPPM_GETCURRENTDIRECTORY, MAX_PATH, (LPARAM)path);
#endif
			::SendMessage(curScintilla, SCI_REPLACETARGET, lenCipher, (LPARAM)bufCipher);
		}
		__finally {
			if (bufPlain != nullptr) {
				delete[]bufPlain;
			}
			if (bufCipher != nullptr) {
				delete[]bufCipher;
			}
		}
	}();
}

void about() {
	tstring text = TEXT("version: ");
	text += TEXT(VERSION);
	::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
}
