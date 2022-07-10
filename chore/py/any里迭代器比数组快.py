import timeit


nums = list(range(100))


def any1() -> bool:
    """any里用迭代器 一般"""
    return any(n % 2 == 0 for n in nums)


def any2() -> bool:
    """any里用list 慢"""
    return any([n % 2 == 0 for n in nums])


def any3() -> bool:
    """不用any 最快"""
    for n in nums:
        if n % 2 == 0:
            return True
    return False


if __name__ == "__main__":
    # !迭代器更快  (不用先创建list)
    print(f"time1: {timeit.timeit(any1, number=1000000)}")  # 0.5180921999999555
    print(f"time2: {timeit.timeit(any2, number=1000000)}")  # 5.973919199997908
    # !不用any更快 (调用函数有额外开销)
    print(f"time3: {timeit.timeit(any3, number=1000000)}")  # 0.1570034000033047
