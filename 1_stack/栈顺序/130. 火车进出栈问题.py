# 每节车厢有两种运动方式，进栈与出栈，问 n 节车厢出栈的可能排列方式有多少种。
# 1≤n≤60000
import math

n = int(input())
A = math.factorial(2 * n)
B = math.factorial(n)
print(A // B // B // (n + 1))
# 卡特兰数经典应用，C(2n,n)/(n+1) 即为答案
# comb(2 * n, n) // (n + 1)
