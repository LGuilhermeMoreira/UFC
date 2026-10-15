//
// Created by guigui on 27/11/2025.
//

#include "kdf.h"

// doc: https://docs.openssl.org/1.0.2/man3/


// vamos implementar kdf utilizando HASH
int kdf_aes_256(const BIGNUM* shared_key,char* result)
{
    // buffer temporário para copiar os bytes da shared_key
    int key_len = BN_num_bytes(shared_key);
    unsigned char* key_bytes = malloc(key_len);
    if (!key_bytes) return EXIT_FAILURE;

    BN_bn2bin(shared_key, key_bytes);
    SHA256_CTX ctx;
    if (!SHA256_Init(&ctx)) {
        free(key_bytes);
        return EXIT_FAILURE;
    }

    if (!SHA256_Update(&ctx, key_bytes, key_len)) {
        free(key_bytes);
        return EXIT_FAILURE;
    }

    if (!SHA256_Final((unsigned char*)result, &ctx)) {
        free(key_bytes);
        return EXIT_FAILURE;
    }

    free(key_bytes);
    return EXIT_SUCCESS;
}

