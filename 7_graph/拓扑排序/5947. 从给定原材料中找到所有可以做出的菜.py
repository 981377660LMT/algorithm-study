from typing import List
from collections import deque, defaultdict

# 此题用字典记录更好，因为键都是字符串
# 拓扑排序两大件:indegree,adjMap
# 起点直接为supplies(原材料入度为0)


class Solution:
    def findAllRecipes(
        self, recipes: List[str], ingredients: List[List[str]], supplies: List[str]
    ) -> List[str]:
        deg = defaultdict(int)
        adjMap = defaultdict(list)
        for cur, deps in zip(recipes, ingredients):
            for dep in deps:
                adjMap[dep].append(cur)
                deg[cur] += 1

        queue = deque(supplies)
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                deg[next] -= 1
                if deg[next] == 0:
                    queue.append(next)

        return [food for food in recipes if deg[food] == 0]


print(
    Solution().findAllRecipes(
        recipes=["bread", "sandwich"],
        ingredients=[["yeast", "flour"], ["bread", "meat"]],
        supplies=["yeast", "flour", "meat"],
    )
)
