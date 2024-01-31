#include "stdafx.h"
#include "offset.h"
#include "c1c2.h"
#include "hook.h"

void rewriteVersion() {
	char* src = (char*)0x0123C754; // S9, "1.xx.xx"
	char* dst = (char*)0x01207020; // S9
	// marshal version string
	for (int i = 0; i < 8; i++) {
		dst[i] = src[i] + i + 1;
	}
}

void __declspec(naked) RewriteVersion() {
	rewriteVersion();
	_asm {
		push 0xC;
		mov eax, 0x00D0ACC8; // S9
		call eax;
		mov eax, 0x004D9B2F; // S9
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

typedef void(*tmuConnectToCS)(char* hostname, int port);
tmuConnectToCS muConnectToCS = (tmuConnectToCS)0x0063CEDA; // S9
void OnConnect()
{
	char * ip = "10.1.2.3";
	int port = 44405;
	muConnectToCS(ip, port);
}

void __declspec(naked) onCsHook()
{
	//004DE303   50               PUSH EAX
	//004DE304   E8 68701500      CALL main_IGC.00635371
	OnConnect();
	_asm
	{
		push eax;
		//call 0x00635371; // it's not needed, 635371 is original Connect function, we use muConnectToCS
		mov edx, 0x004DEA9A; // S9
		jmp edx;
	}
}

void SetHook()
{
	HookThis(PARSE_PACKET_HOOK, 5, (DWORD)&ParsePacket); // S9
	HookThis(SEND_PACKET_HOOK, 8, (DWORD)&SendPacket); // S9
	HookManager.MakeJmpHook(0x004D9B28, 7, RewriteVersion); // S9 rewrite version
	//MemAssign(0x0AA3097E + 3, (BYTE)8); // 1.04R handle version matching 5->8 characters
	//MemAssign(0x0043EBFA + 6, (BYTE)8); // 1.04R login send version 5->8 characters
	//MemAssign(0x004FF24C + 6, (BYTE)8); // 1.04R map server send version 5->8 characters
	MemAssign(0x004FAD40 + 6, (BYTE)8); // S9 ? version 5->8 characters
	MemAssign(0x006E8699 + 6, (BYTE)8); // S9 login send version 5->8 characters
	MemAssign(0x006E8F1E + 6, (BYTE)8); // S9 ? version 5->8 characters
	MemAssign(0x0A17BF30 + 3, (BYTE)8); // S9 handle version matching 5->8 character
	HookManager.MakeJmpHook(0x004DEA94, 6, onCsHook); // S9
	HookThis(0x0043EA60, 6, (DWORD)&OnConnect); // S9
	//MemAssign(0x004D7D42 + 1, (BYTE)0x42); // 1.04R support multi instance
	//MemAssign(0x0067760B + 1, (BYTE)0x2C); // 1.04R fixed item with socket show excellent
	//MemAssign(0x0AF7F02D + 1, (DWORD)0xC8); // 1.04R modify max speed limit of death stab from 300ms to 200ms per attack
	MemSet(CTRL_FREEZE_FIX, 0x02, 1); // S9
	//HookThis(0x004E4F57, 5, (DWORD)&RecordKey); // 1.04R
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