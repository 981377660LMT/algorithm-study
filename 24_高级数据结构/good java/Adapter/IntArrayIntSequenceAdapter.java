package template.string;

import java.util.Arrays;

public class IntArrayIntSequenceAdapter implements IntSequence {
    int[] data;
    int l;
    int r;

    public IntArrayIntSequenceAdapter(int[] data, int l, int r) {
        this.data = data;
        this.l = l;
        this.r = r;
    }

    @Override
    public int get(int i) {
        return data[l + i];
    }

    @Override
    public int length() {
        return r - l + 1;
    }

    @Override
    public IntSequence subsequence(int l, int r) {
        return new IntArrayIntSequenceAdapter(data, this.l + l, this.l + r);
    }

    @Override
    public String toString() {
        return Arrays.toString(Arrays.copyOfRange(data, l, r + 1));
    }
}
