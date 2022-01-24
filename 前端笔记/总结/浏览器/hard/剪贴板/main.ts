/** Available only in secure contexts. */
interface Clipboard extends EventTarget {
  read(): Promise<ClipboardItems>
  readText(): Promise<string>
  write(data: ClipboardItems): Promise<void>
  writeText(data: string): Promise<void>
}

interface ClipboardItem {
  readonly types: ReadonlyArray<string>
  getType(type: string): Promise<Blob>
}

// write 方法除了支持文本数据之外，还支持将图像数据写入到剪贴板
