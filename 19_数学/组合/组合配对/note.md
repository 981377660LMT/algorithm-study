**一遍遍历统计，保证只有后面能看到前面**
`边统计边存`

**全部存起来再统计**

1. 相等情况
2. 不相等只处理一次(默认 x<y)

```Python
# (题目有i<j的索引限制)
for a,b in product(counter,repeat = 2):
   # 1.相等
   if a == b:
     res += freq[a] * (freq[a] - 1) // 2

   # 注意这里防止重复计算的细节
   # 2.两种情况只算一次
   elif a < b and b in freq:
       res += freq[a]*freq[b]
```

```Python
# 3.不重复搜索 (题目有i<j<k的索引限制的话)
for n1, n2, n3 in product(counter.keys(), repeat=3):
   # 不重复搜素排列,非常关键!
   if not n1 <= n2 <= n3:
       continue
```

`1577. 数的平方等于两数乘积的方法数-组合计数.py`
`1711. 大餐计数.py`
