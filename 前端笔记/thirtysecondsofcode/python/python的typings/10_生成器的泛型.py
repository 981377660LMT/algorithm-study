from typing import Generator


def echo_round() -> Generator[int, float, str]:
    sent = yield 0
    while sent >= 0:
        sent = yield round(sent)
    return 'Done'


# 生成器可以由泛型类型 Generator[YieldType, SendType, ReturnType] 注解。例如：
