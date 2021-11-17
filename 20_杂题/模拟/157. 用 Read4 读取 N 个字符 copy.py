class Solution:
    def read(self, buf, n):
        """
        :type buf: Destination buffer (List[str])
        :type n: Number of characters to read (int)
        :rtype: The number of actual characters read (int)
        """
        bi = 0
        for _ in range(0, n, 4):
            tmp = [None] * 4  # 必须先开辟出4个空间
            cur_len = read4(tmp)  # 先读入到tmp
            for j in range(cur_len):
                buf[bi] = tmp[j]  # 从tmp复制到buf
                bi += 1
        return min(bi, n)

