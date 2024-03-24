package template.problem;

import template.math.Factorization;
import template.math.GCDs;
import template.math.LCMs;
import template.primitve.generated.datastructure.IntegerArrayList;
import template.primitve.generated.datastructure.IntegerDequeImpl;
import template.primitve.generated.datastructure.IntegerIterator;
import template.primitve.generated.datastructure.IntegerMultiWayDeque;
import template.primitve.generated.datastructure.IntegerMultiWayStack;
import template.primitve.generated.datastructure.IntegerVersionArray;

import java.util.Arrays;

public class GcdLcmProblem {
    private IntegerMultiWayStack primeFactors;
    private IntegerMultiWayStack factors;
    private IntegerArrayList allFactors = new IntegerArrayList(20);
    private int[] cntOfMultiple;
    private int m;
    private int[] seq;
    public static final int INF = (int) 1e9;
    private int[] coprime;
    private IntegerMultiWayDeque indexesOfSeq;

    public GcdLcmProblem(int m, int[] seq) {
        this.m = m;
        this.seq = seq;
        primeFactors = Factorization.factorizeRangePrime(m);
        factors = Factorization.factorizeRange(m);
        prepareCntOfMultiple();
        coprime = new int[m + 1];
        Arrays.fill(coprime, -1);
        indexesOfSeq = new IntegerMultiWayDeque(m + 1, seq.length);
        for (int i = 0; i < seq.length; i++) {
            indexesOfSeq.addLast(seq[i], i);
        }
    }

    private void prepareCntOfMultiple() {
        cntOfMultiple = new int[m + 1];
        for (int x : seq) {
            cntOfMultiple[x]++;
        }
        for (int i = 1; i <= m; i++) {
            for (int j = i + i; j <= m; j += i) {
                cntOfMultiple[i] += cntOfMultiple[j];
            }
        }
    }

    /**
     * Find how many number in sequence coprime with x. $O(2^{\phi(x)})$ while $\phi(x)$ means the number
     * of distinct prime factors of x
     */
    public int coprime(int x) {
        if (coprime[x] == -1) {
            factorize(x);
            coprime[x] = ie(cntOfMultiple, allFactors.size() - 1, 1, 0);
        }
        return coprime[x];
    }

    private int[] dp;

    /**
     * return the minimum size of subset to ensure all element in subset has gcd | x, INF means impossible.
     */
    public int minimumSizeOfSubsetWhoseGCDDivisibleBy(int x) {
        if (dp == null) {
            dp = new int[m + 1];
            Arrays.fill(dp, INF);
            for (int e : seq) {
                dp[e] = 1;
            }
            for (int i = m; i >= 1; i--) {
                for (int j = i + i; j <= m; j += i) {
                    if (coprime(j / i) > 0) {
                        dp[i] = Math.min(dp[i], dp[j] + 1);
                    }
                }
            }
        }
        return dp[x];
    }

    private void factorize(int x) {
        allFactors.clear();
        allFactors.addAll(primeFactors.iterator(x));
    }

    /**
     * return any pair of indexes (i, j), while i != j and gcd(seq[i], seq[j]) = 1, null for non-existence.
     */
    public int[] findAnyPairIndexesCoprime() {
        return findAnyPairIndexesCoprime(seq);
    }

    private int[] findAnyPairIndexesCoprime(int[] seq) {
        int[] cnts = new int[m + 1];
        for (int i = 0; i < seq.length; i++) {
            factorize(seq[i]);
            int total = ie(cnts, allFactors.size() - 1, 1, 0);
            if (total > 0) {
                for (int j = 0; ; j++) {
                    if (GCDs.gcd(seq[i], seq[j]) == 1) {
                        return new int[]{i, j};
                    }
                }
            }
            add(cnts, seq[i]);
        }
        return null;
    }

    /**
     * <pre>
     * Find a pair of indexes (i, j) satisfy i != j and make gcd(seq[i], seq[j]) as large as possible
     * return null for non-existence.
     *  </pre>
     */
    public int[] findAnyPairIndexesWhileGCDMaximized() {
        if (seq.length < 2) {
            return null;
        }

        int g = -1;
        for (int i = m; i >= 1; i--) {
            if (cntOfMultiple[i] >= 2) {
                g = i;
                break;
            }
        }

        IntegerArrayList ans = new IntegerArrayList(2);
        for (int i = 0; i < seq.length && ans.size() < 2; i++) {
            if (seq[i] % g == 0) {
                ans.add(i);
            }
        }
        return ans.toArray();
    }

    /**
     * Find a pair of indexes (i, j) satisfy i != j and make gcd(seq[i], seq[j]) as small as possible, return null for non-existence.
     */
    public int[] findAnyPairIndexesWhileGCDMinimized() {
        if (seq.length < 2) {
            return null;
        }
        int g = -1;
        for (int i = 1; i <= m; i++) {
            if (minimumSizeOfSubsetWhoseGCDDivisibleBy(i) <= 2) {
                g = i;
                break;
            }
        }

        IntegerArrayList list = new IntegerArrayList(seq.length);
        for (int i = 0; i < seq.length; i++) {
            if (seq[i] % g == 0) {
                list.add(seq[i]);
            }
        }

        int[] pair = findAnyPairIndexesCoprime(list.toArray());
        return new int[]{indexesOfSeq.peekFirst(list.get(pair[0]) * g), indexesOfSeq.peekLast(list.get(pair[1]) * g)};
    }


