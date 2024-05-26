package template.datastructure;

import template.math.DigitUtils;
import template.primitve.generated.datastructure.IntToIntegerFunction;

import java.util.ArrayList;
import java.util.List;

public class PreProd {
    int mod;
    int[] fact;
    int[] invFact;
    int one;
    int n;

    public PreProd(int n, int mod) {
        fact = new int[n];
        invFact = new int[n];
        this.mod = mod;
        this.one = 1 % mod;
    }

    /**
     * O(n+\log_2 mod)
     *
     * @param func
     * @param n
     */
    public void init(IntToIntegerFunction func, int n) {
        this.n = n;
        if (n == 0) {
            return;
        }
        for (int i = 0; i < n; i++) {
            fact[i] = func.apply(i);
            if (fact[i] == 0) {
                fact[i] = 1;
            }
            if (i > 0) {
                fact[i] = (int) ((long) fact[i] * fact[i - 1] % mod);
            }
        }
        invFact[n - 1] = (int) DigitUtils.modInverse(fact[n - 1], mod);
        for (int i = n - 2; i >= 0; i--) {
            invFact[i] = (int) ((long) invFact[i + 1] * func.apply(i + 1) % mod);
        }
    }

    /**
     * O(1)
     */
    public int prefix(int i) {
        if (i >= n) {
            i = n - 1;
        }
        if (i < 0) {
            return one;
        }
        return fact[i];
    }

    /**
     * O(1)
     */
    public int invPrefix(int i) {
        if (i >= n) {
            i = n - 1;
        }
        if (i < 0) {
            return one;
        }
        return invFact[i];
    }

    /**
     * O(1)
     */
    public int interval(int l, int r) {
        if (l < 0) {
            l = 0;
        }
        if (r >= n) {
            r = n - 1;
        }
        if (l > r) {
            return one;
        }
        long ans = fact[r];
        if (l > 0) {
            ans = ans * invFact[l - 1] % mod;
        }
        return (int) ans;
    }

    /**
     * O(1)
     */
    public int intervalInverse(int l, int r) {
        if (l < 0) {
            l = 0;
        }
        if (r >= n) {
            r = n - 1;
        }
        if (l > r) {
            return one;
        }
        long ans = invFact[r];
        if (l > 0) {
            ans = ans * fact[l - 1] % mod;
        }
        return (int) ans;
    }

    @Override
    public String toString() {
        List<Integer> list = new ArrayList<>(n);
        for (int i = 0; i < n; i++) {
            list.add(interval(i, i));
        }
        return list.toString();
    }
}
