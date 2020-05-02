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
	// validate file extension
	TCHAR fileName[MAX_PATH] = { 0 };
	::SendMessage(nppData._nppHandle, NPPM_GETFILENAME, MAX_PATH, (LPARAM)fileName);
	TCHAR* ext = wcsrchr(fileName, TEXT('.'));
	if (ext == nullptr || wcscmp(ext, TEXT(".ozg")) != 0) {
		tstring text = TEXT("[");
		text += fileName;
		text += TEXT("]");
		text += TEXT(" is not a [*.ozg]");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK|MB_ICONWARNING);
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
	auto bufCipher = std::make_unique<unsigned char[]>(lenCipher + 1);
	::SendMessage(curScintilla, SCI_GETTEXT, lenCipher + 1, (LPARAM)bufCipher.get());

	// decrypt
	size_t lenPlain = RandomCipher::decrypt(nullptr, bufCipher.get(), lenCipher);
	// validate cipher length
	if (lenPlain == 0) {
		tstring text = TEXT("convert failed\r\n[");
		text += fileName;
		text += TEXT("] is not a valid [*.ozg]");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK|MB_ICONWARNING);
		return;
	}
	auto bufPlain = std::make_unique<unsigned char[]>(lenPlain);
	RandomCipher::decrypt(bufPlain.get(), bufCipher.get(), lenCipher);

	// marshal(not yet support current)

	// write
	::SendMessageW(nppData._nppHandle, NPPM_MENUCOMMAND, 0, IDM_FILE_NEW);
	which = -1;
	::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
	if (which == -1)
		return;
	curScintilla = (which == 0) ? nppData._scintillaMainHandle : nppData._scintillaSecondHandle;
	::SendMessage(curScintilla, SCI_ADDTEXT, lenPlain, (LPARAM)bufPlain.get());
#if 0
	// write bytes to a file
	TCHAR filePath[MAX_PATH] = { 0 };
	::SendMessage(nppData._nppHandle, NPPM_GETFULLCURRENTPATH, MAX_PATH, (LPARAM)filePath);
	*wcsrchr(filePath, TEXT('.')) = 0;
	wcscat(filePath, TEXT(".gfx"));
	std::unique_ptr<FILE, void(*)(FILE*)> f(_wfopen(filePath, TEXT("wb")), [](FILE* f) {if (f != nullptr) fclose(f);});
	if (f == nullptr || fwrite(bufPlain.get(), 1, lenPlain, f.get()) != lenPlain) {
		tstring text = TEXT("write file to [");
		text += filePath;
		text += TEXT("] failed \r\nerror number [");
		TCHAR err[6] = { 0 };
		text += _itow(GetLastError(), err, 10);
		text += TEXT("]");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK | MB_ICONWARNING);
		return;
	}
	tstring text = TEXT("convert success\r\nthe gfx result has written to [");
	text += filePath;
	text += TEXT("]");
	::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
#endif
}

void gfx2ozg() {
	// validate file extension
	TCHAR fileName[MAX_PATH] = { 0 };
	::SendMessage(nppData._nppHandle, NPPM_GETFILENAME, MAX_PATH, (LPARAM)fileName);
	TCHAR* ext = wcsrchr(fileName, TEXT('.'));
	if (ext == nullptr || wcscmp(ext, TEXT(".gfx")) != 0) {
		tstring text = TEXT("[");
		text += fileName;
		text += TEXT("]");
		text += TEXT(" is not a [*.gfx]");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK | MB_ICONWARNING);
		return;
	}

	// Get the current scintilla
	int which = -1;
	::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
	if (which == -1)
		return;
	HWND curScintilla = (which == 0) ? nppData._scintillaMainHandle : nppData._scintillaSecondHandle;

	// get current text
	int lenPlain = ::SendMessage(curScintilla, SCI_GETLENGTH, 0, 0);
	if (lenPlain == 0) {
		tstring text = TEXT("convert failed\r\n[");
		text += fileName;
		text += TEXT("] is empty");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
		return;
	}
	auto bufPlain = std::make_unique<unsigned char[]>(lenPlain + 1);
	::SendMessage(curScintilla, SCI_GETTEXT, lenPlain + 1, (LPARAM)bufPlain.get());

	// encrypt
	size_t lenCipher = RandomCipher::encrypt(nullptr, bufPlain.get(), lenPlain);
	auto bufCipher = std::make_unique<unsigned char[]>(lenCipher);
	RandomCipher::encrypt(bufCipher.get(), bufPlain.get(), lenPlain);

	// marshal

	// write
	::SendMessageW(nppData._nppHandle, NPPM_MENUCOMMAND, 0, IDM_FILE_NEW);
	which = -1;
	::SendMessage(nppData._nppHandle, NPPM_GETCURRENTSCINTILLA, 0, (LPARAM)&which);
	if (which == -1)
		return;
	curScintilla = (which == 0) ? nppData._scintillaMainHandle : nppData._scintillaSecondHandle;
	::SendMessage(curScintilla, SCI_ADDTEXT, lenCipher, (LPARAM)bufCipher.get());
#if 0
	// write bytes to a file
	TCHAR filePath[MAX_PATH] = { 0 };
	::SendMessage(nppData._nppHandle, NPPM_GETFULLCURRENTPATH, MAX_PATH, (LPARAM)filePath);
	*wcsrchr(filePath, TEXT('.')) = 0;
	wcscat(filePath, TEXT(".ozg"));
	std::unique_ptr<FILE, void(*)(FILE*)> f(_wfopen(filePath, TEXT("wb")), [](FILE* f) {if (f != nullptr) fclose(f);});
	if (f == nullptr || fwrite(bufCipher.get(), 1, lenCipher, f) != lenCipher) {
		tstring text = TEXT("write file to [");
		text += filePath;
		text += TEXT("] failed \r\nerror number [");
		TCHAR err[6] = { 0 };
		text += _itow(GetLastError(), err, 10);
		text += TEXT("]");
		::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK | MB_ICONWARNING);
		return;
	}
	tstring text = TEXT("convert success the ozg result has written to [");
	text += filePath;
	text += TEXT("]");
	::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
#endif
}

void about() {
	tstring text = TEXT("version: ");
	text += TEXT(VERSION);
	::MessageBox(nppData._nppHandle, text.c_str(), TEXT("mu editor"), MB_OK);
}
