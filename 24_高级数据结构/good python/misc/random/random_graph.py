# 生成随机图
# random: 生成随机图
# cycle: 生成环

import enum
from typing import Optional, List, Tuple
import random


class RandomGraphType(enum.Enum):
    random = enum.auto()
    cycle = enum.auto()


class RandomGraph:
    @classmethod
    def build(
        cls,
        n: int,
        m: int,
        typ: RandomGraphType = RandomGraphType.random,
        seed: Optional[int] = None,
    ) -> List[Tuple[int, int]]:
        random.seed(seed)
        if typ == RandomGraphType.random:
            return cls._build_random(n, m)
        if typ == RandomGraphType.cycle:
            return cls._build_cycle(n, m)
        raise ValueError(typ)

    @classmethod
    def _build_cycle(cls, n: int, m: int) -> List[Tuple[int, int]]:
        assert m == n
        cycle = list(range(n))
        random.shuffle(cycle)
        cycle.append(cycle[-1])
        edges = [None] * n
        for i in range(n):
            u, v = cycle[i], cycle[i + 1]
            if random.random() < 0.5:
                edges[i] = (v, u)  # type: ignore
            else:
                edges[i] = (u, v)  # type: ignore
        random.shuffle(edges)
        assert len(edges) == m
        return edges  # type: ignore

    @classmethod
    def _build_random(cls, n: int, m: int) -> List[Tuple[int, int]]:
        assert m <= n * (n - 1) // 2
        edges = set()
        while len(edges) < m:
            u = random.randrange(0, n)
            v = random.randrange(0, n)
            while u == v:
                v = random.randrange(0, n)
            if u > v:
                u, v = v, u
            edges.add((u, v))
        edges = list(edges)
        for i in range(m):
            u, v = edges[i]
            if random.random() < 0.5:
                edges[i] = (v, u)
        random.shuffle(edges)
        assert len(edges) == m
        return edges
