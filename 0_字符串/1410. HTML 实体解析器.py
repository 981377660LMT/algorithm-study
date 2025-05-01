entityMap = {
    "&quot;": '"',
    "&apos;": "'",
    "&gt;": ">",
    "&lt;": "<",
    "&frasl;": "/",
    "&amp;": "&",
}


class Solution:
    def entityParser(self, text: str) -> str:
        n = len(text)
        ptr = 0
        res = []
        while ptr < n:
            isEntity = False
            if text[ptr] == "&":
                for e in entityMap:
                    if text[ptr : ptr + len(e)] == e:
                        res.append(entityMap[e])
                        isEntity = True
                        ptr += len(e)
                        break
            if not isEntity:
                res.append(text[ptr])
                ptr += 1
        return "".join(res)
