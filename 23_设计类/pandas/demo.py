import pandas as pd
import numpy as np
import matplotlib.pyplot as plt

# 表格: DataFrame
# 行索引: index
# 列索引: columns

sheet = {
    "car": ["BMW", "Benz", "Audi"],
    "price": [300000, 400000, 500000],
    "country": ["Germany", "Germany", "Germany"],
}

df = pd.DataFrame(sheet)
print(df.columns)
print(df.index)
print(df.columns.tolist())
print(df.index.tolist())
print(df)
print(df.to_dict())

print(df["car"], type(df["car"]))  # dataframe取一维索引返回的是series
print(df[["car"]], type(df[["car"]]))  # dataframe取二维索引返回的是dataframe
print(df.info())
print(df.describe())
df.so


s = pd.Series([1, 3, 5, np.nan, 6, 8])
print(s)
