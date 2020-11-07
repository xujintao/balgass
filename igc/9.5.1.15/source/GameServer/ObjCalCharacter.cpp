#include "stdafx.h"
#include "ObjCalCharacter.h"
#include "Gamemain.h"
#include "DarkSpirit.h"
#include "SetItemOption.h"
#include "TLog.h"
#include "ItemSocketOptionSystem.h"
#include "BuffEffect.h"
#include "BuffEffectSlot.h"
#include "PentagramSystem.h"
#include "configread.h"
#include "MasterLevelSkillTreeSystem.h"
#include "MuunSystem.h"
#include "ItemOptionTypeMng.h"
#include "StatSpecialize.h"

// GS-N 0.99.60T 0x004A4740
//	GS-N	1.00.18	JPN	0x004C2650	-	Completed
// ADDONS by Dudi (FROM 1.00.90 CS + S5 ITEMS)

CObjCalCharacter gObjCalCharacter;

CObjCalCharacter::CObjCalCharacter() : m_Lua(true)
{

}

CObjCalCharacter::~CObjCalCharacter()
{

}

void CObjCalCharacter::Init()
{
	this->m_Lua.DoFile(g_ConfigRead.GetPath("\\Scripts\\Character\\CalcCharacter.lua"));
}

void CObjCalCharacter::CalcCharacter(int aIndex)
{
	if ( gObj[aIndex].Type != OBJ_USER )
	{
		return;
	}

	LPOBJ lpObj = &gObj[aIndex];
	CItem* Right = &lpObj->pInventory[0];
	CItem* Left  = &lpObj->pInventory[1];
	CItem* Helmet = &lpObj->pInventory[2];
	CItem* Armor = &lpObj->pInventory[3];
	CItem* Pants = &lpObj->pInventory[4];
	CItem* Gloves = &lpObj->pInventory[5];
	CItem* Boots = &lpObj->pInventory[6];
	CItem* Amulet = &lpObj->pInventory[9];
	CItem* Helper = &lpObj->pInventory[8];
	CItem* Wing = &lpObj->pInventory[7];
	CItem* RightRing = &lpObj->pInventory[10];
	CItem* LeftRing = &lpObj->pInventory[11];
	CItem* Pentagram = &lpObj->pInventory[236];

	lpObj->HaveWeaponInHand = true;
	if ( Right->IsItem() == FALSE && Left->IsItem() == FALSE )
	{
		lpObj->HaveWeaponInHand = false;
	}
	else if ( Left->IsItem() == FALSE && Right->m_Type == ITEMGET(4,15) )	// Arrow
	{
		lpObj->HaveWeaponInHand = false;
	}
	else if ( Right->IsItem() == FALSE )
	{
		int iType = Left->m_Type / MAX_SUBTYPE_ITEMS;
		if ( Left->m_Type == ITEMGET(4,7) ) // Bolt
		{
			lpObj->HaveWeaponInHand = false;
		}
		else if ( iType == 6 ) // Shield
		{
			lpObj->HaveWeaponInHand = false;
		}
	}

	for (int iItemIndex=0; iItemIndex<MAX_PLAYER_EQUIPMENT; iItemIndex++)
	{
		lpObj->pInventory[iItemIndex].m_IsValidItem = false;
		if ( lpObj->pInventory[iItemIndex].IsItem() != FALSE )
		{
			lpObj->pInventory[iItemIndex].m_IsValidItem = true;
		}
	}
	lpObj->pInventory[236].m_IsValidItem = false;
	if ( lpObj->pInventory[236].IsItem() != FALSE )
	{
		lpObj->pInventory[236].m_IsValidItem = true;
	}

	// We have to clear all previous buff
	g_BuffEffect.ClearPrevEffectAll(lpObj);
	lpObj->AddLife = 0;
	lpObj->AddMana = 0;
	lpObj->AddBP = 0;
	lpObj->iAddShield = 0;
	lpObj->MonsterDieGetMoney = 0;
	lpObj->MonsterDieGetLife = 0;
	lpObj->MonsterDieGetMana = 0;
	lpObj->AutoHPRecovery = 0;
	lpObj->DamageReflect = 0;
	lpObj->DamageMinus = 0;
	lpObj->AddStrength = 0;
	lpObj->AddDexterity = 0;
	lpObj->AddVitality = 0;
	lpObj->AddEnergy = 0;
	lpObj->AddLeadership = 0;
	lpObj->m_CriticalDamageSuccessRate = 0;
	lpObj->m_ExcellentDamageSuccessRate = 0;
	lpObj->m_PlayerData->SetOpAddDamage = 0;
	lpObj->m_PlayerData->SetOpAddAttackDamage = 0;
	lpObj->m_PlayerData->SetOpAddSkillAttack = 0;
	lpObj->m_PlayerData->SetOpAddMagicPower = 0;
	lpObj->m_PlayerData->SetOpAddMinAttackDamage = 0;
	lpObj->m_PlayerData->SetOpAddMaxAttackDamage = 0;
	lpObj->m_PlayerData->SetOpIncAGValue = 0;
	lpObj->m_PlayerData->SetOpAddCriticalDamageSuccessRate = 0;
	lpObj->m_PlayerData->SetOpAddCriticalDamage = 0;
	lpObj->m_PlayerData->SetOpAddExcellentDamageSuccessRate = 0;
	lpObj->m_PlayerData->SetOpAddExcellentDamage = 0;
	lpObj->m_PlayerData->SetOpAddDefence = 0;
	lpObj->m_PlayerData->SetOpAddDefenceRate = 0;
	lpObj->m_PlayerData->SetOpIgnoreDefense = 0;
	lpObj->m_PlayerData->SetOpDoubleDamage = 0;
	lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage = 0;
	lpObj->m_PlayerData->SetOpImproveSuccessAttackRate = 0;
	lpObj->m_PlayerData->SetOpReflectionDamage = 0;
	lpObj->m_PlayerData->SetOpImproveSheldDefence = 0;
	lpObj->m_PlayerData->SetOpDecreaseAG = 0;
	lpObj->m_PlayerData->SetOpImproveItemDropRate = 0;
	lpObj->m_PlayerData->IsFullSetItem = false;
	lpObj->m_PlayerData->m_WingExcOption.Clear();
	memset(lpObj->m_AddResistance, 0, sizeof(lpObj->m_AddResistance));
	memset(&lpObj->m_PlayerData->m_PentagramOptions, 0, sizeof(lpObj->m_PlayerData->m_PentagramOptions));
	g_StatSpec.ClearUserOption(lpObj);

	// apply previous stat buff
	g_BuffEffect.ApplyPrevEffectStat(lpObj);

	// apply stat of master dimemsions stat
	lpObj->AddEnergy += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncEnergyStat;
	lpObj->AddVitality += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncVitalStat;
	EnterCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
	lpObj->m_PlayerData->AgilityCheckDelay = GetTickCount();
	lpObj->AddDexterity += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncDexStat;
	LeaveCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
	lpObj->AddStrength += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncPowerStat;
	lpObj->AddLeadership += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncLeadershipStat;

	// apply stat of item contributed
	this->CalcSetItemStat(lpObj);
	if (Wing->m_IsValidItem)
	{
		Wing->PlusSpecial(&lpObj->AddLeadership, 105);
	}

	// stat(base+master+buff+item)
	int Strength = lpObj->m_PlayerData->Strength + lpObj->AddStrength;
	int Dexterity = lpObj->m_PlayerData->Dexterity + lpObj->AddDexterity;
	int Vitality = lpObj->m_PlayerData->Vitality + lpObj->AddVitality;
	int Energy = lpObj->m_PlayerData->Energy + lpObj->AddEnergy;
	int Leadership = lpObj->Leadership + lpObj->AddLeadership;

	// damage base
	lpObj->m_CurseDamageMin = 0;
	lpObj->m_CurseDamageMax = 0;
	switch (lpObj->Class)
	{
		case CLASS_ELF:
			if ((Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_CROSSBOW) || (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_BOW)) {
				this->m_Lua.Generic_Call("ElfWithBowDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
					&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			}
			else {
				this->m_Lua.Generic_Call("ElfWithoutBowDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
					&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			}
			this->m_Lua.Generic_Call("ElfMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_KNIGHT:
			this->m_Lua.Generic_Call("KnightDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("KnightMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_MAGUMSA:
			this->m_Lua.Generic_Call("GladiatorDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("GladiatorMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_DARKLORD:
			this->m_Lua.Generic_Call("LordDamageCalc", "iiiii>iiii", Strength, Dexterity, Vitality, Energy, Leadership,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("LordMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_RAGEFIGHTER:
			this->m_Lua.Generic_Call("RageFighterDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("RageFighterMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_SUMMONER:
			this->m_Lua.Generic_Call("SummonerDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("SummonerMagicDamageCalc", "i>iiii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax, &lpObj->m_CurseDamageMin, &lpObj->m_CurseDamageMax);
			break;
		case CLASS_WIZARD:
			this->m_Lua.Generic_Call("WizardDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("WizardMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
			break;
		case CLASS_GROWLANCER:
			this->m_Lua.Generic_Call("GrowLancerDamageCalc", "iiii>iiii", Strength, Dexterity, Vitality, Energy,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight, &lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			this->m_Lua.Generic_Call("GrowLancerMagicDamageCalc", "i>ii", Energy,
				&lpObj->m_MagicDamageMin, &lpObj->m_MagicDamageMax);
	}

	// calc Stat Specialization: increase attack power
	// http://muonline.webzen.com/guides/219/1976/season-9/season-9-character-renewal
	// - (Bonus stat generated by equipping items like weapon, armor, wing, or master skill,
	// etc is not applied to specialization calculation.)
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_ATTACK_POWER);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_MAGIC_DAMAGE);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_CURSE_DAMAGE);

	// calc defense
	this->m_Lua.Generic_Call("CalcDefense", "ii>i", lpObj->Class, Dexterity, &lpObj->m_Defense);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_DEFENSE);

	// attack success rate
	this->m_Lua.Generic_Call("CalcAttackSuccessRate_PvM", "iiiii>i", lpObj->Class, Strength, Dexterity, Leadership, (lpObj->Level + lpObj->m_PlayerData->MasterLevel),
		&lpObj->m_AttackRate);
	this->m_Lua.Generic_Call("CalcDefenseSuccessRate_PvM", "ii>i", lpObj->Class, Dexterity,
		&lpObj->m_DefenseRate);
	this->m_Lua.Generic_Call("CalcAttackSuccessRate_PvP", "iii>d", lpObj->Class, Dexterity, lpObj->Level + lpObj->m_PlayerData->MasterLevel,
		&lpObj->m_PlayerData->m_AttackRatePvP);
	this->m_Lua.Generic_Call("CalcDefenseSuccessRate_PvP", "iii>d", lpObj->Class, Dexterity, lpObj->Level + lpObj->m_PlayerData->MasterLevel,
		&lpObj->m_PlayerData->m_DefenseRatePvP);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_ATTACK_RATE);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_ATTACK_RATE_PVP);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_DEFENSE_RATE);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_DEFENSE_RATE_PVP);

	// attack/magic speed
	this->m_Lua.Generic_Call("CalcAttackSpeed", "ii>ii", lpObj->Class, Dexterity,
		&lpObj->m_AttackSpeed, &lpObj->m_MagicSpeed);
	g_StatSpec.CalcStatOption(lpObj, STAT_OPTION_INC_ATTACK_SPEED);

	// damage(base+master+item+buff)

	// damage of item weapon and addition
	if (Right->m_IsValidItem)
	{
		lpObj->m_AttackDamageMinRight += Right->m_DamageMin;
		lpObj->m_AttackDamageMaxRight += Right->m_DamageMax;
		Right->PlusSpecial(&lpObj->m_AttackDamageMinRight, 80);
		Right->PlusSpecial(&lpObj->m_AttackDamageMaxRight, 80);
		Right->PlusSpecial(&lpObj->m_MagicDamageMin, 81);
		Right->PlusSpecial(&lpObj->m_MagicDamageMax, 81);
		Right->PlusSpecial(&lpObj->m_CurseDamageMin, 113);
		Right->PlusSpecial(&lpObj->m_CurseDamageMax, 113);
	}
	if (Left->m_IsValidItem)
	{
		lpObj->m_AttackDamageMinLeft += Left->m_DamageMin;
		lpObj->m_AttackDamageMaxLeft += Left->m_DamageMax;
		Left->PlusSpecial(&lpObj->m_AttackDamageMinLeft, 80);
		Left->PlusSpecial(&lpObj->m_AttackDamageMaxLeft, 80);
		Left->PlusSpecial(&lpObj->m_MagicDamageMin, 81);
		Left->PlusSpecial(&lpObj->m_MagicDamageMax, 81);
		Left->PlusSpecial(&lpObj->m_CurseDamageMin, 113);
		Left->PlusSpecial(&lpObj->m_CurseDamageMax, 113);
	}

	// damage of item wing addition
	Wing->PlusSpecial(&lpObj->m_AttackDamageMinRight, 80);
	Wing->PlusSpecial(&lpObj->m_AttackDamageMaxRight, 80);
	Wing->PlusSpecial(&lpObj->m_AttackDamageMinLeft, 80);
	Wing->PlusSpecial(&lpObj->m_AttackDamageMaxLeft, 80);
	Wing->PlusSpecial(&lpObj->m_MagicDamageMin, 81);
	Wing->PlusSpecial(&lpObj->m_MagicDamageMax, 81);
	Wing->PlusSpecial(&lpObj->m_CurseDamageMin, 113);
	Wing->PlusSpecial(&lpObj->m_CurseDamageMax, 113);

	// defense of item contributed
	for (int i=1; i<8; i++) {
		if (lpObj->pInventory[i].m_IsValidItem) {
			lpObj->m_Defense += lpObj->pInventory[i].ItemDefense();
			lpObj->pInventory[i].PlusSpecial(&lpObj->m_Defense, 83);
		}
	}

	if (lpObj->m_btInvenPetPos
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_JewelOfHarmonyOption == TRUE
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_Type == ITEMGET(13,4) // Dark Horse
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_Durability > 0.0f)
	{
		int addDefense = 0;
		this->m_Lua.Generic_Call("CalcDarkHorseDefenseBonus", "ii>i",
			Dexterity, lpObj->pInventory[lpObj->m_btInvenPetPos].m_PetItem_Level, &addDefense);
		//int addDefense = Dexterity / 20 + 5 + lpObj->pInventory[lpObj->m_btInvenPetPos].m_PetItem_Level * 2;
		//addDefense = addDefense * g_ConfigRead.calc.DarkHorseDefenseMultiplier / 100.0f;
		lpObj->m_Defense += addDefense;
	}

	// Pet Unicorn
	if (Helper->m_Type == ITEMGET(13, 106) && Helper->IsPeriodItemExpire() == FALSE)
	{
		lpObj->m_Defense += g_ConfigRead.pet.UnicornAddDefense;
	}

	// defense success rate of item contributed
	if (Left->m_IsValidItem != false)
	{
		lpObj->m_DefenseRate += Left->m_DefenseRate;
		lpObj->pInventory[1].PlusSpecial(&lpObj->m_DefenseRate, 82);
	}

	// defense level>=10 bonus of item the same type contributed
	// defense success rate bonus of item the same type contributed
	// calc
	int sameTypeCount = 0;
	int Level10Count = 0;
	int Level11Count = 0;
	int Level12Count = 0;
	int Level13Count = 0;
	int Level14Count = 0;
	int Level15Count = 0;

	if (lpObj->pInventory[6].m_IsValidItem) {
		int rvalue = lpObj->pInventory[6].m_Type % MAX_SUBTYPE_ITEMS;
		for (int i=2; i<=6; i++) {
			if ((i==2 && lpObj->Class == CLASS_MAGUMSA)
			|| (i==5 && lpObj->Class == CLASS_RAGEFIGHTER)) {
				sameTypeCount++;
				Level10Count++;
				Level11Count++;
				Level12Count++;
				Level13Count++;
				Level14Count++;
				Level15Count++;
			}
			if (lpObj->pInventory[i].m_IsValidItem
			&& lpObj->pInventory[i].m_Type % MAX_SUBTYPE_ITEMS == rvalue) {
				sameTypeCount++;
				if ( lpObj->pInventory[i].m_Level > 14 )
					Level15Count++;
				else if ( lpObj->pInventory[i].m_Level > 13 )
					Level14Count++;
				else if ( lpObj->pInventory[i].m_Level > 12 )
					Level13Count++;
				else if ( lpObj->pInventory[i].m_Level > 11 )
					Level12Count++;
				else if ( lpObj->pInventory[i].m_Level > 10 )
					Level11Count++;
				else if ( lpObj->pInventory[i].m_Level > 9 )
					Level10Count++;
			}
		}
	}

	// apply bonus of item the same type
	if (sameTypeCount == 5) {
		if ( Level15Count == 5 )
			lpObj->m_Defense += lpObj->m_Defense * 30 / 100;
		else if ( Level14Count == 5 || (Level14Count + Level15Count) == 5)
			lpObj->m_Defense += lpObj->m_Defense * 25 / 100;
		else if ( Level13Count == 5 || (Level13Count + Level14Count + Level15Count) == 5)
			lpObj->m_Defense += lpObj->m_Defense * 20 / 100;
		else if ( Level12Count == 5 || (Level12Count + Level13Count + Level14Count + Level15Count) == 5 )
			lpObj->m_Defense += lpObj->m_Defense * 15 / 100;
		else if ( Level11Count == 5 || (Level11Count + Level12Count + Level13Count + Level14Count + Level15Count) == 5 )
			lpObj->m_Defense += lpObj->m_Defense * 10 / 100;
		else if ( Level10Count == 5 || (Level10Count + Level11Count + Level12Count + Level13Count + Level14Count + Level15Count) == 5)
			lpObj->m_Defense += lpObj->m_Defense * 5 / 100;
		lpObj->m_DefenseRate += lpObj->m_DefenseRate / 10;
	}

	// calc master passive skill
	g_MasterLevelSkillTreeSystem.SetItemMLPassiveSkill(lpObj, Right->GetDetailItemType());
	g_MasterLevelSkillTreeSystem.SetItemMLPassiveSkill(lpObj, Left->GetDetailItemType());
	g_MasterLevelSkillTreeSystem.SetWingMLPassiveSkill(lpObj, Wing->m_Type);
	g_MasterLevelSkillTreeSystem.SetPetItemMLPassiveSkill(lpObj, Helper->m_Type);
	this->ApplyMLSkillItemOption(lpObj);

	// speed of item weapon contributed
	bool bRight = false;
	bool bLeft = false;
	if ( Right->m_Type != ITEMGET(4,7) && Right->m_Type != ITEMGET(4,15) && Right->m_Type >= ITEMGET(0,0) && Right->m_Type < ITEMGET(6,0) )
	{
		if ( Right->m_IsValidItem != false )
		{
			bRight = true;
		}
	}
	if ( Left->m_Type != ITEMGET(4,7) && Left->m_Type != ITEMGET(4,15) && Left->m_Type >= ITEMGET(0,0) && Left->m_Type < ITEMGET(6,0) )
	{
		if ( Left->m_IsValidItem != false )
		{
			bLeft = true;
		}
	}
	if (bRight && bLeft)
	{
		if (Left->m_IsValidItem != false && Right->m_IsValidItem != false)
		{
			lpObj->m_AttackSpeed += (Right->m_AttackSpeed + Left->m_AttackSpeed) / 2;
			lpObj->m_MagicSpeed += (Right->m_AttackSpeed + Left->m_AttackSpeed) / 2;
		}
	}
	else if (bRight)
	{
		if (Right->m_IsValidItem != false)
		{
			lpObj->m_AttackSpeed += Right->m_AttackSpeed;
			lpObj->m_MagicSpeed += Right->m_AttackSpeed;
		}
	}
	else if (bLeft)
	{
		if (Left->m_IsValidItem != false)
		{
			lpObj->m_AttackSpeed += Left->m_AttackSpeed;
			lpObj->m_MagicSpeed += Left->m_AttackSpeed;
		}
	}

	// speed of item glove contributed
	if ( Gloves->m_Type != -1 )
	{
		if ( Gloves->m_IsValidItem != false )
		{
			lpObj->m_AttackSpeed += Gloves->m_AttackSpeed;
			lpObj->m_MagicSpeed += Gloves->m_AttackSpeed;
		}
	}

	// speed of item helper contributed
	if ( Helper->m_Type != -1 )
	{
		if ( Helper->m_IsValidItem != false )
		{
			lpObj->m_AttackSpeed += Helper->m_AttackSpeed;
			lpObj->m_MagicSpeed += Helper->m_AttackSpeed;
		}
	}

	// speed of item necklace contributed
	if ( Amulet->m_Type != -1 )
	{
		if ( Amulet->m_IsValidItem != false )
		{
			lpObj->m_AttackSpeed += Amulet->m_AttackSpeed;
			lpObj->m_MagicSpeed += Amulet->m_AttackSpeed;
		}
	}

	// misc item
	if ( Wing->IsItem() == TRUE && Wing->m_IsValidItem == true && Wing->m_Type == ITEMGET(12,268) ) // additional bonus of s9 wings
	{
		lpObj->m_AttackDamageMinRight += 100 + (Wing->m_Level * 5);
		lpObj->m_AttackDamageMaxRight += 100 + (Wing->m_Level * 5);
		lpObj->m_AttackDamageMinLeft += 100 + (Wing->m_Level * 5);
		lpObj->m_AttackDamageMaxLeft += 100 + (Wing->m_Level * 5);
		lpObj->m_MagicDamageMin += 100 + (Wing->m_Level * 5);
		lpObj->m_MagicDamageMax += 100 + (Wing->m_Level * 5);
		lpObj->m_CurseDamageMin += 100 + (Wing->m_Level * 5);
		lpObj->m_CurseDamageMax += 100 + (Wing->m_Level * 5);
	}

	// Wizard Ring effect
	if ((lpObj->pInventory[10].IsItem() == TRUE && lpObj->pInventory[10].m_Type == ITEMGET(13,20) && lpObj->pInventory[10].m_Level == 0 && lpObj->pInventory[10].m_Durability > 0.0f)
	|| (lpObj->pInventory[11].IsItem() == TRUE && lpObj->pInventory[11].m_Type == ITEMGET(13,20) && lpObj->pInventory[11].m_Level == 0 && lpObj->pInventory[11].m_Durability > 0.0f))
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.WizardRingAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.WizardRingAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.WizardRingAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.WizardRingAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.WizardRingAddMagicDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.WizardRingAddMagicDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.WizardRingAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.WizardRingAddAttackSpeed;
	}

	// Lethal Wizard's Ring
	if ((lpObj->pInventory[10].IsItem() == TRUE && lpObj->pInventory[10].m_Type == ITEMGET(13,107) && lpObj->pInventory[10].m_Level == 0 && lpObj->pInventory[10].m_Durability > 0.0f)
	|| (lpObj->pInventory[11].IsItem() == TRUE && lpObj->pInventory[11].m_Type == ITEMGET(13,107) && lpObj->pInventory[11].m_Level == 0 && lpObj->pInventory[11].m_Durability > 0.0f))
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.LethalRingAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.LethalRingAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.LethalRingAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.LethalRingAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.LethalRingAddMagicDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.LethalRingAddMagicDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.LethalRingAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.LethalRingAddAttackSpeed;
	}

	// Pet Skeleton
	if (Helper->m_Type == ITEMGET(13, 123) && Helper->IsPeriodItemExpire() == FALSE)
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.EliteSkeletonPetAddDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.EliteSkeletonPetAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.EliteSkeletonPetAddAttackSpeed;
	}

	if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,10) && LeftRing->m_Level == 5) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,10) && RightRing->m_Level == 5)  )
	{
		lpObj->AddLife = ((int)(lpObj->MaxLife * 20.0f))/100;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,39)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,39)) )
	{
		lpObj->m_Defense += (lpObj->m_Defense * g_ConfigRead.pet.EliteSkeletonRingAddDefense / 100);
		lpObj->AddLife += lpObj->Level+lpObj->m_PlayerData->MasterLevel;
	}
	else if ((LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,41)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,41)))
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.ChristmasRingAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.ChristmasRingAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.ChristmasRingAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.ChristmasRingAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.ChristmasRingAddMagicDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.ChristmasRingAddMagicDamage / 100;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,42)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,42)) )
	{
		lpObj->m_AddResistance[R_ICE] += (char)255;
		lpObj->m_AddResistance[R_POISON] += (char)255;
		lpObj->m_AddResistance[R_LIGHTNING] += (char)255;
		lpObj->m_AddResistance[R_FIRE] += (char)255;
		lpObj->m_AddResistance[R_EARTH] += (char)255;
		lpObj->m_AddResistance[R_WIND] += (char)255;
		lpObj->m_AddResistance[R_WATER] += (char)255;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,76)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,76)) )
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.PandaRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.PandaRingAddDamage;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,77)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,77)) )
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.BrownPandaRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.BrownPandaRingAddDamage;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,78)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,78)) )
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.PinkPandaRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.PinkPandaRingAddDamage;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,122)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,122)) )
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.SkeletonRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.SkeletonRingAddDamage;
	}
	else if ((LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13, 163)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13, 163)))
	{
		lpObj->m_Defense += g_ConfigRead.pet.RobotKnightRingAddDefense;
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.RobotKnightRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.RobotKnightRingAddDamage;
	}
	else if ((LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13, 164)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13, 164)))
	{
		lpObj->m_AttackSpeed += g_ConfigRead.pet.MiniRobotRingAddAttackSpeed;
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.MiniRobotRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.MiniRobotRingAddDamage;
	}
	else if ((LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13, 165)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13, 165)))
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.MageRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.MageRingAddDamage;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,169) && LeftRing->m_Durability > 0.0) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,169) && RightRing->m_Durability > 0.0) )
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.DecorationRingAddDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.DecorationRingAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.DecorationRingAddAttackSpeed;
	}
	else if ( (LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13,170) && LeftRing->m_Durability > 0.0) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13,170) && RightRing->m_Durability > 0.0) )
	{
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.BlessedDecorationRingAddDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.BlessedDecorationRingAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.BlessedDecorationRingAddAttackSpeed;
	}
	else if ((LeftRing->IsItem() == TRUE && LeftRing->m_Type == ITEMGET(13, 268)) || (RightRing->IsItem() == TRUE && RightRing->m_Type == ITEMGET(13, 268)))
	{
		lpObj->m_AttackDamageMinRight += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_AttackDamageMaxRight += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_AttackDamageMinLeft += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_AttackDamageMaxLeft += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_MagicDamageMin += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_MagicDamageMax += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_CurseDamageMin += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_CurseDamageMax += g_ConfigRead.pet.DSFRingAddDamage;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.DSFRingAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.DSFRingAddAttackSpeed;
		lpObj->m_Defense += g_ConfigRead.pet.DSFRingAddDefense;
		lpObj->AddLife += g_ConfigRead.pet.DSFRingAddHP;
	}

	// Golden Fenrir Attribute
	if ( lpObj->pInventory[lpObj->m_btInvenPetPos].m_Type == ITEMGET(13,37)
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_NewOption == 4
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_Durability > 0.0f
	&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_JewelOfHarmonyOption == TRUE )
	{
		lpObj->AddLife += 200;	//Life +200
		lpObj->AddMana += 200;	//Mana +200
		
		/*
		//Attack Power Already Coded on zzzitem.cpp
		lpObj->m_AttackDamageMaxLeft += lpObj->m_AttackDamageMaxLeft * 33 / 100; //Attack Power +33
		lpObj->m_AttackDamageMinLeft += lpObj->m_AttackDamageMinLeft * 33 / 100; //Attack Power +33
		lpObj->m_AttackDamageMaxRight += lpObj->m_AttackDamageMaxRight * 33 / 100; //Attack Power +33
		lpObj->m_AttackDamageMinRight += lpObj->m_AttackDamageMinRight * 33 / 100; //Attack Power +33
		*/

		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * 16 / 100; //Wizardry +16
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * 16 / 100; //Wizardry +16
	}

	// Dinorant
	if (lpObj->pInventory[lpObj->m_btInvenPetPos].m_Type == ITEMGET(13,3)
		&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_Durability > 0.0
		&& lpObj->pInventory[lpObj->m_btInvenPetPos].m_JewelOfHarmonyOption == TRUE)
	{
		lpObj->pInventory[lpObj->m_btInvenPetPos].PlusSpecial(&lpObj->AddBP, 103); // AG +50
		lpObj->pInventory[lpObj->m_btInvenPetPos].PlusSpecial(&lpObj->m_AttackSpeed, 104); // Speed 5
		lpObj->pInventory[lpObj->m_btInvenPetPos].PlusSpecial(&lpObj->m_MagicSpeed, 104); // Speed 5
	}

	if ( lpObj->Type == OBJ_USER && lpObj->m_PlayerData->ISBOT == false )
	{
		gDarkSpirit[lpObj->m_Index- g_ConfigRead.server.GetObjectStartUserIndex()].Set(lpObj->m_Index, &lpObj->pInventory[1]);
	}

	if ( lpObj->Class == CLASS_SUMMONER && Left->IsItem() )
	{
		lpObj->m_CurseSpell = Left->IsCurseSpell();
	}

	int nMuunItemEffectValue = 0;

	if (g_CMuunSystem.GetMuunItemValueOfOptType(lpObj, MUUN_INC_ATTACK_POWER, &nMuunItemEffectValue, 0) == 1)
	{
		lpObj->m_AttackDamageMaxLeft += nMuunItemEffectValue;
		lpObj->m_AttackDamageMinLeft += nMuunItemEffectValue;
		lpObj->m_AttackDamageMaxRight += nMuunItemEffectValue;
		lpObj->m_AttackDamageMinRight += nMuunItemEffectValue;
		lpObj->m_MagicDamageMin += nMuunItemEffectValue;
		lpObj->m_MagicDamageMax += nMuunItemEffectValue;
		lpObj->m_CurseDamageMin += nMuunItemEffectValue;
		lpObj->m_CurseDamageMax += nMuunItemEffectValue;
	}

	if (g_CMuunSystem.GetMuunItemValueOfOptType(lpObj, MUUN_INC_MAX_ATTACK_POWER, &nMuunItemEffectValue, 0) == 1)
	{
		lpObj->m_AttackDamageMaxLeft += nMuunItemEffectValue;
		lpObj->m_AttackDamageMaxRight += nMuunItemEffectValue;
		lpObj->m_MagicDamageMax += nMuunItemEffectValue;
		lpObj->m_CurseDamageMax += nMuunItemEffectValue;
	}

	lpObj->pInventory[10].PlusSpecialPercentEx(&lpObj->AddMana, lpObj->MaxMana, 172);
	lpObj->pInventory[11].PlusSpecialPercentEx(&lpObj->AddMana, lpObj->MaxMana, 172);
	lpObj->pInventory[9].PlusSpecialPercentEx(&lpObj->AddBP, lpObj->MaxBP, 173);

	// resistance
	CItem * rItem[3];
	int comparecount = 0;
	
	rItem[0] = &lpObj->pInventory[10];
	rItem[1] = &lpObj->pInventory[11];
	rItem[2] = &lpObj->pInventory[9];

