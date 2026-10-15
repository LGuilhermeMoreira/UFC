cmake_minimum_required(VERSION 3.10)
project(TRAB1)

set(CMAKE_C_STANDARD 11)

set(IMPL "CTR" CACHE STRING "Escolha entra esses modos de operação: CTR ou CBC")

find_package(OpenSSL REQUIRED)

include_directories(include ${OPENSSL_INCLUDE_DIR})

if(IMPL STREQUAL "CTR")
    add_compile_definitions(AES_MODE_CTR)
    set(AES_IMPL
            src/shared/aes_ctr.c
    )
elseif(IMPL STREQUAL "CBC")
    add_compile_definitions(AES_MODE_CBC)
    set(AES_IMPL
            src/shared/aes_cbc.c
    )
else()
    message(FATAL_ERROR "Valor de implementação inválido: ${IMPL}. Possíveis valores são CTR, CBC")
endif()

add_library(shared_lib STATIC
        src/shared/dhke.c
        src/shared/kdf.c
        src/shared/constants.c
        src/shared/rand_gen.c
        ${AES_IMPL}
)

add_library(lib_cbc STATIC
        src/shared/dhke.c
        src/shared/kdf.c
        src/shared/constants.c
        src/shared/rand_gen.c
        src/shared/aes_cbc.c
)

add_library(lib_ctr STATIC
        src/shared/dhke.c
        src/shared/kdf.c
        src/shared/constants.c
        src/shared/rand_gen.c
        src/shared/aes_ctr.c
)

target_link_libraries(shared_lib OpenSSL::Crypto)

add_executable(servidor src/server/main.c)
target_link_libraries(servidor shared_lib OpenSSL::Crypto)

add_executable(cliente src/client/main.c)
target_link_libraries(cliente shared_lib OpenSSL::Crypto)

add_executable(teste_dh src/test/teste_dh.c)
target_link_libraries(teste_dh shared_lib OpenSSL::Crypto)

add_executable(teste_kdf src/test/teste_kdf.c)
target_link_libraries(teste_kdf shared_lib OpenSSL::Crypto)

add_executable(teste_aes_cbc src/test/teste_aes_cbc.c)
target_link_libraries(teste_aes_cbc lib_cbc OpenSSL::Crypto)

add_executable(teste_rand src/test/teste_rand.c)
target_link_libraries(teste_rand shared_lib OpenSSL::Crypto)

add_executable(teste_aes_ctr src/test/teste_aes_ctr.c)
target_link_libraries(teste_aes_ctr lib_ctr OpenSSL::Crypto)