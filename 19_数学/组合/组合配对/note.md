1. 相等情况
2. 不相等只处理一次(默认 x<y)

```Python
# 1.相等
if a == b:
  res += freq[a] * (freq[a] - 1) // 2

# 注意这里防止重复计算的细节
# 2.两种情况只算一次
elif a < b and b in freq:
    res += freq[a]*freq[b]
```

`1577. 数的平方等于两数乘积的方法数-组合计数.py`
`1711. 大餐计数.py`
