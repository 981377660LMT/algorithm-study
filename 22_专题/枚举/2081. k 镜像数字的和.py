# 一个 k 镜像数字 指的是一个在十进制和 k 进制下从前往后读和从后往前读都一样的 没有前导 0 的 正 整数。
# 2 <= k <= 9
# 1 <= n <= 30


palindrome = [1, 2, 3, 4, 5, 6, 7, 8, 9]
# 构造九位数的回文只需要左右部分最多四位 一共有11*10000+9个回文数
for side in range(1, 10000):
    s1 = str(side) + str(side)[::-1]
    palindrome.append(int(s1))
    for mid in range(10):
        s2 = str(side) + str(mid) + str(side)[::-1]
        palindrome.append(int(s2))
palindrome.sort()


class Solution:
    def kMirror(self, k: int, n: int) -> int:
        res = []
        for p in palindrome:
            cand = str(p ** 2)
            if cand == cand[::-1]:
                res.append(cand)
        return len(res)


print(Solution().kMirror(k=3, n=7))
# 输出：499
# 解释：
# 7 个最小的 3 镜像数字和它们的三进制表示如下：
#   十进制       三进制
#     1          1
#     2          2
#     4          11
#     8          22
#     121        11111
#     151        12121
#     212        21212
# 它们的和为 1 + 2 + 4 + 8 + 121 + 151 + 212 = 499 。
