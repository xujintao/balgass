// dllmain.cpp : Defines the entry point for the DLL application.
#include "stdafx.h"
#include "Connection.h"
#include "offsets.h"
#include "HookManager.h"
#include "ItemManagement.h"
#include "DamageExtension.h"
#include "ResourceManager.h"
#include "ServerInfo.h"
#include "ItemLoad.h"
#include "packetsend_hooks.h"
#include "AntiHack.h"
#include "CustomWings.h"
#include "Interface.h"
//#include "Camera.h"
#include "WndHook.h"
#include "Macros.h"
#include "Glow.h"
#include "mdump.h"
#include "ReconnectSystem.h"
#include "Camera.h"
#include "protocol.h"
#include "giocp.h"
#include "Crack.h"
#include "GensExt.h"
#include "Controller.h"
#include "EventEntryLevel.h"

char IP[49];
WORD Port;
char Serial[17];

BYTE hookedrecv[5] = { 0 };
BYTE originalrecvheader[5] = { 0x8B, 0xFF, 0x55, 0x8B, 0xEC };

BYTE hookedwglswap[5] = { 0 };
BYTE originalwglswapheader[5] = { 0x8B, 0xFF, 0x55, 0x8B, 0xEC };

BYTE originalvirtualprotectex[5] = { 0x8B, 0xFF, 0x55, 0x8B, 0xEC };

typedef HMODULE (WINAPI *LoadLibraryFPtr)(LPCTSTR dllName);
LoadLibraryFPtr origFunc = LoadLibraryA;
LoadLibraryFPtr mineFunc = NULL;
DWORD dwAntiHackTime = 0;

DWORD dwHitHackTime = 0;

void __declspec(naked) onCsHook()
{
	//004DE303   50               PUSH EAX
	//004DE304   E8 68701500      CALL main_IGC.00635371
	OnConnect();

	_asm
	{
		//push eax; // 1.04R
		//call 0x00635371; // it's not needed, 635371 is original Connect function, we use muConnectToCS
		mov edx ,0x004E1D3B; // 1.04R
		jmp edx;
	}
}
HMODULE WINAPI LoadLibraryMine(LPCTSTR dllName)
{

}




void WriteMem(HANDLE hProcess, DWORD addr, DWORD mem, DWORD size)
{
	DWORD OldProtect;
	VirtualProtectEx(hProcess, (LPVOID)addr, size, PAGE_EXECUTE_READWRITE, &OldProtect);
	SIZE_T bytes = size;
	WriteProcessMemory(hProcess, (LPVOID)addr, (LPVOID)mem, size, &bytes);

	VirtualProtectEx(hProcess, (LPVOID)addr, size, OldProtect, &OldProtect);
}

