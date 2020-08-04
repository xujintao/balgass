#pragma once

#define NAMESPACE_BEGIN(x) namespace x {
#define NAMESPACE_END }

#define VERSION "1.0"

#ifdef _UNICODE
#define tstring std::wstring
#else
#define tstring std::string
#endif