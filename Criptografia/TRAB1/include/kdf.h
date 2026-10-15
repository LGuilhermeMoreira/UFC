//
// Created by guigui on 25/11/2025.
//
#ifndef TRAB1_KDF_H
#define TRAB1_KDF_H
#include <openssl/bn.h>
#include <openssl/sha.h>

int kdf_aes_256(const BIGNUM* shared_key,char* result);

#endif //TRAB1_KDF_H