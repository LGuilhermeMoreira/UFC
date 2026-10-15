#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <openssl/bn.h>
#include <openssl/aes.h>

#include "dhke.h"
#include "kdf.h"
#include "constants.h"
#include "aes.h"
#include "rand.h"

int main()
{
    int server_fd, client_fd;
    struct sockaddr_in server_addr, client_addr;
    char buffer[BUFFER_SIZE];
    socklen_t addr_len = sizeof(client_addr);

    server_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (server_fd < 0)
    {
        perror("socket");
        exit(1);
    }

    int opt = 1;
    setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));

    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    if (bind(server_fd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0)
    {
        perror("bind");
        exit(1);
    }

    if (listen(server_fd, 1) < 0)
    {
        perror("listen");
        exit(1);
    }

    client_fd = accept(server_fd, (struct sockaddr *)&client_addr, &addr_len);
    if (client_fd < 0)
    {
        perror("accept");
        exit(1);
    }

    printf("--- Iniciando Troca de Chaves DHKE ---\n");
    BIGNUM *P = BN_new();
    BIGNUM *G = BN_new();
    BIGNUM *A = BN_new();
    BIGNUM *b = BN_new();

    memset(buffer, 0, BUFFER_SIZE);
    recv(client_fd, buffer, BUFFER_SIZE, 0);
    BN_hex2bn(&P, buffer);
    send(client_fd, "ACK", 3, 0);

    memset(buffer, 0, BUFFER_SIZE);
    recv(client_fd, buffer, BUFFER_SIZE, 0);
    BN_hex2bn(&G, buffer);
    send(client_fd, "ACK", 3, 0);

    memset(buffer, 0, BUFFER_SIZE);
    recv(client_fd, buffer, BUFFER_SIZE, 0);
    BN_hex2bn(&A, buffer);

    generate_private_key(b, P);
    BIGNUM *B = pow_mod(G, b, P);

    char *b_str = BN_bn2hex(B);
    send(client_fd, b_str, strlen(b_str), 0);
    OPENSSL_free(b_str);

    BIGNUM *shared_secret = pow_mod(A, b, P);

    printf("--- Troca de chaves concluída ---\n\n");

    // printf("A: %s\n\nB: %s\n\nG: %s\n\nShared Key: %s\n\n", BN_bn2hex(A), BN_bn2hex(B), BN_bn2hex(G), BN_bn2hex(shared_secret));
    printf("Segredo compartilhado: %s\n", BN_bn2hex(shared_secret));

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
        // recebe o iv do cliente
        unsigned char iv[AES_BLOCK_SIZE];
        ssize_t iv_bytes = recv(client_fd, iv, AES_BLOCK_SIZE, 0);

        if (iv_bytes <= 0)
        {
            printf("Conexão encerrada.\n");
            continue;
        }
        if (iv_bytes != AES_BLOCK_SIZE)
        {
            printf("Iv veio todo bichado.\n");
            continue;
        }

        // recebe o texto cifrado do cliente
        unsigned char cipher_in[BUFFER_SIZE];
        ssize_t cbytes = recv(client_fd, cipher_in, sizeof(cipher_in), 0);

        if (cbytes <= 0)
        {
            printf("Cliente desconectou.\n");
            continue;
        }

        printf("Texto cifrado recebido: ");
        for(int i = 0; i < cbytes; i++) printf("%02x", cipher_in[i]);
        printf("\n");

        unsigned char *plain_text = NULL;
        int plain_len = 0;

        // decifra a mensagem
        if (aes_decrypt(cipher_in, (int)cbytes,
                        &plain_text, &plain_len,
                        (unsigned char *)derived_key, iv) != 0)
        {
            printf("Deu zebra para descriptografar.\n");
            free(plain_text);
            continue;
        }

        char *msg = malloc(plain_len + 1);
        memcpy(msg, plain_text, plain_len);
        msg[plain_len] = '\0';

        printf("Cliente disse: %s\n", msg);

        if (strcmp(msg, "exit") == 0)
        {
            free(msg);
            free(plain_text);
            continue;
        }

        free(plain_text);
        free(msg);

        char server_input[BUFFER_SIZE];
        printf("Sua resposta: ");
        if (fgets(server_input, sizeof(server_input), stdin) == NULL)
        {
            continue;
        }
        server_input[strcspn(server_input, "\n")] = 0;

#ifdef AES_MODE_CTR
        // gera novo nouce
        if (rand_gen(iv) == EXIT_FAILURE)
        {
            printf("Erro ao gerar novo IV/Nonce.\n");
            continue;
        }
        send(client_fd, iv, AES_BLOCK_SIZE, 0);
#endif

        unsigned char *cipher_out = NULL;
        int cipher_out_len = 0;

        // encriptando a mensagem do servidor
        if (aes_encrypt((unsigned char *)server_input, strlen(server_input),
                        &cipher_out, &cipher_out_len,
                        (unsigned char *)derived_key, iv) != 0)
        {
            printf("Deu zebra para criptografia.\n");
            free(cipher_out);
            continue;
        }

        printf("Texto cifrado a enviar: ");
        for(int i = 0; i < cipher_out_len; i++) printf("%02x", cipher_out[i]);
        printf("\n");
        // 6) Enviar ciphertext de volta
        ssize_t sent = send(client_fd, cipher_out, cipher_out_len, 0);
        if (sent != cipher_out_len)
        {
            printf("Deu zebra ao enviar resposta criptografada.\n");
            free(cipher_out);
            continue;
        }

        free(cipher_out);
    }

    free(derived_key);
    BN_free(P);
    BN_free(G);
    BN_free(A);
    BN_free(B);
    BN_free(b);
    BN_free(shared_secret);
    close(client_fd);
    close(server_fd);

    return 0;
}