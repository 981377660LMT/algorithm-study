// 如果没有 key，你就只能通过 Plugin 实例本身来访问其状态。但在很多情况下，代码的不同部分（例如，一个 UI 组件）需要访问某个插件的状态，但它并没有对该插件实例的直接引用。
// 通过 key，任何代码都可以通过 myPluginKey.getState(editorState) 来安全、可靠地获取该插件的状态，实现了模块间的解耦。

const keys = Object.create(null)

function createKey(name: string) {
  if (name in keys) return name + '$' + ++keys[name]
  keys[name] = 0
  return name + '$'
}

class PluginKey<PluginState = any> {
  /// @internal
  key: string

  constructor(name = 'key') {
    this.key = createKey(name)
  }

  /// Get the active plugin with this key, if any, from an editor
  /// state.
  get(state: EditorState): Plugin<PluginState> | undefined {
    return state.config.pluginsByKey[this.key]
  }

  /// Get the plugin's state from an editor state.
  getState(state: EditorState): PluginState | undefined {
    return (state as any)[this.key]
  }
}

export {}
