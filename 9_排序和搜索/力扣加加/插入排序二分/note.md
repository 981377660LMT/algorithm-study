```Python
d = []
for a in A:
    i = bisect.bisect_left(d, a)
    if d and i < len(d):
        d[i] = a
    else:
        d.append(a)
```

题目条件很明显:
限定了**i<j**
或者
**i>j**(此时需要反转数组)
在归并的 merge 阶段统计
