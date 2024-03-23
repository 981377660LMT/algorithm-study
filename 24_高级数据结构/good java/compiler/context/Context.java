package template.compiler.context;

import java.util.HashMap;
import java.util.Map;

public class Context {
    private Map<String, Variable> map = new HashMap<>();
    private Context parent;

    public Context() {
    }

    public Context(Context parent) {
        this.parent = parent;
    }

    public Variable get(String s, boolean lookUpParent, boolean createIfNotExist) {
        Variable ans = map.get(s);
        if (ans == null && lookUpParent && parent != null) {
            ans = parent.get(s, true, false);
        }
        if (ans == null && createIfNotExist) {
            ans = new Variable(s, null);
            map.put(ans.getName(), ans);
        }
        return ans;
    }
}
