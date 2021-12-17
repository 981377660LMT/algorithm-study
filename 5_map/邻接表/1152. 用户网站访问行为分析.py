from typing import List
from collections import defaultdict
from itertools import combinations

# 我们需要找到用户访问网站时的 『`共性行为路径`』，也就是有最多的用户都 `至少按某种次序访问过一次 的三个页面路径`
class Solution:
    def mostVisitedPattern(
        self, username: List[str], timestamp: List[int], website: List[str]
    ) -> List[str]:
        userHistory = defaultdict(list)
        patternCounter = defaultdict(set)

        for name, _, site in sorted(zip(username, timestamp, website)):
            userHistory[name].append(site)
        for user, visited in userHistory.items():
            for triplet in combinations(visited, 3):
                # 一个人加一次
                patternCounter[triplet].add(user)

        return sorted(patternCounter.items(), key=lambda item: (-len(item[1]), item[0]))[0][0]


print(
    Solution().mostVisitedPattern(
        username=["joe", "joe", "joe", "james", "james", "james", "james", "mary", "mary", "mary"],
        timestamp=[1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
        website=[
            "home",
            "about",
            "career",
            "home",
            "cart",
            "maps",
            "home",
            "home",
            "about",
            "career",
        ],
    )
)
# 输出：["home","about","career"]
# 解释：
# 由示例输入得到的记录如下：
# ["joe", 1, "home"]
# ["joe", 2, "about"]
# ["joe", 3, "career"]
# ["james", 4, "home"]
# ["james", 5, "cart"]
# ["james", 6, "maps"]
# ["james", 7, "home"]
# ["mary", 8, "home"]
# ["mary", 9, "about"]
# ["mary", 10, "career"]
# 有 2 个用户至少访问过一次 ("home", "about", "career")。
# 有 1 个用户至少访问过一次 ("home", "cart", "maps")。
# 有 1 个用户至少访问过一次 ("home", "cart", "home")。
# 有 1 个用户至少访问过一次 ("home", "maps", "home")。
# 有 1 个用户至少访问过一次 ("cart", "maps", "home")。

