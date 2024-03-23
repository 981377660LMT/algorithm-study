package template.string.re;

public class CharacterState extends AbstractState {
    Transfer invalid;
    Transfer next;
    int character;

    public CharacterState(Transfer invalid, Transfer next, int character) {
        this.invalid = invalid;
        this.next = next;
        this.character = character;
    }

    @Override
    public Transfer next(int c) {
        return c == character ? next : invalid;
    }

    @Override
    public String toString() {
        return super.toString() + id() + "-" + ((char) character) + "->" + next.toString() + "\n";
    }
}
