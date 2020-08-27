// stdafx.h : include file for standard system include files,
// or project specific include files that are used frequently, but
// are changed infrequently
//

#pragma once

#include "targetver.h"

#define DLL_VERSION "9.5.1.13"

#define WIN32_LEAN_AND_MEAN             // Exclude rarely-used stuff from Windows headers
// Windows Header Files:
#include <windows.h>
#include <mmsystem.h>
#include <string>
#include <sstream>
#include <vector>
#include <list>
#include <stdio.h>
#include <math.h>
#include <WindowsX.h>
#include <ShellAPI.h>
#include <tlhelp32.h>
#include <tchar.h>
#include <wingdi.h>
#include <map>

// socket
#include <WinSock2.h>
#pragma comment(lib, "ws2_32.lib")

// psapi: The process status API
#include <Psapi.h>
#pragma comment(lib, "psapi.lib")

// cryptopp
#pragma comment(lib, "cryptlib.lib")

// opengl
#include <gl\GL.h>
#pragma comment(lib, "Opengl32.lib")

// glu: The OpenGL Utility (GLU) library
//#include <gl\GLU.h>
//#pragma comment(lib, "glu32.lib")

// GLee: OpenGL Easy Extension library)
//#include <GLee.h>
//#pragma comment(lib, "Glee.lib")

// mpir: Multiple Precision Integers and Rationals
//#include <gmp.h>
//#include <mpir.h>
//#pragma comment(lib, "mpir.lib")

//#include "SecureEngineSDK.h"
//#pragma comment(lib, "SecureEngineSDK32.lib")		// wait
// TODO: reference additional headers your program requires here

