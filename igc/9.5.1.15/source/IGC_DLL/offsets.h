#define MU_MAP_ADDR 0x1206ED0 // S9
#define LOAD_ITEMS_CALLBACK 0x006185BE // S9

#define VERSION_HOOK1 0x0043EC09+3 // 1.04R
#define VERSION_HOOK2 0x004FF25B+3 // 1.04R
//#define VERSION_HOOK3 0x09ED4BAF+3 // S9
#define SERIAL_HOOK1 0x0043EC4C+3 // 1.04R
#define SERIAL_HOOK2 0x004FF29E+3 // 1.04R


#define CHARSET_ADDRESS 0x1205B4C // S9
#define XSHOP_FIX 0x005C43E4 // S9


#define MU_WINDOW_SIZE 0x1206F10 // S9
#define MU_DRAW_VERSION 0x004DF445 // S9

extern DWORD* g_Map;

#define MAKE_NUMBERW(x,y)  ( (WORD)(((BYTE)((y)&0xFF)) |   ((BYTE)((x)&0xFF)<<8 ))  )

#define MU_GET_CHARR_POS1 0x00969F14 // S9
#define MU_GET_CHARR_POS2 0x00829DDE // S9
#define MU_GET_CHARR_POS3 0x00840443 // S9
#define SETGLOW_EX 0x00678694 // 1.04R
#define SETGLOW_BACKJMP_EX702 0x0067869A // 1.04R

#define PROTOCOL_SOCKET_ADDR 0x08C88FF0+0xC // 1.04R

#define CTRL_FREEZE_FIX 0x005359B7+1 // 1.04R
#define MAPSRV_DELACCID_FIX 0x004FAC1D // S9

#define JEWEL_OF_BLESS_PRICE 0x59C295+3 // S9
#define JEWEL_OF_SOUL_PRICE 0x59C2B1+3  // S9
#define JEWEL_OF_CHAOS_PRICE 0x59C2CD+3 // S9
#define JEWEL_OF_LIFE_PRICE 0x59C2E9+3 // S9
#define JEWEL_OF_CREATION_PRICE 0x59C305+3 // S9
#define CREST_OF_MONARCH_PRICE 0x0059C409+3 // S9
#define LOCH_FEATHER_PRICE 0x0059C3FC+3 // S9
#define JEWEL_OF_GUARDIAN_PRICE 0x0059CBE7+3 // S9
#define WERERABBIT_EGG_PRICE 0x0059DF54+3 // S9
#define ITEM_SELL_DIVISOR 0x0059DA9F+1 // S9


//123CA8C

#define pGameWindow				*(HWND*)(0x123CA8C)

//00958F37   55               PUSH EBP


#define pGetPreviewStruct		((DWORD(__thiscall*)(LPVOID This, int ViewportID)) 0x958F37)
#define pPreviewThis			((LPVOID(*)()) 0x437357)


// UI
#define pMessageBoxThis			((DWORD(__cdecl*)()) 0x00969F14) // S9

enum ObjState
{
	SelectServer	= 2,
	SwitchCharacter = 4,
	GameProcess		= 5,
};

#define pPlayerState			*(int*)0x1205338 // s9e2


#define PARSE_PACKET_HOOK 0x00760292 // 1.04R
#define PARSE_PACKET_STREAM 0x006BDC33 // 1.04R
#define PROTOCOL_CORE1 0x00735DC7 // 1.04R
#define PROTOCOL_CORE2 0x0075C3B2 // 1.04R

#define SEND_PACKET_HOOK 0x00439420 // 1.04R
#define MU_SEND_PACKET 0x004397E3 // 1.04R
#define MU_SENDER_CLASS 0x012E4034 // 1.04R

#define ACCOUNT_ID 0x08C88F78 // 1.04R

#define SPRINTF_HOOK 0x00DE817A // 1.04R

#define DMG_SEND_HOOK 0x006CABB6 // 1.04R
#define DMG_SEND_HOOK1 0x0089C526 // 1.04R
#define DMG_SEND_HOOK_RET 0x006CABDE // 1.04R

