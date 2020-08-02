#include <windows.h>

void handleHooks();

BOOL APIENTRY DllMain(HMODULE hModule, DWORD ul_reason_for_call, LPVOID lpReserved)
{
	switch (ul_reason_for_call)
	{
	case DLL_PROCESS_ATTACH:
	{
		handleHooks();
		break;
	}
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}

enum DIR {
	DIR_MEMSET,
	DIR_MEMCOPY,
	DIR_MEMASSIGN
};

typedef struct {
	DIR directive;
	DWORD addrDst;
	DWORD value;
	DWORD valueSize;
	DWORD addrSrc;
}hook;

hook hooks[] = {
	{ DIR_MEMASSIGN, 0x00405B22 + 1, (DWORD)"设备未初始化 请安装最新的图形驱动程序"/*gbk encode*/, 4, 0 },
	{ DIR_MEMASSIGN, 0x00405BD0 + 1, (DWORD)"分辨率已重置。请按“选项设置”按钮以指定新的分辨率。"/*gbk encode*/, 4, 0 },
};

void handleHooks() {
	for (auto h : hooks) {
		switch (h.directive) {
		case DIR_MEMSET:break;
		case DIR_MEMCOPY:break;
		case DIR_MEMASSIGN:
			DWORD OldProtect;	
			switch (h.valueSize) {
			case 1:
				VirtualProtect((LPVOID)h.addrDst, 1, PAGE_EXECUTE_READWRITE, &OldProtect);
				*(BYTE*)h.addrDst = (BYTE)h.value;
				VirtualProtect((LPVOID)h.addrDst, 1, OldProtect, &OldProtect);
				break;
			case 2:
				VirtualProtect((LPVOID)h.addrDst, 2, PAGE_EXECUTE_READWRITE, &OldProtect);
				*(WORD*)h.addrDst = (WORD)h.value;
				VirtualProtect((LPVOID)h.addrDst, 2, OldProtect, &OldProtect);
				break;
			case 4:
				VirtualProtect((LPVOID)h.addrDst, 4, PAGE_EXECUTE_READWRITE, &OldProtect);
				*(DWORD*)h.addrDst = (DWORD)h.value;
				VirtualProtect((LPVOID)h.addrDst, 4, OldProtect, &OldProtect);
				break;
			default:__asm int 3;
			}
			switch (h.valueSize) {
			case 1:
				if (*(BYTE*)h.addrDst != (BYTE)h.value) {
					__asm int 3;
				}
				break;
			case 2:
				if (*(WORD*)h.addrDst != (WORD)h.value) {
					__asm int 3;
				}
				break;
			case 4:
				if (*(DWORD*)h.addrDst != (DWORD)h.value) {
					__asm int 3
				}
				break;
			default:__asm int 3;
			}
		default:
			break;
		}
	}
}