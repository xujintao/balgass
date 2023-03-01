#include <windows.h>

void f1() {}
void f2() {}

typedef void(*hookFunc)();
hookFunc hooks[] = {
	f1,
	f2,
};

void handleHooks() {
	for (auto h : hooks) {
		h();
	}
}

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