package main

import (
  "bytes"
  "io"
  "syscall/js"

  lz4 "github.com/pierrec/lz4/v4"
  "lz4-hc-wasm/utils"
)

// -------------------------------------------------------------------
// frame-level data format

func compressFrame(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "compressFrame() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  var depthNumber int
  if len(args) < 2 || !utils.ValidateJsNumber(args[1]) {
    depthNumber = 0
  } else {
    depthNumber = utils.ReadJsInt(args[1])
  }

  // resolve to Uint8Array, or reject Error
  dstBytes, err := compressFrameWithReader(srcBytes, depthNumber)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
}

func uncompressFrame(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "uncompressFrame() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  // resolve to Uint8Array, or reject Error
  dstBytes, err := uncompressFrameWithReader(srcBytes)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
}

// -------------------------------------------------------------------
// frame-level helpers

func compressFrameWithReader(src []byte, depth int) ([]byte, error) {
  level := getCompressionLevel(depth)
  r := io.NopCloser(bytes.NewReader(src))
  zcomp := lz4.NewCompressingReader(r)
  if err := zcomp.Apply(lz4.CompressionLevelOption(level)); err != nil {
    return nil, err
  }
  return io.ReadAll(zcomp)
}

func uncompressFrameWithReader(src []byte) ([]byte, error) {
  zr := lz4.NewReader(bytes.NewReader(src))
  return io.ReadAll(zr)
}

// -------------------------------------------------------------------
// block-level data format

func compressBlockBound(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsNumber(args[0]) {
    js.Global().Get("console").Call("error", "compressBlockBound() called incorrectly. parameter #1: requires Integer")
    return -1
  }

  return lz4.CompressBlockBound(
    utils.ReadJsInt(args[0]),
  )
}

func compressBlock(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "compressBlock() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  var dstSize int
  if len(args) < 2 || !utils.ValidateJsNumber(args[1]) {
    dstSize = lz4.CompressBlockBound(len(srcBytes))
  } else {
    dstSize = utils.ReadJsInt(args[1])
  }
  dstBytes := make([]byte, dstSize)

  // resolve to Uint8Array, or reject Error
  n, err := lz4.CompressBlock(srcBytes, dstBytes, nil)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, n, err)
}

func compressBlockHC(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "compressBlockHC() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  // Validate input
  if len(args) < 2 || !utils.ValidateJsNumber(args[1]) {
    js.Global().Get("console").Call("error", "compressBlockHC() called incorrectly. parameter #2 (depth): requires Integer")
    return nil
  }
  depthNumber := utils.ReadJsInt(args[1])

  var dstSize int
  if len(args) < 3 || !utils.ValidateJsNumber(args[2]) {
    dstSize = lz4.CompressBlockBound(len(srcBytes))
  } else {
    dstSize = utils.ReadJsInt(args[2])
  }
  dstBytes := make([]byte, dstSize)

  // resolve to Uint8Array, or reject Error
  level := getCompressionLevel(depthNumber)
  n, err := lz4.CompressBlockHC(srcBytes, dstBytes, level, nil, nil)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, n, err)
}

func uncompressBlock(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "uncompressBlock() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  var dstSize int
  if len(args) < 2 || !utils.ValidateJsNumber(args[1]) {
    dstSize = len(srcBytes) * 255
  } else {
    dstSize = utils.ReadJsInt(args[1])
  }
  dstBytes := make([]byte, dstSize)

  // resolve to Uint8Array, or reject Error
  n, err := lz4.UncompressBlock(srcBytes, dstBytes)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, n, err)
}

func uncompressBlockWithDict(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "uncompressBlockWithDict() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  // Validate input
  if len(args) < 2 || !utils.ValidateJsUint8Array(args[1]) {
    js.Global().Get("console").Call("error", "uncompressBlockWithDict() called incorrectly. parameter #2 (dict): requires Uint8Array")
    return nil
  }
  dictBytes := utils.ReadJsBytes(args[1])

  var dstSize int
  if len(args) < 3 || !utils.ValidateJsNumber(args[2]) {
    dstSize = len(srcBytes) * 255
  } else {
    dstSize = utils.ReadJsInt(args[2])
  }
  dstBytes := make([]byte, dstSize)

  // resolve to Uint8Array, or reject Error
  n, err := lz4.UncompressBlockWithDict(srcBytes, dstBytes, dictBytes)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, n, err)
}

// -------------------------------------------------------------------
// common helpers

func getCompressionLevel(depth int) (lz4.CompressionLevel) {
  if depth < 0 {
    depth = 0
  }
  if depth > 9 {
    depth = 9
  }

  switch depth {
  case 0:
    return lz4.Fast
  case 1:
    return lz4.Level1
  case 2:
    return lz4.Level2
  case 3:
    return lz4.Level3
  case 4:
    return lz4.Level4
  case 5:
    return lz4.Level5
  case 6:
    return lz4.Level6
  case 7:
    return lz4.Level7
  case 8:
    return lz4.Level8
  case 9:
    return lz4.Level9
  }
  return lz4.Fast
}

// -------------------------------------------------------------------
// JS API

func main() {
  // frame-level data format
  js.Global().Set("compressFrame",           js.FuncOf(compressFrame))
  js.Global().Set("uncompressFrame",         js.FuncOf(uncompressFrame))

  // block-level data format
  js.Global().Set("compressBlockBound",      js.FuncOf(compressBlockBound))
  js.Global().Set("compressBlock",           js.FuncOf(compressBlock))
  js.Global().Set("compressBlockHC",         js.FuncOf(compressBlockHC))
  js.Global().Set("uncompressBlock",         js.FuncOf(uncompressBlock))
  js.Global().Set("uncompressBlockWithDict", js.FuncOf(uncompressBlockWithDict))

  // Keep the Go program running indefinitely (essential for Wasm)
  <-make(chan bool)
}
