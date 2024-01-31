#pragma once
class CXor {
public:
	static void Enc(unsigned char* buff, int headerSize, int len);
	static void Dec(unsigned char* buff, int headerSize, int len);
	static void BuxConvert(char* buf, int size);
};
extern CXor* Xor;