    /**
     * Find a pair of indexes (i, j) satisfy i != j and make lcm(seq[i], seq[j]) as large as possible, return null for non-existence.
     */
    public int[] findAnyPairIndexesWhileLCMMaximized() {
        if (seq.length < 2) {
            return null;
        }
        int size = 0;
        for (int i = 1; i <= m; i++) {
            size += m / i;
        }

        IntegerMultiWayDeque deque = new IntegerMultiWayDeque(m + 1, size);
        for (int i = 1; i <= m; i++) {
            for (int j = i; j <= m; j += i) {
                if (!indexesOfSeq.isEmpty(j)) {
                    deque.addLast(i, j / i);
                }
            }
        }

        long lcm = -1;
        int v1 = 0;
        int v2 = 0;

        for (int i = m; i >= 1; i--) {
            if (!indexesOfSeq.isEmpty(i) && indexesOfSeq.peekFirst(i) != indexesOfSeq.peekLast(i)) {
                lcm = v1 = v2 = i;
                break;
            }
        }

        IntegerVersionArray iva = new IntegerVersionArray(m + 1);
        IntegerDequeImpl stack = new IntegerDequeImpl(m);
        for (int i = 1; i <= m; i++) {
            iva.clear();
            stack.clear();
            while (!deque.isEmpty(i)) {
                int last = deque.removeLast(i);
                factorize(last);
                int total = ie(iva, allFactors.size() - 1, 1, 0);
                if (total > 0) {
                    int pop = 0;
                    while (total > 0) {
                        pop = stack.removeLast();
                        if (GCDs.gcd(pop, last) == 1) {
                            total--;
                        }
                        add(iva, pop, -1);
                    }
                    long l = LCMs.lcm(last, pop) * i;
                    if (l > lcm) {
                        lcm = l;
                        v1 = last * i;
                        v2 = pop * i;
                    }
                }
                add(iva, last, 1);
                stack.addLast(last);
            }
        }

        return new int[]{indexesOfSeq.peekFirst(v1), indexesOfSeq.peekLast(v2)};
    }

    /**
     * Find a pair of indexes (i, j) satisfy i != j and make lcm(seq[i], seq[j]) as small as possible, return null for non-existence.
     */
    public int[] findAnyPairIndexesWhileLCMMinimized() {
        if (seq.length < 2) {
            return null;
        }
        int size = 0;
        for (int i = 1; i <= m; i++) {
            size += m / i;
        }

        IntegerMultiWayDeque deque = new IntegerMultiWayDeque(m + 1, size);
        for (int i = 1; i <= m; i++) {
            for (int j = i; j <= m; j += i) {
                if (!indexesOfSeq.isEmpty(j)) {
                    deque.addLast(i, j / i);
                }
            }
        }

        long lcm = INF;
        int v1 = 0;
        int v2 = 0;

        for (int i = 1; i <= m; i++) {
            if (!indexesOfSeq.isEmpty(i) && indexesOfSeq.peekFirst(i) != indexesOfSeq.peekLast(i)) {
                lcm = v1 = v2 = i;
                break;
            }
        }

        IntegerVersionArray iva = new IntegerVersionArray(m + 1);
        IntegerDequeImpl stack = new IntegerDequeImpl(m);
        for (int i = 1; i <= m; i++) {
            iva.clear();
            stack.clear();
            while (!deque.isEmpty(i)) {
                int last = deque.removeFirst(i);
                factorize(last);
                int total = ie(iva, allFactors.size() - 1, 1, 0);
                if (total > 0) {
                    int pop = 0;
                    while (total > 0) {
                        pop = stack.removeLast();
                        if (GCDs.gcd(pop, last) == 1) {
                            total--;
                        }
                        add(iva, pop, -1);
                    }
                    long l = LCMs.lcm(last, pop) * i;
                    if (l < lcm) {
                        lcm = l;
                        v1 = last * i;
                        v2 = pop * i;
                    }
                }
                add(iva, last, 1);
                stack.addLast(last);
            }
        }

        return new int[]{indexesOfSeq.peekFirst(v1), indexesOfSeq.peekLast(v2)};
    }

    private void add(int[] cnts, int x) {
        for (IntegerIterator iterator = factors.iterator(x); iterator.hasNext(); ) {
            int d = iterator.next();
            cnts[d]++;
        }
    }

    private void add(IntegerVersionArray cnts, int x, int c) {
        for (IntegerIterator iterator = factors.iterator(x); iterator.hasNext(); ) {
            int d = iterator.next();
            cnts.modify(d, c);
        }
    }

    private int ie(int[] cnts, int i, int divisor, int numberOfFactors) {
        if (i < 0) {
            int ans = cnts[divisor];
            if ((numberOfFactors & 1) == 1) {
                ans = -ans;
            }
            return ans;
        }
        return ie(cnts, i - 1, divisor, numberOfFactors) +
                ie(cnts, i - 1, divisor * allFactors.get(i),
                        numberOfFactors + 1);
    }

    private int ie(IntegerVersionArray cnts, int i, int divisor, int numberOfFactors) {
        if (i < 0) {
            int ans = cnts.get(divisor);
            if ((numberOfFactors & 1) == 1) {
                ans = -ans;
            }
            return ans;
        }
        return ie(cnts, i - 1, divisor, numberOfFactors) +
                ie(cnts, i - 1, divisor * allFactors.get(i),
                        numberOfFactors + 1);
    }
}
