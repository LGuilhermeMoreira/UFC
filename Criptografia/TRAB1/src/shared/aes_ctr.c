//
// Created by guigui on 05/01/2026.
//


#include <openssl/evp.h>
#include <string.h>
#include <stdlib.h>

// incrementa o nonce em 1
void increment_counter(unsigned char *counter) {
    for (int i = 15; i >= 0; i--) {
        if (++counter[i] != 0)
            break;
    }
}

int aes_ctr_process(const unsigned char *input,
                   int input_len,
                   unsigned char **output,
                   int *output_len,
                   const unsigned char *key,
                   const unsigned char *iv) {
    
    EVP_CIPHER_CTX *ctx;
    unsigned char counter[16];
    unsigned char keystream[16];
    int out_tmp_len;

    // aloca memória para saida
    *output = malloc(input_len);
    if (!*output) return 1;
    *output_len = input_len;

    memcpy(counter, iv, 16);

    //inicializa o contexto
    if (!(ctx = EVP_CIPHER_CTX_new())) return 1;

    for (int i = 0; i < input_len; i += 16) {
        // seta o contexto para utilizar aes_256_ecb sem padding (funciona como aes padrão)
        if (!EVP_EncryptInit_ex(ctx, EVP_aes_256_ecb(), NULL, key, NULL)) return 1;
        EVP_CIPHER_CTX_set_padding(ctx, 0); 
        
        if (!EVP_EncryptUpdate(ctx, keystream, &out_tmp_len, counter, 16)) return 1;

        // faz o xor entre a chave e o bloco de entrada
        int block_size = (input_len - i < 16) ? (input_len - i) : 16;
        for (int j = 0; j < block_size; j++) {
            (*output)[i + j] = input[i + j] ^ keystream[j];
        }

        // incrementa o contador para o proximo bloco
        increment_counter(counter);
    }

    EVP_CIPHER_CTX_free(ctx);
    return 0;
}

int aes_encrypt(const unsigned char *plaintext, int plaintext_len, unsigned char **ciphertext, int *ciphertext_len, const unsigned char *key, const unsigned char *iv) {
    return aes_ctr_process(plaintext, plaintext_len, ciphertext, ciphertext_len, key, iv);
}

int aes_decrypt(const unsigned char *ciphertext, int ciphertext_len, unsigned char **plaintext, int *plaintext_len, const unsigned char *key, const unsigned char *iv) {
    return aes_ctr_process(ciphertext, ciphertext_len, plaintext, plaintext_len, key, iv);
}