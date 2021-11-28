# 返回表示最终状态的字符串。
# 如果同时有多米诺骨牌落在一张垂直竖立的多米诺骨牌的两边，由于受力平衡， 该骨牌仍然保持不变。
class Solution:
    # 模拟
    def pushDominoes1(self, dominoes: str) -> str:
        pre, cur = "", dominoes

        while cur != pre:
            pre = cur
            cur = cur.replace("R.L", "T")
            cur = cur.replace(".L", "LL")
            cur = cur.replace("R.", "RR")
            cur = cur.replace("T", "R.L")
        return cur

    # 计算合力
    # 我们可以对每个多米诺骨牌计算净受力。我们关心的受力取决于一个多米诺骨牌和最近的左侧 'R' 右侧 'L' 的距离：哪边近，就受哪边力更多。
    # 从左向右扫描，我们的力每轮迭代减少 1
    def pushDominoes(self, dominoes: str) -> str:
        n = len(dominoes)
        force = [0] * n

        f_to_right = 0
        for i in range(n):
            if dominoes[i] == 'R':
                f_to_right = n
            elif dominoes[i] == 'L':
                f_to_right = 0
            else:
                f_to_right = max(f_to_right - 1, 0)
            force[i] += f_to_right

        f_to_left = 0
        for i in range(n - 1, -1, -1):
            if dominoes[i] == 'L':
                f_to_left = n
            elif dominoes[i] == 'R':
                f_to_left = 0
            else:
                f_to_left = max(f_to_left - 1, 0)
            force[i] -= f_to_left
        return ''.join('.' if f == 0 else 'R' if f > 0 else 'L' for f in force)


print(Solution().pushDominoes(".L.R...LR..L.."))
