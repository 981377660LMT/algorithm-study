import re
from collections import deque
from typing import List

PATTERN = re.compile(r"%([A-Z])%")


class Solution:
    def applySubstitutions(self, replacements: List[List[str]], text: str) -> str:
        mapping = {key: value for key, value in replacements}
        indeg = {key: 0 for key in mapping}
        revGraph = {key: [] for key in mapping}
        for key, value in mapping.items():
            for dep in PATTERN.findall(value):
                indeg[key] += 1
                revGraph[dep].append(key)

        queue = deque([key for key, deg in indeg.items() if deg == 0])
        order = []
        while queue:
            cur = queue.popleft()
            order.append(cur)
            for next_ in revGraph[cur]:
                indeg[next_] -= 1
                if indeg[next_] == 0:
                    queue.append(next_)

        resolved = {}
        for key in order:

            def repl(match):
                inner = match.group(1)
                return resolved[inner]

            trans = PATTERN.sub(repl, mapping[key])
            resolved[key] = trans

        return PATTERN.sub(lambda m: resolved[m.group(1)], text)