#define GET_MAX_RESISTANCE(x,y,z) ( ( ( ( (x) > (y) ) ? (x) : (y) ) > (z) ) ? ( ( (x) > (y) ) ? (x) : (y) ) : (z) )	

	lpObj->m_Resistance[1] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[1], rItem[1]->m_Resistance[1], rItem[2]->m_Resistance[1]);
	lpObj->m_Resistance[0] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[0], rItem[1]->m_Resistance[0], rItem[2]->m_Resistance[0]);
	lpObj->m_Resistance[2] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[2], rItem[1]->m_Resistance[2], rItem[2]->m_Resistance[2]);
	lpObj->m_Resistance[3] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[3], rItem[1]->m_Resistance[3], rItem[2]->m_Resistance[3]);
	lpObj->m_Resistance[4] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[4], rItem[1]->m_Resistance[4], rItem[2]->m_Resistance[4]);
	lpObj->m_Resistance[5] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[5], rItem[1]->m_Resistance[5], rItem[2]->m_Resistance[5]);
	lpObj->m_Resistance[6] = GET_MAX_RESISTANCE(rItem[0]->m_Resistance[6], rItem[1]->m_Resistance[6], rItem[2]->m_Resistance[6]);

	lpObj->m_Resistance[R_POISON] += lpObj->m_PlayerData->m_MPSkillOpt.iMpsResistancePoison;
	lpObj->m_Resistance[R_LIGHTNING] += lpObj->m_PlayerData->m_MPSkillOpt.iMpsResistanceThunder;
	lpObj->m_Resistance[R_ICE] += lpObj->m_PlayerData->m_MPSkillOpt.iMpsResistanceIce;

	lpObj->m_PlayerData->m_Resistance_Stun = 0;
	lpObj->m_PlayerData->m_Resistance_Double = 0;
	lpObj->m_PlayerData->m_Resistance_Perfect = 0;
	lpObj->m_PlayerData->m_Resistance_SD = 0;
	lpObj->m_PlayerData->m_Resistance_Excellent = 0;
	lpObj->m_PlayerData->m_Resistance_Critical = 0;

	for (int i = 0; i < 3; i++)
	{
		switch (rItem[i]->m_Type)
		{
		case ITEMGET(13, 171):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_Perfect += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						lpObj->m_PlayerData->m_Resistance_Perfect += 2;
					}
				}
			}
			break;
		case ITEMGET(13, 172):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_SD += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						lpObj->m_PlayerData->m_Resistance_SD += 2;
					}
				}
			}
			break;
		case ITEMGET(13, 173):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_Double += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						lpObj->m_PlayerData->m_Resistance_Double += 2;
					}
				}
			}
			break;
		case ITEMGET(13, 174):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_Stun += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						if (rItem[i]->m_Level)
						{
							lpObj->m_PlayerData->m_Resistance_Stun += 2;
						}
						else
						{
							lpObj->m_PlayerData->m_Resistance_Stun++;
						}
					}
				}
			}
			break;
		case ITEMGET(13, 175):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_Excellent += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						lpObj->m_PlayerData->m_Resistance_Excellent += 2;
					}
				}
			}
			break;
		case ITEMGET(13, 176):
			for (int j = 0; j < 6; j++)
			{
				if (rItem[i]->m_Type == ITEMGET(g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemType, g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].ItemIndex))
				{
					lpObj->m_PlayerData->m_Resistance_Critical += g_ItemOptionTypeMng.m_AccessoryItemOptionValue[j].LevelValue[rItem[i]->m_Level];

					if (rItem[i]->IsExtItem() == TRUE)
					{
						lpObj->m_PlayerData->m_Resistance_Critical += 2;
					}
				}
			}
			break;
		}
	}

	g_kItemSystemFor380.ApplyFor380Option(lpObj);
	g_kJewelOfHarmonySystem.SetApplyStrengthenItem(lpObj);
	g_SocketOptionSystem.SetApplySocketEffect(lpObj);
	gObjInventoryEquipment(lpObj);

	lpObj->AddLife += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxHP;
	lpObj->AddBP += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxBP;
	lpObj->iAddShield += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxSD;
	lpObj->AddMana += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxMana;
	lpObj->AddMana += lpObj->MaxMana * lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxManaRate / 100.0;

	g_ItemOptionTypeMng.CalcExcOptionEffect(lpObj);
	g_ItemOptionTypeMng.CalcWingOptionEffect(lpObj);

	// calc option set item contributed
	this->CalcSetItemOption(lpObj);
	this->SetItemApply(lpObj);

	this->PremiumItemApply(lpObj);

	// calc critical damage success rate of luck item
	for (int i=0; i<8; i++) {
		if (lpObj->pInventory[i].m_IsValidItem) {
			lpObj->pInventory[i].PlusSpecial(&lpObj->m_CriticalDamageSuccessRate, 84);
		}
	}
	lpObj->m_CriticalDamageSuccessRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncCriticalDamageRate;
	lpObj->m_ExcellentDamageSuccessRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncExcellentDamageRate;

	// calc pentagram
	if (lpObj->pInventory[236].IsItem() == TRUE)
	{
		if (this->ValidItem(lpObj, &lpObj->pInventory[236], 236))
		{
			lpObj->pInventory[236].m_IsValidItem = true;
			g_PentagramSystem.CalcPentagramItem(aIndex, &lpObj->pInventory[236]);
		}
		else
		{
			lpObj->pInventory[236].m_IsValidItem = false;
		}
	}
	else
	{
		g_PentagramSystem.ClearPentagramItem(aIndex);
	}

	gObjNextExpCal(lpObj);

	if (lpObj->Class)
		lpObj->m_DetectSpeedHackTime = (gAttackSpeedTimeLimit - lpObj->m_AttackSpeed * gDecTimePerAttackSpeed);
	else
		lpObj->m_DetectSpeedHackTime = (gAttackSpeedTimeLimit - (2 * lpObj->m_MagicSpeed) * gDecTimePerAttackSpeed);

	if (lpObj->m_DetectSpeedHackTime < gMinimumAttackSpeedTime)
		lpObj->m_DetectSpeedHackTime = gMinimumAttackSpeedTime;

	if ( gObjCheckUsedBuffEffect(lpObj, BUFFTYPE_BERSERKER_PRO) == TRUE )
	{
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_CurseDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
		lpObj->m_CurseDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncValueBerserker3;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_BOW)
	{
		if ( Right->m_Type == ITEMGET(4,15) && Right->m_Level == 1 )
		{
			lpObj->m_AttackDamageMinLeft += (WORD)(lpObj->m_AttackDamageMinLeft * 0.03f + 1.0f);
			lpObj->m_AttackDamageMaxLeft += (WORD)(lpObj->m_AttackDamageMaxLeft * 0.03f + 1.0f);
		}
		else if ( Right->m_Type == ITEMGET(4,15) && Right->m_Level == 2 )
		{
			lpObj->m_AttackDamageMinLeft += (WORD)(lpObj->m_AttackDamageMinLeft * 0.05f + 1.0f);
			lpObj->m_AttackDamageMaxLeft += (WORD)(lpObj->m_AttackDamageMaxLeft * 0.05f + 1.0f);
		}
		else if ( Right->m_Type == ITEMGET(4,15) && Right->m_Level == 3 )
		{
			lpObj->m_AttackDamageMinLeft += (WORD)(lpObj->m_AttackDamageMinLeft * 0.07f + 1.0f);
			lpObj->m_AttackDamageMaxLeft += (WORD)(lpObj->m_AttackDamageMaxLeft * 0.07f + 1.0f);
		}
	}
	else if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_CROSSBOW)
	{
		if ( Left->m_Type == ITEMGET(4,7) && Left->m_Level == 1 )
		{
			lpObj->m_AttackDamageMinRight += (WORD)(lpObj->m_AttackDamageMinRight * 0.03f + 1.0f);
			lpObj->m_AttackDamageMaxRight += (WORD)(lpObj->m_AttackDamageMaxRight * 0.03f + 1.0f);
		}
		else if ( Left->m_Type == ITEMGET(4,7) && Left->m_Level == 2 )
		{
			lpObj->m_AttackDamageMinRight += (WORD)(lpObj->m_AttackDamageMinRight * 0.05f + 1.0f);
			lpObj->m_AttackDamageMaxRight += (WORD)(lpObj->m_AttackDamageMaxRight * 0.05f + 1.0f);
		}
		else if ( Left->m_Type == ITEMGET(4,7) && Left->m_Level == 3 )
		{
			lpObj->m_AttackDamageMinRight += (WORD)(lpObj->m_AttackDamageMinRight * 0.07f + 1.0f);
			lpObj->m_AttackDamageMaxRight += (WORD)(lpObj->m_AttackDamageMaxRight * 0.07f + 1.0f);
		}
	}
	
	// convert right/left damage to damage
	switch (lpObj->Class) {
	case CLASS_KNIGHT:
	case CLASS_MAGUMSA:
		if (Right->m_IsValidItem && Left->m_IsValidItem) {
			if (Right->m_Type == Left->m_Type) {
				this->m_Lua.Generic_Call("CalcTwoSameWeaponBonus", "iiii>iiii",
					lpObj->m_AttackDamageMinLeft, lpObj->m_AttackDamageMaxLeft,
					lpObj->m_AttackDamageMinRight, lpObj->m_AttackDamageMaxRight,
					&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight,
					&lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			}
			else {
				this->m_Lua.Generic_Call("CalcTwoDifferentWeaponBonus", "iiii>iiii",
					lpObj->m_AttackDamageMinLeft, lpObj->m_AttackDamageMaxLeft,
					lpObj->m_AttackDamageMinRight, lpObj->m_AttackDamageMaxRight,
					&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight,
					&lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			}
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft;
		}
		else if (Right->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight;
		}
		else if (Left->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinLeft;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxLeft;
		}
		else {
			lpObj->m_AttackDamageMin = (lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft)/2;
			lpObj->m_AttackDamageMax = (lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft)/2;
		}
		break;
	case CLASS_RAGEFIGHTER:
		if (Right->m_IsValidItem && Left->m_IsValidItem) {
			this->m_Lua.Generic_Call("CalcRageFighterTwoWeaponBonus", "iiii>iiii",
				lpObj->m_AttackDamageMinLeft, lpObj->m_AttackDamageMaxLeft,
				lpObj->m_AttackDamageMinRight, lpObj->m_AttackDamageMaxRight,
				&lpObj->m_AttackDamageMinLeft, &lpObj->m_AttackDamageMinRight,
				&lpObj->m_AttackDamageMaxLeft, &lpObj->m_AttackDamageMaxRight);
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft;
		}
		else if (Right->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight;
		}
		else if (Left->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinLeft;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxLeft;
		}
		else {
			lpObj->m_AttackDamageMin = (lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft)/2;
			lpObj->m_AttackDamageMax = (lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft)/2;
		}
		break;
	case CLASS_DARKLORD:
		if (Right->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight;
		}
		else {
			lpObj->m_AttackDamageMin = (lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft)/2;
			lpObj->m_AttackDamageMax = (lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft)/2;
		}
		break;
	case CLASS_ELF:
		if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_CROSSBOW) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight;
		}
		else if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_BOW) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinLeft;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxLeft;
		}
		else {
			lpObj->m_AttackDamageMin = (lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft)/2;
			lpObj->m_AttackDamageMax = (lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft)/2;
		}
		break;
	default: // wizard summoner
		if (Right->m_IsValidItem) {
			lpObj->m_AttackDamageMin = lpObj->m_AttackDamageMinRight;
			lpObj->m_AttackDamageMax = lpObj->m_AttackDamageMaxRight;
		}
		else {
			lpObj->m_AttackDamageMin = (lpObj->m_AttackDamageMinRight + lpObj->m_AttackDamageMinLeft)/2;
			lpObj->m_AttackDamageMax = (lpObj->m_AttackDamageMaxRight + lpObj->m_AttackDamageMaxLeft)/2;
		}
	}

	// Now calculation is done, we have to recovery all previous buff.
	g_BuffEffect.ClearPrevEffectStat(lpObj);
	g_BuffEffect.ApplyPrevEffectAll(lpObj);

	if (gObjSatanSprite(lpObj) == TRUE){
		lpObj->m_AttackDamageMin += lpObj->m_AttackDamageMin * g_ConfigRead.pet.SatanAddDamage / 100;
		lpObj->m_AttackDamageMax += lpObj->m_AttackDamageMax * g_ConfigRead.pet.SatanAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.SatanAddDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.SatanAddDamage / 100;
		lpObj->m_CurseDamageMin += lpObj->m_CurseDamageMin * g_ConfigRead.pet.SatanAddDamage / 100;
		lpObj->m_CurseDamageMax += lpObj->m_CurseDamageMax * g_ConfigRead.pet.SatanAddDamage / 100;
	}
	else if (gObjDemonSprite(lpObj) == TRUE)
	{
		lpObj->m_AttackDamageMin += lpObj->m_AttackDamageMin * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_AttackDamageMax += lpObj->m_AttackDamageMax * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_CurseDamageMin += lpObj->m_CurseDamageMin * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_CurseDamageMax += lpObj->m_CurseDamageMax * g_ConfigRead.pet.DemonAddDamage / 100;
		lpObj->m_AttackSpeed += g_ConfigRead.pet.DemonAddAttackSpeed;
		lpObj->m_MagicSpeed += g_ConfigRead.pet.DemonAddAttackSpeed;
	}
	else if (gObjAngelSprite(lpObj) == TRUE) {
		lpObj->AddLife += g_ConfigRead.pet.AngelAddHP;
	}
	else if (gObjSafeGuardSprite(lpObj) == TRUE)
	{
		lpObj->AddLife += g_ConfigRead.pet.SafeGuardAddHP;
	}
	else if (Helper->m_IsValidItem && Helper->m_Type == ITEMGET(13, 80) && lpObj->pInventory[8].m_Durability > 0.0)
	{
		lpObj->m_Defense += g_ConfigRead.pet.PandaPetAddDefense;
		lpObj->m_MagicDefense += g_ConfigRead.pet.PandaPetAddDefense;
	}

	this->CalcShieldPoint(lpObj);

	if ( lpObj->iShield > ( lpObj->iMaxShield + lpObj->iAddShield ) )
	{
		lpObj->iShield = lpObj->iMaxShield + lpObj->iAddShield ;
		GSProtocol.GCReFillSend(lpObj->m_Index, lpObj->Life, 0xFF, 0, lpObj->iShield);
	}

	int nEffectValue = 0;
	int iTotalMaxShield = lpObj->iMaxShield + lpObj->iAddShield;
	int iTotalMaxMana = lpObj->MaxMana + lpObj->AddMana;
	int iTotalMaxLife = lpObj->MaxLife + lpObj->AddLife;

	if (gObjGetValueOfBuffIndex(lpObj, BUFFTYPE_ACHERON_FIRE, &nEffectValue, 0))
	{
		iTotalMaxShield /= 2.0;
		lpObj->iAddShield -= iTotalMaxShield;

		if ( lpObj->iShield > ( lpObj->iMaxShield + lpObj->iAddShield ) )
		{
			lpObj->iShield = lpObj->iMaxShield + lpObj->iAddShield ;
			GSProtocol.GCReFillSend(lpObj->m_Index, lpObj->Life, 0xFF, 0, lpObj->iShield);
		}
	}

	if (gObjGetValueOfBuffIndex(lpObj, BUFFTYPE_ACHERON_FROST, &nEffectValue, 0))
	{
		iTotalMaxMana /= 2.0;
		lpObj->AddMana -= iTotalMaxMana;

		if ( (lpObj->MaxMana + lpObj->AddMana ) < lpObj->Mana )
		{
			lpObj->Mana = lpObj->MaxMana + lpObj->AddMana;
			GSProtocol.GCManaSend(lpObj->m_Index, lpObj->Mana, 0xFF, 0, lpObj->BP);
		}
	}

	if (gObjGetValueOfBuffIndex(lpObj, BUFFTYPE_ACHERON_BIND, &nEffectValue, 0))
	{
		iTotalMaxLife /= 2.0;
		lpObj->AddLife -= iTotalMaxLife;

		if ( (lpObj->MaxLife + lpObj->AddLife ) < lpObj->Life )
		{
			lpObj->Life = lpObj->MaxLife + lpObj->AddLife;
			GSProtocol.GCReFillSend(lpObj->m_Index, lpObj->Life, 0xFF, 0, lpObj->iShield);
		}
	}

	GSProtocol.GCReFillSend(lpObj->m_Index, iTotalMaxLife, 0xFE, 0, iTotalMaxShield);
	GSProtocol.GCManaSend(lpObj->m_Index, iTotalMaxMana, 0xFE, 0, lpObj->MaxBP + lpObj->AddBP);
	g_StatSpec.SendOptionList(lpObj);
	GSProtocol.GCSendAttackSpeed(lpObj->m_Index);
}

