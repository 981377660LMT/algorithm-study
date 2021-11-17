type CanvasImageSource = HTMLOrSVGImageElement | HTMLVideoElement | HTMLCanvasElement | ImageBitmap

const video = document.querySelector('#video') as HTMLVideoElement
const canvas = document.createElement('canvas') as HTMLCanvasElement
const img = document.createElement('img')
img.crossOrigin = ''
const ctx = canvas.getContext('2d')!

function captureVideo() {
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight
  ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
  img.src = canvas.toDataURL()
  document.body.append(img)
}

export {}
