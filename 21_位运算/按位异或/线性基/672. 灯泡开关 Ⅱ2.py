# 672. 灯泡开关 Ⅱ
# https://leetcode.cn/problems/bulb-switcher-ii/description/
#
# 房间中有 n 只已经打开的灯泡，
# !编号从 1 到 n 。墙上挂着 4 个开关 。
#
# 这 4 个开关各自都具有不同的功能，其中：
#
# 开关 1 ：反转当前所有灯的状态（即开变为关，关变为开）
# 开关 2 ：反转编号为偶数的灯的状态（即 2, 4, ...）
# 开关 3 ：反转编号为奇数的灯的状态（即 1, 3, ...）
# 开关 4 ：反转编号为 j = 3k + 1 的灯的状态，其中 k = 0, 1, 2, ...（即 1, 4, 7, 10, ...）
# 你必须 恰好 按压开关 presses 次。每次按压，你都需要从 4 个开关中选出一个来执行按压操作。
#
# !给你两个整数 n 和 presses ，执行完所有按压之后，返回 不同可能状态 的数量。

import math
from typing import List


class Solution:
    def flipLights(self, n: int, presses: int) -> int:
        """
        线性基解法（F₂ 空间）：
        1. 只看前 L=min(n,6) 盏灯，因为状态周期为 6。
        2. 把 4 种开关动作分别用 L 位掩码表示，插入到线性基中，
           计算出秩 r。
        3. 若 presses > r，则任何子空间向量都可达，答案 = 2^r。
        4. 否则只能选恰好 presses 次复制（即选 k 个基向量且 k≡presses mod2),
           状态数 = ∑_{0≤k≤r, k≤presses, k≡presses (mod2)} C(r,k)。
        """
        # 按压 0 次时，只有初始全开状态
        if presses == 0:
            return 1

        # L 是我们要考虑的灯泡数
        L = min(n, 6)
        masks: List[int] = []

        # 开关 1：翻转所有
        masks.append((1 << L) - 1)

        # 开关 2：翻转偶数号（1-based）→ 0-based 下标为奇数
        m2 = 0
        for i in range(L):
            if (i + 1) % 2 == 0:
                m2 |= 1 << i
        masks.append(m2)

        # 开关 3：翻转奇数号（1-based）→ 0-based 下标为偶数
        m3 = 0
        for i in range(L):
            if (i + 1) % 2 == 1:
                m3 |= 1 << i
        masks.append(m3)

        # 开关 4：翻转下标 i 满足 i%3==0 的灯（1-based 3k+1）
        m4 = 0
        for i in range(L):
            if i % 3 == 0:
                m4 |= 1 << i
        masks.append(m4)

        # --- 构建线性基 ---
        basis = [0] * L  # basis[i] 存“最高位为 i” 的基向量
        r = 0  # 当前基向量个数（秩）
        for v in masks:
            x = v
            # 自高位向低位消元
            for i in range(L - 1, -1, -1):
                if (x >> i) & 1 == 0:
                    continue
                if basis[i] != 0:
                    # 消掉 i 位
                    x ^= basis[i]
                else:
                    # 找到新的基向量
                    basis[i] = x
                    r += 1
                    break
        # --- 基构建完毕，秩为 r ---

        # 若按压次数超过秩，则整个子空间都可达到
        if presses > r:
            return 1 << r

        # 否则只统计恰好 presses 次按键对应的子集大小
        ans = 0
        for k in range(0, r + 1):
            # k ≤ presses 且 k 与 presses 同奇偶
            if k <= presses and (k % 2) == (presses % 2):
                ans += math.comb(r, k)
        return ans


if __name__ == "__main__":
    sol = Solution()
    tests = [
        (1, 0, 1),
        (1, 5, 2),
        (2, 1, 3),
        (2, 2, 4),
        (3, 0, 1),
        (3, 1, 4),
        (3, 2, 7),
        (3, 3, 8),
        (10, 100, 8),
    ]
    for n, presses, expect in tests:
        res = sol.flipLights(n, presses)
        print(f"n={n}, presses={presses} -> {res} (expected {expect})")
