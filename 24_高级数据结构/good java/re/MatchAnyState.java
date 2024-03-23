package template.string.re;

public class MatchAnyState extends AbstractState {
    Transfer next;

    public MatchAnyState(Transfer next) {
        this.next = next;
    }

    @Override
    public Transfer next(int c) {
        return next;
    }

    @Override
    public String toString() {
        return super.toString() + id() + "-*->" + next.toString() + "\n";
    }
}
