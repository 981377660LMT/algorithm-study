from typing import List
from pprint import pprint
from json import dumps, loads

Pairs = List[List[int]]
AdjList = Pairs


class Neuron:
    def __init__(self, id: int, children: List['Neuron'] = None) -> None:
        self.id = id
        self.children = [] if children is None else children


# 时间复杂度为O(V+E)。
def build(num: int, pairs: Pairs) -> List['Neuron']:
    """[summary]

    Args:
        num (int): [神经元个数]
        pairs (Pairs): [神经元的pairs]

    Returns:
        List[Neuron]: [神经元数组]
    """

    def getAdjList(num: int, pairs: Pairs) -> AdjList:
        adjList = [[] for _ in range(num + 1)]
        for pre, next in pairs:
            adjList[pre].append(next)
        return adjList

    def getNeuronById(id: int) -> Neuron:
        return neurons[id - 1]

    def dfs(curId: int) -> None:
        if visited[curId]:
            return
        visited[curId] = True

        for nextId in adjList[curId]:
            curNeuron = getNeuronById(curId)
            nextNeuron = getNeuronById(nextId)
            curNeuron.children.append(nextNeuron)
            dfs(nextId)

    neurons = [Neuron(i + 1) for i in range(num)]
    adjList = getAdjList(num, pairs)
    visited = [False for _ in range(num + 1)]

    for id in range(1, num + 1):
        if visited[id]:
            continue
        dfs(id)

    return neurons


if __name__ == '__main__':
    pairs = [
        [1, 2],
        [1, 3],
        [2, 3],
        [4, 5],
        [5, 6],
        [7, 8],
        [9, 10],
    ]
    neurons = build(10, pairs)
    pprint(
        [loads(dumps(neuron, default=lambda obj: obj.__dict__)) for neuron in neurons],
        sort_dicts=False,
    )