void SetHook()
{
	HookThis(PARSE_PACKET_HOOK, 5, (DWORD)&ParsePacket); // 1.04R
	HookThis(SEND_PACKET_HOOK, 8, (DWORD)&SendPacket); // 1.04R
	HookThis(SPRINTF_HOOK, 5, (DWORD)sprintf); // 1.04R, why hook sprintf?
	//HookThis((DWORD)&Value, 0x0059E682);
	//HookThis((DWORD)&SetPrice, 0x00594618); // S8 E3 
	HookThis(DMG_SEND_HOOK, 6, (DWORD)&dmgSendHook); // 1.04R
	HookThis(DMG_SEND_RF_HOOK, 6, (DWORD)&dmgSendRFHook); // 1.04R
	HookThis(HP_SEND_HOOK, 7, (DWORD)&HPSendHook); // 1.04R
	HookThis(MANA_SEND_HOOK, 7, (DWORD)&manaSendHook); // 1.04R
	HookThis(HP_SEND_DK_HOOK, 7, (DWORD)&HPSendDKHook); // 1.04R
	HookThis(HP_SEND_C_HOOK, 5, (DWORD)&HPSendCHook); // 1.04R
	HookThis(MANA_SEND_DK_HOOK, 7, (DWORD)&ManaSendDKHook); // 1.04R
	HookThis(MANA_SEND_C_HOOK, 5, (DWORD)&ManaSendCHook); // 1.04R
	HookThis(SKILL_65K_HOOK, 9, (DWORD)&FixSkill65k); // 1.04R
	HookThis(ERRTEL_MIX_STAFF_HOOK, 8, (DWORD)&ErrtelMixStaffFix); // 1.04R
	HookThis(ITEM_LEVEL_REQ_HOOK, 6, (DWORD)&ItemLevelReqFix); // 1.04R
	HookThis(MOVE_WND_HOOK, 6, (DWORD)&MoveWndFix); // 1.04R
	//AlterPShopDisplayCurrency()
	HookThis(ALTER_PSHOP_DISPLAY_HOOK, 7, (DWORD)&AlterPShopDisplayCurrency); // 1.04R
	//g_Camera->Init();
	MemAssign(MU_WND_PROC_HOOK, FPTR(WndProc));

	// HD Resolution
	//HookManager.MakeJmpHook(RESOLUTION_HOOK, MuLoadResolutionAsm);
	//GCSetCharSet(1255);
	//00591530 wing
	//0099500B
	//0063D3E3   51               PUSH ECX

	//0059E68B
}
void SpeedCheckThread()
{

	while(true){
	
//	static DWORD tick = GetTickCount();
/*	while ((GetTickCount() - tick) < 20000) 
		{
			Sleep(10);
		}
	/*
		while ((GetTickCount() - tick) < 20000) 
		{
			Sleep(10);
		}
		if(*g_Map != -1 || *g_Map != 0x4a || *g_Map != 0x49 || *g_Map != 0x4E )
		{
			CGSendAgi(0);
		}

		tick = GetTickCount();*/
	}

}
/*void LoadNewItems()
{
//	CreateThread(NULL, NULL, (LPTHREAD_START_ROUTINE)SpeedCheckThread, NULL, NULL, NULL);
}*/

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

