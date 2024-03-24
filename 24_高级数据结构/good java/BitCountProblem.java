package template.problem;

import template.binary.Log2;

public class BitCountProblem {
    public static final long INF = (long) 2e18;

    public static long countZero(long n) {
        return countLength(n) - countOne(n);
    }

    public static long countLength(long n) {
        if (n == 0) {
            return 1;
        }
        long ans = 0;
        int len = Log2.floorLog(n);
        for (int i = 0; i <= len; i++) {
            long l = 1L << i;
            long r = Math.min(l * 2 - 1, n);
            ans += (i + 1) * (r - l + 1);
        }
        return ans;
    }

    /**
     * count how many 1 occur in binary form of 0, 1, ..., n
     */
    public static long countOne(long n) {
        if (n == 0) {
            return 0;
        }
        long ans = 0;
        int len = Log2.floorLog(n);
        for (int i = 0; i <= len; i++) {
            long bit = 1L << i;
            long remain = n;
            long block = remain / (2 * bit);
            ans += block * bit;
            remain %= (2 * bit);
            ans += Math.max(0, remain - bit + 1);
        }
        return ans;
    }

    public static long countOne(long l, long r) {
        long ans = countOne(r);
        if (l > 0) {
            ans -= countOne(l - 1);
        }
        return ans;
    }

    public static long countLength(long l, long r) {
        long ans = countLength(r);
        if (l > 0) {
            ans -= countLength(l - 1);
        }
        return ans;
    }

    public static long countZero(long l, long r) {
        long ans = countZero(r);
        if (l > 0) {
            ans -= countZero(l - 1);
        }
        return ans;
    }

}