#define DMG_SEND_RF_HOOK 0x006CA7D9 // 1.04R
#define DMG_SEND_RF_HOOK_RET 0x006CA7EA // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK1 0x00528A00 // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK2 0x00528BDF // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK3 0x00528C41 // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK4 0x00528CA3 // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK5 0x00528D1A // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK6 0x00528D58 // 1.04R
#define RF_SKILL_DMG_DISP_FIX_HOOK7 0x00528D80 // 1.04R

#define HP_SEND_HOOK 0x00AAA2F5 // 1.04R
#define HP_SEND_HOOK1 0x00A3A4F2 // 1.04R
#define HP_SEND_HOOK_RET 0x00AAA32B // 1.04R

#define MANA_SEND_HOOK 0x00AAA34C // 1.04R
#define MANA_SEND_HOOK1 0x00A3A4F2 // 1.04R

#define HP_SEND_DK_HOOK 0x00A8A786 // 1.04R
#define HP_SEND_DK_HOOK1 0x08C88E58 // 1.04R
#define HP_SEND_DK_HOOK2 0x086105EC // 1.04R
#define HP_SEND_DK_HOOK3 0x08610600 // 1.04R
#define HP_SEND_DK_HOOK4 0x00436DA8 // 1.04R
#define HP_SEND_DK_HOOK_RET 0x00A8A7AA // 1.04R

#define HP_SEND_C_HOOK 0x00A8A7BC // 1.04R
#define HP_SEND_C_HOOK1 0x08610600 // 1.04R
#define HP_SEND_C_HOOK2 0x00436DA8 // 1.04R
#define HP_SEND_C_HOOK_RET 0x00A8A7E5 // 1.04R

#define MANA_SEND_DK_HOOK 0x00A8AAB0 // 1.04R
#define MANA_SEND_DK_HOOK1 0x08610600 // 1.04R
#define MANA_SEND_DK_HOOK2 0x00436DA8 // 1.04R
#define MANA_SEND_DK_HOOK_RET 0x00A8AAD4 // 1.04R

#define MANA_SEND_C_HOOK 0x00A8AAE6 // 1.04R
#define MANA_SEND_C_HOOK1 0x086105EC // 1.04R
#define MANA_SEND_C_HOOK2 0x08610600 // 1.04R
#define MANA_SEND_C_HOOK3 0x00436DA8 // 1.04R
#define MANA_SEND_C_HOOK_RET 0x00A8AB0F // 1.04R

#define SKILL_65K_HOOK 0x00A0A03B // 1.04R
#define SKILL_65K_HOOK1 0x00A0A1E7 // 1.04R
#define SKILL_65K_HOOK_RET 0x00A0A044 // 1.04R

#define ERRTEL_MIX_STAFF_HOOK 0x00AC2B9F // 1.04R
#define ERRTEL_MIX_STAFF_HOOK1 0x008E2BBA // 1.04R
#define ERRTEL_MIX_STAFF_HOOK_RET1 0x00AC2BA7 // 1.04R
#define ERRTEL_MIX_STAFF_HOOK_RET2 0x00AC2BB3 // 1.04R

#define ITEM_LEVEL_REQ_HOOK 0x00A28654 // 1.04R
#define ITEM_LEVEL_REQ_HOOK_RET 0x00A2865A // 1.04R

#define MOVE_WND_HOOK 0x00D52D0B // 1.04R
#define MOVE_WND_HOOK_RET1 0x00D52D21 // 1.04R
#define MOVE_WND_HOOK_RET2 0x00D52D11 // 1.04R
#define MOVE_WND_HOOK1 0x00D52D18 // 1.04R

#define ALTER_PSHOP_DISPLAY_HOOK 0x008E4492 // 1.04R
#define ALTER_PSHOP_DISPLAY_HOOK_RET 0x008E4499 // 1.04R

#define MU_WND_PROC_HOOK 0x004D6FDC+3 // 1.04R
#define MU_WND_PROC 0x004D5F98 // 1.04R