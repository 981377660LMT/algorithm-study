# 遍历区间的差分字典.


from collections import defaultdict


class DiffMapIntervals:
    __slots__ = ("_mp",)

    def __init__(self):
        self._mp = defaultdict(int)

    def add(self, left: int, right: int, x: int):
        """閉区間 [l,r] に x を加算する."""
        if left <= right:
            self._mp[left] += x
            self._mp[right + 1] -= x

    def enumerate(self, since: int, until: int):
        """累積和を求める.

        [Output]
        (y, l, r) という形のリスト. ただし, (y, l, r) は l<=x<=y の範囲では y であるということを意味する.
        """
        curSum = 0
        pre = since
        mp = self._mp
        for t in sorted(mp):
            if t > until:
                break
            if mp[t] == 0:
                continue
            if pre <= t - 1:
                yield (curSum, pre, t - 1)
            curSum += mp[t]
            pre = t
        if pre <= until:
            yield (curSum, pre, until)


if __name__ == "__main__":

    def demo():
        D = DiffMapIntervals()

        # 添加闭区间 [left, right] 加上值 x
        D.add(1, 3, 10)  # 在区间 [1,3] 加 10
        D.add(2, 5, -5)  # 在区间 [2,5] 加 -5
        D.add(4, 6, 3)  # 在区间 [4,6] 加 3

        # 枚举从 1 到 7 的累积和区间
        print("累积和区间 (value, start, end):")
        for value, start, end in D.enumerate(1, 7):
            print(f"({value}, {start}, {end})")

    demo()
