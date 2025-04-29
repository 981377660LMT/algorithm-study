import re
from typing import List


class Solution:
    def removeComments(self, source: List[str]) -> List[str]:
        return list(filter(None, re.sub("//.*|/\*(.|\n)*?\*/", "", "\n".join(source)).split("\n")))
