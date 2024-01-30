#include "stdafx.h"
#include "offset.h"
#include "c1c2.h"
#include "hook.h"

void rewriteVersion() {
	char* src = (char*)0x01319A44; // 1.04R, "1.xx.xx"
	char* dst = (char*)0x012E4018; // 1.04R
								   // marshal version string
	for (int i = 0; i < 8; i++) {
		dst[i] = src[i] + i + 1;
	}
}

void __declspec(naked) RewriteVersion() {
	rewriteVersion();
	_asm {
		push 1;
		mov eax, 0x00DE852F;
		call eax;
		mov eax, 0x004D80F4;
		jmp eax;
	}
}

void __declspec(naked) LoginHook1() {
	__asm {
		push 0x15;
		push 0;
		lea eax, ss:[ebp - 0x30];
		push eax;
		mov edx, 0x00DE8100; // 1.04R
		call edx;
		add esp, 0xC;
		push 0x15;
		lea eax, ss:[ebp - 0x30];
		push eax;
		mov edx, 0x0043EA62; // 1.04R
		jmp edx;
	}
}

void __declspec(naked) LoginHook2() {
	__asm {
		lea eax, ss:[ebp - 0x30];
		push eax;
		mov edx, 0x0043EE13; // 1.04R
		call edx;
		mov edx, 0x0043EAAB; // 1.04R
		jmp edx;
	}
}

void __declspec(naked) LoginHook3() {
	__asm {
		push 0xA;
		lea eax, ss:[ebp - 0x18];
		push eax;
		mov edx, 0x0043B750; // 1.04R
		call edx;
		pop ecx;
		pop ecx;
		push 0x14;
		lea eax, ss:[ebp - 0x30];
		push eax;
		mov edx, 0x0043B750; // 1.04R
		call edx;
		pop ecx;
		pop ecx;
		push 1;
		push 0xA;
		lea eax, ss:[ebp - 0x18];
		push eax;
		lea ecx, ss:[ebp - 0x14BC];
		mov edx, 0x00439298; // 1.04R
		call edx;
		push 1;
		push 0x14;
		lea eax, ss:[ebp - 0x30];
		push eax;
		mov edx, 0x0043EBC7; // 1.04R
		jmp edx;
	}
}

void __declspec(naked) RecordKey() {
	__asm {
		call dword ptr ds : [0x01149744] // GetActiveWindow
			cmp dword ptr ds : [0x01319D6C], eax
			jne skip
			mov eax, 0x008AF00D
			call eax
			mov ecx, eax
			mov eax, 0x008AF06A
			call eax
			skip :
		mov eax, 0x004E4F63
			jmp eax
	}
}

