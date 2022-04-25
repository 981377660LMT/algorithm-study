# 给你一个字符串 num
# 其中，num 表示一个很大的整数，字符串中的每个字符依次对应整数上的各个 数位 。
# 你可以交换这个整数相邻数位的数字 最多 k 次。
# 请你返回你能得到的最小整数，并以字符串形式返回。

# 1 <= num.length <= 30000
# 1 <= k <= 10^9

# 首先要使首位尽可能小，从 0 到 9 遍历对应的位置列表，
# 找到第一个距离小于等于 k 的位置 idx，即应该交换到首位。
# 剩下的同理。显然每轮只需要 0 到 9 最近的位置，
# 因此用 A 维护剩下的 0 到 9 分别对应的有序位置列表，每轮最多检查 10 个位置。

# 难点在于候选位置 idx 和当前位置 i 的距离是变化的。
# 当 idx 后有位置被交换时，idx 要后移一位。
# 因此需要求出 idx 后有多少个被交换的位置。

# 可以用树状数组，也可以通过深入分析证明：
# idx 后被交换的位置数量=所有已交换位置中小于 num[idx] 的数量 - num[:idx] 中小于 num[idx] 的数量，从而 O(N) 求解。
# 不过最简单粗暴的方法还是用 sortedcontainers.SortedList 维护所有已交换位置的有序集合，二分查找比 idx 大的数量即可。

# https://leetcode-cn.com/problems/minimum-possible-integer-after-at-most-k-adjacent-swaps-on-digits/comments/1156207
from collections import defaultdict, deque
from sortedcontainers import SortedList

# 输出可以包含前导 0 ，但输入保证不会有前导 0 。


class Solution:
    def minInteger(self, num: str, k: int) -> str:
        n = len(num)
        indexMap = defaultdict(deque)
        for index, char in enumerate(num):
            indexMap[char].append(index)

        res = []
        sl = SortedList()

        for i in range(n):
            for digit in range(10):
                indexes = indexMap[str(digit)]
                if not indexes:
                    continue

                pos = indexes[0]
                dist = pos - i  # 距离
                preCount = len(sl) - sl.bisect_right(pos)  # 原来在pos右边的数,换到前面了,看右边多少个数比他大
                cost = dist + preCount
                if cost <= k:
                    k -= cost
                    res.append(str(digit))
                    sl.add(indexes.popleft())
                    break

        return ''.join(res)


print(Solution().minInteger("4321", 4))
