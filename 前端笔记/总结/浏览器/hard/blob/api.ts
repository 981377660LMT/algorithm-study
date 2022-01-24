interface Blob {
  readonly size: number
  readonly type: string
  arrayBuffer(): Promise<ArrayBuffer>
  slice(start?: number, end?: number, contentType?: string): Blob
  stream(): NodeJS.ReadableStream
  text(): Promise<string>
}
