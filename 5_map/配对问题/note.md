遍历 `sorted(counter)` ，配对元素减去当前需要的元素数

```Python
for key in sorted(counter):
    # 配对元素不足
    if counter[key + 2 * k] < counter[key]:
        break
    res.extend([key + k] * counter[key])
    counter[key + 2 * k] -= counter[key]
else:
    return res
```

困难题多了一个枚举 k
`5966. 还原原数组-配对问题`
