# 题意:给定1≤L<R≤1e18，求一对(x,y)满足L<x ≤y<R使得gcd(x,y) = 1。

# 由于在题目约束下，素数间隔至多为K= 1500 左右，
# 所以直接从Ⅰ右边和R左边暴力找答案即可。时间复杂度O(K2 log R)。

# 时间复杂度O(K^2 logR)。

# prime gap 质数间隙
# 1e9以内，两个相邻素数距离不超过400

# Prime Gap:
# https://primes.utm.edu/notes/GapsTable.html
