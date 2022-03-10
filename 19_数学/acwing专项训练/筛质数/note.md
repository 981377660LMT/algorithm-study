```Python
def getPrimes(upper: int) -> List[int]:
    """筛选出1-upper中的质数 nloglogn"""
    visited = [False] * (upper + 1)
    res = []
    for num in range(2, upper + 1):
        if visited[num]:
            continue
        res.append(num)
        for multi in range(num * num, upper + 1, num):
            visited[multi] = True

    return res
```
