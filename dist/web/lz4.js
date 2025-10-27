window.LZ4_WASM = {
  currentScript: document.currentScript.src
}

window.LZ4_WASM.init = async () => {
  const go = new Go()

  const baseUrl = new URL(window.LZ4_WASM.currentScript)
  const wasmUrl = new URL("../lz4.wasm", baseUrl)

  const {instance} = await WebAssembly.instantiateStreaming(
    fetch(wasmUrl),
    go.importObject
  )
  go.run(instance)

  window.LZ4_WASM.compressFrame           = globalThis.compressFrame
  window.LZ4_WASM.uncompressFrame         = globalThis.uncompressFrame
  window.LZ4_WASM.compressBlockBound      = globalThis.compressBlockBound
  window.LZ4_WASM.compressBlock           = globalThis.compressBlock
  window.LZ4_WASM.compressBlockHC         = globalThis.compressBlockHC
  window.LZ4_WASM.uncompressBlock         = globalThis.uncompressBlock
  window.LZ4_WASM.uncompressBlockWithDict = globalThis.uncompressBlockWithDict

  delete globalThis.Go
  delete globalThis.compressFrame
  delete globalThis.uncompressFrame
  delete globalThis.compressBlockBound
  delete globalThis.compressBlock
  delete globalThis.compressBlockHC
  delete globalThis.uncompressBlock
  delete globalThis.uncompressBlockWithDict
}
