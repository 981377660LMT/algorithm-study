from typing import List

# 使用字典记录｛共同喜欢的商店：索引和｝，返回索引和并列最小的商店名
class Solution:
    def findRestaurant(self, list1: List[str], list2: List[str]) -> List[str]:
        dic = {name: list1.index(name) + list2.index(name) for name in set(list1) & set(list2)}
        return [name for name in dic if dic[name] == min(dic.values())]


print(
    Solution().findRestaurant(
        ["Shogun", "Tapioca Express", "Burger King", "KFC"],
        ["Piatti", "The Grill at Torrey Pines", "Hungry Hunter Steakhouse", "Shogun"],
    )
)

