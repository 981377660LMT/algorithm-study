明了胜于晦涩 简洁胜于复杂 扁平胜于嵌套 稀疏胜于拥挤 可读性很重要

Python 有用特殊的语法能够让 else 语块在循环体结束的时候立刻得到执行。
循环体后的 else 语块只有`在循环体没有触发 break 语句的时候才会执行。`
避免在循环体的后面使用 else 语块，因为这样的表达不直观，而且容易误导读者。

使用 try/except/else 块可以使得代码中对哪种一场要被处理，哪种异常要往上抛出的处理变得更加清晰。

```Python
def load_json_key(data, key):
    try:
        result_dict = json.loads(data)  # 可能引发ValueError异常
    except ValueError as e:
        raise KeyError from e
    else:
        return result_dict[key]       # 可能引发KeyError
```

else 块经常会被用于在 try 块成功运行后添加额外的行为，但是要确保代码会在 finally 块之前得到运行。