void SetHook()
{
	HookThis(PARSE_PACKET_HOOK, 5, (DWORD)&ParsePacket); // 1.04R
	HookThis(SEND_PACKET_HOOK, 8, (DWORD)&SendPacket); // 1.04R
	MemSet(0x004D8102, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "config.ini read error"
	MemSet(0x004D8145, 0xEB, 1); // 1.04R, 75->EB, jne->jmp, disable GameGuard, search "gg init error"
	MemSet(0x004DD639, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	MemSet(0x006B7B2A, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	MemSet(0x006B7C63, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	char disableRGonLoadCharacterInfo[6] = { 0xE9, 0xA0, 0x4C, 0xDB, 0xF5, 0x90 };
	MemCpy(0x0A9503D2, disableRGonLoadCharacterInfo, sizeof(disableRGonLoadCharacterInfo)); // 1.04R, je->jmp, disable GameGuard, search "ResourceGuard"
	MemAssign(0x00701DF4 + 1, (DWORD)0x0A4F6808); // 1.04R, jump over shell logic: check integrity on load character info
	MemSet(0x00B38653, 0xC3, 1); // 1.04R, 55->C3, push ebp->ret, disable log encode
	MemAssign(0x00E1A4AC, (WORD)0x90C3); // 1.04R, 6A 02->C3 90, push 02->ret, fix r6602 floating point support not loaded, search cmd "mov dword ptr ds:[ebx+0xC], esi"
	//MemSet(0x0ABD696B, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable main.exe version match
	HookManager.MakeJmpHook(0x004D80ED, 7, RewriteVersion); // rewrite version
	MemAssign(0x0AA3097E + 3, (BYTE)8); // handle version matching 5->8 characters
	MemAssign(0x0043EBFA + 6, (BYTE)8); // login send version 5->8 characters
	MemAssign(0x004FF24C + 6, (BYTE)8); // map server send version 5->8 characters
	MemSet(0x0AD33703, 0xC3, 1); // 1.04R, 55->C3, push ebp->ret, disable some encode
	//HookThis(0x0043EA4C, 7, (DWORD)&LoginHook1); // 1.04R, extend pwd length, fill pwd buf
	//HookThis(0x0043EAA2, 9, (DWORD)&LoginHook2); // 1.04R, extend pwd length, validate pwd length
	//HookThis(0x0043EB4E, 7, (DWORD)&LoginHook3); // 1.04R, extend pwd length, xor pwd buf
	MemAssign(0x0051B429 + 1, (DWORD)0x0A87D703); // 1.04R, handleState2, disable anti-temper with backup code
	MemAssign(0x004E17B9 + 1, (DWORD)0x099DAE9F); // 1.04R, handleState4, disable anti-temper with backup code
	MemAssign(0x005B5D6C + 1, (DWORD)0x09A1DB88); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x004E0E03 + 1, (DWORD)0x0A36BF12); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x006C18B6 + 1, (DWORD)0x0A33476B); // 1.04R, handleF6, jump over complicated shell logic, which would send trap message
	MemAssign(0x006EAB21 + 1, (DWORD)0x0A8283B4); // 1.04R, handleA9petItemInfo, jump over complicated shell logic that would send trap message
	//MemAssign(0x005B0252 + 6, (BYTE)0xD4); // 1.04R, position set, D7->D4, compatible with s9 server
	MemAssign(0x0075FF6A + 0xD4, *(BYTE*)(0x0075FF6A + 0xD7)); // 1.04R, position set, D7->D4, compatible with s9 server
	//MemAssign(0x00A16091 + 1, (BYTE)0x15); // 1.04R, position set, DA->15, compatible with s9 server
	MemAssign(0x0075FF6A + 0x15, *(BYTE*)(0x0075FF6A + 0xDA)); // 1.04R, position get, DA->15, compatible with s9 server
	DWORD inlineAddr[] = {
		0x004E70D9, 0x005ABA01, 0x0A8FB8B9, 0x09FC2882, 0x005CEB05, 0x005D0607, 0x005D1B5C, 0x005D38B5,
		0x005D4974, 0x0060594F, 0x0060C8E1, 0x006144CF, 0x00615B39, 0x006171A3, 0x006185CB, 0x00619761,
		0x0061AC07, 0x0061BE96, 0x0061D198, 0x0061E43E, 0x0061F725, 0x00620BE2, 0x00621E89, 0x00623130,
		0x006243D7, 0x0062571C, 0x006269EA, 0x00627D25, 0x00629028, 0x00638A4C, 0x00639EE5, 0x0063B39E,
		0x00650BE9, 0x00651E00, 0x0A8F9C57, 0x006D026C, 0x09EB7E04, 0x006D3CCA, 0x0AF1145B, 0x0A9D549B,
		0x0A892721, 0x09FB6290, 0x0A891849, 0x0A05DCAB, 0x0A3A1A6A, 0x0A43AABA, 0x0070AB31, 0x09FB93F0,
		0x0AF824F4, 0x0A4ECC2F, 0x0AD70413, 0x0A5D32BC, 0x0AA06551, 0x0AF8A8D1, 0x0A4E848A, 0x09F8469D,
	};
	for (auto addr : inlineAddr) {
		MemAssign(addr, (BYTE)0x90); // inc eax->nop, inline send disable encrypt
	}
	//MemAssign(0x00612007 + 6, (BYTE)0x11); // 1.04R, normal attack send, D9->11, compatible with s9 server
	MemAssign(0x0075FF6A + 0x11, *(BYTE*)(0x0075FF6A + 0xD9)); // 1.04R,  normal attack send recv, D9->11, compatible with s9 server
	MemAssign(0x00A0A5E1 + 1, (DWORD)0x0941B4CA); // 1.04R, jump over complicated shell logic, which would send trap message
	//MemAssign(0x0A0C438C + 6, (BYTE)0xDF); // 1.04R, be attacked send, 1D->DF, compatible with s9 server
	MemAssign(0x004D7D42 + 1, (BYTE)0x42); // support multi instance
	// MemAssign(0x00B0F1C1 + 1, (BYTE)0x00); // discard high dmDisplayFrequency, fixed resolution mismatching
	MemAssign(0x0067760B + 1, (BYTE)0x2C); // fixed item with socket show excellent
	MemAssign(0x0AF7F02D + 1, (DWORD)0xC8); // modify max speed limit of death stab from 300ms to 200ms per attack
	MemSet(CTRL_FREEZE_FIX, 0x02, 1);
	HookThis(0x004E4F57, 5, (DWORD)&RecordKey);
}

BOOL APIENTRY DllMain(HMODULE hModule,
	DWORD  ul_reason_for_call,
	LPVOID lpReserved
	)
{
	switch (ul_reason_for_call)
	{
	case DLL_PROCESS_ATTACH:
		SetHook();
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}