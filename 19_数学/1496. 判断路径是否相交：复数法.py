from typing import Set


record = {'N': 1j, 'S': -1j, 'E': 1, 'W': -1}


class Solution:
    def isPathCrossing(self, path: str) -> bool:
        visited: Set[complex] = set([0])
        pos = 0
        for dir in path:
            pos += record[dir]
            if pos in visited:
                return True
            visited.add(pos)
        return False

