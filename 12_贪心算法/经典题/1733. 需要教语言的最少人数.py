from typing import List
from collections import Counter

# 你可以选择 一门 语言并教会一些用户，使得所有好友之间都可以相互沟通。请返回你 最少 需要教会多少名用户。
# 好友关系`没有传递性`，也就是说如果 x 和 y 是好友，且 y 和 z 是好友， x 和 z 不一定是好友。

# 总结：
# 1.要教的学生是没有公共语言的
# 2. 要使教授的人数最少，需要教授最流行的语言(他们之中会的最多)
class Solution:
    def minimumTeachings(
        self, n: int, languages: List[List[int]], friendships: List[List[int]]
    ) -> int:
        lan = list(map(set, languages))
        students = set(p for u, v in friendships for p in (u, v) if not lan[u - 1] & lan[v - 1])
        counter = Counter()
        for s in students:
            counter += Counter(lan[s - 1])
        return len(students) - max(counter.values(), default=0)


print(
    Solution().minimumTeachings(
        n=2, languages=[[1], [2], [1, 2]], friendships=[[1, 2], [1, 3], [2, 3]]
    )
)
# 输出：1
# 解释：你可以选择教用户 1 第二门语言，也可以选择教用户 2 第一门语言。
