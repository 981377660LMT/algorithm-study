from typing import List

# 1 <= n <= 20
# 整数 1 在序列中只出现一次。
# 2 到 n 之间每个整数都恰好出现两次。
# 对于每个 2 到 n 之间的整数 i ，两个 i 之间出现的距离恰好为 i 。


# 总结：
# 有点像langford数列 => 回溯
# 1.剪枝：取数的时候，从大到小取
# 2.回溯带上bool返回值表示找到一个解 此时不必再向下回溯寻找，而是层层返回
class Solution:
    def constructDistancedSequence(self, n: int) -> List[int]:
        def bt(cur: int) -> bool:
            if cur >= 2 * n - 1:
                return True

            # 这个位置之前已经放过了
            if res[cur] != -1:
                return bt(cur + 1)

            # 字典序最大的序列:优先看大的
            for num in range(n, 0, -1):
                if num not in res:
                    gap = num if num > 1 else 0
                    # 可以填
                    if cur + gap < 2 * n - 1 and res[cur] == res[cur + gap] == -1:
                        res[cur] = res[cur + gap] = num
                        if bt(cur + 1):
                            return True
                        res[cur] = res[cur + gap] = -1

            return False

        res = [-1] * (2 * n - 1)
        bt(0)
        return res


print(Solution().constructDistancedSequence(n=3))
# 输出：[3,1,2,3,2]
# 解释：[2,3,2,1,3] 也是一个可行的序列，但是 [3,1,2,3,2] 是字典序最大的序列。
