https://www.pypandas.cn/docs/getting_started/overview.html#%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84
https://github.com/zhouyanasd/or-pandas/blob/master/articles/Pandas%E6%95%99%E7%A8%8B_01%E6%95%B0%E6%8D%AE%E5%88%86%E6%9E%90%E5%85%A5%E9%97%A8.md

Pandas 的主要数据结构是 Series（一维数据）与 DataFrame（二维数据）
Series:带标签的一维同构数组
DataFrame 带标签的，大小可变的，二维异构表格

```py
for col in df.columns:
    series = df[col]
    # do something with series
```

索引 / 选择
索引基础用法如下：

| 操作             | 句法              | 结果             |
| ---------------- | ----------------- | ---------------- |
| 选择列           | df[col]/df[[col]] | Series/DataFrame |
| 用标签选择行     | df.loc[label]     | Series           |
| 用整数位置选择行 | df.iloc[loc]      | Series           |
| 行切片           | df[5:10]          | DataFrame        |
| 用布尔向量选择行 | df[bool_vec]      | DataFrame        |
| 选择多行多列     | df.iloc[0:2,0:2]  | DataFrame        |
| 选择单元格       | df.at[row,col]    | 值               |
