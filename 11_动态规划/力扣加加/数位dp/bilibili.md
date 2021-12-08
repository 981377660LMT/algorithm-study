`res[L,R]=res[1,R]-res[1,L-1]`
数字大小维度
数字大小->数位字典序

**windy 数**
不含前导零且相邻两个数字之差至少为 2 的正整数被称为 windy 数
5 36 192 是 windy 数 10 21 不是

dfs 求解:

1. 现在枚举到了哪一位 `curPos`
2. 前面一位的数字是多少 `preNum`
3. 这一位可以填哪些数(布尔值等) `flag`
   ![图 1](../../../images/7a0f49718214cb49a19d84bab1252a2a5d2da3cd832f6e65bc469b50102b5c32.png)

dfs(curPos,preNum,flag)
dfs**注意@lru_cache**

**RoundNumber**
![图 2](../../../images/4ddf1f1c0f458d333420758f5e039fb89d57440611ac8e9c3c0962cd980d6cdc.png)  
dp[pos][zeros][ones]

**好字符串的数目**
![图 3](../../../images/4f0b4dba46d282b2fcb4ea8055d3b91426787e98e90a41fe1e4448e4d2b01c5f.png)
dp[pos][len]

总结模板
![图 4](../../../images/77c584ef6d009af7ba1e967023b4d5e1ce9d6cb6151fd1b3f1518d14b7f34cce.png)
