# 693. 交替位二进制数
class Solution:
    def hasAlternatingBits(self, n: int) -> bool:
        """
        位运算最优解：
        对于二进制交替位数 n，n ^ (n>>1) 会得到一串全 1 的二进制，比如
          n=10101(2)=21, n>>1=01010(2)=10, 异或结果 = 11111(2)=31。
        如果 x = n ^ (n>>1) 是全 1，那么 x+1 必然是 100000... 的形式，
        因此 x & (x+1) == 0。否则就有相邻两位相同。
        时间 O(1)，空间 O(1)。
        """
        x = n ^ (n >> 1)
        return (x & (x + 1)) == 0


if __name__ == "__main__":
    sol = Solution()
    tests = [
        (5, True),  # 101
        (7, False),  # 111
        (10, True),  # 1010
        (11, False),  # 1011
        (1, True),  # 1
        (3, False),  # 11
    ]
    for n, expect in tests:
        res = sol.hasAlternatingBits(n)
        print(f"{n} -> {res} (expected {expect})")
