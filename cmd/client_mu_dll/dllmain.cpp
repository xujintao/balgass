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
	{ DIR_MEMASSIGN, 0x00405B22 + 1, (DWORD)"�豸δ��ʼ�� �밲װ���µ�ͼ����������"/*gbk encode*/, 4, 0 },
	{ DIR_MEMASSIGN, 0x00405BD0 + 1, (DWORD)"�ֱ��������á��밴��ѡ�����á���ť��ָ���µķֱ��ʡ�"/*gbk encode*/, 4, 0 },
	// { DIR_MEMASSIGN, 0x0040309D + 3, 0x30, 1, 0 }, // dmDisplayFrequency limit to >= 48
	// { DIR_MEMASSIGN, 0x004030A6 + 1, 0x8C, 1, 0 }, // dmDisplayFrequency jne -> jl
	// { DIR_MEMASSIGN, 0x0040314C + 0, 0xEB, 1, 0 }, // discard high dmDisplayFrequency, fixed resolution mismatching
	{ DIR_MEMASSIGN, 0x0040AC4A + 1, 0x17, 1, 0 }, // do not discard the first tag when version tags are all higher then current
	{ DIR_MEMASSIGN, 0x0040AC60, 0xE8EB, 2, 0 }, // do not discard the first tag when version tags are all higher then current
	{ DIR_MEMASSIGN, 0x0040AA07, 0x05EB, 2, 0 }, // keep status of failure on connecting failed
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
			break;
		default:
			break;
		}
	}
}