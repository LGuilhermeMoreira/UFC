#include "dhke.h"
#include <math.h>

// doc: https://docs.openssl.org/1.0.2/man3/

int generate_large_prime(BIGNUM* prime,int bits)
{
   return BN_generate_prime_ex(prime,bits,1,NULL,NULL,NULL);
}

// range [2,prime-2]
// lógica: gerar [0,prime-4] e somar 2
int generate_private_key(BIGNUM* private_key,BIGNUM* prime)
{
    BIGNUM* max = BN_dup(prime);
    if (!max)
    {
        BN_free(max);
        return EXIT_FAILURE;
    }

    if (!BN_sub_word(max,4))
    {
        BN_free(max);
        return EXIT_FAILURE;
    };

    if (!BN_rand_range(private_key,max))
    {
        BN_free(max);
        return EXIT_FAILURE;
    }

    BIGNUM* min = BN_new();
    if (!min)
    {
        BN_free(min);
        BN_free(max);
        return EXIT_FAILURE;
    }

    if (BN_set_word(min,2))
    {
        BN_free(max);
        BN_free(min);
        return EXIT_FAILURE;
    }

    if (!BN_add(private_key,private_key,min))
    {
        BN_free(max);
        BN_free(min);
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}


BIGNUM* pow_mod(const BIGNUM* base,const  BIGNUM* exponent,const  BIGNUM* mod)
{
    BIGNUM* result = BN_new();

    if (!result)
    {
        BN_free(result);
        return NULL;
    }
    BN_CTX* ctx = BN_CTX_new();
    if (!ctx)
    {
        BN_free(result);
        BN_CTX_free(ctx);
        return NULL;
    }

    //r=a^p % m
    if (!BN_mod_exp(result, base, exponent, mod, ctx))
    {
        BN_free(result);
        BN_CTX_free(ctx);
        return NULL;
    }

    BN_CTX_free(ctx);

    return result;
}
