"a.1", 2) == 2
        assert persist.query("a.4", 1) == 1
        assert persist.query("a.4", 3) == 3
        assert persist.query(