void SetValues()
{
	DWORD OldProtect;

	//ResourceManager.Preloader(); // 1.04R has no glow.bmd
	//g_Controller.Load(); 1.04R
	g_ServerInfo->Load(".\\Data\\Local\\ServerInfo.bmd");

	//MemSet(0x004D7E2A, 0xEB, 1); // 1.04R, 73->EB, jae->jmp, disable mu.exe, search "mu.exe"
	MemAssign(0x004D7E2A, (BYTE)0xEB); // MemSet can not work, why?
	MemSet(0x004D8102, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "config.ini read error"
	MemSet(0x004D8145, 0xEB, 1); // 1.04R, 75->EB, jne->jmp, disable GameGuard, search "gg init error"
	MemSet(0x004DD639, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	MemSet(0x006B7B2A, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	MemSet(0x006B7C63, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable GameGuard, search "ResourceGuard"
	MemSet(0x00B38653, 0xC3, 1); // 1.04R, 55->C3, push ebp->ret, disable log encode
	MemAssign(0x00E1A4AC, (WORD)0x90C3); // 1.04R, 6A 02->C3 90, push 02->ret, fix r6602 floating point support not loaded, search cmd "mov dword ptr ds:[ebx+0xC], esi"
	//MemSet(0x0ABD696B, 0xEB, 1); // 1.04R, 74->EB, je->jmp, disable main.exe version match
	MemSet(0x0AD33703, 0xC3, 1); // 1.04R, 55->C3, push ebp->ret, disable some encode
	HookThis(0x0043EA4C, 7, (DWORD)&LoginHook1); // 1.04R, extend pwd length, fill pwd buf
	HookThis(0x0043EAA2, 9, (DWORD)&LoginHook2); // 1.04R, extend pwd length, validate pwd length
	HookThis(0x0043EB4E, 7, (DWORD)&LoginHook3); // 1.04R, extend pwd length, xor pwd buf
	MemAssign(0x0051B429 + 1, (DWORD)0x0A87D703); // 1.04R, handleState2, disable anti-temper with backup code
	MemAssign(0x004E17B9 + 1, (DWORD)0x099DAE9F); // 1.04R, handleState4, disable anti-temper with backup code
	MemAssign(0x005B5D6C + 1, (DWORD)0x09A1DB88); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x004E0E03 + 1, (DWORD)0x0A36BF12); // 1.04R, handleState5, disable anti-temper with backup code
	MemAssign(0x006C18B6 + 1, (DWORD)0x0A33476B); // 1.04R, handleF6, jump over complicated shell logic, which would send trap message
	//MemAssign(0x005B0252 + 6, (BYTE)0xD4); // 1.04R, position send, D7->D4, compatible with s9 server
	MemAssign(0x0075FF6A + 0xD4, *(BYTE*)(0x0075FF6A + 0xD7)); // 1.04R, position recv, D7->D4, compatible with s9 server
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


	MemSet(0x00AF4B68+3, 7, 1); // 1.04R, Option +28
	// GCSetCharSet(g_ServerInfo->GetCharset());
	//HookManager.MakeJmpHook(CONNECT_HOOK1, 5, onCsHook);
	HookManager.MakeCallback(CONNECT_HOOK1, OnConnect, 0, 5, false, false);
	HookManager.MakeCallback(CONNECT_HOOK2, OnConnect, 0, 5, false, false);
	MemAssign(VERSION_HOOK1, (DWORD)g_ServerInfo->GetVersion());
	MemAssign(VERSION_HOOK2, (DWORD)g_ServerInfo->GetVersion());
	MemAssign(VERSION_HOOK3, (DWORD)g_ServerInfo->GetVersion());
	MemAssign(SERIAL_HOOK1, (DWORD)g_ServerInfo->GetSerial());
	MemAssign(SERIAL_HOOK2, (DWORD)g_ServerInfo->GetSerial());
	MemAssign(CHATSERVER_PORT,(DWORD)g_ServerInfo->GetChatServerPort());
	HookManager.MakeCallback(GSCONNECT_HOOK1, OnGsConnect, 0, 10, false);
	HookManager.MakeCallback(GSCONNECT_HOOK2, OnGsConnect, 0, 10, false);
	MemSet(CTRL_FREEZE_FIX, 0x02, 1);
	//MemSet(MAPSRV_DELACCID_FIX, 0x90, 5); // nop call to memcpy in mapserver, because buff source is empty, due to no account id from webzen http (remove new login system)
	//MemSet(MAPSRV_DELACCID_FIX, 0x90, 5); // this may be annotated in 1.04R
	HookManager.MakeCallback(RENDER_HP_BAR_CALLBACK, RenderLoadingBar, 0, 5, true); // 1.04R
	

	HookManager.MakeJmpHook(0x004E56C8, 10, HookDCFunc); // 1.04R, reconnect
	//HookManager.MakeJmpHook(0x004E1B32, PreventSwapBufferHack); //s9
	//HookManager.MakeJmpHook(0x004E1B1E, PreventSwapBufferHack2); //s9
	//PreventSwapBufferHack2
	//009C6990   68 E0C60401      PUSH main.0104C6E0                       ; ASCII "> Menu - Exit game. "
	//009EDC9B   68 00D80901      PUSH main.0109D800                       ; ASCII "> Menu - Exit game. "


	HookManager.MakeJmpHook(0x00ACD10F, 5, HookExitFunc); // 1.04R
	HookManager.MakeJmpHook(0x0043F7EA, 5, HookExitCharSelectFunc); // 1.04R
	ResourceManager.LoadGraphics();

	// HookManager.MakeJmpHook(SETGLOW_EX, 6, OnSetGlowEx); // 1.04R has no glow.bmd

	MemSet(0x00A29AC3, 0x90, 2); // 1.04R, Exc socket item FIX
	MemSet(0x00A29448, 0xEB, 1); // 1.04R, Exc Socket Visual Bug Option FIX

	MemSet(0x006C8E82, 0x90, 3); // 1.04R, Fenrir exc visual FIX 1
	MemSet(0x006C8E9E, 0x90, 2); // 1.04R, Fenrir exc visual FIX 2
	MemSet(0x006C8EB6, 0xEB, 1); // 1.04R, Fenrir exc visual FIX 3

	MemSet(0x00AC525B,0xEB,1); // 1.04R -- Cherry Blossom Opening MuRuumy Windows Fix, JNZ -> JMP
	MemSet(0x008E55FC,0x90,6); // 1.04R -- Joh option display on Ancients
	MemSet(0x00B0F0F4+3, 0x30, 1); // 1.04R -- Decreased limit of screen refresh frequency, 60 -> 48
	MemSet(0x00B0F0F8, 0x7D, 1); // 1.04R, JE -> JGE -- limit check == -> >=
	// MemSet(0x006440A7 + 2, 0x26, 1); // 1.04R may have fixed it, Box of Kundun 2-5 fix
	MemSet(0x00644068 + 3, 0x07, 1); // 1.04R, Box of Heaven fix

	RFSkillDamageDisplayFix();
	FixTraps();

	GensInitData();
	HookManager.MakeJmpHook(0x00640FF7, 6, GensIsBattleZoneMap); // 1.04R
	HookManager.MakeJmpHook(0x00ABBF4C, 6, GensWarpMenuFix); // 1.04R
	MemSet(0x00ABBE96+1, 0x85, 1); // 1.04R, je->jne

	// LoadBCEntryLevel(); // private custom !!!
	// LoadDSEntryLevel(); // private custom !!!

	// Launcher init
	SYSTEMTIME Date;
	GetLocalTime(&Date);
	DWORD dwParam = (Date.wDay * Date.wMonth * Date.wYear) ^ 0x12894376 + 0xFFFFFFF;
	if (g_ServerInfo->GetLauncher())
	{
		LPSTR cmdline = GetCommandLine();

		int len = strlen(cmdline);
		int i;
		for (i = (len - 1); i >= 0; i--)
		{
			if (cmdline[i] == ' ')
				break;
		}

		char* parameter = &cmdline[i + 1];
		int paramLen = strlen(parameter);
		bool isNum = true;
		for (int i = 0; i < paramLen; i++)
		{
			if (parameter[i] < '0' || parameter[i] > '9')
				isNum = false;
		}
		DWORD dwParameter = 0;
		if (isNum)
		{
			dwParameter = atoi(parameter);
		}

		if (dwParam != dwParameter)
		{
			WinExec(g_ServerInfo->GetLauncherFile().c_str(), SW_SHOW);
			ExitProcess(0);
		}
	}

	if (g_ServerInfo->GetPatcher())
	{
		LPSTR cmdline = GetCommandLine();

		int len = strlen(cmdline);


		int i;
		for (i = (len - 1); i >= 0; i--)
		{
			if (cmdline[i] == ' ')
				break;
		}

		char* parameter = &cmdline[i + 1];
		int paramLen = strlen(parameter);

		bool isNum = true;
		for (int i = 0; i < paramLen; i++)
		{
			if (parameter[i] < '0' || parameter[i] > '9')
				isNum = false;
		}
		DWORD dwParameter = 0;

		if (!isNum || paramLen == 0)
		{
			WinExec(g_ServerInfo->GetPatcherFile().c_str(), SW_SHOW);
			ExitProcess(0);
		}
	}

	// item

	HookManager.MakeJmpHook(TOOLTIP_FIX, 9, ToolTipFix); // 1.04R
}

void RecvCheckThread()
{
	
	BYTE check[5];
	PMSG_ANTIHACK_CHECK pMsg;
	
	while (true)
	{
		DWORD dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("kernel32.dll"), "VirtualProtectEx");

		if (dwAddr == NULL)
		{
			int Error = GetLastError();
			char szTemp[128];
			sprintf(szTemp, "ERROR :%d", Error);
			MessageBox(0, szTemp, "ERROR", MB_OK);
		}

		ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, check, 5, 0);

		if (memcmp(originalvirtualprotectex, check, 5) != 0) // If recv is hooked earlier
		{
			DWORD OldProtect;
			VirtualProtect((LPVOID)(0x0116D0BB), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
			BYTE FuckYou[6] = { 0xEB, 0x44, 0x27, 0xA1, 0x04, 0xDE };
			memcpy((void*)0x0116D0BB, &FuckYou, 6); // Trap :: It will crash main on send packet :D
			GCSendAntihackBreach(40000);
		} 

		/* dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("opengl32.dll"), "wglSwapBuffers");
		ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, check, 5, 0);

		 if (memcmp(originalwglswapheader, check, 5) != 0) // If recv is hooked by some shit
		{
			DWORD OldProtect;
			VirtualProtect((LPVOID)(0x0043937C), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
			BYTE FuckYou[6] = { 0xE8, 0x9B, 0xE9, 0x98, 0x22, 0xC1 };
			memcpy((void*)0x0043937C, &FuckYou, 6); // Trap :: It will crash main on send packet
			GCSendAntihackBreach(60000);
		} */

		AHCheckSwapBufferMod();
		if (GetTickCount() - dwHitHackTime >= 3000 && g_Connection == GS_CONNECTED)
		{
			if( pPlayerState == GameProcess )
			{
				CheckHitHack();
			}
			dwHitHackTime = GetTickCount();
		}

		if (GetTickCount() - dwAntiHackTime >= 20000 && g_Connection == GS_CONNECTED)
		{
			
			pMsg.h.c = 0xC1;
			pMsg.h.size = sizeof(pMsg);
			pMsg.h.headcode = 0xFA;
			pMsg.subcode = 0x11;
			
			dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("ws2_32.dll"), "recv");
			ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, pMsg.checkdata, 5, 0);

			DataSend((LPBYTE)&pMsg, pMsg.h.size);
			
			dwAntiHackTime = GetTickCount();
		}

		Sleep(1000);
	}

	
}

