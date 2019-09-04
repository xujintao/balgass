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