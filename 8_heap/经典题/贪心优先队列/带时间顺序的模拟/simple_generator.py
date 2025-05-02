def simple_generator():
    try:
        # Step 1: 初始阶段，等待第一次调用 `next()` 或 `send()`
        print("Generator started.")
        value = yield "Step 1: Please send a value"

        # Step 2: 处理 `send()` 传入的值
        print(f"Received value from send(): {value}")
        value = yield "Step 2: Please send another value"

        # Step 3: 处理第二个 `send()` 传入的值
        print(f"Received value from send(): {value}")
        value = yield "Step 3: Please send a final value"

        # Step 4: 处理错误
        print(f"Received value from send(): {value}")
        raise ValueError("An error occurred!")

    except ValueError as e:
        print(f"Caught an exception in generator: {e}")
        yield "Error handled, generator is ending."

    finally:
        # Step 5: 结束生成器并返回值
        print("Generator is returning a value.")
        return "Generator has finished."


# 使用生成器
gen = simple_generator()

# Step 1: 使用 `next()` 启动生成器
print(next(gen))  # 输出: "Step 1: Please send a value"

# Step 2: 使用 `send()` 向生成器传入一个值
print(gen.send(10))  # 输出: "Received value from send(): 10" -> "Step 2: Please send another value"

# Step 3: 使用 `send()` 向生成器传入另一个值
print(gen.send(20))  # 输出: "Received value from send(): 20" -> "Step 3: Please send a final value"

# Step 4: 使用 `send()` 向生成器传入一个值并触发错误
print(gen.send(30))  # 输出: "Received value from send(): 30" -> 触发 ValueError

# Step 5: 错误处理和结束
try:
    print(next(gen))  # 输出: "Error handled, generator is ending." -> "Generator has finished."
except StopIteration as e:
    print(f"Generator finished with message: {e.value}")
