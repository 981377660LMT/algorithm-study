# # 配列の初期状態
# A = [2, 0, 1, 9]
# # 更新操作
# operations = [('a.0', 2, 2, 'a.1'), ('a.1', 0, 1, 'a.2'), ('a.2', 1, 1, 'a.3'),
#               ('a.1', 3, 0, 'b.2'), ('a.3', 3, 3, 'a.4')]
# # クエリ
# queries = [('a.1', 2), ('a.4', 1), ('a.4', 3), ('b.2', 3)]
# https://qiita.com/wotsushi/items/72e7f8cdd674741ffd61#%E5%8F%82%E8%80%83%E8%A8%98%E4%BA%8B

# 可持久化数组-离线(バッチ処理)
# !把版本看成节点，那么根据版本的继承关系这些节点会构成一棵树。
# !多个版本组成了一棵树,本质上是树上的离线查询

from collections import defaultdict
from typing import List, Tuple


def offlineQuery(
    nums: List[int], operations: List[Tuple[str, int, int, str]], queries: List[Tuple[str, int]]
) -> List[int]:
    """预处理出查询组,dfs的过程中输出所有查询结果

    Args:
        nums (List[int]): 初始数组
        operations (List[Tuple[str, int, int, str]]): 更新操作,也叫action
        queries (List[Tuple[str, int]]): 查询操作
    """

    def dfs(version: str, curNums: List[int]) -> None:
        for qi, qv in queryGroup[version]:
            res[qi] = curNums[qv]
        for nextVersion, pos, newValue in adjMap[version]:
            oldValue = curNums[pos]
            curNums[pos] = newValue
            dfs(nextVersion, curNums)
            curNums[pos] = oldValue

    queryGroup = defaultdict(list)  # {'a.1': [2], 'a.4': [1, 3], 'b.2': [3]}  预处理每个版本的查询
    for qi, (vesion, qv) in enumerate(queries):
        queryGroup[vesion].append((qi, qv))

    adjMap = defaultdict(list)  # {'a.0': [('a.1', 2, 2)],'a.1': [('a.2', 0, 1), ('b.2', 3, 0)]}
    for curVersion, pos, newValue, nextVersion in operations:
        adjMap[curVersion].append((nextVersion, pos, newValue))

    res = [-1] * len(queries)
    dfs("a.0", nums)
    return res


if __name__ == "__main__":
    assert offlineQuery(
        [2, 0, 1, 9],
        [
            ("a.0", 2, 2, "a.1"),
            ("a.1", 0, 1, "a.2"),
            ("a.2", 1, 1, "a.3"),
            ("a.1", 3, 0, "b.2"),
            ("a.3", 3, 3, "a.4"),
        ],
        [("a.1", 2), ("a.4", 1), ("a.4", 3), ("b.2", 3)],
    ) == [2, 1, 3, 0]
