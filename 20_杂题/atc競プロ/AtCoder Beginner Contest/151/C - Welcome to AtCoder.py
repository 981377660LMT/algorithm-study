from typing import List, Tuple


# 计算AC和WA的数量之和
def countACAndWA(n: int, submissions: List[Tuple[int, str]]) -> Tuple[int, int]:
    AC, WA = [False] * (n + 1), [0] * (n + 1)
    for id, result in submissions:
        if result == "AC":
            AC[id] = True
        elif not AC[id]:  # 已经AC的题目不再计算WA
            WA[id] += 1
    return sum(AC), sum(wa for ac, wa in zip(AC, WA) if ac)  # 只有AC的题目才计算WA


n, m = map(int, input().split())
submissions = []
for _ in range(m):
    count, result = input().split()
    submissions.append((int(count), result))
print(*countACAndWA(n, submissions))
