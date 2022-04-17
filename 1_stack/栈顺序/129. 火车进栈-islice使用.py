# 按《字典序》输出前 20 种可能的出栈方案。
# 字典序就是dfs

# // 面对任何一个状态, 我们只有两种选择:
# // 1. 把下一个数进栈(如果还有下一个数)
# // 2. 把当前栈顶的数出栈(如果栈非空)
# 因为要顺序输出字典序，所以要优先出栈(越早出栈字典序越小)


from itertools import islice
from typing import Generator, List


def main(n: int) -> None:
    def dfs(index: int, inTrains: List[int], outTrains: List[int]) -> Generator[str, None, None]:
        """要按字典序输出，肯定是优先出栈"""

        if len(outTrains) == n:
            yield ''.join(map(str, outTrains))

        if inTrains:
            outTrains.append(inTrains.pop())
            yield from dfs(index, inTrains, outTrains)
            inTrains.append(outTrains.pop())

        if index <= n:
            inTrains.append(index)
            yield from dfs(index + 1, inTrains, outTrains)
            inTrains.pop()

    gen = dfs(1, [], [])
    slice = islice(gen, 20)  # 按《字典序》输出前 20 种可能的出栈方案。
    for s in slice:
        print(s)


n = int(input())
main(n)

