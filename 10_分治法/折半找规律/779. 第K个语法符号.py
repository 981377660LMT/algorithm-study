# 在第一行我们写上一个 0。接下来的每一行，将前一行中的0替换为01，1替换为10。
# 给定行数 N 和序数 K，返回第 N 行中第 K个字符。（K从1开始）


class Solution:
    def kthGrammar(self, n: int, k: int) -> int:
        if n == 1:
            return 0
        length = 1 << (n - 1)
        if k > length // 2:
            k -= length // 2
            return 0 if self.kthGrammar(n - 1, k) == 1 else 1
        else:
            return self.kthGrammar(n - 1, k)


print(Solution().kthGrammar(n=2, k=1))
# 解释:
# 第一行: 0
# 第二行: 01
# 第三行: 0110
# 第四行: 01101001
