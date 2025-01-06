import time
import sys


# Python by Steve Hanov
# Suffix array construction from:
# From Ge Nong, Sen Zhang and Wai Hong Chan, Two Efficient Algorithms for Linear Suffix Array Construction, 2008
# longestMatches from: http://stevehanov.ca/blog/index.php?id=146
# Released to the public domain.
class SuffixArray:
    def __init__(self, A, B):
        self.length1 = len(A) + 1
        self.length2 = len(B) + 1
        self.seq = A + b"\001" + B + b"\000"

        # find the start or end of each bucket
        def getBuckets(s, bkt, n, K, end, os):
            sum = 0
            # clear all buckets
            for i in range(K + 1):
                bkt[i] = 0

            # compute the size of each bucket
            for i in range(n):
                bkt[s[i + os]] += 1

            for i in range(K + 1):
                sum += bkt[i]
                if end:
                    bkt[i] = sum
                else:
                    bkt[i] = sum - bkt[i]

        # compute SAl
        def induceSAl(t, SA, s, bkt, n, K, end, os, oa):
            # find starts of buckets
            getBuckets(s, bkt, n, K, end, os)
            for i in range(n):
                j = SA[i + oa] - 1
                if j >= 0 and not t[j]:
                    SA[bkt[s[j + os]] + oa] = j
                    bkt[s[j + os]] += 1

        # compute SAs
        def induceSAs(t, SA, s, bkt, n, K, end, os, oa):
            # find ends of buckets
            getBuckets(s, bkt, n, K, end, os)
            for i in range(n - 1, -1, -1):
                j = SA[i + oa] - 1
                if j >= 0 and t[j]:
                    bkt[s[j + os]] -= 1
                    SA[bkt[s[j + os]] + oa] = j

        # find the suffix array SA of s[0..n-1] in {1..K}^n
        # require s[n-1]=0 (the sentinel!), n>=2
        # use a working space (excluding s and SA) of
        # at most 2.25n+O(1) for a constant alphabet
        def SA_IS(s, SA, n, K, os=0, oa=0):
            def isLms(i):
                return i > 0 and t[i] and not t[i - 1]

            # print("Stage 1");
            t = bytearray(n)
            t[n - 2] = 0
            t[n - 1] = 1
            for i in range(n - 3, -1, -1):
                if s[i + os] < s[i + 1 + os] or s[i + os] == s[i + 1 + os] and t[i + 1] == 1:
                    t[i] = 1
                else:
                    t[i] = 0

            # print("Stage 2");
            bkt = [0] * (K + 1)
            getBuckets(s, bkt, n, K, True, os)

            for i in range(n):
                SA[i + oa] = -1

            for i in range(1, n):
                if isLms(i):
                    bkt[s[i + os]] -= 1
                    SA[bkt[s[i + os]] + oa] = i

            # print("Stage 3");
            induceSAl(t, SA, s, bkt, n, K, False, os, oa)

            induceSAs(t, SA, s, bkt, n, K, True, os, oa)

            bkt = None

            n1 = 0
            for i in range(n):
                if isLms(SA[i + oa]):
                    SA[n1 + oa] = SA[i + oa]
                    n1 += 1

            for i in range(n1, n):
                SA[i + oa] = -1

            # print("Stage 4");
            name = 0
            prev = -1
            for i in range(n1):
                pos = SA[i + oa]
                diff = False
                for d in range(n):
                    if (
                        prev == -1
                        or s[pos + d + os] != s[prev + d + os]
                        or t[pos + d] != t[prev + d]
                    ):
                        diff = True
                        break
                    elif d > 0 and (isLms(pos + d) or isLms(prev + d)):
                        break
                if diff:
                    name += 1
                    prev = pos
                pos = int(pos / 2)
                SA[n1 + pos + oa] = name - 1

            j = n - 1
            for i in range(n - 1, n1 - 1, -1):
                if SA[i + oa] >= 0:
                    SA[j + oa] = SA[i + oa]
                    j -= 1

            s1 = SA1 = SA
            s1off = n - n1
            if name < n1:
                SA_IS(s1, SA1, n1, name - 1, s1off + oa, oa)
            else:
                for i in range(n1):
                    SA1[s1[i + s1off + oa] + oa] = i

            bkt = [0] * (K + 1)

            # print("Stage 5");
            getBuckets(s, bkt, n, K, True, os)

            j = 0
            for i in range(1, n):
                if isLms(i):
                    s1[s1off + oa + j] = i
                    j += 1

            for i in range(n1):
                SA1[i + oa] = s1[oa + s1off + SA1[oa + i]]

            for i in range(n1, n):
                SA[i + oa] = -1

            for i in range(n1 - 1, -1, -1):
                j = SA[oa + i]
                SA[oa + i] = -1
                bkt[s[os + j]] -= 1
                SA[bkt[s[os + j]] + oa] = j

            # print("Stage 6");
            induceSAl(t, SA, s, bkt, n, K, False, os, oa)
            induceSAs(t, SA, s, bkt, n, K, True, os, oa)

        # Given string s and suffix array sa, compute the longest common prefix
        # information in O(N) time.
        def calculateLcp(s, sa):
            n = len(sa)
            k = 0
            lcp = [0] * n
            rank = [0] * n

            for i in range(n):
                rank[sa[i]] = i

            i = 0
            while True:
                if i >= n:
                    break

                if rank[i] == n - 1:
                    k = 0
                    i += 1
                    continue

                j = sa[rank[i] + 1]

                while i + k < n and j + k < n and s[i + k] == s[j + k]:
                    k += 1

                lcp[rank[i]] = k

                i += 1
                if k:
                    k -= 1

            return lcp

        def remap(a, b):
            # Fast SA construction algorithms assume the sequence is
            # numeric, and has an upper value due to the bucket sort.

            # reserve 0 and 1 for the separator and end character.
            nextName = 2

            # Get the set of all characters used.
            items = set(a)
            items.update(b)

            # map each character to its index in the sorted list
            mapping = {}
            for item in sorted(items):
                mapping[item] = nextName
                nextName += 1

            # create a new string that is numeric.
            newString = []
            newString.extend([mapping[item] for item in a])
            newString.append(1)
            newString.extend([mapping[item] for item in b])
            newString.append(0)

            return newString, nextName - 1

        # Here is the code to initialize the suffix array object
        remapped, K = remap(A, B)
        self.sa = [0] * len(remapped)

        if self.length1 + self.length2 > 2:
            SA_IS(remapped, self.sa, len(self.sa), K)
            # self.sa = sais_native(remapped, len(remapped), K)
        self.lcp = calculateLcp(self.seq, self.sa)

    def show(self):
        for i in range(len(self.sa)):
            self.showPosition(i)

    def showPosition(self, saIndex):
        i = saIndex
        p = self.sa[i]
        if self.sa[i] < self.length1:
            s = "A  "
        else:
            s = "  B"
            p -= self.length1
        print(
            "{3} {0:8} {1:5} |{2}".format(
                p, self.lcp[i], json.dumps(self.seq[self.sa[i] : self.sa[i] + 20]), s
            )
        )

    def longestMatches(self):
        # returns, for every position in B, a tuple with the longest matching
        # position in A and the length of that match.
        result = [None] * self.length2

        # forward pass
        lcp = 0
        aIndex = 0
        for i in range(len(self.sa)):
            if self.sa[i] < self.length1:
                # string in A
                lcp = self.lcp[i]
                aIndex = self.sa[i]
            else:
                # string in B.
                result[self.sa[i] - self.length1] = (aIndex, lcp)
                lcp = min(lcp, self.lcp[i])

        # reverse pass
        lcp = 0
        aIndex = 0
        for i in range(len(self.sa) - 1, -1, -1):
            if self.sa[i] < self.length1:
                # string in A
                aIndex = self.sa[i]
                if i > 0:
                    lcp = self.lcp[i - 1]
            else:
                # string in B.
                lcp = min(lcp, self.lcp[i])
                bIndex = self.sa[i] - self.length1
                oldAIndex, oldLcp = result[bIndex]
                if lcp > oldLcp:
                    result[bIndex] = (aIndex, lcp)

        return result


