set -e

echo "Qual modo de criptografia deseja utilizar?"
echo "1) CBC"
echo "2) CTR"
read -p "Escolha (1 ou 2): " choice

if [[ "$choice" == "1" ]]; then
    MODE="CBC"
elif [[ "$choice" == "2" ]]; then
    MODE="CTR"
else
    echo "Opção inválida! Use 1 para CBC ou 2 para CTR."
    exit 1
fi

echo "Modo selecionado: $MODE"

echo "[0/5] Removendo pasta build antiga, se existir..."
rm -rf build

echo "[1/5] Criando pasta build..."
mkdir -p build

echo "[2/5] Entrando na pasta build..."
cd build

echo "[3/5] Rodando CMake com modo $MODE..."
cmake -DIMPL=$MODE ..

echo "[4/5] Rodando Make..."
make

echo "[5/5] Build concluído com sucesso!"
