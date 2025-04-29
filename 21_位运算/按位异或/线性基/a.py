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
        if presses == 0:
            return 1

        L = min(n, 6)
        masks: List[int] = []

        masks.append((1 << L) - 1)

        m2 = 0
        for i in range(L):
            if (i + 1) % 2 == 0:
                m2 |= 1 << i
        masks.append(m2)

        m3 = 0
        for i in range(L):
            if (i + 1) % 2 == 1:
                m3 |= 1 << i
        masks.append(m3)

        m4 = 0
        for i in range(L):
            if i % 3 == 0:
                m4 |= 1 << i
        masks.append(m4)

        basis = [0] * L
        r = 0
        for v in masks:
            x = v
            for i in range(L - 1, -1, -1):
                if (x >> i) & 1 == 0:
                    continue
                if basis[i] != 0:
                    x ^= basis[i]
                else:
                    basis[i] = x
                    r += 1
                    break

        if presses > r:
            return 1 << r

        ans = 0
        for k in range(0, r + 1):
            if k <= presses and (k % 2) == (presses % 2):
                ans += math.comb(r, k)
        return ans
