public class Foo {
 
  public static interface Modify {
    int index();
}

public static interface Query {
    public int left();
    public int right();
}

public static interface VersionQuery extends Query {
    public int version();
}

public static interface State<Q extends Query> {
    public void answer(Q q);

    public void add(int i);

    public void remove(int i);
}

public static interface AddOnlyState<Q extends Query> extends State<Q> {
    public void save();

    public void rollback();

}

public static interface RemoveOnlyState<Q extends Query> extends State<Q> {
    public void save();
    public void rollback();

}

public static interface ModifiableState<Q extends VersionQuery, M extends Modify> extends State<Q> {
    public void apply(M m);
    public void revoke(M m);
}
}