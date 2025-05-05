https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/math_fwt.go
https://maspypy.github.io/library/setfunc/hadamard.hpp

把卷积变成点乘?

---

FWT（快速沃尔什-哈达玛变换）是一种高效计算**位运算卷积**（XOR/AND/OR卷积）的算法，复杂度O(n log n)，常用于子集相关的计数、DP优化等。

**核心思想**：  
把序列变换到新空间，**卷积变为点乘**，逆变换回原空间，类似FFT。

---

## 典型用法：异或卷积

**问题**：给定两个长度为2^k的数组A、B，求C[i]=∑A[j]\*B[k]，其中j^k=i。

**步骤**：

1. 对A、B分别FWT变换
2. 对应元素相乘
3. 逆FWT变换

**代码（XOR卷积）**：

```cpp
void FWT(int *a, int n, int op) {
    for (int d = 1; d < n; d <<= 1)
        for (int i = 0; i < n; i += d << 1)
            for (int j = 0; j < d; ++j) {
                int x = a[i + j], y = a[i + j + d];
                a[i + j] = x + y;
                a[i + j + d] = x - y;
                if (op == -1) { // 逆变换
                    a[i + j] /= 2;
                    a[i + j + d] /= 2;
                }
            }
}
```

**使用**：

```cpp
FWT(A, n, 1);
FWT(B, n, 1);
for (int i = 0; i < n; ++i) C[i] = A[i] * B[i];
FWT(C, n, -1);
```

---

## 总结

- FWT适合处理**子集卷积**、**位运算相关计数**等问题
- 常见变种：XOR、AND、OR卷积
- 步骤：FWT变换 → 点乘 → 逆FWT

**一句话**：FWT就是位运算卷积的FFT！
