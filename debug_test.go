package wazero_test

import (
	"context"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func TestDebug(t *testing.T) {
	ctx := context.Background()

	config := wazero.NewRuntimeConfigCompiler().WithDebugInfoEnabled(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, config)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	binpath := "/Users/julienfabre/code/github.com/stealthrocket/testdata/rust/fib/target/wasm32-wasi/debug/fib.wasm"

	fh, err := os.Open(binpath)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		t.Fatal(err)
	}

	_, err = runtime.CompileModule(ctx, data)
	if err != nil {
		t.Fatal(err)
	}

	//d := compiled.DWARFLines().Data()
	//r := d.Reader()

	//for {
	//	println("--- entry ---")
	//	entry, err := r.Next()
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	//	println(entry.Offset, entry.Tag.String(), entry.Children)

	//	for _, f := range entry.Field {
	//		str := fmt.Sprintf("field: attr=%s class=%s", f.Attr.String(), f.Class.String())
	//		switch f.Class {
	//		case dwarf.ClassString:
	//			str += fmt.Sprintf(" val=%s", f.Val.(string))
	//		case dwarf.ClassConstant:
	//			str += fmt.Sprintf(" val=%d", f.Val.(int64))
	//		case dwarf.ClassAddress:
	//			str += fmt.Sprintf(" val=%d", f.Val.(uint64))
	//		case dwarf.ClassBlock:
	//			str += fmt.Sprintf(" val=%d", f.Val.([]byte))

	//		}
	//		if f.Class == dwarf.ClassString {
	//			println("field: ", f.Attr.String(), f.Class.String(), f.Val.(string))
	//		} else {
	//			println("field: ", f.Attr.String(), f.Class.String(), f.Val)
	//		}
	//		println(str)
	//	}
	//	println("--- end entry ---")
	//	println()
	//}

}

func newJITDumpFile() (*os.File, error) {

	return os.Create("jitdump.log")
}
