const image = document.querySelector('#previewContainer') as HTMLImageElement
fetch('https://avatars3.githubusercontent.com/u/4220799')
  .then(response => response.blob())
  .then(blob => {
    const objectURL = URL.createObjectURL(blob)
    image.src = objectURL
  })

/** This Fetch API interface represents a resource request. */
interface Request extends Body {}
interface Body {
  readonly body: ReadableStream<Uint8Array> | null
  readonly bodyUsed: boolean
  arrayBuffer(): Promise<ArrayBuffer>
  blob(): Promise<Blob>
  formData(): Promise<FormData>
  json(): Promise<any>
  text(): Promise<string>
}
