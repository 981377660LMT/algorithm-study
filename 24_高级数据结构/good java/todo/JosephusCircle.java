package template.problem;

import template.math.DigitUtils;

/**
 * O(min(n,k ln(n/k)))
 */
public class JosephusCircle {
    public static int survivorBF(int n, int k) {
        if (n == 1) {
            return 0;
        }
        return (survivorBF(n - 1, k) + k) % n;
    }

    /**
     * There are n people form a circle, numbered with 0, 1, 2, ..., n -1. And we will start with 0 in
     * order, if some guy is the k-th from the start, then kill him, start with the next guy. <br>
     * This method find who will die last
     */
    public static int survivor(int n, int k) {
        if (k == 1) {
            return n - 1;
        }
        if (k >= n) {
            return survivorBF(n, k);
        }
        int survivor = survivorBF(k, k);
        for (int i = k + 1; i <= n; i++) {
            int t = DigitUtils.ceilDiv(i - survivor - 1, k - 1);
            // i + t - 1 <= n => t <= n + 1 - i
            t = Math.min(t, n + 1 - i);
            i += t - 1;
            survivor = (survivor + k * t) % i;
        }
        return survivor;
    }

    public static int dieAtRoundBF(int n, int k, int round) {
        if (round == 1) {
            return (k - 1) % n;
        }
        int who = dieAtRoundBF(n - 1, k, round - 1);
        return (who + k) % n;
    }

    /**
     * There are n people form a circle, numbered with 0, 1, 2, ..., n -1. And we will start with 0 in
     * order, if some guy is the k-th from the start, then kill him, start with the next guy. <br>
     * This method find who will die on round <br>
     */
    public static int dieAtRound(int n, int k, int round) {
        int alive = n - round + 1;
        if (k == 1) {
            return round - 1;
        }
        int who = (k - 1) % alive;
        for (int i = alive + 1; i <= n; i++) {
            int t = DigitUtils.ceilDiv(i - who - 1, k - 1);
            t = Math.min(t, n + 1 - i);
            i += t - 1;
            who = (who + k * t) % i;
        }
        return who;
    }

    public static int dieTimeBF(int n, int k, int who) {
        if ((k - 1) % n == who) {
            return 1;
        }
        return dieTimeBF(n - 1, k, DigitUtils.mod(who - k, n)) + 1;
    }

    /**
     * There are n people form a circle, numbered with 0, 1, 2, ..., n -1. And we will start with 0 in
     * order, if some guy is the k-th from the start, then kill him, start with the next guy. <br>
     * This method find who will die on which round <br>
     * O(min(n, kln(n/k)))
     */
    public static int dieTime(int n, int k, int who) {
        if ((who + 1) % k == 0) {
            return (who + 1) / k;
        }
        int turn = n / k;
        if (turn <= 1) {
            return dieTimeBF(n, k, who);
        }
        int next;
        if (who >= turn * k) {
            next = who - turn * k;
        } else {
            next = n + who - (who + 1) / k - turn * k;
        }
        return dieTime(n - turn, k, DigitUtils.mod(next, n)) + turn;
    }
}
