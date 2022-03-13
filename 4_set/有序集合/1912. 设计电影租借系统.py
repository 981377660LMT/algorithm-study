from typing import List
from sortedcontainers import SortedList
from collections import defaultdict

# 你有一个电影租借公司和 n 个电影商店。
# 你想要实现一个电影租借系统，它支持查询、预订和返还电影的操作。
# 同时系统还能生成一份当前被借出电影的报告。

# 所有电影用二维整数数组 entries 表示，其中 entries[i] = [shopi, moviei, pricei]

# 总结:
# 1.题目限制了`商店需要按照 价格 升序排序，如果价格相同，则 shopi 较小 的商店排在前面。`
# 说明需要用sortedList维护电影信息
# 2.每个商店都有movie=>price信息 用dict维护
# 3.renting也需要sortedList来生成报告:res 中的电影需要按 价格 升序排序；如果价格相同，则 shopj 较小 的排在前面；如果仍然相同，则 moviej 较小 的排在前面。


class MovieRentingSystem:
    def __init__(self, n: int, entries: List[List[int]]):
        self.n = n
        self.movies = defaultdict(SortedList)
        self.shops = defaultdict(dict)
        self.renting = SortedList([])
        for shop, movie, price in entries:
            self.movies[movie].add((price, shop))
            self.shops[shop][movie] = price

    # 找到拥有指定电影且 未借出 的商店中 最便宜的 5 个商店
    def search(self, movie: int) -> List[int]:
        return [info[1] for info in list(self.movies[movie].islice(stop=5))]

    # 从指定商店借出指定电影，题目保证指定电影在指定商店 未借出 。
    def rent(self, shop: int, movie: int) -> None:
        price = self.shops[shop][movie]
        self.movies[movie].discard((price, shop))
        self.renting.add((price, shop, movie))

    # 在指定商店返还 之前已借出 的指定电影。
    def drop(self, shop: int, movie: int) -> None:
        price = self.shops[shop][movie]
        self.movies[movie].add((price, shop))
        self.renting.discard((price, shop, movie))

    # 返回 最便宜的 5 部已借出电影
    def report(self) -> List[List[int]]:
        return [[shop, movie] for _, shop, movie in self.renting.islice(stop=5)]


# Your MovieRentingSystem object will be instantiated and called as such:
# obj = MovieRentingSystem(n, entries)
# param_1 = obj.search(movie)
# obj.rent(shop,movie)
# obj.drop(shop,movie)
# param_4 = obj.report()
