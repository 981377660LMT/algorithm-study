from typing import List
from collections import defaultdict
from sortedcontainers import SortedList


# !总结:Dict+SortedList 可以解决所有这类问题
# !因为怕KeyError所以全用defaultdict了


class FoodRatings:
    def __init__(self, foods: List[str], cuisines: List[str], ratings: List[int]):
        """
        foods[i] 是第 i 种食物的名字。
        cuisines[i] 是第 i 种食物的烹饪方式。
        ratings[i] 是第 i 种食物的最初评分。
        """
        self.foodToScore = defaultdict(int)
        self.foodToCuision = defaultdict(str)
        self.cuisionRank = defaultdict(lambda: SortedList(key=lambda x: (-x[0], x[1])))
        for food, cuision, score in zip(foods, cuisines, ratings):
            self.foodToScore[food] = score
            self.foodToCuision[food] = cuision
            self.cuisionRank[cuision].add((score, food))

    def changeRating(self, food: str, newRating: int) -> None:
        """修改名字为 food 的食物的评分。删除旧的，添加新的"""
        preScore = self.foodToScore[food]
        cuision = self.foodToCuision[food]
        self.cuisionRank[cuision].discard((preScore, food))

        self.foodToScore[food] = newRating
        self.foodToCuision[food] = cuision
        self.cuisionRank[cuision].add((newRating, food))

    def highestRated(self, cuisine: str) -> str:
        """返回指定烹饪方式 cuisine 下评分最高的食物的名字。如果存在并列，返回 字典序较小 的名字。"""
        if not self.cuisionRank[cuisine]:
            return ""
        return self.cuisionRank[cuisine][0][1]


# Your FoodRatings object will be instantiated and called as such:
# obj = FoodRatings(foods, cuisines, ratings)
# obj.changeRating(food,newRating)
# param_2 = obj.highestRated(cuisine)
