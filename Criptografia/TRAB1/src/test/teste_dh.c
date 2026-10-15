#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <openssl/bn.h>
#include "dhke.h"

int main(int argc, char const *argv[])
{
    BIGNUM* P = BN_new();
    BIGNUM* G = BN_new();

    BIGNUM* a = BN_new();

    BIGNUM* b = BN_new();

    BN_set_word(G, 2);

    printf("=== Iniciando Teste DHKE ===\n");
    printf("Gerando primo P (2048 bits)... isso pode demorar um pouco...\n");

    if (!generate_large_prime(P, 2048)) {
        fprintf(stderr, "Erro: Falha ao gerar primo P\n");
        return EXIT_FAILURE;
    }

    char *p_str = BN_bn2dec(P);
    printf("P gerado (inicio): %.20s...\n", p_str); // Mostra só o começo para não poluir
    OPENSSL_free(p_str);

    if (!generate_private_key(a, P)) {
        fprintf(stderr, "Erro: Falha ao gerar chave privada de Alice\n");
        return EXIT_FAILURE;
    }

    BIGNUM* A = pow_mod(G, a, P);
    if (!A) {
        fprintf(stderr, "Erro: Falha ao calcular chave pública de Alice (A)\n");
        return EXIT_FAILURE;
    }

    if (!generate_private_key(b, P)) {
        fprintf(stderr, "Erro: Falha ao gerar chave privada de Bob\n");
        return EXIT_FAILURE;
    }

    BIGNUM* B = pow_mod(G, b, P);
    if (!B) {
        fprintf(stderr, "Erro: Falha ao calcular chave pública de Bob (B)\n");
        return EXIT_FAILURE;
    }

    printf("\n--- Chaves Publicas Geradas ---\n");
    char *A_str = BN_bn2dec(A);
    char *B_str = BN_bn2dec(B);
    printf("Publica Alice (A): %.20s...\n", A_str);
    printf("Publica Bob   (B): %.20s...\n", B_str);
    OPENSSL_free(A_str);
    OPENSSL_free(B_str);

    printf("\n--- Calculando Segredos Compartilhados ---\n");

    BIGNUM* secret_alice = pow_mod(B, a, P);

    BIGNUM* secret_bob = pow_mod(A, b, P);

    if (!secret_alice || !secret_bob) {
        fprintf(stderr, "Erro: Falha no cálculo dos segredos\n");
        return EXIT_FAILURE;
    }

    if (BN_cmp(secret_alice, secret_bob) == 0) {
        printf("\n[SUCESSO] O protocolo funcionou!\n");

        char *s_str = BN_bn2dec(secret_alice);
        printf("Segredo Compartilhado: %.30s...\n", s_str);
        OPENSSL_free(s_str);
    } else {
        printf("\n[FALHA] Os segredos sao diferentes!\n");

        char *sa_str = BN_bn2dec(secret_alice);
        char *sb_str = BN_bn2dec(secret_bob);
        printf("Segredo Alice: %s\n", sa_str);
        printf("Segredo Bob:   %s\n", sb_str);
        OPENSSL_free(sa_str);
        OPENSSL_free(sb_str);
    }

    // --- LIMPEZA DE MEMÓRIA ---
    BN_free(P);
    BN_free(G);
    BN_free(a);
    BN_free(b);
    BN_free(A);
    BN_free(B);
    BN_free(secret_alice);
    BN_free(secret_bob);

    return EXIT_SUCCESS;
}