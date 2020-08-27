#include "stdafx.h"
#include "EventEntryLevel.h"
#include "HookManager.h"
// warning
// it's private custom
// don't add it to public release

void LoadBCEntryLevel()
{
	char szTemp[128];

	FILE * hFile = fopen("./Data/Local/EventEntry.bmd", "r"); // check if file exists

	if (hFile == NULL)
	{
		return;
	}

	fclose(hFile);

	for (int i = 0; i < 7; i++)
	{
		sprintf(szTemp, "BC%dLevelMin_Normal", i + 1);
		g_BloodCastleEnterLevel_Normal[i][0] = GetPrivateProfileInt("BloodCastle", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "BC%dLevelMax_Normal", i + 1);
		g_BloodCastleEnterLevel_Normal[i][1] = GetPrivateProfileInt("BloodCastle", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "BC%dLevelMin_MGDLRF", i + 1);
		g_BloodCastleEnterLevel_Magumsa[i][0] = GetPrivateProfileInt("BloodCastle", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "BC%dLevelMax_MGDLRF", i + 1);
		g_BloodCastleEnterLevel_Magumsa[i][1] = GetPrivateProfileInt("BloodCastle", szTemp, 0, "./Data/Local/EventEntry.bmd");
	}

	HookManager.MakeJmpHook(0x00AD37C8, 6, HookBCLevel); // 1.04R
}

void LoadDSEntryLevel()
{
	char szTemp[128];

	FILE * hFile = fopen("./Data/Local/EventEntry.bmd", "r"); // check if file exists

	if (hFile == NULL)
	{
		return;
	}

	fclose(hFile);

	for (int i = 0; i < 6; i++)
	{
		sprintf(szTemp, "DS%dLevelMin_Normal", i + 1);
		g_DevilSquareEnterLevel_Normal[i][0] = GetPrivateProfileInt("DevilSquare", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "DS%dLevelMax_Normal", i + 1);
		g_DevilSquareEnterLevel_Normal[i][1] = GetPrivateProfileInt("DevilSquare", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "DS%dLevelMin_MGDLRF", i + 1);
		g_DevilSquareEnterLevel_Magumsa[i][0] = GetPrivateProfileInt("DevilSquare", szTemp, 0, "./Data/Local/EventEntry.bmd");

		sprintf(szTemp, "DS%dLevelMax_MGDLRF", i + 1);
		g_DevilSquareEnterLevel_Magumsa[i][1] = GetPrivateProfileInt("DevilSquare", szTemp, 0, "./Data/Local/EventEntry.bmd");
	}

	ModifyDSLevel();
}

void __declspec(naked) HookBCLevel()
{
	BC_LEVEL_INFO_MAIN m_BCInfo;

	for (int i = 0; i < 7; i++)
	{
		m_BCInfo.iBloodCastleLevel = i + 1;
		m_BCInfo.iMinLevel = g_BloodCastleEnterLevel_Normal[i][0];
		m_BCInfo.iMaxLevel = g_BloodCastleEnterLevel_Normal[i][1];
		m_BCInfo.bIsMagumsa = false;

		__asm
		{
			LEA EAX, m_BCInfo;
			PUSH EAX;
			MOV ECX, DWORD PTR SS : [EBP - 0x74];
			ADD ECX, 0x74;
			MOV EDX, 0x00AD40B5; // 1.04R
			CALL EDX;
		}
	}

	for (int i = 0; i < 7; i++)
	{
		m_BCInfo.iBloodCastleLevel = i + 1;
		m_BCInfo.iMinLevel = g_BloodCastleEnterLevel_Magumsa[i][0];
		m_BCInfo.iMaxLevel = g_BloodCastleEnterLevel_Magumsa[i][1];
		m_BCInfo.bIsMagumsa = true;

		__asm
		{
			LEA EAX, m_BCInfo;
			PUSH EAX;
			MOV ECX, DWORD PTR SS : [EBP - 0x74];
			ADD ECX, 0x74;
			MOV EDX, 0x00AD40B5; // 1.04R
			CALL EDX;
		}
	}

	__asm
	{
		MOV EDX, 0x00AD38A6; // 1.04R
		JMP EDX;
	}
}

void ModifyDSLevel()
{
	MemAssign(0x0091D058 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[0][0]); // 1.04R
	MemAssign(0x0091D065 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[0][1]); // 1.04R
	MemAssign(0x0091D072 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[1][0]); // 1.04R
	MemAssign(0x0091D07F + 6, (DWORD)g_DevilSquareEnterLevel_Normal[1][1]); // 1.04R
	MemAssign(0x0091D08C + 6, (DWORD)g_DevilSquareEnterLevel_Normal[2][0]); // 1.04R
	MemAssign(0x0091D099 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[2][1]); // 1.04R
	MemAssign(0x0091D0A6 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[3][0]); // 1.04R
	MemAssign(0x0091D0B3 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[3][1]); // 1.04R
	MemAssign(0x0091D0C0 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[4][0]); // 1.04R
	MemAssign(0x0091D0CD + 6, (DWORD)g_DevilSquareEnterLevel_Normal[4][1]); // 1.04R
	MemAssign(0x0091D0DA + 6, (DWORD)g_DevilSquareEnterLevel_Normal[5][0]); // 1.04R
	MemAssign(0x0091D0E7 + 6, (DWORD)g_DevilSquareEnterLevel_Normal[5][1]); // 1.04R

	MemAssign(0x0091D108 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[0][0]); // 1.04R
	MemAssign(0x0091D115 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[0][1]); // 1.04R
	MemAssign(0x0091D122 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[1][0]); // 1.04R
	MemAssign(0x0091D12F + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[1][1]); // 1.04R
	MemAssign(0x0091D13C + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[2][0]); // 1.04R
	MemAssign(0x0091D149 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[2][1]); // 1.04R
	MemAssign(0x0091D156 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[3][0]); // 1.04R
	MemAssign(0x0091D163 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[3][1]); // 1.04R
	MemAssign(0x0091D170 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[4][0]); // 1.04R
	MemAssign(0x0091D17D + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[4][1]); // 1.04R
	MemAssign(0x0091D18A + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[5][0]); // 1.04R
	MemAssign(0x0091D197 + 6, (DWORD)g_DevilSquareEnterLevel_Magumsa[5][1]); // 1.04R
}