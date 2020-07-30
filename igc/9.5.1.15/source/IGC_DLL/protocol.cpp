#include "stdafx.h"
#include "protocol.h"
#include "DamageExtension.h"
#include "AntiHack.h"
#include "ReconnectSystem.h"
#include "offsets.h"
#include "packetsend_hooks.h"
#include "HookManager.h"
#include "EncDec.h"
#include "giocp.h"
#include "ItemManagement.h"
#include "GensExt.h"
#include "Interface.h"
#include "ServerInfo.h"

// global variables
char ChatServerIp[17];

//start of protocolcore
bool CliProtocolCore(LPBYTE aRecv, BYTE ProtoNum, int len, bool Encrypt)
{
	switch(ProtoNum)
	{
		/*
		case 0xF1:
			{
				PMSG_DEFAULT2 * lpMsg2 = (PMSG_DEFAULT2 *)aRecv;
				switch (lpMsg2->subcode)
				{
					case 0x00:
						GCJoinResultSend((PMSG_JOINRESULT*)aRecv);
						break;
				}
			}
			break;
		*/
		case 0xF3:
			{
				PMSG_DEFAULT2 * lpMsg2 = (PMSG_DEFAULT2 *)aRecv;
				switch (lpMsg2->subcode)
				{
				case 0x03:
					GCCharMapJoinResult((PMSG_CHARMAPJOINRESULT*)aRecv);
					break;
				case 0x04:
					GCCharRegen((PMSG_CHARREGEN*)aRecv);
					break;
				case 0x05:
					GCLevelUpSend((PMSG_LEVELUP*)aRecv);
					break;
				case 0x06:
					GCLevelUpPointAdd((PMSG_LVPOINTADDRESULT*)aRecv);
					break;
				case 0x51:
					GCMasterLevelUpSend((PMSG_MASTER_LEVEL_UP_SEND*)aRecv);
					break;
				}
			}
			break;
		case 0xBF:
		{
			PMSG_DEFAULT2 * lpMsg2 = (PMSG_DEFAULT2*)aRecv;
			switch (lpMsg2->subcode)
			{
			case 0x51:
				GCAnsMuBotRecvUse((PMSG_MUBOT_USE*)aRecv);
				break;
			}
		}
		break;
		/*
		case 0xF4:
			{
				// for reconnect must disable 0xF4 protocol or client crash
				BYTE internalProto = aRecv[3];
				if (g_ReconnectProcess >= 2 && g_ReconnectProcess != 6)
				{
					switch (internalProto)
					{
					case 0x03:
					{
						reconnectBuf = aRecv;
						g_ReconnectProcess = 3;

					}
					break;

					}
					return false;
				}
			}
			break;
		*/
		case 0xCA:
			GCFriendRoomCreate((PMSG_FRIEND_ROOMCREATE_RESULT*)aRecv);
			break;
		case 0xD4:
			AHCheckGetTickHook();
			break;
		case 0xFA:
		{
			if (aRecv[0] == 0xC1)
			{
				PMSG_DEFAULT2 * lpMsg2 = (PMSG_DEFAULT2 *)aRecv;
				switch (lpMsg2->subcode)
				{
				case 0x0D:
					GCOffTradeReq(aRecv);
					break;
				case 0x10:
				case 0x11:
					GCCustomPost((PMSG_POST_DATA *)aRecv);
					break;
				case 0x12:
					GCSetChatColors((PMSG_SET_CHAT_COLOR *)aRecv);
					break;
				case 0xA0:
					GCEnableSiegeSkills();
					break;
				case 0xA2:
					GCSetAgilityFix((PMSG_SETAGILITYBUG*)aRecv);
					break;
				case 0xA3:
					GCSetCashItemMoveEnable();
					break;
				case 0xA4:
					HookThis(0x00AF1790, 7, (DWORD)&LuckyItemToolTipFix); // 1.04R
					HookThis(0x00AF3484, 7, (DWORD)&LuckyItemEnhancementFix); // 1.04R
					break;
				case 0xA5:
					GCDisableReconnect();
					break;
				case 0xA6:
					ExitProcess(0);
					break;
				case 0xA7:
					GCDisableReconnectSystem();
					break;
				case 0xA8:{
					PMSG_ALTERPSHOPVAULT * lpMsg = (PMSG_ALTERPSHOPVAULT*)aRecv;
					GCAlterPShopVault(lpMsg->type);}
					break;
				case 0x90:
					GCIGCStatsAdd((PMSG_ADDSTATS*)aRecv);
					break;
				case 0xB0:
					GCDropSellMod();
					break;
					
				}
				return false;
			}
			else if (aRecv[0] == 0xC2)
			{
				PWMSG_DEFAULT2 * lpMsg2 = (PWMSG_DEFAULT2*)aRecv;
				switch (lpMsg2->subcode)
				{
					case 0xA4:
						GCRareItemSetPrices((PMSG_RAREITEMPIRCESEND*)aRecv);
						break;
					case 0xF8:
						GCSetBattleZoneData((_tagPMSG_SEND_GENS_BATTLEZONE_DATA *)aRecv);
						break;
				}
				return false;
			}
		}
			break;
		case 0x11:
			GCAttackResult((PMSG_ATTACKRESULT*)aRecv);
			break;
		case 0x17:
			GCDiePlayerSend((PMSG_DIEPLAYER*)aRecv);
			break;
		case 0x26:
			GCRefillSend((PMSG_REFILL*)aRecv);
			break;
		case 0x27:
			GCManaSend((PMSG_MANASEND*)aRecv);
			break;
		case 0x44:
			handlePartyHPMP((PMSG_PARTY_COUNT*)aRecv);
			break;
	}
	return true;
}

