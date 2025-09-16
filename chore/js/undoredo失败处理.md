# `undoredo失败处理`

好的，我们来讨论基于“一个线性历史记录（如 `std::vector` 或数组）+ 一个当前位置指针（ptr，通常是索引）”实现的 `undo/redo` 系统。

这种模型通常是这样工作的：

- 一个 `std::vector<Command*>` 或 `std::vector<std::unique_ptr<Command>>` 存储了所有已执行的命令。
- 一个整数 `current_position_ptr`（或叫 `history_index`）指向**当前状态**在历史记录中的位置。它指向的是**最后一个被成功执行或重做的命令**。
- **Undo**：调用 `history[current_position_ptr]` 的 `undo()` 方法，然后将 `current_position_ptr` 向后移动（`--`）。
- **Redo**：将 `current_position_ptr` 向前移动（`++`），然后调用 `history[current_position_ptr]` 的 `execute()` 方法。
- **新命令**：在 `current_position_ptr` 之后添加新命令，并丢弃该位置之后的所有旧命令（即 `redo` 历史），然后将 `current_position_ptr` 移动到新命令的位置。

在这种模型下，当 `undo` 或 `redo` 操作失败时，对指针 `current_position_ptr` 的处理是关键。

### 核心原则

**`current_position_ptr` 必须始终指向一个代表了当前应用程序真实状态的有效历史点。** 如果操作失败，指针绝不能移动到代表一个未达成状态的位置。

---

### 场景一：Undo 失败

假设历史记录为 `[Cmd1, Cmd2, Cmd3]`，`current_position_ptr` 为 `2`（指向 `Cmd3`）。

1.  **用户点击 Undo。**
2.  系统尝试执行 `history[2]->undo()`（即 `Cmd3->undo()`）。
3.  **操作失败**，并抛出异常或返回错误。

**错误处理：**

- **指针 `current_position_ptr` 不能移动！** 它必须保持在 `2`。因为 `undo` 失败了，应用程序的当前状态仍然是 `Cmd3` 执行完毕后的状态。如果将指针移动到 `1`，历史记录就会谎称当前状态是 `Cmd2` 完成后的状态，这是错误的。
- **处理未来的历史（Redo 历史）：** 此时，`Cmd3` 这个命令已经变得不可信赖。我们无法保证它能被再次 `redo`。因此，最安全的操作是**截断历史记录**。
  - 丢弃 `Cmd3` 以及它之后的所有命令（虽然此时它后面没有命令）。
  - 将 `history` 的大小调整为 `current_position_ptr`（即 `history.resize(2)`）， effectively 丢弃 `Cmd3`。
  - `current_position_ptr` 需要更新为新的末尾位置，即 `current_position_ptr--`（变为 `1`）。

**结论：当 `undo` 失败时，将失败的命令及其之后的所有历史记录全部清除。`current_position_ptr` 指向新的历史记录末尾。**

#### 伪代码示例 (Undo 失败)

```cpp
// std::vector<std::unique_ptr<Command>> history;
// int current_position_ptr = -1; // -1 表示初始状态

void onUndo() {
    if (current_position_ptr < 0) return;

    try {
        history[current_position_ptr]->undo();

        // 如果成功，指针后移
        current_position_ptr--;

    } catch (const std::exception& e) {
        std::cerr << "Undo failed: " << e.what() << std::endl;

        // 失败处理：
        // 1. 获取失败命令的索引
        int failed_cmd_index = current_position_ptr;

        // 2. 截断历史记录，删除失败的命令及其之后的所有内容
        //    (从 failed_cmd_index 开始的所有元素)
        history.resize(failed_cmd_index);

        // 3. 更新指针到新的历史末尾
        current_position_ptr = failed_cmd_index - 1;
    }
    updateUI();
}
```

---

### 场景二：Redo 失败

假设历史记录为 `[Cmd1, Cmd2, Cmd3]`，`current_position_ptr` 为 `0`（指向 `Cmd1`）。用户已经 `undo` 了两次。

1.  **用户点击 Redo。**
2.  系统首先移动指针：`current_position_ptr++`（现在是 `1`）。
3.  然后尝试执行 `history[1]->execute()`（即 `Cmd2->execute()`）。
4.  **操作失败**，并抛出异常或返回错误。

**错误处理：**

- **指针 `current_position_ptr` 必须回滚！** 因为 `redo` 失败，应用程序的真实状态仍然是 `Cmd1` 执行完毕后的状态。指针必须移回到 `0`。
- **处理未来的历史（Redo 历史）：** `Cmd2` 这个命令已经不可信。它以及它之后的所有 `redo` 历史（`Cmd3`）都必须被认为是无效的。
  - 截断历史记录，丢弃 `Cmd2` 和 `Cmd3`。
  - 将 `history` 的大小调整为 `current_position_ptr` 的原始值加一（即 `history.resize(1)`）。

**结论：当 `redo` 失败时，将失败的命令及其之后的所有历史记录全部清除。`current_position_ptr` 恢复到操作前的位置。**

#### 伪代码示例 (Redo 失败)

```cpp
void onRedo() {
    if (current_position_ptr >= history.size() - 1) return;

    // 记录原始指针位置
    int original_ptr = current_position_ptr;

    try {
        // 1. 指针前移
        current_position_ptr++;
        // 2. 执行命令
        history[current_position_ptr]->execute();

    } catch (const std::exception& e) {
        std::cerr << "Redo failed: " << e.what() << std::endl;

        // 失败处理：
        // 1. 截断历史记录，删除失败的命令及其之后的所有内容
        //    history 的新大小应该是 original_ptr + 1
        history.resize(original_ptr + 1);

        // 2. 将指针恢复到操作前的位置
        current_position_ptr = original_ptr;
    }
    updateUI();
}
```

### 总结

无论是哪种场景，核心思想都是**当操作失败时，销毁所有未来的、不可靠的历史记录，并将 `current_position_ptr` 调整到能准确反映当前程序状态的最后一个有效位置**。这可以防止历史记录被污染，保证系统的稳定性。
