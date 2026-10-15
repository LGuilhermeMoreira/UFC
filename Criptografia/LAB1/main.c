#include<stdio.h>
#include<string.h>
#include <stdlib.h>

const size_t ASCII_LEN = 256;

//link: https://danieldonda.com/a-matematica-da-cifra-de-vigenere/


//Ci=Pi+Ki (mod 256)
void encrypt(char *plain_text,char *cypher_text, char *key){
    size_t plain_text_len = strlen(plain_text);
    size_t key_len = strlen(key);

    if(key_len == 0){
        strcpy(cypher_text,plain_text);
        return;
    }

    for (size_t i = 0; i < plain_text_len; i++)
    {
        size_t key_value = key[i % key_len];
        cypher_text[i] = (plain_text[i] + key_value) % ASCII_LEN; 
    }
    
    cypher_text[plain_text_len] = '\0';
}

//Pi=Ci-Ki+256 (mod 256)
void decrypt(char *cypher_text, char *plain_text,char *key){
    size_t cypher_text_len = strlen(cypher_text);
    size_t key_len = strlen(key);

    if(key_len == 0){
        strcpy(plain_text,cypher_text);
        return;
    }

    for (size_t i = 0; i < cypher_text_len; i++)
    {
        size_t key_value = key[i % key_len];
        plain_text[i] = (cypher_text[i] - key_value + ASCII_LEN) % ASCII_LEN;
    }
    
    plain_text[cypher_text_len] = '\0';
}

void run_tests(const char *filename) {
    FILE *file = fopen(filename, "r");
    if (file == NULL) {
        printf("Erro: Nao foi possivel abrir o arquivo '%s'\n", filename);
        return;
    }

    char line[512]; 
    int test_count = 1;

    printf("--- Iniciando Testes ---\n\n");

    while (fgets(line, sizeof(line), file)) {        
        line[strcspn(line, "\n")] = 0;

        char *plain_text_original = strtok(line, ",");
        char *key = strtok(NULL, ",");

        
        if (plain_text_original == NULL || key == NULL) {
            continue;
        }
        
        
        char *cypher_text = malloc(strlen(plain_text_original) + 1);
        char *decrypted_text = malloc(strlen(plain_text_original) + 1);

        
        encrypt(plain_text_original, cypher_text, key);
        decrypt(cypher_text, decrypted_text, key);

        printf("--- Teste %d ---\n", test_count++);
        printf("frase Original: '%s'\n", plain_text_original);
        printf("Key: '%s'\n", key);
        printf("Cifrado: '%s'\n", cypher_text);
        printf("Decifrado: '%s'\n", decrypted_text);

        
        if (strcmp(plain_text_original, decrypted_text) == 0) {
            printf("OK\n\n");
        } else {
            printf("NOT OK\n\n");
        }

        
        free(cypher_text);
        free(decrypted_text);
    }
    
    printf("--- FIM DOS TESTES ---\n");
    fclose(file);
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Uso: %s <nome_do_arquivo>\n", argv[0]);
        return 1;
    }

    char *nome_arquivo = argv[1];
    run_tests(nome_arquivo);
    return 0;
}

