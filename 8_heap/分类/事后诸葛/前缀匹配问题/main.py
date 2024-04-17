# 前缀匹配问题
# https://taodaling.github.io/blog/2019/10/28/%E4%B8%80%E4%BA%9B%E8%B4%AA%E5%BF%83%E9%97%AE%E9%A2%98/


from itertools import permutations
from random import randint

from heapq import heappop, heappush
from typing import List


# 问题1：电影院有n个人和n把椅子，椅子编号为0到n-1，对于第i个人，
# 他希望能坐在前[0,ai]把椅子中的一把椅子上，不然他会不满意。
# 现在要求我们安排这n个人的座次，要求不满意的人数最少。
#
# 很容易看出，将人和椅子建立二分图，就是一个简单的二分图匹配问题。
# 但是，假如我告诉你n=1e6的时候，是不是就发现这问题不简单了。
# 观察一下这个问题与经典的二分图匹配问题有哪些区别，最重要的区别就是每个人能够匹配的椅子都是`某个前缀`。
# 因此我们可以贪心地解决这个问题。
# !因此算法就得出来了，我们把所有人按ai进行排序，之后从前往后遍历每一张座位，
# 如果有座位在，我们让所有能坐在该座位上的人中ai属性最小的人坐上去。
# （假如最优方案由编号更大的人坐在这里，我们交换两人座位方案依旧可行且最大匹配不变）
def solve1(n: int, prefer: List[int]) -> int:
    m = len(prefer)
    prefer = sorted(prefer)
    i, j = 0, 0  # 椅子，人
    res = 0
    while i < n and j < m:
        if prefer[j] >= i:
            i += 1
            j += 1
        else:
            j += 1
            res += 1
    res += m - j
    return res


# Renovation (排序贪心+反悔堆，带权前缀匹配问题)
# https://www.luogu.com.cn/problem/CF883J
# 问题2(带生气值)：电影院有n个人和n把椅子，椅子编号为0到n-1，对于第i个人，
# 他希望能坐在前ai把椅子中的一把椅子上，不然他的不满意值会变成angry[i]（如果满意则为0）
# 现在要求我们安排这n个人的座次，要求不满意的人数最少。
#
# !先对所有人按照ai进行排序。之后扫描所有椅子，让所有能坐上这把椅子的人中ai最小的那个人坐上。
# 但是不同的是，当处理完第i把椅子后，我们需要考虑所有最后能匹配的椅子为i的人，他们该何去何从。
# 我们会试图用他们替换已经坐上椅子中的人中ui最小的那个，假如比我们现在手头的候选人的ui更小，那么就替换它，否则我们的候选人将永远失去坐上椅子的机会。
def solve2(n: int, prefer: List[int], angry: List[int]) -> int:
    m = len(prefer)
    order = sorted(range(m), key=lambda x: prefer[x])
    pq = []  # 生气值最小堆
    res = 0
    i, j = 0, 0  # 椅子，人
    while i < n and j < m:
        pid = order[j]
        if prefer[pid] >= i:
            heappush(pq, angry[pid])
            i += 1
            j += 1
        else:
            if pq and pq[0] < angry[pid]:
                res += heappop(pq)
                heappush(pq, angry[pid])
            else:
                res += angry[pid]
            j += 1
    while j < m:
        if pq and pq[0] < angry[order[j]]:
            res += heappop(pq)
            heappush(pq, angry[order[j]])
        else:
            res += angry[order[j]]
        j += 1
    return res


# 形态3：电影院有n个人和n把椅子，椅子编号为0到n-1，对于第i个人，
# 他希望能坐在前ai把椅子中的一把椅子上，或者坐在后bi把椅子中的一把椅子上，不然他会不满意。
# 现在要求我们安排这n个人的座次，要求不满意的人数最少。
#
# 这个问题不仅存在前缀，还存在后缀，怎么解决呢。
# !我们先分配前缀座位，再对剩下的没有坐的椅子分配后缀座位。
# 对于一个人，如果bi越大，那么在分配前缀的阶段中，即使被淘汰，成本也越小。
# !因此我们可以规定一个人被淘汰的费用为bi，那么分配前缀的过程实际上就是一个最小费用最大匹配，而分配后缀的过程就是一个普通的最大匹配。
def solve3(n: int, prefer1: List[int], prefer2: List[int]) -> int:
    ...


# 形态4：电影院有n个人和m把椅子，椅子编号为0到m-1，对于第i个人，
# 他希望能坐在前ai把椅子中的连续的ci把椅子上，不然他会不满意。
# 现在要求我们安排这n个人的座次，要求不满意的人数最少。
# !一个人可以分配到不止一把椅子。我可以把第i个人占用的ci个座位理解成为让他没有椅子坐的费用为ci。之后就是一个带权的前缀匹配问题了。
def solve4(n: int, m: int, prefer: List[int], count: List[int]) -> int:
    ...


if __name__ == "__main__":

    def check1ByBruteForce(n: int, prefer: List[int]) -> int:
        m = len(prefer)
        res = int(1e18)
        for perm in permutations(range(m)):
            curSum = 0
            for i, p in enumerate(perm):
                if i >= n:
                    curSum += 1
                    continue
                if prefer[p] < i:
                    curSum += 1
            if curSum < res:
                res = curSum
        return res

    def check2ByBruteForce(n: int, prefer: List[int], angry: List[int]) -> int:
        m = len(prefer)
        res = int(1e18)
        for perm in permutations(range(m)):
            curSum = 0
            for i, p in enumerate(perm):
                if i >= n:
                    curSum += angry[p]
                    continue
                if prefer[p] < i:
                    curSum += angry[p]
            if curSum < res:
                res = curSum
        return res

    for _ in range(100):
        n = 5
        m = randint(1, 7)
        prefer = [randint(0, n) for _ in range(m)]
        res1, res2 = check1ByBruteForce(n, prefer), solve1(n, prefer)
        if res1 != res2:
            print(n, prefer)
            print(res1, res2)
            raise ValueError("error 1")
    print("pass 1")

    for _ in range(1000):
        n = 6
        m = randint(1, 7)
        prefer = [randint(0, n) for _ in range(m)]
        angry = [randint(0, 10) for _ in range(m)]
        res1, res2 = check2ByBruteForce(n, prefer, angry), solve2(n, prefer, angry)
        if res1 != res2:
            print(n, prefer, angry)
            print(res1, res2)
            raise ValueError("error 2")
    print("pass 2")
