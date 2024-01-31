#pragma once
void HookThis(DWORD Addr, DWORD size, DWORD newAddr);
void MemSet(DWORD Addr, int mem, int size);
void MemCpy(DWORD Addr, void* mem, int size);
void MemAssign(DWORD Addr, BYTE value);
void MemAssign(DWORD Addr, WORD value);
void MemAssign(DWORD Addr, DWORD value);

struct HOOK_STRUCT {
	DWORD Id;
	DWORD Addr;
	std::string Type;
};

class CHookManager {
	DWORD m_LastEntry;
	DWORD m_Count;
	std::list<HOOK_STRUCT> m_List;
public:
	CHookManager();
	void XorLine(char* buff, int len);
	void DumpList(std::string file);
	void RegisterHook(DWORD Addr, DWORD Type);

	void MakeJmpHook(DWORD Addr, DWORD size, void* pProc);
	void MakeCallback(DWORD Addr, void* pFunc, int Args, int MemSize, bool saveEcx, bool keep = true);
	void MakeClassCallback(DWORD Addr, void* pMethod, void* pObject, int Args, int MemSize, bool saveEcx);
};

extern CHookManager HookManager;