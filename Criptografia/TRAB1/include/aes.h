//
// Created by guigui on 03/12/2025.
//

#ifndef TRAB1_AES_H
#define TRAB1_AES_H
int aes_encrypt(const unsigned char *plaintext,
    int plaintext_len,
    unsigned char **ciphertext,
    int *ciphertext_len,
    const unsigned char *key,
    const unsigned char *iv
    );

int aes_decrypt(const unsigned char *ciphertext,
    int ciphertext_len,
    unsigned char **plaintext,
    int *plaintext_len,
    const unsigned char *key,
    const unsigned char *iv
    );


#endif //TRAB1_AES_H