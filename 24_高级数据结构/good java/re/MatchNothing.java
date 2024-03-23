package template.string.re;

public class MatchNothing extends AbstractState {
    Transfer invalid;

    public MatchNothing(Transfer invalid) {
        this.invalid = invalid;
    }

    @Override
    public Transfer next(int c) {
        return invalid;
    }
}
