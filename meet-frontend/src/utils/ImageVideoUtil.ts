const blurPxAmount = 10

const convertImageFileToBlob = async (file: File, resizeWidth: number, applyResize: boolean = false, applyBlur: boolean = false): Promise<Blob> => {
  return await new Promise<Blob>((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)

    reader.onload = async () => {
      const img = new Image()
      img.src = reader.result as string

      img.onload = () => {
        // Calcular la relación de aspecto de la imagen original
        const aspectRatio = img.width / img.height

        // Calcular la altura correspondiente al ancho dado manteniendo la relación de aspecto
        const resizeHeight = resizeWidth / aspectRatio

        // Crear un elemento canvas en memoria
        const canvas = document.createElement('canvas')
        const ctx = canvas.getContext('2d')

        // Establecer las dimensiones del canvas
        canvas.width = applyResize ? resizeWidth : img.width
        canvas.height = applyResize ? resizeHeight : img.height

        // Aplicar un efecto de desenfoque a la imagen
        if (applyBlur && ctx !== null) {
          ctx.filter = `blur(${blurPxAmount}px)`
        }

        // Dibujar la imagen en el canvas con las dimensiones deseadas
        ctx?.drawImage(img, 0, 0, canvas.width, canvas.height)

        canvas.toBlob((blob) => {
          if (blob != null) {
            resolve(blob)
          } else {
            reject(new Error('Failed to resize image'))
          }
        }, 'image/png')
      }

      img.onerror = () => {
        reject(new Error('Failed to load image'))
      }
    }

    reader.onerror = () => {
      reject(new Error('Failed to read file'))
    }
  })
}

const convertVideoFileToBlob = async (file: File, resizeWidth: number, applyResize: boolean = false, applyBlur: boolean = false): Promise<Blob> => {
  return await new Promise<Blob>((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsArrayBuffer(file)

    reader.onerror = (e) => {
      console.error(e)
      reject(new Error('Failed to read file'))
    }

    reader.onload = async (eReader) => {
      const buffer = eReader.target?.result as ArrayBuffer

      let type = ''
      if (file.type.includes('mp4')) {
        type = 'video/mp4'
      } else if (file.type.includes('ogg')) {
        type = 'video/ogg'
      }

      const videoBlob = new Blob([new Uint8Array(buffer)], { type })

      if (!applyResize) {
        resolve(videoBlob)
      } else {
        const videoPreview = document.createElement('video')
        videoPreview.src = URL.createObjectURL(videoBlob)

        videoPreview.addEventListener('loadeddata', function () {
          if (!isNaN(videoPreview.duration)) {
            const rand = Math.round(Math.random() * videoPreview.duration * 1000) + 1
            videoPreview.currentTime = rand / 1000
          }
        }, false)

        videoPreview.addEventListener('seeked', function () {
        // Calcular la relación de aspecto de la imagen original
          const aspectRatio = videoPreview.videoWidth / videoPreview.videoHeight

          // Calcular la altura correspondiente al ancho dado manteniendo la relación de aspecto
          const resizeHeight = resizeWidth / aspectRatio

          // Crear un elemento canvas en memoria
          const canvas = document.createElement('canvas')
          const ctx = canvas.getContext('2d')

          // Establecer las dimensiones del canvas
          canvas.width = applyResize ? resizeWidth : videoPreview.videoWidth
          canvas.height = applyResize ? resizeHeight : videoPreview.videoHeight

          // Aplicar un efecto de desenfoque a la imagen
          if (applyBlur && ctx !== null) {
            ctx.filter = `blur(${blurPxAmount}px)`
          }

          // Dibujar la imagen en el canvas con las dimensiones deseadas
          ctx?.drawImage(videoPreview, 0, 0, canvas.width, canvas.height)

          canvas.toBlob((blob) => {
            if (blob != null) {
              resolve(blob)
            } else {
              reject(new Error('Failed to resize image'))
            }
          }, 'image/png')
        }, false)

        videoPreview.onerror = () => {
          reject(new Error('Failed to load image'))
        }
      }
    }
  })
}

const isVideoSupportedContentType = (contentType: string): boolean => {
  return ['video/mp4', 'video/ogg', 'audio/ogg'].some(ct => ct === contentType)
}

const isImageSupportedContentType = (contentType: string): boolean => {
  return ['image/jpg', 'image/jpeg', 'image/png'].some(ct => ct === contentType)
}

export const ImageVideoUtil = {
  convertImageFileToBlob,
  convertVideoFileToBlob,
  isVideoSupportedContentType,
  isImageSupportedContentType
}
