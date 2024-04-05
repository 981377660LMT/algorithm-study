# 缝隙二分
# https://taodaling.github.io/blog/2019/09/12/%E4%BA%8C%E5%88%86/
#
# !二分还能作用在非单调函数上查找缝隙。
# !所谓的缝隙是指这样一个整数x，满足check(x−1)≠check(x)。
# 要执行缝隙二分的前提是，一开始给定的l和r满足check(l)≠check(r)。
# 在执行缝隙而二分的过程，利用l和r算出mid=(l+r)/2后，如果check(mid)=check(l)，那么就令l=mid，否则令r=mid。
# 容易发现这样做始终能保证check(l)≠check(r)。
# 由于区间在不断缩小，因此最终一定能找到一个缝隙。当然我们无法确定找到的是哪个缝隙，但是这不重要。
# !交互题，让你01序列上找一个跳跃点


# E - Majority of Balls
# https://atcoder.jp/contests/ddcc2020-qual/tasks/ddcc2020_qual_e
# https://www.luogu.com.cn/problem/solution/AT_ddcc2020_qual_e
# 这是一道交互题。有 2n 个球，球有红色和蓝色。
# 每次你可以询问n个球，交互库会回答这n个球中红色球的数量和蓝色球的数量谁更多。
# 你需要在210 次询问内猜出答案。
# !n<=99且n为奇数
#
# !1.如果我们找到了一组有 N−1 个球的球堆，满足红球数量等于蓝球数量，就可以快速求出每个球的颜色了。
#    具体的，对于不属于这 N−1 个球中的球，询问时询问这 N−1 个球与该球共同组成的 N 个球中数量较多的颜色，
#    对于属于这 N−1 个球中的球，用这N-1个球的补集去除A[left]和A[left+n]即可.
#    得到的颜色即为该球颜色。
# !2.如何快速找到那一组满足条件的球堆呢？
#    !如果[left,left+n)的答案为RED，[left+1,left+n+1)的答案为BLUE，则[left+1,left+n)这n-1个球中红球数量等于蓝球数量
#    且A[left]一定为RED，A[left+n]一定为BLUE
#    又因为[1,n+1)和[n+1,2n+1)的答案一定不同，所以我们可以用间隙二分找到最左边的left满足[left+1,left+n)这n-1个球中红球数量等于蓝球数量.

from typing import List


def majorityOfBalls() -> None:
    def output(res: List[str]) -> None:
        print("!", "".join(res), flush=True)

    def query(interval: List[int]) -> str:
        print("?", *interval, flush=True)
        return input()

    def findLeft() -> int:
        """找到最靠左的left满足[left+1,left+n)这n-1个球中红球数量等于蓝球数量."""
        x = query(list(range(1, n + 1)))  # !此时 [n+1,2*n+1) 的答案与x一定不同
        left, right = 1, n
        while left <= right:
            mid = (left + right) // 2
            # !如果[left,left+n)的答案为RED，[left+1,left+n+1)的答案为BLUE，则[left+1,left+n)这n-1个球中红球数量等于蓝球数量
            # !且A[left]一定为RED，A[left+n]一定为BLUE
            answer = query(list(range(mid + 1, mid + n + 1)))  # 缝隙二分求间隙
            if answer != x:
                right = mid - 1
            else:
                left = mid + 1
        return left

    n = int(input())
    left = findLeft()
    res = []
    for i in range(1, 2 * n + 1):
        if left + 1 <= i < left + n:
            answer = query([i] + list(range(1, left)) + list(range(left + n + 1, 2 * n + 1)))
        else:
            answer = query([i] + list(range(left + 1, left + n)))
        res.append("R" if answer == "Red" else "B")

    output(res)


if __name__ == "__main__":
    majorityOfBalls()
