package main

import (
  "bytes"
  "io"
  "reflect"
  "syscall/js"
  "unsafe"

  lz4 "github.com/pierrec/lz4/v4"
  "lz4-hc-wasm/utils"
)

// -------------------------------------------------------------------

func uncompressBlock(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "uncompressBlock() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  // resolve to Uint8Array, or reject Error
  dstBytes, err := uncompressBlockWithReader(srcBytes, nil)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
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

  // resolve to Uint8Array, or reject Error
  dstBytes, err := uncompressBlockWithReader(srcBytes, dictBytes)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
}

func compressBlock(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "compressBlock() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  // resolve to Uint8Array, or reject Error
  dstBytes, err := compressBlockWithReader(srcBytes, 0)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
}

func compressBlockHC(this js.Value, args []js.Value) interface{} {
  // Validate input
  if len(args) < 1 || !utils.ValidateJsUint8Array(args[0]) {
    js.Global().Get("console").Call("error", "compressBlockHC() called incorrectly. parameter #1 (src): requires Uint8Array")
    return nil
  }
  srcBytes := utils.ReadJsBytes(args[0])

  var depthNumber int
  if len(args) < 2 || !utils.ValidateJsNumber(args[1]) {
    depthNumber = 9
  } else {
    depthNumber = utils.ReadJsInt(args[1])
  }

  // resolve to Uint8Array, or reject Error
  dstBytes, err := compressBlockWithReader(srcBytes, depthNumber)
  return utils.ReturnJsUint8ArrayAsPromise(dstBytes, -1, err)
}

// -------------------------------------------------------------------

func uncompressBlockWithReader(src, dict []byte) ([]byte, error) {
  zr := lz4.NewReader(bytes.NewReader(src))

  if dict != nil {
    rv := reflect.ValueOf(zr).Elem()
    dictField := rv.FieldByName("dict")

    if dictField.IsValid() {
      ptrToDict := (*[]byte)(unsafe.Pointer(dictField.UnsafeAddr()))
      *ptrToDict = dict
    }
  }

  return io.ReadAll(zr)
}

func compressBlockWithReader(src []byte, depth int) ([]byte, error) {
  var level lz4.CompressionLevel
  switch depth {
  case 0:
    level = lz4.Fast
  case 1:
    level = lz4.Level1
  case 2:
    level = lz4.Level2
  case 3:
    level = lz4.Level3
  case 4:
    level = lz4.Level4
  case 5:
    level = lz4.Level5
  case 6:
    level = lz4.Level6
  case 7:
    level = lz4.Level7
  case 8:
    level = lz4.Level8
  case 9:
    level = lz4.Level9
  default:
    level = lz4.Level9
  }

  r := io.NopCloser(bytes.NewReader(src))
  zcomp := lz4.NewCompressingReader(r)
  if err := zcomp.Apply(lz4.CompressionLevelOption(level)); err != nil {
    return nil, err
  }
  return io.ReadAll(zcomp)
}

// -------------------------------------------------------------------

func main() {
  js.Global().Set("uncompressBlock",         js.FuncOf(uncompressBlock))
  js.Global().Set("uncompressBlockWithDict", js.FuncOf(uncompressBlockWithDict))
  js.Global().Set("compressBlock",           js.FuncOf(compressBlock))
  js.Global().Set("compressBlockHC",         js.FuncOf(compressBlockHC))

  // Keep the Go program running indefinitely (essential for Wasm)
  <-make(chan bool)
}