void CObjCalCharacter::GetSetItemOption(LPOBJ lpObj, LPBYTE pSetOptionTable, LPBYTE pSetOptionCountTable, int * pSetOptionCount)
{
	*pSetOptionCount = 0;
	int RightSetOptionIndex = -1;
	int RightRingSetOptionIndex = -1;

	for ( int i=0;i<INVETORY_WEAR_SIZE;i++)
	{
		//1.01.00 fix
		if (lpObj->pInventory[i].m_Durability <= 0.0)
		{
			continue;
		}

		int iSetItemType = lpObj->pInventory[i].IsSetItem();

		if ( iSetItemType )
		{
			if ( lpObj->pInventory[i].m_IsValidItem )
			{
				int iSetItemOption = gSetItemOption.GetOptionTableIndex(lpObj->pInventory[i].m_Type, iSetItemType );

				if ( iSetItemOption != -1 )
				{
					if ( i == 0 )
					{
						RightSetOptionIndex = iSetItemOption;
					}
					else if ( i == 1  )
					{
						if (  RightSetOptionIndex == iSetItemOption )
						{
							if ( lpObj->pInventory[0].GetWeaponType() == lpObj->pInventory[1].GetWeaponType() && lpObj->pInventory[0].IsSetItem() )
							{
								continue;
							}
						}
					}
					
					if ( i == 10 )
					{
						RightRingSetOptionIndex = iSetItemOption;
					}
					else if ( i == 11   )
					{
						if ( lpObj->pInventory[10].m_Type == lpObj->pInventory[11].m_Type && lpObj->pInventory[10].IsSetItem() )
						{
							continue;
						}
					}

					BOOL bFound = FALSE;

					for ( int n=0;n<*pSetOptionCount;n++)
					{
						if ( pSetOptionTable[n] == iSetItemOption )
						{
							pSetOptionCountTable[n]++;
							bFound = TRUE;

							break;
						}
					}

					if ( bFound == FALSE )
					{
						pSetOptionTable[*pSetOptionCount] = iSetItemOption;
						pSetOptionCountTable[*pSetOptionCount]++;

						*pSetOptionCount += 1;
					}
				}
			}
		}
	}
}

