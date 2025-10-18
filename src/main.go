package main

import (
  lz4 "github.com/pierrec/lz4/v4"
  "lz4-hc-wasm/utils"
  "syscall/js"
)

// -------------------------------------------------------------------

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
  depthNumber := lz4.CompressionLevel(
    utils.ReadJsUint32(args[1]),
  )

  var dstSize int
  if len(args) < 3 || !utils.ValidateJsNumber(args[2]) {
    dstSize = lz4.CompressBlockBound(len(srcBytes))
  } else {
    dstSize = utils.ReadJsInt(args[2])
  }
  dstBytes := make([]byte, dstSize)

  // resolve to Uint8Array, or reject Error
  n, err := lz4.CompressBlockHC(srcBytes, dstBytes, depthNumber, nil, nil)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, n, err)
}

// -------------------------------------------------------------------

func main() {
  js.Global().Set("compressBlockBound",      js.FuncOf(compressBlockBound))
  js.Global().Set("uncompressBlock",         js.FuncOf(uncompressBlock))
  js.Global().Set("uncompressBlockWithDict", js.FuncOf(uncompressBlockWithDict))
  js.Global().Set("compressBlock",           js.FuncOf(compressBlock))
  js.Global().Set("compressBlockHC",         js.FuncOf(compressBlockHC))

  // Keep the Go program running indefinitely (essential for Wasm)
  <-make(chan bool)
}
