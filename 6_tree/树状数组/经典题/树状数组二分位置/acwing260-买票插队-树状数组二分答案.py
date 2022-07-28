# 告诉你每一个排队的人的编号，和他进入队列时的具体位置。
# 请你确定最终的队列顺序。
# 每个测试用例，输出一行包含 N 个整数（表示每个人的编号）的结果，表示最终的人员队列顺序。

# 这道题是牛的复刻版 根据身高和前面的人数重建队列
# 树状数组维护一个01序列的前缀和；从后往前考虑,因为后面的去除后不影响前面的；快速找到第 k 个 1 所在的位置，使用二分即可
from collections import defaultdict


class BIT:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError("index 必须是正整数")
        while index <= self.size:
            self.tree[index] += delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


# 倒序分析：
# - 69：前面有 2 个人
# - 33：除去 69 之后，前面有 1 个人
# - 51：除去 69、33 之后，前面有 1 个人
def main() -> None:
    n = int(input())
    res = [0] * n
    bit = BIT(n + 10)

    people = []

    for i in range(n):
        preCount, id = map(int, input().split())  # 插队后他前面的人数 个人编号
        people.append((preCount, id))
        bit.add(i + 1, 1)  # 离散化了

    # 从后往前考虑,因为后面的去除后不影响前面的
    # 二分第i个人在树状数组中排的位置
    for i in range(n - 1, -1, -1):
        preCount, id = people[i]
        preCount += 1  # 注意这里

        left, right = 1, n
        while left <= right:
            mid = (left + right) >> 1
            count = bit.query(mid)
            if preCount <= count:
                right = mid - 1
            else:
                left = mid + 1

        res[left - 1] = id
        bit.add(left, -1)

    print(" ".join(map(str, res)))


while True:
    try:
        main()
    except EOFError:
        break