void GCDiePlayerSend(PMSG_DIEPLAYER *lpMsg)
{
	int aIndex = *(DWORD*)0x08C88E0C; // aIndex 1.04R

	int index = MAKE_NUMBERW(lpMsg->NumberH, lpMsg->NumberL);

	if (index == aIndex)
	{
		hpDisplay = 0;
	}
}

void GCAttackResult(PMSG_ATTACKRESULT *lpMsg)
{
	

	int iIndex = *(DWORD*)0x08C88E0C; // aIndex 1.04R
	dmgDisplay = lpMsg->IGCDamage;

	int IndexH = lpMsg->NumberH;
	int IndexL = lpMsg->NumberL;

	IndexH &= 0x7F; // Delete MSB Flag

	int AttackIndex = MAKE_NUMBERW(IndexH, IndexL);

	if (AttackIndex == iIndex)
	{
		if (hpDisplay - dmgDisplay <= 0)
		{
			hpDisplay = 0;
		}
		else
		{
			hpDisplay -= dmgDisplay;
		}
	}

	

	

	BYTE check[5];
	DWORD dwAddr = (DWORD)GetProcedureAddress(GetModuleHandle("kernel32.dll"), "QueryPerformanceCounter");
	ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, check, 5, 0);

	if (memcmp(originalqueryperformancecounter, check, 5) != 0) // If recv is hooked by some shit
	{
		DWORD OldProtect;
		VirtualProtect((LPVOID)(0x00635438), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
		BYTE FuckYou[6] = { 0x48, 0xC1, 0x46, 0x55, 0xEB, 0x22 };
		memcpy((void*)0x00635438, &FuckYou, 6); // Trap :: It will crash main on send packet :D
		GCSendAntihackBreach(10000);
	}

	
}

