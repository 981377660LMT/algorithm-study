package template.algo;

public interface UndoOperation {
  static final UndoOperation NIL = new UndoOperation() {
    @Override
    public void apply() {

    }

    @Override
    public void undo() {

    }
  };

  void apply();

  void undo();
}
