package template.primitve.generated.datastructure;

public class LongSubArray {
    private long[] data;
    private int l;
    private int r;

    public LongSubArray(long[] data, int l, int r) {
        this.data = data;
        this.l = l;
        this.r = r;
    }

    public int length() {
        return r - l + 1;
    }

    public long get(int i) {
        return data[l + i];
    }

    public void set(int i, long x) {
        data[i] = x;
    }

    public LongSubArray subArray(int ll, int rr) {
        return new LongSubArray(data, ll + l, rr + l);
    }
}
