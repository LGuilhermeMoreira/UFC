//
// Created by guigui on 03/12/2025.
//

// doc: https://docs.openssl.org/1.0.2/man3/

#include <stdlib.h>
#include <string.h>
#include <openssl/evp.h>
#include <openssl/aes.h>
#include "aes.h"

// lógica:
/*
 * 1. Fazer o padding
 *
 * 2. Para o primeiro bloco (i = 0):
 *       temp = P0^IV
 *       C0 = encrypt(temp,key)
 *
 * 3. Para cada próximo bloco (i >= 1):
 *       temp = Pi^C(i-1)
 *       Ci = encrypt(temp,key)
 */
int aes_encrypt(const unsigned char *plaintext,
    const int plaintext_len,
    unsigned char **ciphertext,
    int *ciphertext_len,
    const unsigned char *key,
    const unsigned char *iv
    )
{
    EVP_CIPHER_CTX *ctx = NULL;
    unsigned char *temp_plaintext = NULL;

    temp_plaintext = malloc(plaintext_len);
    if (!temp_plaintext) return EXIT_FAILURE;

    int temp_plaintext_len = plaintext_len;
    memcpy(temp_plaintext, plaintext, plaintext_len);

    // fazendp padding
    if (plaintext_len % AES_BLOCK_SIZE != 0)
    {
        int padded_len = AES_BLOCK_SIZE - (plaintext_len % AES_BLOCK_SIZE);
        unsigned char *new_ptr = realloc(temp_plaintext, plaintext_len + padded_len);

        if (!new_ptr)
        {
            free(temp_plaintext);
            return EXIT_FAILURE;
        }

        temp_plaintext = new_ptr;
        temp_plaintext_len = plaintext_len + padded_len;

        for (int i = plaintext_len; i < temp_plaintext_len; i++)
        {
            temp_plaintext[i] = padded_len;
        }
    }else
    {
        int padded_len = AES_BLOCK_SIZE;
        unsigned char *new_ptr = realloc(temp_plaintext, plaintext_len + padded_len);

        if (!new_ptr)
        {
            free(temp_plaintext);
            return EXIT_FAILURE;
        }

        temp_plaintext = new_ptr;
        temp_plaintext_len = plaintext_len + padded_len;

        for (int i = plaintext_len; i < temp_plaintext_len; i++)
        {
            temp_plaintext[i] = padded_len;
        }
    }

    // alocando memória para ciphertext
    *ciphertext_len = temp_plaintext_len;
    *ciphertext = malloc(temp_plaintext_len);
    if (!*ciphertext)
    {
        free(temp_plaintext);
        return EXIT_FAILURE;
    }

    // criando buffers
    unsigned char temp[AES_BLOCK_SIZE];
    unsigned char enc_block[AES_BLOCK_SIZE];
    unsigned char prev_block[AES_BLOCK_SIZE];
    memcpy(prev_block, iv, AES_BLOCK_SIZE);

    // inicializando contexto
    ctx = EVP_CIPHER_CTX_new();
    if (!ctx)
    {
        free(temp_plaintext);
        free(*ciphertext);
        return EXIT_FAILURE;
    }

    EVP_EncryptInit_ex(ctx, EVP_aes_256_ecb(), NULL, key, NULL);
    EVP_CIPHER_CTX_set_padding(ctx, 0);

    // serve para escrever o tamanho de saida do ecb, mas como não
    int aux = 0;

    for (int i = 0; i < temp_plaintext_len; i += AES_BLOCK_SIZE)
    {
        for (int j = 0; j < AES_BLOCK_SIZE; j++)
            temp[j] = prev_block[j] ^ temp_plaintext[i + j];

        EVP_EncryptUpdate(ctx, enc_block, &aux, temp, AES_BLOCK_SIZE);

        memcpy(&(*ciphertext)[i], enc_block, AES_BLOCK_SIZE);
        memcpy(prev_block, enc_block, AES_BLOCK_SIZE);
    }

    EVP_CIPHER_CTX_free(ctx);
    free(temp_plaintext);

    return EXIT_SUCCESS;
}


int aes_decrypt(const unsigned char *ciphertext,
    int ciphertext_len,
    unsigned char **plaintext,
    int *plaintext_len,
    const unsigned char *key,
    const unsigned char *iv
    )
{
    EVP_CIPHER_CTX *ctx = NULL;

    // aloca plaintext
    *plaintext = malloc(ciphertext_len);
    if (*plaintext == NULL)
        return EXIT_FAILURE;

    *plaintext_len = ciphertext_len;

    unsigned char temp[AES_BLOCK_SIZE];
    unsigned char dec_block[AES_BLOCK_SIZE];
    unsigned char prev_block[AES_BLOCK_SIZE];

    memcpy(prev_block, iv, AES_BLOCK_SIZE);

    // cria contexto
    ctx = EVP_CIPHER_CTX_new();
    EVP_DecryptInit_ex(ctx, EVP_aes_256_ecb(), NULL, key, NULL);
    EVP_CIPHER_CTX_set_padding(ctx, 0);

    int aux = 0;

    for (int i = 0; i < ciphertext_len; i += AES_BLOCK_SIZE)
    {
        // decripta o bloco atual
        EVP_DecryptUpdate(ctx, dec_block, &aux, ciphertext + i, AES_BLOCK_SIZE);

        // faz xor com bloco anterior
        for (int j = 0; j < AES_BLOCK_SIZE; j++)
            temp[j] = dec_block[j] ^ prev_block[j];

        memcpy(*plaintext + i, temp, AES_BLOCK_SIZE);
        memcpy(prev_block, ciphertext + i, AES_BLOCK_SIZE);
    }

    EVP_CIPHER_CTX_free(ctx);

    // removendo padding
    int pad = (*plaintext)[(*plaintext_len) - 1];

    if (pad <= 0 || pad > AES_BLOCK_SIZE)
    {
        // padding inválido → corrompido
        free(*plaintext);
        return EXIT_FAILURE;
    }

    // reduzindo tamanho do plaintext
    *plaintext_len -= pad;

    return EXIT_SUCCESS;
}


