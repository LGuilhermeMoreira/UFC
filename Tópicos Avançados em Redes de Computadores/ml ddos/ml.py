import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, OneHotEncoder
from sklearn.tree import DecisionTreeClassifier
from sklearn.metrics import confusion_matrix, classification_report, accuracy_score
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
from sklearn.tree import plot_tree

# --- ETAPA 1 e 2: Carregamento e Pré-processamento dos Dados ---

# Tente carregar o dataset.
try:
    df = pd.read_csv("DDoS_dataset.csv")
    print("Dataset carregado com sucesso!")
    print(f"O dataset tem {df.shape[0]} linhas e {df.shape[1]} colunas.")
except FileNotFoundError:
    print(
        "Erro: Arquivo 'DDoS_dataset.csv' não encontrado. Certifique-se de que ele está na mesma pasta que o seu script."
    )
    exit()

# Para este primeiro modelo, vamos simplificar e focar nas features mais diretas.
# IPs podem ser complexos de tratar (existem milhões). Vamos removê-los por enquanto
# para criar um modelo base sólido e fácil de explicar.
df = df.drop(["Source IP", "Dest IP"], axis=1)
print("Colunas de IP removidas para simplificar o modelo inicial.")

# Separar as features (X) do nosso alvo (y)
X = df.drop("target", axis=1)
y = df["target"]

# Identificar quais colunas são categóricas e quais são numéricas
categorical_features = ["Highest Layer", "Transport Layer"]
numeric_features = ["Source Port", "Dest Port", "Packet Length", "Packets/Time"]

# Criar um 'pipeline' de pré-processamento. Isso organiza os passos de transformação.
# 1. Para features numéricas: Vamos padronizar a escala (StandardScaler).
# 2. Para features categóricas: Vamos convertê-las em colunas numéricas (OneHotEncoder).
preprocessor = ColumnTransformer(
    transformers=[
        ("num", StandardScaler(), numeric_features),
        ("cat", OneHotEncoder(handle_unknown="ignore"), categorical_features),
    ]
)

# Dividir os dados em conjunto de TREINO e conjunto de TESTE.
# 80% para treinar, 20% para testar. random_state garante que a divisão seja sempre a mesma.
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42, stratify=y
)
print(f"Dados divididos em {len(X_train)} amostras de treino e {len(X_test)} de teste.")


# --- ETAPA 3: Definição e Treinamento do Modelo ---

# Criar o pipeline final que inclui o pré-processamento E o modelo.
# O modelo escolhido é a Árvore de Decisão (DecisionTreeClassifier).
# max_depth=5 limita a profundidade da árvore para evitar que ela fique complexa demais e para facilitar a visualização.
model_pipeline = Pipeline(
    steps=[
        ("preprocessor", preprocessor),
        ("classifier", DecisionTreeClassifier(max_depth=5, random_state=42)),
    ]
)

# Treinar o modelo! O comando .fit() é onde a "mágica" do aprendizado acontece.
print("\nIniciando o treinamento do modelo...")
model_pipeline.fit(X_train, y_train)
print("Modelo treinado com sucesso!")


# --- ETAPA 4: Avaliação do Modelo ---

# Agora, vamos usar o modelo treinado para fazer previsões no conjunto de teste.
print("\nRealizando previsões no conjunto de teste...")
y_pred = model_pipeline.predict(X_test)

# Medir a acurácia
accuracy = accuracy_score(y_test, y_pred)
print(f"\nAcurácia do Modelo: {accuracy * 100:.2f}%")

# Imprimir o Relatório de Classificação (com Precisão e Recall)
# - Precisão: Das vezes que o modelo disse "Ataque", quantas eram de fato um ataque? (Evita falsos alarmes)
# - Recall: De todos os ataques reais, quantos o modelo conseguiu pegar? (Mede a capacidade de detecção)
print("\nRelatório de Classificação:")
print(classification_report(y_test, y_pred, target_names=["Normal (0)", "DDoS (1)"]))

# Gerar e exibir a Matriz de Confusão
print("Gerando a Matriz de Confusão...")
cm = confusion_matrix(y_test, y_pred)
plt.figure(figsize=(8, 6))
sns.heatmap(
    cm,
    annot=True,
    fmt="d",
    cmap="Blues",
    xticklabels=["Previsto Normal", "Previsto DDoS"],
    yticklabels=["Real Normal", "Real DDoS"],
)
plt.title("Matriz de Confusão")
plt.ylabel("Verdadeiro")
plt.xlabel("Previsto")
plt.show()  # Exibe o gráfico


# --- ETAPA 5: Visualização da Árvore de Decisão ---

# Este é um ótimo gráfico para sua apresentação!
print("\nGerando a visualização da Árvore de Decisão...")

# Pegar os nomes das features após o OneHotEncoding para usar no gráfico
feature_names = numeric_features + list(
    model_pipeline.named_steps["preprocessor"]
    .named_transformers_["cat"]
    .get_feature_names_out(categorical_features)
)

plt.figure(figsize=(20, 10))
plot_tree(
    model_pipeline.named_steps["classifier"],
    feature_names=feature_names,
    class_names=["Normal", "DDoS"],
    filled=True,
    rounded=True,
    fontsize=10,
)
plt.title("Visualização da Árvore de Decisão")
plt.show()  # Exibe a árvore
