from typing import Set


class Solution:
    def isPathCrossing(self, path: str) -> bool:
        record = {'N': 1j, 'S': -1j, 'E': 1, 'W': -1}
        visited: Set[complex] = set([0])
        position = 0
        for direction in path:
            position += record[direction]
            if position in visited:
                return True
            visited.add(position)
        return False

