package template.problem;

/**
 * There are two row of numbers, each row contains 1, ...,n.
 * <br>
 * Given m, for row 1 number x and row 2 number y, if x + y <= m, there is an edge (x,y) between them.
 * <br>
 * You're asked to find a largest matching on this bipartite graph.
 */
public class NumberSumLeqMatch {
    long n;
    long m;
    long used;
    long type1;

    public NumberSumLeqMatch(long n, long m) {
        this.n = n;
        this.m = m;
        used = Math.min(n, m - 1);
        type1 = m - used - 1;
    }

    /**
     * Get the maximum matching
     */
    public long maxMatching() {
        return used;
    }

    /**
     * Get the partner of x, if not exist, return -1
     */
    public long partner(long x) {
        if (x > used) {
            return -1;
        }
        if (x <= type1) {
            return x;
        }
        return m - x;
    }

    /**
     * Get the maximum matching if remove a from row 1 and remove b from row 2
     */
    public long maxMatching(long a, long b) {
        if (a >= m && b >= m) {
            return used;
        }
        if (a + b < m && type1 <= 1) {
            return used - 2;
        }
        return used - 1;
    }
}
