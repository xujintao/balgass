#include "stdafx.h"
#include "exception.h"
#include "hook.h"

void HookThis(DWORD Addr, DWORD size, DWORD newAddr)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)Addr, size, PAGE_EXECUTE_READWRITE, &OldProtect);
	memset((LPVOID)Addr, 0x90, size);
	*(BYTE*)(Addr) = 0xE9;
	*(DWORD*)(Addr + 1) = newAddr - (Addr + 5);
	VirtualProtect((LPVOID)Addr, size, OldProtect, &OldProtect);
}

void MemSet(DWORD Addr, int mem, int size)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)Addr, size, PAGE_EXECUTE_READWRITE, &OldProtect);
	memset((void*)Addr, mem, size);
	VirtualProtect((LPVOID)Addr, size, OldProtect, &OldProtect);
	for (int i = 0; i < size; i++) {
		if (0 != memcmp((void*)(Addr + i), (void*)&mem, 1)) {
			__asm int 3;
		}
	}
}

void MemCpy(DWORD Addr, void* mem, int size)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)Addr, size, PAGE_EXECUTE_READWRITE, &OldProtect);
	memcpy((void*)Addr, mem, size);
	VirtualProtect((LPVOID)Addr, size, OldProtect, &OldProtect);
	for (int i = 0; i < size; i++) {
		if (0 != memcmp((void*)(Addr + i), (void*)(uintptr_t(mem) + i), 1)) {
			__asm int 3;
		}
	}
}

void MemAssign(DWORD addr, BYTE value)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)addr, 1, PAGE_EXECUTE_READWRITE, &OldProtect);
	*(BYTE*)addr = value;
	VirtualProtect((LPVOID)addr, 1, OldProtect, &OldProtect);
	if (*(BYTE*)addr != value) {
		__asm int 3;
	}
}

void MemAssign(DWORD addr, WORD value)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)addr, 2, PAGE_EXECUTE_READWRITE, &OldProtect);
	*(WORD*)addr = value;
	VirtualProtect((LPVOID)addr, 2, OldProtect, &OldProtect);
	if (*(WORD*)addr != value) {
		__asm int 3;
	}
}

void MemAssign(DWORD addr, DWORD value)
{
	DWORD OldProtect;
	VirtualProtect((LPVOID)addr, 4, PAGE_EXECUTE_READWRITE, &OldProtect);
	*(DWORD*)addr = value;
	VirtualProtect((LPVOID)addr, 4, OldProtect, &OldProtect);
	if (*(DWORD*)addr != value) {
		__asm int 3;
	}
}

CHookManager HookManager;

CHookManager::CHookManager()
{
	m_Count = 0;
	m_LastEntry = -1;
}

void CHookManager::RegisterHook(DWORD Addr, DWORD Type)
{
	HOOK_STRUCT hook;
	hook.Addr = Addr;
	hook.Id = m_Count;
	m_Count++;
	switch (Type)
	{
	case 0:
		hook.Type = "Jmp";
		break;
	case 1:
		hook.Type = "Call";
		break;
	case 2:
		hook.Type = "Class";
		break;
	}
	m_List.push_back(hook);

}

void CHookManager::MakeJmpHook(DWORD Addr, DWORD size, void* pProc)
{
	try {
		DWORD memlen = 10;		// MOV DWORD PTR...
		memlen += 5;				// JMP pProc;

		DWORD NewAddr = (DWORD)VirtualAlloc(NULL, memlen, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
		if (!NewAddr)
		{
			throw CException("Virtual Alloc error: %d", GetLastError());
		}

		DWORD offs = 0;

		BYTE aMovMem[] = { 0xC7, 0x05 };
		memcpy((LPVOID)(NewAddr + offs), aMovMem, 2);
		offs += 2;

		DWORD* pEntry = &m_LastEntry;
		memcpy((LPVOID)(NewAddr + offs), &pEntry, 4);
		offs += 4;

		memcpy((LPVOID)(NewAddr + offs), &m_Count, 4);
		offs += 4;

		HookThis(NewAddr + offs, 5, (DWORD)pProc);

		HookThis(Addr, size, NewAddr);

		RegisterHook(Addr, 0);
	}
	catch (std::exception& e)
	{
		//		CErrorLog::OutException(e.what());
	}
}