void CObjCalCharacter::CalcSetItemStat(LPOBJ lpObj)
{
	BYTE SetOptionTable[29];
	BYTE SetOptionCountTable[29];
	int SetOptionCount = 0;
	int op1;
	int op2;
	int op3;
	int op4;
	int op5;
	int opvalue1;
	int opvalue2;
	int opvalue3;
	int opvalue4;
	int opvalue5;

	if (lpObj->Type != OBJ_USER)
	{
		return;
	}

	for ( int i=0;i<INVETORY_WEAR_SIZE;i++)
	{
		if ( lpObj->pInventory[i].IsSetItem() && lpObj->pInventory[i].m_IsValidItem && lpObj->pInventory[i].m_Durability != 0.0f )
		{
			int AddStatType = lpObj->pInventory[i].GetAddStatType();

			switch ( AddStatType )
			{
				case 1:
					lpObj->pInventory[i].SetItemPlusSpecialStat(&lpObj->AddStrength, 196);
					break;
				case 2:
					EnterCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
					lpObj->m_PlayerData->AgilityCheckDelay = GetTickCount();
					lpObj->pInventory[i].SetItemPlusSpecialStat(&lpObj->AddDexterity, 197);
					LeaveCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
					break;
				case 3:
					lpObj->pInventory[i].SetItemPlusSpecialStat(&lpObj->AddEnergy, 198);
					break;
				case 4:
					lpObj->pInventory[i].SetItemPlusSpecialStat(&lpObj->AddVitality, 199);
					break;
			}
		}
	}

	memset(SetOptionTable, 0, sizeof(SetOptionTable));
	memset(SetOptionCountTable, 0, sizeof(SetOptionCountTable));
	this->GetSetItemOption(lpObj, SetOptionTable, SetOptionCountTable, &SetOptionCount);

	for ( int optioncount=0;optioncount<SetOptionCount;optioncount++)
	{
		int OptionTableIndex = SetOptionTable[optioncount];
		int OptionTableCount = SetOptionCountTable[optioncount];

		if ( OptionTableCount >= 2 )
		{
			gSetItemOption.IsRequireClass(OptionTableIndex, lpObj->Class, lpObj->m_PlayerData->ChangeUP);

			for ( int tablecnt=0;tablecnt<(OptionTableCount-1);tablecnt++)
			{
				if ( gSetItemOption.GetSetOption(OptionTableIndex, tablecnt,
					op1, op2, opvalue1, opvalue2, lpObj->Class, lpObj->m_PlayerData->ChangeUP) )
				{
					this->SetItemStatPlusSpecial(lpObj, op1, opvalue1);
					this->SetItemStatPlusSpecial(lpObj, op2, opvalue2);
				}
			}

			if ( gSetItemOption.GetMaxSetOptionCount(OptionTableIndex) < OptionTableCount )
			{
				lpObj->m_PlayerData->IsFullSetItem = true;

				gSetItemOption.GetGetFullSetOption(OptionTableIndex,
					op1, op2, op3, op4, op5,
					opvalue1, opvalue2, opvalue3, opvalue4, opvalue5,
					lpObj->Class, lpObj->m_PlayerData->ChangeUP);

				this->SetItemStatPlusSpecial(lpObj, op1, opvalue1);
				this->SetItemStatPlusSpecial(lpObj, op2, opvalue2);
				this->SetItemStatPlusSpecial(lpObj, op3, opvalue3);
				this->SetItemStatPlusSpecial(lpObj, op4, opvalue4);
				this->SetItemStatPlusSpecial(lpObj, op5, opvalue5);
			}		
		}
	}
}

