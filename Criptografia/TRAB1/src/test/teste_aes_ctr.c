//
// Created by guigui on 05/01/2026.
//

#include <string.h>

#include "kdf.h"
#include "dhke.h"
#include "aes.h"
#include "openssl/aes.h"

int main(void)
{
    const char message[] = "Ola mundo 123";
    const int message_len = strlen(message);
    unsigned char iv[AES_BLOCK_SIZE] = {
        0x00,0x01,0x02,0x03,0x04,0x05,0x06,0x07,
        0x08,0x09,0x0a,0x0b,0x0c,0x0d,0x0e,0x0f
    };

    printf("Alocando memória para content...\n");
    char* content = malloc(32);
    if (!content) return EXIT_FAILURE;
    BIGNUM* P = BN_new();

    printf("Gerando o primo\n");
    generate_large_prime(P,2048);

    if (kdf_aes_256(P,content))
    {
        printf("Deu zebra na função de derivação\n");
        return EXIT_FAILURE;
    }

    char *cipher = NULL;
    int cipher_len = 0;
    if (aes_encrypt((unsigned char*)message,message_len,(unsigned char**)&cipher,&cipher_len,(unsigned char*)content,iv))
    {
        printf("Deu zebra na encriptação\n");
        return EXIT_FAILURE;
    }

    char *plain =NULL;
    int plain_len = 0;

    if (aes_decrypt((unsigned char*)cipher,cipher_len,(unsigned char**)&plain,&plain_len,(unsigned char*)content,iv))
    {
        printf("Deu zebra na decriptação\n");
        return EXIT_FAILURE;
    }

    printf("Menssagem original: %s\n",message);
    printf("Menssagem encriptada: %s\n",cipher);
    printf("Menssagem decriptada: %s\n",plain);

    return EXIT_SUCCESS;
}