# 在[A,B]间选两个不同的数x<y 使得gcd(x,y)最大 求最大公约数的最大值
# 1≤A<B≤2×1e5
# !性质:最大公约数不可能超过(B-A)

A, B = map(int, input().split())
diff = B - A
res = 1
for f in range(diff, 0, -1):
    if B // f - (A - 1) // f > 1:
        print(f)
        exit(0)