class SAFindLongest:
    def __init__(self, A, B):
        self.matches = SuffixArray(A, B).longestMatches()

    def findLongest(self, positionInB):
        return self.matches[positionInB]


class BruteFindLongest:
    def __init__(self, A, B):
        self.A = A
        self.B = B

    def findLongest(self, positionInB):
        a = self.A

        def prefixLength(position, text):
            l = 0
            while l < len(text) and position + l < len(a) and text[l] == a[position + l]:
                l += 1
            return l

        text = self.B[positionInB:]

        longest = 0
        longestPos = 0
        for p in range(len(self.A)):
            l = prefixLength(p, text)
            if l > longest:
                longest = l
                longestPos = p

        return longestPos, longest


class TableFindLongest:
    def __init__(self, A, B):
        self.A = A
        self.B = B
        minlen = 4
        self.table = {}
        for i in range(len(A) - minlen + 1):
            code = A[i : i + minlen]
            if True or code not in self.table:
                self.table[code] = [i]
            else:
                self.table[code].append(i)

    def findLongest(self, positionInB):
        minlen = 4
        code = self.B[positionInB : positionInB + minlen]
        if code not in self.table:
            return 0, 0

        a = self.A
        lcp = 0
        pos = 0

        def prefixLength(position, text):
            l = 0
            while l < len(text) and position + l < len(a) and text[l] == a[position + l]:
                l += 1
            return l

        text = self.B[positionInB:]
        for index in self.table[code]:
            l = prefixLength(index, text)
            if l > lcp:
                pos = index
                lcp = l

        return pos, lcp