void CObjCalCharacter::CalcSetItemOption(LPOBJ lpObj)
{
	BYTE SetOptionTable[29];
	BYTE SetOptionCountTable[29];
	int SetOptionCount = 0;
	int op1;
	int op2;
	int op3;
	int op4;
	int op5;
	int opvalue1;
	int opvalue2;
	int opvalue3;
	int opvalue4;
	int opvalue5;

	memset(SetOptionTable, 0, sizeof(SetOptionTable));
	memset(SetOptionCountTable, 0, sizeof(SetOptionCountTable));
	this->GetSetItemOption(lpObj, SetOptionTable, SetOptionCountTable, &SetOptionCount);

	for (int optioncount=0;optioncount<SetOptionCount;optioncount++)
	{
		int OptionTableIndex = SetOptionTable[optioncount];
		int OptionTableCount = SetOptionCountTable[optioncount];

		if ( OptionTableCount >= 2 )
		{
			for ( int tablecnt=0;tablecnt<(OptionTableCount-1);tablecnt++)
			{
				if ( gSetItemOption.GetSetOption(OptionTableIndex, tablecnt,
					op1, op2, opvalue1, opvalue2, lpObj->Class, lpObj->m_PlayerData->ChangeUP) )
				{
					this->SetItemPlusSpecial(lpObj, op1, opvalue1);
					this->SetItemPlusSpecial(lpObj, op2, opvalue2);
				}
			}
				
			if ( gSetItemOption.GetMaxSetOptionCount(OptionTableIndex) < OptionTableCount )
			{
				lpObj->m_PlayerData->IsFullSetItem = true;

				gSetItemOption.GetGetFullSetOption(OptionTableIndex,
					op1, op2, op3, op4, op5,
					opvalue1, opvalue2, opvalue3, opvalue4, opvalue5,
					lpObj->Class, lpObj->m_PlayerData->ChangeUP);

				this->SetItemPlusSpecial(lpObj, op1, opvalue1);
				this->SetItemPlusSpecial(lpObj, op2, opvalue2);
				this->SetItemPlusSpecial(lpObj, op3, opvalue3);
				this->SetItemPlusSpecial(lpObj, op4, opvalue4);
				this->SetItemPlusSpecial(lpObj, op5, opvalue5);
			}
		}
	}
}

