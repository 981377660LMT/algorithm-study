# 黒色・白色・灰色のカードが 1 枚ずつあります。

# 以下の条件のうち一つ以上を満たすように、各カードに 1 以上 N 以下の整数を書き込む方法が何通りあるかを求めてください。

# 黒色と白色のカードに書かれている整数の差の絶対値は K 以上
# 黒色と灰色のカードに書かれている整数の差の絶対値は K 以上
# 灰色と白色のカードに書かれている整数の差の絶対値は K 以上
# 1≤N≤100000
# 1≤K≤min(5,N−1)
# !减去都不满足的情况

n, k = map(int, input().split())

res = n ** 3
for a in range(1, n + 1):
    bl, br = max(1, a - k + 1), min(n, a + k - 1)
    for b in range(bl, br + 1):
        cl, cr = max(1, b - k + 1), min(n, b + k - 1)
        for c in range(cl, cr + 1):
            if abs(a - b) < k and abs(b - c) < k and abs(a - c) < k:
                res -= 1
print(res)

