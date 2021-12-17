from typing import List
from collections import Counter


class Solution:
    def sortFeatures(self, features: List[str], responses: List[str]) -> List[str]:
        #### word和index的对应关系
        idByWord = {word: id for id, word in enumerate(features)}

        #### 统计每个单词的出现freq
        freq = Counter()
        for resp in responses:
            for word in set(resp.split(' ')):  # 注意，在一条string出现多次计1次。如例1中的cooler
                if word in idByWord:
                    freq[word] += 1

        return sorted(features, key=lambda word: -freq[word])


print(
    Solution().sortFeatures(
        features=["cooler", "lock", "touch"],
        responses=["i like cooler cooler", "lock touch cool", "locker like touch"],
    )
)

