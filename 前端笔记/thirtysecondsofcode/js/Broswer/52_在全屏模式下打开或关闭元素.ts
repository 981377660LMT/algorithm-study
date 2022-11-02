const fullscreen = (mode = true, el = 'body') =>
  mode ? document.querySelector(el)!.requestFullscreen() : document.exitFullscreen()

fullscreen() // Opens `body` in fullscreen mode
fullscreen(false) // Exits fullscreen mode
export {}
