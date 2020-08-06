#pragma once
#include "config.h"
// include headers that implement a archive in simple text format
#include <boost/archive/text_oarchive.hpp>
#include <iostream>
//#include <boost/archive/text_iarchive.hpp>
NAMESPACE_BEGIN(WTFCipher)
/*
class msg {
	friend class boost::serialization::access;
	template<class Archive>
	void serialize(Archive& ar, const unsigned int version) {
		ar << degress;
		ar << minutes;
		ar << seconds;
	}
};
*/
size_t encrypt(unsigned char* out, const unsigned char* in, size_t len) {
	return 0;
}

size_t decrypt(unsigned char* out, const unsigned char* in, size_t len) {
	if (in == nullptr || len == 0) {
		return 0;
	}
	int flag = in[0];

	boost::archive::text_oarchive oa(std::cout);
	return 0;
}

NAMESPACE_END