void CheckMainFunctions()
{
	

	BYTE check[5];

	DWORD dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("kernel32.dll"), "VirtualProtectEx");
	ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, check, 5, 0);

	if (memcmp(originalvirtualprotectex, check, 5) != 0) // If recv is hooked earlier
	{
		DWORD OldProtect;
		VirtualProtect((LPVOID)(0x0116D0BB), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
		BYTE FuckYou[6] = { 0xEB, 0x44, 0x27, 0xA1, 0x04, 0xDE };
		memcpy((void*)0x0116D0BB, &FuckYou, 6); // Trap :: It will crash main on send packet
	}

	dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("ws2_32.dll"), "recv");
	ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, check, 5, 0);

	if (memcmp(originalrecvheader, check, 5) != 0) // If recv is hooked earlier
	{
		DWORD OldProtect;
		VirtualProtect((LPVOID)(0x0043937C), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
		BYTE FuckYou[6] = {0xE8, 0x9B, 0xE9, 0x98, 0x22, 0xC1};
		memcpy((void*)0x0043937C, &FuckYou, 6); // Trap :: It will crash main on send packet
	}

	hCheckThread = CreateThread(NULL, NULL, (LPTHREAD_START_ROUTINE)RecvCheckThread, NULL, NULL, NULL);

	dwAntiHackTime = GetTickCount();
	dwHitHackTime = GetTickCount();
	PrintModules(GetCurrentProcessId());

	
}

BOOL APIENTRY DllMain (HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved)
{
	switch (ul_reason_for_call)
	{
		case DLL_PROCESS_ATTACH:
		{
			CMiniDump::Begin();
			//DWORD OldProtect;
			//VirtualProtect((LPVOID)0x00401000, 12681214, PAGE_EXECUTE_READ, &OldProtect);
			CheckMainFunctions();
			// DoHooks();
			SetHook();
			SetValues();
			//g_Controller.m_Instance = hModule;
		}
		break;

		case DLL_PROCESS_DETACH:
			CMiniDump::End();
			break;
	}

	return TRUE;
}