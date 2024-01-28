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
	HookThis(PARSE_PACKET_HOOK, 5, (DWORD)&ParsePacket); // S9
	HookThis(SEND_PACKET_HOOK, 8, (DWORD)&SendPacket); // S9
	HookManager.MakeJmpHook(0x004D80ED, 7, RewriteVersion); // 1.04R rewrite version
	MemAssign(0x0AA3097E + 3, (BYTE)8); // 1.04R handle version matching 5->8 characters
	MemAssign(0x0043EBFA + 6, (BYTE)8); // 1.04R login send version 5->8 characters
	MemAssign(0x004FF24C + 6, (BYTE)8); // 1.04R map server send version 5->8 characters
	MemSet(0x0AD33703, 0xC3, 1); // 1.04R, 55->C3, push ebp->ret, disable some encode
	HookThis(0x0043EA4C, 7, (DWORD)&LoginHook1); // 1.04R, extend pwd length, fill pwd buf
	HookThis(0x0043EAA2, 9, (DWORD)&LoginHook2); // 1.04R, extend pwd length, validate pwd length
	HookThis(0x0043EB4E, 7, (DWORD)&LoginHook3); // 1.04R, extend pwd length, xor pwd buf
	MemAssign(0x0051B429 + 1, (DWORD)0x0A87D703); // 1.04R, handleState2, disable anti-temper with backup code
	MemAssign(0x004E17B9 + 1, (DWORD)0x099DAE9F); // 1.04R, handleState4, disable anti-temper with backup code
	MemAssign(0x005B5D6C + 1, (DWORD)0x09A1DB88); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x004E0E03 + 1, (DWORD)0x0A36BF12); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x004D7D42 + 1, (BYTE)0x42); // 1.04R support multi instance
	// MemAssign(0x00B0F1C1 + 1, (BYTE)0x00); // 1.04R discard high dmDisplayFrequency, fixed resolution mismatching
	MemAssign(0x0067760B + 1, (BYTE)0x2C); // 1.04R fixed item with socket show excellent
	MemAssign(0x0AF7F02D + 1, (DWORD)0xC8); // 1.04R modify max speed limit of death stab from 300ms to 200ms per attack
	MemSet(CTRL_FREEZE_FIX, 0x02, 1); // S9
	HookThis(0x004E4F57, 5, (DWORD)&RecordKey); // 1.04R
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