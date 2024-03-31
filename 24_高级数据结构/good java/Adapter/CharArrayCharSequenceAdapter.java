package template.string;

public class CharArrayCharSequenceAdapter implements CharSequence {
    char[] s;
    int l;
    int r;

    public CharArrayCharSequenceAdapter(char[] s, int l, int r) {
        this.s = s;
        this.l = l;
        this.r = r;
    }

    @Override
    public int length() {
        return r - l + 1;
    }

    @Override
    public char charAt(int index) {
        return s[index + l];
    }

    @Override
    public CharSequence subSequence(int start, int end) {
        return new CharArrayCharSequenceAdapter(s, l + start, l + end - 1);
    }
}
