一般解决 gcd 计数问题，三个步骤：

1. zeta：用 c1 表示 counter, c2 表示每个因子(作为 gcd)有多少个倍数，累加所有倍数的个数;
2. 根据题意，对 c2 进行变换(每个数作为 gcd 的答案);
3. moebius 反演：`容斥`，减去每个数的倍数的答案;

```python
upper = max(nums) + 1
c1, c2 = [0] * upper, [0] * upper  # !c2[i]表示gcd为i的二元组对数
for v in nums:
    c1[v] += 1
for f in range(1, upper):
    for m in range(f, upper, f):
        c2[f] += c1[m]
for i in range(1, upper):
    c2[i] = c2[i] * (c2[i] - 1) // 2
for f in range(upper - 1, 0, -1):
    for m in range(2 * f, upper, f):
        c2[f] -= c2[m]
```

https://aprilganmo.hatenablog.com/entry/2020/07/24/190816