void GCRefillSend(PMSG_REFILL* lpMsg)
{

	if (lpMsg->IPos == 0xFE)
	{
		maxHPDisplay = lpMsg->IGCLife;
	}
	else if (lpMsg->IPos == 0xFF)
	{
		hpDisplay = lpMsg->IGCLife;
	}

	

	BYTE checksh[3];
	BYTE original[] = { 0x83, 0xE8, 0x28 };
	DWORD dwAddr = 0x004E50ED; // 1.04R
	ReadProcessMemory(GetCurrentProcess(), (LPVOID)dwAddr, checksh, 3, 0);

	if (memcmp(checksh, original, 3) != 0)
	{
		DWORD OldProtect;
		VirtualProtect((LPVOID)(0x00552B5C), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
		BYTE trapArray[6] = { 0xE8, 0x22, 0xA9, 0x99, 0x09, 0xEb };
		memcpy((void*)0x00552B5C, &trapArray, 6); // Trap :: crash on using potion
		GCSendAntihackBreach(60000);
	}
	
}

void GCManaSend(PMSG_MANASEND *lpMsg)
{
		if (lpMsg->IPos == 0xFE)
		{
			maxmanaDisplay = lpMsg->IGCMana;
		}
		else if (lpMsg->IPos == 0xFF)
		{
			manaDisplay = lpMsg->IGCMana;
		}

		/*
		DWORD dwExitCode;
		GetExitCodeThread(hCheckThread, &dwExitCode);

		if (dwExitCode != STILL_ACTIVE)
		{
			DWORD OldProtect;
			VirtualProtect((LPVOID)(0x009C51496), 6, PAGE_EXECUTE_READWRITE, &OldProtect);
			BYTE trapArray[6] = { 0xEB, 0x28, 0x16, 0x77, 0xAF, 0x4B };
			memcpy((void*)0x009C51496, &trapArray, 6); // Trap :: It will crash main on send packet :D
			GCSendAntihackBreach(20000);
		}
		*/
}

void GCJoinResultSend(PMSG_JOINRESULT *lpMsg)
{
	
	if (g_ReconnectProcess >= 4 && g_ReconnectProcess != 6)
	{
		g_ReconnectProcess = 5;
	}
	
}

void GCCharMapJoinResult(PMSG_CHARMAPJOINRESULT *lpMsg)
{
	hpDisplay = lpMsg->IGCLife;
	maxHPDisplay = lpMsg->IGCMaxLife;
	manaDisplay = lpMsg->IGCMana;
	maxmanaDisplay = lpMsg->IGCMaxMana;

	CGSendFilesCrc();
}

void GCCharRegen(PMSG_CHARREGEN * lpMsg)
{
	hpDisplay = lpMsg->IGCLife;
	manaDisplay = lpMsg->IGCMana;
}

void GCLevelUpSend(PMSG_LEVELUP *lpMsg)
{
	maxHPDisplay = lpMsg->IGCMaxLife;
	maxmanaDisplay = lpMsg->IGCMaxMana;
}

void GCLevelUpPointAdd(PMSG_LVPOINTADDRESULT *lpMsg)
{
	if (lpMsg->ResultType == 0x12)
	{
		maxHPDisplay = lpMsg->IGCMaxLifeAndMana;
	}

	else if (lpMsg->ResultType == 0x13)
	{
		maxmanaDisplay = lpMsg->IGCMaxLifeAndMana;
	}
}

void GCIGCStatsAdd(PMSG_ADDSTATS* lpMsg)
{
	WORD stats = lpMsg->Points;

	switch (lpMsg->statstype)
	{
	case 0:
		*(WORD*)((*(DWORD*)IGC_STAT) + 0x118) = stats; // 1.04R
		break;
	case 1:
		*(WORD*)((*(DWORD*)IGC_STAT) + 0x11A) = stats; // 1.04R
		break;
	case 2:
		*(WORD*)((*(DWORD*)IGC_STAT) + 0x11C) = stats; // 1.04R
		maxHPDisplay = lpMsg->MaxLifeAndMana;
		break;
	case 3:
		*(WORD*)((*(DWORD*)IGC_STAT) + 0x11E) = stats; // 1.04R
		maxmanaDisplay = lpMsg->MaxLifeAndMana;
		break;
	case 4:
		*(WORD*)((*(DWORD*)IGC_STAT) + 0x120) = stats; // 1.04R
		break;

	}

	*(WORD*)((*(DWORD*)IGC_STAT) + 0x174) = lpMsg->LUP; // 1.04R
}

void __fastcall GCOffTradeReq(BYTE* aRecv)
{
	g_UserQuit = TRUE;
	closesocket(*g_MuSocket);
}

void GCSetAgilityFix(PMSG_SETAGILITYBUG* lpMsg)
{
	MemAssign(0x00566F71 + 6, lpMsg->value); // 1.04R
}

void GCSetCashItemMoveEnable()
{
	return;

	BYTE btWareEnable[6] = { 0xE9, 0xDA, 0x00, 0x00, 0x00, 0x90 };
	BYTE btTradeEnable[6] = { 0xE9, 0x94, 0x00, 0x00, 0x00, 0x90 };

	DWORD OldProtect;
	VirtualProtect((LPVOID)0x009B5648, 6, PAGE_EXECUTE_READWRITE, &OldProtect);// addr wrong, not s9 or 1.04R
	memcpy((LPVOID)0x009B5648, btWareEnable, sizeof(btWareEnable));

	VirtualProtect((LPVOID)0x009B4F86, 6, PAGE_EXECUTE_READWRITE, &OldProtect);// addr wrong, not s9 or 1.04R
	memcpy((LPVOID)0x009B4F86, btTradeEnable, sizeof(btTradeEnable));

}
void GCMasterLevelUpSend(PMSG_MASTER_LEVEL_UP_SEND *lpMsg)
{
	maxHPDisplay = lpMsg->IGCMaxLife;
	maxmanaDisplay = lpMsg->IGCMaxMana;
}
void GCFriendRoomCreate(PMSG_FRIEND_ROOMCREATE_RESULT* lpMsg)
{
	strcpy(ChatServerIp, (const char*)lpMsg->ServerIp);
}

void GCEnableSiegeSkills()
{
	MemAssign(0x005A6226 + 1, (BYTE)0x4A); // 1.04R
	MemAssign(0x005A622C + 1, (BYTE)0x44); // 1.04R
	MemAssign(0x005A6232 + 1, (BYTE)0x3E); // 1.04R
	MemAssign(0x005A6238 + 1, (BYTE)0x38); // 1.04R
	MemAssign(0x005A623E + 1, (BYTE)0x32); // 1.04R
	MemAssign(0x005A6244 + 1, (BYTE)0x2C); // 1.04R
	MemAssign(0x005A624D, (BYTE)0xEB); // 1.04R
}

void CGSendFilesCrc()
{
	PMSG_SEND_CRC msg;
	msg.h.c = 0xC1;
	msg.h.headcode = 0xFA;
	msg.h.size = sizeof(PMSG_SEND_CRC);
	msg.m_SubCode = 0x0D;


	char szPath[MAX_PATH];
	GetModuleFileName(GetCurrentModule(), szPath, MAX_PATH);
	msg.m_Crc[1] = g_EncDec->ComputeFileCRC("IGC.dll");

	GetModuleFileName(NULL, szPath, MAX_PATH);
	msg.m_Crc[0] = g_EncDec->ComputeFileCRC(szPath);

	//msg.m_Crc[0] = g_EncDec->ComputeFileCRC("Main.exe");
	//msg.m_Crc[1] = g_EncDec->ComputeFileCRC("IGC.dll");
	msg.m_Crc[2] = g_EncDec->ComputeFileCRC("Data\\Player\\Player.bmd");
	msg.m_Crc[4] = g_EncDec->ComputeFileCRC("Data\\Local\\Eng\\skill_eng.bmd");
	//msg.m_Crc[3] = g_ProcessScan->GetCrc();
	DataSend(&msg, sizeof(PMSG_SEND_CRC));
}

//inline void GCSendAntihackBreach(DWORD dwErrorCode)
void GCSendAntihackBreach(DWORD dwErrorCode)
{
	
	PMSG_ANTIHACK_BREACH pMsg;

	pMsg.h.c = 0xC1;
	pMsg.h.size = sizeof(pMsg);
	pMsg.h.headcode = 0xFA;
	pMsg.subcode = 0x0A;
	pMsg.dwErrorCode = dwErrorCode;

	DataSend(&pMsg, pMsg.h.size);
	
}

void GCSetCharSet(WORD CharSet)
{
	BYTE CHAT_FIX[] = { 0xB8, LOBYTE(CharSet), HIBYTE(CharSet), 0x00, 0x00, 0x90 };
	MemCpy(0x00A18E07, CHAT_FIX, sizeof(CHAT_FIX));
	BYTE CHAT_ASIANFIX[] = { LOBYTE(CharSet), HIBYTE(CharSet) };
	MemCpy(0x00DF46BB, CHAT_ASIANFIX, sizeof(CHAT_ASIANFIX));
	MemCpy(0x00A18FCE, CHAT_ASIANFIX, sizeof(CHAT_ASIANFIX)); // s9
	MemCpy(0x00A1912D, CHAT_ASIANFIX, sizeof(CHAT_ASIANFIX)); // s9
	MemAssign(0x00A19128, (DWORD)(0x100000000 - CharSet)); // s9
	MemCpy(0x00A18D3E, CHAT_ASIANFIX, sizeof(CHAT_ASIANFIX)); // s9
	MemAssign(0x0051B6AA, (BYTE)0x6D); // s9
	MemAssign(CHARSET_ADDRESS, CharSet);
	MemSet(0x0045277A, 0x74, 1); // s9
	MemSet(0x00A18D0A, 0xEB, 1); // s9
	MemSet(0x00D5AE3E, 0xEB, 1); // s9
	MemSet(0x00438528, 0xEB, 1); // S9 -- JE -> JMP, Opens support of Asian Symbols in character creation
	//MemSet(0x00452685, 0xEB, 1); // S9 -- IME and Virtual inputs support
	//MemSet(0x00452686, 0x48, 1); // S9 -- IME and Virtual inputs support
	MemAssign(0x00452685, (WORD)0x48EB); // S9 -- IME and Virtual inputs support
	MemSet(0x00452744, 0xEB, 1); // S9 -- Asian support for range of input boxes
	MemCpy(0x004B391C, CHAT_ASIANFIX, 2);
}

void GCRareItemSetPrices(PMSG_RAREITEMPIRCESEND * lpMsg)
{
	ModifyValueInt(JEWEL_OF_CHAOS_PRICE, lpMsg->iJewelOfChaosPrice);
	ModifyValueInt(JEWEL_OF_BLESS_PRICE, lpMsg->iJewelOfBlessPrice);
	ModifyValueInt(JEWEL_OF_SOUL_PRICE, lpMsg->iJewelOfSoulPrice);
	ModifyValueInt(JEWEL_OF_LIFE_PRICE, lpMsg->iJewelOfLifePrice);
	ModifyValueInt(JEWEL_OF_CREATION_PRICE, lpMsg->iJewelOfCreationPrice);
	ModifyValueInt(CREST_OF_MONARCH_PRICE, lpMsg->iCrestOfMonarchPrice);
	ModifyValueInt(LOCH_FEATHER_PRICE, lpMsg->iLochFeatherPrice);
	ModifyValueInt(JEWEL_OF_GUARDIAN_PRICE, lpMsg->iJewelOfGuardianPrice);
	ModifyValueInt(WERERABBIT_EGG_PRICE, lpMsg->iWereRabbitEggPrice);
	MemSet(ITEM_SELL_DIVISOR, lpMsg->btItemSellPriceDivisor,1 );
}

void GCDisableReconnect()
{
	g_UserQuit = 1;
	closesocket(*g_MuSocket);
}

void GCDisableReconnectSystem()
{
	g_UserQuit = 1;
}

void GCAnsMuBotRecvUse(PMSG_MUBOT_USE * aRecv)
{
	g_MuBotEnabled = aRecv->Status;
}

void CGEnableMUBot(int state)
{
	PMSG_MUBOT_SWITCH pMsg;

	pMsg.h.set((LPBYTE)&pMsg,0xBF,0x51,sizeof(pMsg));

	pMsg.Status = state;

	DataSend(&pMsg,sizeof(pMsg));
}

void GCAlterPShopVault(int type)
{
	switch(type)
	{
	case 0:
		pShopPriceTag = "Zen: %s";
		break;
	case 1:
		pShopPriceTag = "WCoinC: %s";
		break;
	case 2:
		pShopPriceTag = "WCoinP: %s";
		break;
	case 3:
		pShopPriceTag = "Goblin Point: %s";
		break;
	}
}

/* You are not allowed to drop this expensive item
005AFF17   E9 92000000      JMP main.005AFFAE
005AFF1C   90               NOP
*/

/* Reinforced item can't be sold.
00A1426E   E9 97000000      JMP main.00A1430A
00A14273   90               NOP
*/

void GCDropSellMod()
{
	BYTE btDropItem[6] = { 0xE9, 0xFB, 0x00, 0x00, 0x00, 0x90 };
	MemCpy(0x00637209, btDropItem, sizeof(btDropItem)); // 1.04R
	BYTE btReinforcedSell[6] = {  0xE9, 0x97, 0x00, 0x00, 0x00, 0x90 };
	MemCpy(0x00AF245A, btReinforcedSell, sizeof(btReinforcedSell)); // 1.04R
}

void GCCustomPost(PMSG_POST_DATA *lpMsg)
{
	if (lpMsg == NULL)
		return;

	int tmpClass = pMessageBoxThis();

	if (tmpClass == 0)
		return;

	char* pText = &lpMsg->szMessage[1];

	char* tmpWrongSymbolCheck = strstr(pText, "<");

	if (tmpWrongSymbolCheck != NULL) 
	{
		return;
	}

	char szName[11];
	szName[10] = '\0';
	memcpy(szName, lpMsg->szName, 10);

	char szMsg[512];
	std::string szPostFormat;

	if (lpMsg->h.subcode == 0x10)
	{
		szPostFormat = "<FONT color='#FF%02X%02X%02X'>%s %s: %s</FONT><BR>";
		sprintf(szMsg, szPostFormat.c_str(),
			lpMsg->btColorRGB[0],
			lpMsg->btColorRGB[1],
			lpMsg->btColorRGB[2],
			lpMsg->szTag,
			szName,
			pText);
	}

	else if (lpMsg->h.subcode == 0x11)
	{
		szPostFormat = "<FONT color='#FF%02X%02X%02X'>[%s]%s %s: %s</FONT><BR>";
		sprintf(szMsg, szPostFormat.c_str(),
			lpMsg->btColorRGB[0],
			lpMsg->btColorRGB[1],
			lpMsg->btColorRGB[2],
			lpMsg->szServerName,
			lpMsg->szTag,
			szName,
			pText);
	}

	 DWORD dwNum = MultiByteToWideChar(g_ServerInfo->GetCharset(), 0, szMsg, -1, NULL, 0);

	if (!dwNum)
		return;

	wchar_t *pwText;
	pwText = new wchar_t[dwNum];
	MultiByteToWideChar(g_ServerInfo->GetCharset(), 0, szMsg, -1, pwText, dwNum);

	dwNum = WideCharToMultiByte(CP_UTF8, 0, pwText, -1, NULL, 0, NULL, NULL);
	char *psText;
	psText = new char[dwNum];
	WideCharToMultiByte(CP_UTF8, 0, pwText, -1, psText, dwNum, NULL, NULL);

	ScaleFormAction(*(DWORD*)(tmpClass + 0x124), "SetChatLogText", "%s %d", psText, 2);

	delete[] pwText;
	delete[] psText; 
}

void GCSetChatColors(PMSG_SET_CHAT_COLOR * lpMsg)
{
	MemCpy(0x00AA0DC8 + 3, lpMsg->btWhisperMsg, sizeof(lpMsg->btWhisperMsg)); // 1.04R
	MemCpy(0x00AA0DF2 + 3, lpMsg->btInfoMsg, sizeof(lpMsg->btInfoMsg)); // 1.04R
	MemCpy(0x00AA0E01 + 3, lpMsg->btErrorMsg, sizeof(lpMsg->btErrorMsg)); // 1.04R
	MemCpy(0x00AA0E10 + 3, lpMsg->btChatMsg, sizeof(lpMsg->btChatMsg)); // 1.04R
	MemCpy(0x00AA0E1F + 3, lpMsg->btPartyMsg, sizeof(lpMsg->btPartyMsg)); // 1.04R
	MemCpy(0x00AA0E2E + 3, lpMsg->btGuildMsg, sizeof(lpMsg->btGuildMsg)); // 1.04R
	MemCpy(0x00AA0E3D + 3, lpMsg->btAllianceMsg, sizeof(lpMsg->btAllianceMsg)); // 1.04R
	MemCpy(0x00AA0E4C + 3, lpMsg->btGMChatMsg, sizeof(lpMsg->btGMChatMsg)); // 1.04R
	MemCpy(0x00AA0E5B + 3, lpMsg->btGensMsg, sizeof(lpMsg->btGensMsg)); // 1.04R
}


void handlePartyHPMP(PMSG_PARTY_COUNT* lpMsg) {
	int size = lpMsg->Count;
	if (lpMsg->h.size <= sizeof(PMSG_PARTY_COUNT) + size * sizeof(PMSG_PARTYLIFEALL))
		return; // there is no HP/MP value extension

	// range msg
	for (int i = 0; i < size; i++) {
		PMSG_PARTYHPMP* partyMember = (PMSG_PARTYHPMP*)((uintptr_t)lpMsg + sizeof(PMSG_PARTY_COUNT) + size * sizeof(PMSG_PARTYLIFEALL) + i * sizeof(PMSG_PARTYHPMP));
		char name[11] = { 0 };
		strncpy(name, partyMember->Name, 10);
		
		// partyFrame = f00A49798game().m184partyFrame
		void* partyFrame = nullptr;
		__asm {
			mov eax, 0x00A49798; // 1.04R f00A49798game()
			call eax;
			mov eax, dword ptr ds:[eax+0x184]; // 1.04R mPartyFrame
			mov partyFrame, eax;
		}

		// bRet = partyFrame.f00A83122validateName(name)
		bool bRet = false;
		__asm {
			mov ecx, partyFrame;
			lea eax, name;
			push eax;
			mov eax, 0x00A83122; // 1.04R partyFrame.f00A83122validateName(char* name) __stdcall
			call eax;
			mov bRet, al;
		}
		if (bRet == false)
			continue;

		// range PartyFrame
		int partyFrameSize = *(int*)((uintptr_t)partyFrame + 0x74);
		for (int index = 0; index < partyFrameSize; index++) {
			void* partyFrameMember = (int*)((uintptr_t)partyFrame + 0x78 + index * 0x54);
			// bRet = partyFrame.f00A83171matchIndexByName(index, name)
			__asm {
				mov ecx, partyFrame;
				lea eax, name;
				push eax;
				push index;
				mov eax, 0x00A83171; // 1.04R partyFrame.f00A83171matchIndexByName(int index, char* name) __stdcall
				call eax;
				mov bRet, al;
			}
			if (bRet) {
				typedef int(*ScaleformGFX)(void*, char*, char*, ...);
				ScaleformGFX scaleformGFX = (ScaleformGFX)0x00A3A4F2; // 1.04R
				
				int HP = partyMember->HP;
				int HPMax = partyMember->HPMax;
				int* partyFrameMemberHP = (int*)((uintptr_t)partyFrameMember + 0x14);
				int* partyFrameMemberHPMax = (int*)((uintptr_t)partyFrameMember + 0x18);
				if (*partyFrameMemberHP != HP || *partyFrameMemberHPMax != HPMax) {
					*partyFrameMemberHP = HP;
					*partyFrameMemberHPMax = HPMax;
					scaleformGFX(partyFrame, "SetPartyHPValue", "%d, %d, %d", index, HP, HPMax);
				}

				int MP = partyMember->MP;
				int MPMax = partyMember->MPMax;
				int* partyFrameMemberMP = (int*)((uintptr_t)partyFrameMember + 0x1C);
				int* partyFrameMemberMPMax = (int*)((uintptr_t)partyFrameMember + 0x20);
				if (*partyFrameMemberMP != MP || *partyFrameMemberMPMax != MPMax) {
					*partyFrameMemberMP = MP;
					*partyFrameMemberMPMax = MPMax;
					scaleformGFX(partyFrame, "SetPartyMPValue", "%d, %d, %d", index, MP, MPMax);
				}
				break;
			}
		}
	}
}