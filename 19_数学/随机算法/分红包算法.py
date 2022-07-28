# 让每个人公平

# 自己想出来的算法
# 红包在一开始创建的时候，分配方案就订好了 生成一组加权随机数即可

from random import randint
from typing import List


def randomSplit(people: int, money: float) -> List[int]:
    """
    为了避免浮点数误差,
    使用`分`来表示金额单位,
    且使用`Fraction(python)`或者`math.BigDecimal(java)`来计算
    """
    money = int(money * 100)
    weights = [randint(5, 15) for _ in range(people)]
    sum_ = sum(weights)
    res = [int(money * weight / sum_) for weight in weights]
    remain = money - sum(res)
    res[randint(0, people - 1)] += remain
    return res


# 20个人 666.66元红包
print(randomSplit(20, 666.66))
