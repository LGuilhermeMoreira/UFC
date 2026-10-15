#ifndef DHKE_H
#define DHKE_H
#include <openssl/bn.h>

int generate_large_prime(BIGNUM* prime,int bits);
int generate_private_key(BIGNUM* private_key,BIGNUM* prime);
BIGNUM* pow_mod(const BIGNUM* base, const BIGNUM* exponent,const BIGNUM* mod);

#endif
