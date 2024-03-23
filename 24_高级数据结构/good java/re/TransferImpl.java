package template.string.re;

public class TransferImpl implements Transfer {
    State state;

    public TransferImpl(State state) {
        this.state = state;
    }

    public TransferImpl() {
    }

    @Override
    public State get() {
        return state;
    }

    @Override
    public void set(State state) {
        assert this.state == null;
        this.state = state;
    }

    @Override
    public String toString() {
        if(state == null){
            return null;
        }
        return "" + state.id();
    }
}
