package template.string;

public interface IntSequence extends Comparable<IntSequence> {
    int get(int i);

    int length();

    IntSequence subsequence(int l, int r);

    @Override
    default int compareTo(IntSequence o) {
        int n = length();
        int m = o.length();
        for (int i = 0, to = Math.min(n, m); i < to; i++) {
            int ans = Integer.compare(get(i), o.get(i));
            if (ans != 0) {
                return ans;
            }
        }
        return Integer.compare(n, m);
    }
}
