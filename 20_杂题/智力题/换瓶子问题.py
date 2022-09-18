# !牛奶5元1瓶 两个空瓶加1元换1瓶 喝n瓶需要多少钱


def calMoney(n: int) -> int:
    def check(mid: int) -> bool:
        """mid元能不能喝n瓶"""
        start, money = divmod(mid, 5)
        res = start
        empty = start
        while empty >= 2 and money >= 1:
            count = min(empty // 2, money)
            money -= count
            empty -= count * 2
            empty += count
            res += count
        return res >= n

    left, right = 0, int(1e9)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1
    return left


if __name__ == "__main__":
    assert [calMoney(i) for i in range(10)] == [0, 5, 10, 11, 16, 21, 22, 23, 28, 37]
