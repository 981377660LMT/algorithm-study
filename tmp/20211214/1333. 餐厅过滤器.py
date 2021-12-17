from typing import List

# 过滤后返回餐馆的 id，按照 rating 从高到低排序。如果 rating 相同，那么按 id 从高到低排序
class Solution:
    def filterRestaurants(
        self, restaurants: List[List[int]], veganFriendly: int, maxPrice: int, maxDistance: int
    ) -> List[int]:
        return [
            r[0]
            for r in sorted(
                [
                    r
                    for r in restaurants
                    if r[2] >= veganFriendly and r[3] <= maxPrice and r[4] <= maxDistance
                ],
                key=lambda r: (r[1], r[0]),
                reverse=True,
            )
        ]


print(
    Solution().filterRestaurants(
        restaurants=[
            [1, 4, 1, 40, 10],
            [2, 8, 0, 50, 5],
            [3, 8, 1, 30, 4],
            [4, 10, 0, 10, 3],
            [5, 1, 1, 15, 1],
        ],
        veganFriendly=1,
        maxPrice=50,
        maxDistance=10,
    )
)

# 输出：[3,1,5]

