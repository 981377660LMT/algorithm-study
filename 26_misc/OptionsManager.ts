// 设计一个通用的 OptionsManager，不与任何业务逻辑耦合，可用于各种需要管理选项的场景。
//
// 1. **完全解耦**：不依赖任何业务逻辑或特定组件
// 2. **类型安全**：使用泛型支持不同类型的选项
// 3. **功能丰富**：提供完整的选项管理API，如添加、删除、选中等
// 4. **高可复用性**：可用于任何需要管理选项的场景，如下拉框、多选框、复选框等
// 5. **统一选项管理**：一致地处理系统选项和自定义选项

export interface Option {
  label: string
  value: string
}

/**
 * 通用选项管理器.
 * 负责管理选项的模型层，提供选项处理的通用能力.
 */
export class OptionsManager<T extends Option = Option> {
  private _presetOptions: T[]
  private _customOptions: T[] = []
  private _selectedValues: Set<string>
  private readonly _createOptionFn: (label: string, value: string) => T

  /**
   * @param presetOptions 预设的选项.
   * @param selectedValues 已选中的值.
   * @param createOptionFn 创建自定义选项的工厂函数.
   */
  constructor(presetOptions: T[], selectedValues: string[], createOptionFn: (label: string, value: string) => T) {
    this._presetOptions = [...presetOptions]
    this._selectedValues = new Set(selectedValues)
    this._createOptionFn = createOptionFn

    this._updateMissingCustomOptions()
  }

  getAllOptions(): T[] {
    return [...this._customOptions, ...this._presetOptions]
  }

  getPresetOptions(): T[] {
    return [...this._presetOptions]
  }

  getCustomOptions(): T[] {
    return [...this._customOptions]
  }

  getSelectedValues(): string[] {
    return [...this._selectedValues]
  }

  isSelected(value: string): boolean {
    return this._selectedValues.has(value)
  }

  setSelectedValues(values: string[]): void {
    this._selectedValues = new Set(values)
    this._updateMissingCustomOptions()
  }

  addSelectedValue(value: string): void {
    if (!this._selectedValues.has(value)) {
      this._selectedValues.add(value)
      this._updateMissingCustomOptions()
    }
  }

  removeSelectedValue(value: string): void {
    if (this._selectedValues.has(value)) {
      this._selectedValues.delete(value)
      this._updateMissingCustomOptions()
    }
  }

  toggleValue(value: string): void {
    if (this.isSelected(value)) {
      this.removeSelectedValue(value)
    } else {
      this.addSelectedValue(value)
    }
  }

  setPresetOptions(options: T[]): void {
    this._presetOptions = [...options]
    this._updateMissingCustomOptions()
  }

  addCustomOption(label: string, value: string): void {
    if (!this.hasOption(value)) {
      const newOption = this._createOptionFn(label, value)
      this._customOptions.push(newOption)
    }
  }

  hasOption(value: string): boolean {
    return this._presetOptions.some(opt => opt.value === value) || this._customOptions.some(opt => opt.value === value)
  }

  getOption(value: string): T | undefined {
    return this._presetOptions.find(opt => opt.value === value) || this._customOptions.find(opt => opt.value === value)
  }

  getLabel(value: string): string {
    const option = this.getOption(value)
    return option ? option.label : value
  }

  getSelectedLabels(): string[] {
    return [...this._selectedValues].map(value => this.getLabel(value))
  }

  clearSelectedValues(): void {
    this._selectedValues = new Set()
  }

  clearCustomOptions(): void {
    this._customOptions = []
    this._updateMissingCustomOptions()
  }

  reset(): void {
    this._presetOptions = []
    this._customOptions = []
    this._selectedValues = new Set()
  }

  /**
   * 同步选中的选项.
   * 确保所有选中的值都能在选项列表中找到.
   */
  private _updateMissingCustomOptions(): void {
    // 获取所有已有选项的值集合
    const existingValues = new Set([...this._presetOptions.map(opt => opt.value), ...this._customOptions.map(opt => opt.value)])

    // 找出在选中值中但不在选项中的值
    const missingValues = [...this._selectedValues].filter(value => !existingValues.has(value))

    // 为缺失的值创建自定义选项
    missingValues.forEach(value => {
      this.addCustomOption(value, value)
    })
  }
}

if (require.main === module) {
  // 测试代码
  const manager = new OptionsManager(
    [
      { label: '系统选项1', value: 'sys1' },
      { label: '系统选项2', value: 'sys2' }
    ],
    ['sys1'],
    (label, value) => ({ label, value })
  )

  console.log('所有选项:', manager.getAllOptions())
  console.log('选中的值:', manager.getSelectedValues())
  console.log('选中的标签:', manager.getSelectedLabels())

  manager.addCustomOption('自定义选项1', 'custom1')
  console.log('添加自定义选项后:', manager.getAllOptions())
  console.log('选中的值:', manager.getSelectedValues())
  console.log('选中的标签:', manager.getSelectedLabels())
}