void CObjCalCharacter::SetItemStatPlusSpecial(LPOBJ lpObj, int option, int ivalue)
{
	if (lpObj->Type != OBJ_USER)
		return;

	if ( option == -1 )
		return;

	switch ( option )
	{
		case AT_SET_OPTION_IMPROVE_STRENGTH:
			lpObj->AddStrength += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_DEXTERITY:
			EnterCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
			lpObj->m_PlayerData->AgilityCheckDelay = GetTickCount();
			lpObj->AddDexterity += ivalue;
			LeaveCriticalSection(&lpObj->m_PlayerData->AgiCheckCriti);
			break;
		case AT_SET_OPTION_IMPROVE_ENERGY:
			lpObj->AddEnergy += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_VITALITY:
			lpObj->AddVitality += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_LEADERSHIP:
			lpObj->AddLeadership += ivalue;
			break;
	}
}

void CObjCalCharacter::SetItemPlusSpecial(LPOBJ lpObj, int option, int ivalue)
{
	if ( option == -1 )
		return;
	if( ivalue == 0 )
	{
		g_Log.Add("Error iValue == 0, option %d",option);
		return;
	}
	switch ( option )
	{
		case AT_SET_OPTION_DOUBLE_DAMAGE:
			lpObj->m_PlayerData->SetOpDoubleDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_SHIELD_DEFENCE:
			lpObj->m_PlayerData->SetOpImproveSheldDefence += ivalue;
			break;
		case AT_SET_OPTION_TWO_HAND_SWORD_IMPROVE_DAMAGE:
			lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MINATTACKDAMAGE:
			lpObj->m_PlayerData->SetOpAddMinAttackDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAXATTACKDAMAGE:
			lpObj->m_PlayerData->SetOpAddMaxAttackDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAGICDAMAGE:
			lpObj->m_PlayerData->SetOpAddMagicPower += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_DAMAGE:
			lpObj->m_PlayerData->SetOpAddDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_ATTACKRATE:
			lpObj->m_PlayerData->SetOpImproveSuccessAttackRate += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_DEFENCE:
			lpObj->m_PlayerData->SetOpAddDefence += ivalue;
			break;
		case AT_SET_OPTION_DEFENCE_IGNORE:
			lpObj->m_PlayerData->SetOpIgnoreDefense += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAXLIFE:
			lpObj->AddLife += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAXMANA:
			lpObj->AddMana += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAXAG:
			lpObj->AddBP += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_AG_VALUE:
			lpObj->m_PlayerData->SetOpIncAGValue += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_CRITICALDAMAGE:
			lpObj->m_PlayerData->SetOpAddCriticalDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_CRITICALDAMAGE_SUCCESS:
			lpObj->m_PlayerData->SetOpAddCriticalDamageSuccessRate += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_EX_DAMAGE:
			lpObj->m_PlayerData->SetOpAddExcellentDamage += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_EX_DAMAGE_SUCCESS:
			lpObj->m_PlayerData->SetOpAddExcellentDamageSuccessRate += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_SKILLDAMAGE:
			lpObj->m_PlayerData->SetOpAddSkillAttack += ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_ATTACKDAMAGE_WITH_DEX:
			lpObj->m_PlayerData->SetOpAddAttackDamage += ( lpObj->m_PlayerData->Dexterity + lpObj->AddDexterity ) / ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_ATTACKDAMAGE_WITH_STR:
			lpObj->m_PlayerData->SetOpAddAttackDamage += ( lpObj->m_PlayerData->Strength + lpObj->AddStrength ) / ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_MAGICDAMAGE_WITH_ENERGY:
			lpObj->m_PlayerData->SetOpAddMagicPower += ( lpObj->m_PlayerData->Energy + lpObj->AddEnergy ) / ivalue;
			break;
		case AT_SET_OPTION_IMPROVE_DEFENCE_WITH_DEX:
			lpObj->m_PlayerData->SetOpAddDefence += lpObj->m_PlayerData->Dexterity + ( lpObj->AddDexterity  / ivalue );
			break;
		case AT_SET_OPTION_IMPROVE_DEFENCE_WITH_VIT:
			lpObj->m_PlayerData->SetOpAddDefence += ( lpObj->m_PlayerData->Vitality + lpObj->AddVitality ) / ivalue;
			break;
		case AT_SET_OPTION_FIRE_MASTERY:
			lpObj->m_AddResistance[R_FIRE] += ivalue;
			break;
		case AT_SET_OPTION_ICE_MASTERY:
			lpObj->m_AddResistance[R_ICE] += ivalue;
			break;
		case AT_SET_OPTION_THUNDER_MASTERY:
			lpObj->m_AddResistance[R_LIGHTNING] += ivalue;
			break;
		case AT_SET_OPTION_POSION_MASTERY:
			lpObj->m_AddResistance[R_POISON] += ivalue;
			break;
		case AT_SET_OPTION_WATER_MASTERY:
			lpObj->m_AddResistance[R_WATER] += ivalue;
			break;
		case AT_SET_OPTION_WIND_MASTERY:
			lpObj->m_AddResistance[R_WIND] += ivalue;
			break;
		case AT_SET_OPTION_EARTH_MASTERY:
			lpObj->m_AddResistance[R_EARTH] += ivalue;
			break;
	}
}

