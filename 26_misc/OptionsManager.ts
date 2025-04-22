// 设计一个通用的 OptionManager，不与任何业务逻辑耦合，可用于各种需要管理选项的场景。
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
  private readonly _systemOptions: T[] = []
  private readonly _customOptions: T[] = []
  private readonly selectedValues: string[] = []
  private readonly createOptionFn: (label: string, value: string) => T

  /**
   * @param systemOptions 系统提供的选项.
   * @param selectedValues 已选中的值.
   * @param createOptionFn 创建自定义选项的工厂函数.
   */
  constructor(systemOptions: T[], selectedValues: string[] = [], createOptionFn?: (label: string, value: string) => T) {
    this._systemOptions = [...systemOptions]
    this.selectedValues = [...selectedValues]

    // 默认的创建选项函数
    this.createOptionFn =
      createOptionFn ||
      ((label, value) => {
        return {
          label,
          value,
          isCustom: true
        } as T
      })

    // 同步选项
    this.syncCustomOptions()
  }

  /**
   * 获取所有选项（包括系统和自定义选项）
   */
  public getAllOptions(): T[] {
    return [...this._customOptions, ...this._systemOptions]
  }

  /**
   * 获取系统选项
   */
  public getSystemOptions(): T[] {
    return [...this._systemOptions]
  }

  /**
   * 获取自定义选项
   */
  public getCustomOptions(): T[] {
    return [...this._customOptions]
  }

  /**
   * 获取选中的值
   */
  public getSelectedValues(): string[] {
    return [...this.selectedValues]
  }

  /**
   * 判断值是否被选中
   */
  public isSelected(value: string): boolean {
    return this.selectedValues.includes(value)
  }

  /**
   * 设置选中的值
   */
  public setSelectedValues(values: string[]): void {
    this.selectedValues = [...values]
    this.syncCustomOptions()
  }

  /**
   * 添加选中值
   */
  public addSelectedValue(value: string): void {
    if (!this.selectedValues.includes(value)) {
      this.selectedValues.push(value)
      this.syncCustomOptions()
    }
  }

  /**
   * 移除选中值
   */
  public removeSelectedValue(value: string): void {
    this.selectedValues = this.selectedValues.filter(v => v !== value)
  }

  /**
   * 切换选中状态
   */
  public toggleValue(value: string): void {
    if (this.isSelected(value)) {
      this.removeSelectedValue(value)
    } else {
      this.addSelectedValue(value)
    }
  }

  /**
   * 设置系统选项
   */
  public setSystemOptions(options: T[]): void {
    this._systemOptions = [...options]
    this.syncCustomOptions()
  }

  /**
   * 添加自定义选项
   */
  public addCustomOption(label: string, value: string): void {
    // 检查是否已存在
    if (!this.hasOption(value)) {
      const newOption = this.createOptionFn(label, value)
      this._customOptions.push(newOption)
    }
  }

  /**
   * 检查选项是否存在
   */
  public hasOption(value: string): boolean {
    return this._systemOptions.some(opt => opt.value === value) || this._customOptions.some(opt => opt.value === value)
  }

  /**
   * 获取选项
   */
  public getOption(value: string): T | undefined {
    return this._systemOptions.find(opt => opt.value === value) || this._customOptions.find(opt => opt.value === value)
  }

  /**
   * 获取选项的标签
   */
  public getLabel(value: string): string {
    const option = this.getOption(value)
    return option ? option.label : value
  }

  /**
   * 获取所有选中项的标签
   */
  public getSelectedLabels(): string[] {
    return this.selectedValues.map(value => this.getLabel(value))
  }

  /**
   * 清空选中的值
   */
  public clearSelectedValues(): void {
    this.selectedValues = []
  }

  /**
   * 清空自定义选项
   */
  public clearCustomOptions(): void {
    this._customOptions = []
    this.syncCustomOptions()
  }

  /**
   * 重置所有数据
   */
  public reset(): void {
    this._systemOptions = []
    this._customOptions = []
    this.selectedValues = []
  }

  /**
   * 同步选中的选项
   * 确保所有选中的值都能在选项列表中找到
   */
  private syncCustomOptions(): void {
    // 获取所有已有选项的值集合
    const existingValues = new Set([...this._systemOptions.map(opt => opt.value), ...this._customOptions.map(opt => opt.value)])

    // 找出在选中值中但不在选项中的值
    const missingValues = this.selectedValues.filter(value => !existingValues.has(value))

    // 为缺失的值创建自定义选项
    missingValues.forEach(value => {
      this.addCustomOption(value, value)
    })
  }
}
