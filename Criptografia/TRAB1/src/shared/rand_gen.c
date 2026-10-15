#include <stdlib.h>
#include <openssl/rand.h>
#include <openssl/aes.h>
#include "rand.h"

int rand_gen(unsigned char *buf) {
    if (RAND_bytes(buf, AES_BLOCK_SIZE) != 1) {
        return EXIT_FAILURE;
    }
    return EXIT_SUCCESS;
}