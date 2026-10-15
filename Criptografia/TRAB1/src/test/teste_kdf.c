//
// Created by guigui on 27/11/2025.
//

#include "kdf.h"
#include "dhke.h"

int main(int argc, char* argv[])
{
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

    for (int i = 0; i < 32; i++)
    {
        printf("%02x",(unsigned char)content[i]);
    }
    printf("\n");
    return EXIT_SUCCESS;
}
