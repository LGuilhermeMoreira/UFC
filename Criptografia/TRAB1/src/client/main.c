#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <openssl/aes.h>
#include <openssl/bn.h>

#include "aes.h"
#include "kdf.h"
#include "dhke.h"
#include "constants.h"
#include "rand.h"

int main(int argc, char const *argv[])
{
    int sockfd;
    struct sockaddr_in server_addr;
    char buffer[BUFFER_SIZE];

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0)
    {
        perror("socket");
        exit(1);
    }

    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT);
    server_addr.sin_addr.s_addr = inet_addr("127.0.0.1");

    if (connect(sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
    {
        perror("connect");
        exit(1);
    }

    printf("--- Iniciando Troca de Chaves DHKE ---\n");

    BIGNUM *P = BN_new();
    BIGNUM *G = BN_new();
    BIGNUM *a = BN_new();
    BIGNUM *B = BN_new();

    BN_set_word(G, 2);
    generate_large_prime(P, DHKE_BIT_SIZE);
    generate_private_key(a, P);

    BIGNUM *A = pow_mod(G, a, P);

    char *p_str = BN_bn2hex(P);
    send(sockfd, p_str, strlen(p_str), 0);
    OPENSSL_free(p_str);

    recv(sockfd, buffer, BUFFER_SIZE, 0);

    char *g_str = BN_bn2hex(G);
    send(sockfd, g_str, strlen(g_str), 0);
    OPENSSL_free(g_str);

    recv(sockfd, buffer, BUFFER_SIZE, 0);

    char *a_pub_str = BN_bn2hex(A);
    send(sockfd, a_pub_str, strlen(a_pub_str), 0);
    OPENSSL_free(a_pub_str);

    memset(buffer, 0, BUFFER_SIZE);
    int bytes = recv(sockfd, buffer, BUFFER_SIZE, 0);
    if (bytes > 0)
    {
        BN_hex2bn(&B, buffer);
    }

    BIGNUM *shared_secret = pow_mod(B, a, P);

    printf("--- Troca de chaves concluída ---\n\n");
    printf("Segredo compartilhado: %s\n", BN_bn2hex(shared_secret));
    // printf("A: %s\n\nB: %s\n\nG: %s\n\nShared Key: %s\n\n",BN_bn2hex(A),BN_bn2hex(B),BN_bn2hex(G),BN_bn2hex(shared_secret));

    // chave do aes
    char *derived_key = malloc(DERIVED_KEY_SIZE);

    if (kdf_aes_256(shared_secret, derived_key))
    {
        printf("Deu zebra na função de derivação\n");
        return EXIT_FAILURE;
    }

    // for (int i = 0; i < DERIVED_KEY_SIZE; ++i)
    // {
    //     printf("%02x", (unsigned char)derived_key[i]);
    // }
    // printf("\n");

    while (1)
    {
        // unsigned char *iv = malloc(AES_BLOCK_SIZE);
        // if (!iv)
        // {
        //     perror("malloc iv");
        //     continue;
        // }
        // gerando iv
        unsigned char iv[AES_BLOCK_SIZE];

        if (rand_gen(iv) == EXIT_FAILURE)
        {
            printf("Deu zebra pra gerar o iv\n");
            continue;
        }

        // envio o iv
        ssize_t s = send(sockfd, iv, AES_BLOCK_SIZE, 0);
        if (s != AES_BLOCK_SIZE)
        {
            printf("Deu zebra ao enviar o iv\n");
            continue;
        }

        // pega a mensagem
        char plain_input[BUFFER_SIZE];
        printf("Digite uma mensagem: ");
        if (!fgets(plain_input, sizeof(plain_input), stdin))
        {
            printf("deu zebra ao ler a mensagem\n");
            continue;
        }
        plain_input[strcspn(plain_input, "\n")] = 0;
        if (strcmp(plain_input, "exit") == 0)
        {
            printf("Encerrando conexão.\n");
            continue;
        }

        // encripta a mensagem
        unsigned char *cipher_text = NULL;
        int cipher_text_len = 0;
        if (aes_encrypt((unsigned char *)plain_input,
                        (int)strlen(plain_input),
                        &cipher_text,
                        &cipher_text_len,
                        (unsigned char *)derived_key,
                        iv) != 0)
        {
            printf("Deu zebra ao cifrar a mensagem\n");
            free(cipher_text);
            continue;
        }

        printf("Texto cifrado: ");
        for(int i = 0; i < cipher_text_len; i++) printf("%02x", cipher_text[i]);
        printf("\n");

        // envio mensagem encriptada
        ssize_t sent = send(sockfd, cipher_text, cipher_text_len, 0);
        if (sent != cipher_text_len)
        {
            printf("Deu zebra ao enviar o texto cifrado\n");
            free(cipher_text);
            continue;
        }

        free(cipher_text);

#ifdef AES_MODE_CTR
        // recebe o nouce do servidor
        ssize_t iv_recvd = recv(sockfd, iv, AES_BLOCK_SIZE, 0);
        if (iv_recvd != AES_BLOCK_SIZE)
        {
            printf("Erro ao receber o novo IV do servidor\n");
            continue;
        }
#endif

        // recebo mensagem encriptada de volta
        unsigned char recvbuf[BUFFER_SIZE];
        ssize_t n = recv(sockfd, recvbuf, sizeof(recvbuf), 0);
        if (n <= 0)
        {
            if (n == 0)
                printf("conexão fechada pelo servidor\n");
            else
                printf("deu zebra ao receber a resposta do servidor\n");
            continue;
        }

        printf("Texto cifrado recebido: ");
        for(int i = 0; i < n; i++) printf("%02x", recvbuf[i]);
        printf("\n");

        // decripto a mensagem
        unsigned char *decrypted = NULL;
        int decrypted_len = 0;
        if (aes_decrypt(recvbuf, (int)n, &decrypted, &decrypted_len, (unsigned char *)derived_key, iv) != 0)
        {
            fprintf(stderr, "erro em aes_decrypt\n");
            free(decrypted);
            continue;
        }

        // deixando a saida válida para o printf
        char *out = malloc(decrypted_len + 1);
        if (out)
        {
            memcpy(out, decrypted, decrypted_len);
            out[decrypted_len] = '\0';
            printf("Servidor respondeu: %s\n", out);
            free(out);
        }
        else
        {
            printf("Servidor respondeu (binário, %d bytes)\n", decrypted_len);
        }

        free(decrypted);
    }

    free(derived_key);
    memset(buffer, 0, BUFFER_SIZE);
    BN_free(P);
    BN_free(G);
    BN_free(a);
    BN_free(A);
    BN_free(B);
    BN_free(shared_secret);
    close(sockfd);
    return 0;
}