def longestPrefix(A, B):
    l = 0
    m = min(len(A), len(B))
    while l < m and A[l] == B[l]:
        l += 1

    return l


def longestSuffix(A, B):
    l = 0
    la = len(A)
    lb = len(B)
    m = min(la, lb)
    while l < m and A[la - l - 1] == B[lb - l - 1]:
        l += 1

    return l


def DiffEncode(a, b, Algorithm):
    start = time.time()
    cmds = []
    j = 0
    insertFrom = -1
    check = bytearray()
    totalCost = 0
    minimumMatch = 6

    def copy(start, length):
        nonlocal totalCost, check
        totalCost += 9
        cmds.append("COPY {0}, {1}".format(start, length))
        check += a[start : start + length]

    def insert(start, length):
        nonlocal totalCost, check
        totalCost += 1 + 5 + length
        cmds.append("INSERT {0}, {1}".format(start, length))
        check += b[start : start + length]

    j = prefix = longestPrefix(a, b)
    if prefix > 0:
        copy(0, prefix)

    if prefix != len(b):
        suffix = min(longestSuffix(a, b), len(b) - prefix)
        limit = len(b) - suffix

        if suffix + prefix != len(b):
            if len(b) - suffix > prefix:
                algorithm = Algorithm(a[prefix : len(a) - suffix], b[prefix:limit])

                while j < limit:
                    position, length = algorithm.findLongest(j - prefix)
                    position += prefix
                    if length < minimumMatch:
                        if insertFrom == -1:
                            insertFrom = j
                        j += 1
                    else:
                        if insertFrom >= 0:
                            insert(insertFrom, j - insertFrom)
                            insertFrom = -1
                        copy(position, length)
                        j += length

                if insertFrom >= 0:
                    insert(insertFrom, j - insertFrom)

        if limit < len(b):
            copy(len(a) - suffix, suffix)

    if check != b:
        print("\n".join(cmds))
        print("ERROR")
        raise (Exception("Error"))

    totalTime = time.time() - start

    return "\n".join(cmds), totalTime, totalCost


import os


def performance():
    # Output a CSV file with the N and the time to construct a suffix
    # array of a random string of length N.
    for i in range(1000, 1000000, 1000):
        s = os.urandom(i // 2)
        t = os.urandom(i // 2)
        start = time.time()
        s = SuffixArray(s, t)
        print("{0},{1}".format(i, time.time() - start))


if "--performance" in sys.argv:
    performance()
    sys.exit(0)

Algorithms = {
    "suffixarray": SAFindLongest,
    "bruteforce": BruteFindLongest,
    "table": TableFindLongest,
}

import argparse

parser = argparse.ArgumentParser()
parser.add_argument("file1", type=str, help="First file to compare")
parser.add_argument("file2", type=str, help="Second file to compare")
parser.add_argument("--performance", help="Performance test SA construction")
parser.add_argument(
    "-a",
    "--algorithm",
    type=str,
    help="Algorithm to use",
    choices=Algorithms.keys(),
    default="suffixarray",
)
args = parser.parse_args()

f1 = open(args.file1, "rb").read()
f2 = open(args.file2, "rb").read()

cmds, time, cost = DiffEncode(f1, f2, Algorithms[args.algorithm])
print("Cmds:" + cmds)
print("Total time: {0}".format(time))
print("Instruction cost: {0}".format(cost))
