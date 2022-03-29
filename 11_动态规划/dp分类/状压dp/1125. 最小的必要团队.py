# 691. 贴纸拼词.py


from typing import List
from functools import lru_cache

# 1 <= req_skills.length <= 16
# 1 <= people.length <= 60
# 你规划了一份需求的技能清单 req_skills，并打算从备选人员名单 people 中选出些人组成一个「必要团队」
# 请你返回 任一 规模最小的必要团队，团队成员用人员编号表示。
INF = 0x7FFFFFFF


class Solution:
    def smallestSufficientTeam(self, req_skills: List[str], people: List[List[str]]) -> List[int]:
        idBySkill = {skill: id for id, skill in enumerate(req_skills)}  # 每个技能赋予唯一ID

        # 每位候选人的状态
        cand = []
        for skills in people:
            state = 0
            for skill in skills:
                state |= 1 << idBySkill[skill]
            cand.append(state)

        n = len(req_skills)
        target = (1 << n) - 1

        @lru_cache(None)
        def dfs(person: int, state: int) -> List:
            # 不需要人
            if state == target:
                return []

            # impossible
            if person >= len(people):
                return [0] * 100

            # 要不要当前的人
            return min(
                dfs(person + 1, state), [person] + dfs(person + 1, state | cand[person]), key=len
            )

        return dfs(0, 0)


print(
    Solution().smallestSufficientTeam(
        req_skills=["java", "nodejs", "reactjs"],
        people=[["java"], ["nodejs"], ["nodejs", "reactjs"]],
    )
)
# 输出：[0,2]
