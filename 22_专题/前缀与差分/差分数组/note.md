差分数组的用处：`延迟的范围更新`

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

注意如果开差分数组会超空间限制，就用`哈希表`来存 (而且也更方便)
