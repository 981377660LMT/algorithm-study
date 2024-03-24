package template.string;

import template.algo.DivisionDescending;
import template.string.KMPAutomaton;

public class StringPeriod {
    /**
     * <pre>
     * 查询字符串的最小周期p，满足s[i+p]=s[i]，如果i+p<n
     * 时间复杂度：O(n)
     * 空间复杂度：O(n)
     * </pre>
     */
    public static int minPeriod(char[] s, int n) {
        KMPAutomaton kmp = new KMPAutomaton(n);
        for (int i = 0; i < n; i++) {
            kmp.build(s[i]);
        }
        return n - kmp.maxBorder(n);
    }

    /**
     * <pre>
     * 查询字符串的最小旋转周期p，满足s[(i+p)%n]=s[i]
     * 时间复杂度：O(n\log_2n)
     * 空间复杂度：O(\log_2n)
     * </pre>
     */
    public static int minRotatePeriod(char[] s, int n) {
        return (int) DivisionDescending.find(n, m -> {
            for (int i = 0, j = (int) m; i < n; i++, j++) {
                if (j >= n) {
                    j -= n;
                }
                if (s[i] != s[j]) {
                    return false;
                }
            }
            return true;
        });
    }
    /**
     * <pre>
     * 查询字符串的最小回文旋转周期p，满足s[p..n)+s[0..p)是回文
     * 前置条件：s[0..n)是回文
     * 时间复杂度：O(n\log_2n)
     * 空间复杂度：O(\log_2n)
     * </pre>
     */
    public static int minPalindromeRotatePeriod(char[] s, int n) {
        assert isPalindrome(s, n);
        int rp = minRotatePeriod(s, n);
        if (rp % 2 == 0) {
            rp /= 2;
        }
        return rp;
    }

    /**
     * <pre>
     * 判断字符串s是否是回文
     * 时间复杂度：O(n)
     * 空间复杂度：O(1)
     * </pre>
     */
    public static boolean isPalindrome(char[] s, int n){
        int l = 0;
        int r = n - 1;
        while(l < r){
            if(s[l] != s[r]){
                return false;
            }
            l++;
            r--;
        }
        return true;
    }
}
