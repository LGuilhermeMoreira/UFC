import pandas as pd
dataframe = pd.read_csv('veiculos.csv')
print(dataframe)
try:
    dataframe2 = pd.read_excel('veiculos.xlsx')
    print(dataframe2)
except OSError:
    print(OSError)

