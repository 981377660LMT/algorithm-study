# 映射，求最长上升子序列
# 思考：北岸的城市和南岸的城市是一对自变量与因变量的关系
# 因此，我们把北岸城市位置做一个排序，然后按照相应南岸城市的位置做一个最长上升子序列

# 不相交：上升子序列
from typing import List


class Solution:
    def friendly_city(self, cities: List[tuple]):
        cities.sort(key=lambda x: x[0])
        south = []
        for city in cities:
            south.append(city[1])

        def calcLis(arr) -> List[int]:
            dp = [1 for _ in range(len(arr))]
            for i in range(1, len(dp)):
                for j in range(i):
                    if arr[i] > arr[j]:
                        dp[i] = max(dp[i], dp[j] + 1)
            return dp

        dp = calcLis(south)
        return max(dp)


if __name__ == "__main__":
    solution = Solution()
    nums = int(input())
    cities = []
    for i in range(nums):
        north, south = list(map(int, input().split()))
        cities.append((north, south))
    res = solution.friendly_city(cities)
    print(res)
