# 3141. 最大海明距离(最大汉明距离)
#
# 给定一个数组 nums 和一个整数 m，每个元素 nums[i] 满足 0 <= nums[i] < 2^m，
# !返回数组 answer，每个元素 answer[i] 表示 nums[i] 和数组中其它任何元素 nums[j] 的最大 海明距离。
# m<=17
#
# 我们把每个数视作一个二进制点，例如6表示110号点，改变其中一个二进制位视为一条边，如110号到111号有一条边，
# !则原题变为找到每个nums[i]对应的反码在与nums[j]的最短路径长度，如100的反码011与111最短距离为1，因此100的最大汉明距离点为111，最大距离为3-1=2。
# 最短路径采用01bfs算法，即记录所有可到达点及其距离，每次遍历时找到当前未遍历过的最近可达点，将其连接的点更新至集合中，重复该步骤至集合为空。
#
# !异或最大 => bfs求与反码的最短路


from collections import deque
from typing import List

INF = int(1e18)


class Solution:
    def maxHammingDistances(self, nums: List[int], m: int) -> List[int]:
        bit = max(nums).bit_length()
        dist = [INF] * (1 << bit)
        queue = deque()
        for num in nums:
            queue.append(num)
            dist[num] = 0

        while queue:
            cur = queue.popleft()
            for b in range(bit):
                next_ = cur ^ (1 << b)
                cand = dist[cur] + 1
                if cand < dist[next_]:
                    dist[next_] = cand
                    queue.append(next_)

        mask = (1 << bit) - 1
        return [bit - dist[mask ^ v] for v in nums]


# nums = [9,12,9,11], m = 4
if __name__ == "__main__":
    print(Solution().maxHammingDistances([9, 12, 9, 11], 4))