void CObjCalCharacter::SetItemApply(LPOBJ lpObj)
{
	lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->SetOpAddAttackDamage;
	lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->SetOpAddAttackDamage;
	lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->SetOpAddAttackDamage;
	lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->SetOpAddAttackDamage;

	lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->SetOpAddMinAttackDamage;
	lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->SetOpAddMinAttackDamage;

	lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->SetOpAddMaxAttackDamage;
	lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->SetOpAddMaxAttackDamage;

	lpObj->m_MagicDamageMin += lpObj->m_MagicDamageMin * lpObj->m_PlayerData->SetOpAddMagicPower / 100;
	lpObj->m_MagicDamageMax += lpObj->m_MagicDamageMax * lpObj->m_PlayerData->SetOpAddMagicPower / 100;

	lpObj->AddLife += INT(lpObj->AddVitality * DCInfo.DefClass[lpObj->Class].VitalityToLife);
	lpObj->AddMana += INT(lpObj->AddEnergy * DCInfo.DefClass[lpObj->Class].EnergyToMana);

	lpObj->m_CriticalDamageSuccessRate += lpObj->m_PlayerData->SetOpAddCriticalDamageSuccessRate;
	lpObj->m_ExcellentDamageSuccessRate += lpObj->m_PlayerData->SetOpAddExcellentDamageSuccessRate;

	if ( lpObj->pInventory[10].IsSetItem() )
	{
		lpObj->pInventory[10].PlusSpecialSetRing((LPBYTE)lpObj->m_AddResistance);
	}

	if ( lpObj->pInventory[11].IsSetItem() )
	{
		lpObj->pInventory[11].PlusSpecialSetRing((LPBYTE)lpObj->m_AddResistance);
	}

	if ( lpObj->pInventory[9].IsSetItem() )
	{
		lpObj->pInventory[9].PlusSpecialSetRing((LPBYTE)lpObj->m_AddResistance);
	}

	if ( lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage )
	{
		if ( !lpObj->pInventory[1].IsItem() && !lpObj->pInventory[0].IsItem() )
		{
			lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage = 0;
		}

		if ( lpObj->pInventory[1].IsItem() && lpObj->pInventory[1].m_TwoHand == FALSE )
		{
			lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage = 0;
		}

		if ( lpObj->pInventory[0].IsItem() && lpObj->pInventory[0].m_TwoHand == FALSE )
		{
			lpObj->m_PlayerData->SetOpTwoHandSwordImproveDamage = 0;
		}
	}

	if ( (lpObj->MaxLife + lpObj->AddLife ) < lpObj->Life )
	{
		lpObj->Life = lpObj->MaxLife + lpObj->AddLife;
		GSProtocol.GCReFillSend(lpObj->m_Index, lpObj->Life, 0xFF, 0, lpObj->iShield);
	}

	gObjSetBP(lpObj->m_Index);

	if ( (lpObj->MaxMana + lpObj->AddMana ) < lpObj->Mana )
	{
		lpObj->Mana = lpObj->MaxMana + lpObj->AddMana;
		GSProtocol.GCManaSend(lpObj->m_Index, lpObj->Mana, 0xFF, 0, lpObj->BP);
	}

	lpObj->m_Defense += lpObj->m_PlayerData->SetOpAddDefence * 10 / 20;
	lpObj->m_Defense += lpObj->m_Defense * lpObj->m_PlayerData->SetOpAddDefenceRate / 100;

	if ( lpObj->pInventory[1].m_Type >= ITEMGET(6,0) && lpObj->pInventory[1].m_Type < ITEMGET(7,0) )
	{
		lpObj->m_Defense += lpObj->m_Defense * lpObj->m_PlayerData->SetOpImproveSheldDefence / 100;
	}
}

