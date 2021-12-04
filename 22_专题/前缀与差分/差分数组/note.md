```Python
res = [0] * (length + 1)
for left, right, delta in updates:
    res[left] += delta
    res[right + 1] -= delta
for i in range(1, length+1):
    res[i] += res[i - 1]
return res[:-1]
```

注意返回结果时，`差分数组最后一个数要去掉`
