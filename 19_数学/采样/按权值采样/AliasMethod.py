import random
from typing import List


class AliasMethod:
    """O(1)时间复杂度的采样方法."""

    def __init__(self, probs: List[float]):
        n = len(probs)
        self.prob = [0] * n
        self.alias = [0] * n
        scaled = [p * n for p in probs]
        small, large = [], []
        for i, p in enumerate(scaled):
            (small if p < 1 else large).append(i)
        while small and large:
            s, l = small.pop(), large.pop()
            self.prob[s] = scaled[s]
            self.alias[s] = l
            scaled[l] = scaled[l] + scaled[s] - 1
            (small if scaled[l] < 1 else large).append(l)
        for i in large + small:
            self.prob[i] = 1

    def pick(self):
        i = random.randint(0, len(self.prob) - 1)
        return i if random.random() < self.prob[i] else self.alias[i]


if __name__ == "__main__":

    def test_alias_method():

        probs = [0.1, 0.3, 0.5, 0.1]

        sampler = AliasMethod(probs)

        samples = int(1e6)
        results = [0] * len(probs)

        for _ in range(samples):
            results[sampler.pick()] += 1

        print("元素  理论概率  实际频率")
        for i in range(len(probs)):
            frequency = results[i] / samples
            print(f"{i}    {probs[i]:.2f}      {frequency:.4f}")

    test_alias_method()