void CObjCalCharacter::ApplyMLSkillItemOption(LPOBJ lpObj)
{
	if (lpObj->Type != OBJ_USER)
		return;

	CItem* Right = &lpObj->pInventory[0];
	CItem* Left  = &lpObj->pInventory[1];
	CItem* Helmet = &lpObj->pInventory[2];
	CItem* Armor = &lpObj->pInventory[3];
	CItem* Pant = &lpObj->pInventory[4];
	CItem* Glove = &lpObj->pInventory[5];
	CItem* Boot = &lpObj->pInventory[6];
	CItem* Wing = &lpObj->pInventory[7];

	// pvp defense rate
	lpObj->m_PlayerData->m_DefenseRatePvP += lpObj->m_PlayerData->m_MPSkillOpt.iMpsPVPBlockingRate;

	// defense
	lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsDefence;

	// defense rate
	lpObj->m_DefenseRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsBlockingSuccessRate;
	
	BOOL nFullSet = FALSE;
	if (lpObj->Class == CLASS_MAGUMSA)
	{
		if (Armor->m_IsValidItem
		&& Glove->m_IsValidItem
		&& Pant->m_IsValidItem
		&& Boot->m_IsValidItem
		&& (Armor->m_Type%MAX_SUBTYPE_ITEMS == Glove->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Glove->m_Type%MAX_SUBTYPE_ITEMS == Pant->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Pant->m_Type%MAX_SUBTYPE_ITEMS == Boot->m_Type%MAX_SUBTYPE_ITEMS))
			nFullSet = TRUE;
	}
	else if (lpObj->Class == CLASS_RAGEFIGHTER)
	{
		if (Helmet->m_IsValidItem
		&& Armor->m_IsValidItem
		&& Pant->m_IsValidItem
		&& Boot->m_IsValidItem
		&& (Helmet->m_Type%MAX_SUBTYPE_ITEMS == Armor->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Armor->m_Type%MAX_SUBTYPE_ITEMS == Pant->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Pant->m_Type%MAX_SUBTYPE_ITEMS == Boot->m_Type%MAX_SUBTYPE_ITEMS))
			nFullSet = TRUE;
	}
	else
	{
		if (Helmet->m_IsValidItem
		&& Armor->m_IsValidItem
		&& Glove->m_IsValidItem
		&& Pant->m_IsValidItem
		&& Boot->m_IsValidItem
		&& (Helmet->m_Type%MAX_SUBTYPE_ITEMS == Armor->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Armor->m_Type%MAX_SUBTYPE_ITEMS == Glove->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Glove->m_Type%MAX_SUBTYPE_ITEMS == Pant->m_Type%MAX_SUBTYPE_ITEMS)
		&& (Pant->m_Type%MAX_SUBTYPE_ITEMS == Boot->m_Type%MAX_SUBTYPE_ITEMS))
			nFullSet = TRUE;
	}
	if (nFullSet == TRUE)
	{
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncSetItemDefense;
	}

	if (lpObj->pInventory[lpObj->m_btInvenPetPos].m_Type == ITEMGET(13, 4)) {
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncDarkHorseDefense;
	}

	// apply master wing
	if (Wing->m_IsValidItem) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_CurseDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_CurseDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDamage;
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddWingDefense;
	}

	// attack rate
	lpObj->m_AttackRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAttackSuccessRate;

	// apply master weapon Mastery
	lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddPhysicDamage;
	lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddPhysicDamage;
	lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddPhysicDamage;
	lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddPhysicDamage;
	lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddMagicDamage;
	lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddMagicDamage;
	lpObj->m_CurseDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncMagicPower;
	lpObj->m_CurseDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncMagicPower;

	// pvp attack rate
	lpObj->m_PlayerData->m_AttackRatePvP += lpObj->m_PlayerData->m_MPSkillOpt.iMpsPVPAttackDmgRate;

	// apply master weapon Strengthener
	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_TWOHAND_SWORD) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandSwordAddDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandSwordAddDamage;
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsPvPAttack;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsPvPAttack;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_TWOHAND_SWORD) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandSwordAddDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandSwordAddDamage;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsPvPAttack;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsPvPAttack;
	}
	
	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_ONEHAND_SWORD) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSwordAddDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSwordAddDamage;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_ONEHAND_SWORD) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSwordAddDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSwordAddDamage;
	}

	if ((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_ONEHAND_SWORD)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_ONEHAND_SWORD)) {
		lpObj->m_AttackSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddAttackSpeed;
		lpObj->m_MagicSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddAttackSpeed;
	}
	
	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_MACE) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsScultAddDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsScultAddDamage;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_MACE) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsScultAddDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsScultAddDamage;
	}

	if (!((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_MACE)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_MACE))) {
		lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddSturn = 0.0;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_SPEAR) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSpearAddDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSpearAddDamage;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_SPEAR) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSpearAddDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsSpearAddDamage;
	}

	if (!((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_SPEAR)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_SPEAR))) {
		lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddDoubleDamage = 0.0;
	}

	if ((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_ONEHAND_STAFF)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_ONEHAND_STAFF)) {
		lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsStaffMagicPower;
		lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsStaffMagicPower;
		lpObj->m_AttackSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsStaffAttackSpeed;
		lpObj->m_MagicSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsStaffAttackSpeed;
	}

	if ((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_TWOHAND_STAFF)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_TWOHAND_STAFF)) {
		lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsMagicPower;
		lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsMagicPower;
		lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsStaffManAttack;
		lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsStaffManAttack;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_BOW) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsBowAddDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsBowAddDamage;
		lpObj->m_AttackSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsBowAddAttackSpeed;
		lpObj->m_MagicSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsBowAddAttackSpeed;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_CROSSBOW) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsBowAddDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsBowAddDamage;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsBowManAttack;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsTwoHandsBowManAttack;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_DARKLORD_SCEPTER) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddCeptorDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddCeptorDamage;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncCeptorManAttack;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncCeptorManAttack;
	} else {
		lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddPetAttack = 0;
		lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddControlAttack = 0;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_SUMMONER_STAFF) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddMagicPowerStick;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddMagicPowerStick;
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddManAttackStick;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddManAttackStick;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_SUMMONER_BOOK) {
		lpObj->m_CurseDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddCurseMagicBook;
		lpObj->m_CurseDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddCurseMagicBook;
		lpObj->m_AttackSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddAttackSpeedMagicBook;
		lpObj->m_MagicSpeed += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddAttackSpeedMagicBook;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_RAGEFIGHTER_WEAPON) {
		lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMonkWeaponDamage;
		lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMonkWeaponDamage;
	}

	if (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_RAGEFIGHTER_WEAPON) {
		lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMonkWeaponDamage;
		lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMonkWeaponDamage;
	}

	if (!((Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_RAGEFIGHTER_WEAPON)
	|| (Right->m_IsValidItem && Right->GetDetailItemType() == ITEM_RAGEFIGHTER_WEAPON))) {
		lpObj->m_PlayerData->m_MPSkillOpt.iMpsMonkDoubleDamage = 0;
	}

	if (Left->m_IsValidItem && Left->GetDetailItemType() == ITEM_SHIELD) {
		// wizard
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncShield;
		lpObj->m_DefenseRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsShieldBlockRate;

		// elf
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsElfAddShield;
		lpObj->m_DefenseRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsElfShieldBlockRate;
		
		// lord
		lpObj->m_Defense += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddShieldDefense;
		lpObj->m_DefenseRate += lpObj->m_PlayerData->m_MPSkillOpt.iMpsAddShieldBlockingRage;
	}

	// calc master min/max attack damage
	lpObj->m_AttackDamageMinRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMinAttackDamage;
	lpObj->m_AttackDamageMaxRight += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxAttackDamage;
	lpObj->m_AttackDamageMinLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMinAttackDamage;
	lpObj->m_AttackDamageMaxLeft += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxAttackDamage;
	lpObj->m_MagicDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMinMagicDamage;
	lpObj->m_MagicDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsMaxMagicDamage;
	lpObj->m_CurseDamageMin += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncMinMagicDamage;
	lpObj->m_CurseDamageMax += lpObj->m_PlayerData->m_MPSkillOpt.iMpsIncMaxMagicDamage;
}

BOOL CObjCalCharacter::ValidItem(LPOBJ lpObj, CItem * lpItem, int pos)
{
	if (g_ConfigRead.data.common.CheckValidItem == false)
	{
		return TRUE;
	}

	LPITEM_ATTRIBUTE p = &ItemAttribute[lpItem->m_Type];
	
	if ( p->RequireStrength != 0 )
	{
		if ( (lpObj->m_PlayerData->Strength + lpObj->AddStrength) < ( lpItem->m_RequireStrength - lpItem->m_HJOpStrength ) )
		{
			return FALSE;
		}
	}

	if ( p->RequireDexterity != 0 )
	{
		if ( (lpObj->m_PlayerData->Dexterity + lpObj->AddDexterity) < ( lpItem->m_RequireDexterity - lpItem->m_HJOpDexterity ) )
		{
			return FALSE;
		}
	}


	if ( p->RequireEnergy != 0 )
	{
		if ( (lpObj->m_PlayerData->Energy + lpObj->AddEnergy) < lpItem->m_RequireEnergy )
		{
			return FALSE;
		}
	}

	if ( p->RequireLevel != 0 )
	{
		if ( (lpObj->Level) < lpItem->m_RequireLevel )
		{
			return FALSE;
		}
	}

	if ( lpItem->m_Type >= ITEMGET(0,0) && lpItem->m_Type <= ITEMGET(11, 0) )
	{
		if ( p->RequireVitality )
		{
			if ( (lpObj->m_PlayerData->Vitality + lpObj->AddVitality ) < lpItem->m_RequireVitality )
			{
				return FALSE;
			}
		}

		if ( p->RequireLeadership )
		{
			if ( (lpObj->Leadership + lpObj->AddLeadership ) < lpItem->m_RequireLeaderShip )
			{
				return FALSE;
			}
		}
	}

	return TRUE;
}
	


void CObjCalCharacter::CalcShieldPoint(LPOBJ lpObj)
{
	int iMaxShieldPoint = 0;
	int iExpressionA = ( lpObj->m_PlayerData->Strength + lpObj->AddStrength ) + ( lpObj->m_PlayerData->Dexterity + lpObj->AddDexterity ) + ( lpObj->m_PlayerData->Vitality + lpObj->AddVitality ) + (lpObj->m_PlayerData->Energy + lpObj->AddEnergy );

	if ( lpObj->Class == CLASS_DARKLORD )
	{
		iExpressionA += lpObj->Leadership + lpObj->AddLeadership;
	}

	if ( g_iShieldGageConstB == 0 )
		g_iShieldGageConstB = 30;

	int Level = lpObj->Level + lpObj->m_PlayerData->MasterLevel;

	int iExpressionB = (Level * Level) / g_iShieldGageConstB;

	if ( g_iShieldGageConstA == 0 )
		g_iShieldGageConstA = 12;

	iMaxShieldPoint = ( iExpressionA * g_iShieldGageConstA ) / 10 + iExpressionB  + lpObj->m_Defense;
	lpObj->iMaxShield = iMaxShieldPoint;
}

void CObjCalCharacter::PremiumItemApply(LPOBJ lpObj)
{
	for (int i=9;i<12;++i)
	{
		if(lpObj->pInventory[i].m_Durability <= 0.0 || lpObj->pInventory[i].IsPeriodItemExpire() == TRUE)
			continue;

		switch ( lpObj->pInventory[i].m_Type ) 
		{
			case ITEMGET(13,109):
				lpObj->AutoHPRecovery += 3;
				lpObj->AddMana += (short)(lpObj->MaxMana * 4 / 100.0f);
				break;
			case ITEMGET(13,110):
				lpObj->AutoHPRecovery += 3;
				lpObj->AddLife += (short)(lpObj->MaxLife * 4 / 100.0f);
				break;
			case ITEMGET(13,111):
				lpObj->AutoHPRecovery += 3;
				lpObj->MonsterDieGetMoney += 30;
				break;
			case ITEMGET(13,112):
				lpObj->AutoHPRecovery += 3;
				lpObj->DamageMinus += 4;
				break;
			case ITEMGET(13,113):
				lpObj->AutoHPRecovery += 3;
				lpObj->m_ExcellentDamageSuccessRate += 10;
				break;
			case ITEMGET(13,114):
				lpObj->AutoHPRecovery += 3;
				lpObj->m_AttackSpeed += 7;
				lpObj->m_MagicSpeed += 7;
				break;
			case ITEMGET(13,115):
				lpObj->AutoHPRecovery += 3;
				lpObj->MonsterDieGetMana += 1;
				break;
			case ITEMGET(13,128):
				lpObj->m_CriticalDamageSuccessRate += 10;
				break;
			case ITEMGET(13,129):
				lpObj->m_ExcellentDamageSuccessRate += 10;
				break;
			case ITEMGET(13,130):
				lpObj->AddLife += 50;
				break;
			case ITEMGET(13,131):
				lpObj->AddMana += 50;
				break;
			case ITEMGET(13,132):
				lpObj->AddLife += 100;
				lpObj->iAddShield += 500;
				break;
			case ITEMGET(13,133):
				lpObj->AddMana += 100;
				lpObj->AddBP += 500;
				break;
			case ITEMGET(13,134):
				break;
		}
	}
}