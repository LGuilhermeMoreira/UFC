//
// Created by guigui on 11/12/2025.
//

#include <openssl/aes.h>
#include "rand.h"
#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
    unsigned char* iv = malloc(AES_BLOCK_SIZE);
    const int status = rand_gen(iv);
    if (status == EXIT_SUCCESS) {
        printf("iv: \n");
        for (int i = 0; i < AES_BLOCK_SIZE; i++) {
            printf("%02x", (unsigned char)iv[i]);
        }
        printf("\n");
    }else {
        printf("Deu zebra!\n");
    }

    free(iv);
}
