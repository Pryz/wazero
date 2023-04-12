package wazero_test

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func TestTracepoint(t *testing.T) {
	ctx := context.Background()

	config := wazero.NewRuntimeConfigCompiler().WithDebugInfoEnabled(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, config)
	defer runtime.Close(ctx)

	builder := runtime.NewHostModuleBuilder("tracepoint")

	addParams := []uint64{42, 24}
	addResult := []uint64{66}
	var lat time.Duration
	var start time.Time

	builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
		start = time.Now()
		t.Logf("tracepoint.call_add: %v", stack)

		if !reflect.DeepEqual(stack, addParams) {
			t.Log("got", stack)
			t.Log("expected", addParams)
			t.Fatal("unexpected call stack")
		}
	}), []api.ValueType{api.ValueTypeI32, api.ValueTypeI32}, []api.ValueType{}).Export("call_add")

	builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
		lat = time.Since(start)
		t.Logf("tracepoint.return_add: %v", stack)

		if !reflect.DeepEqual(stack, addResult) {
			t.Log("got", stack)
			t.Log("expected", addResult)
			t.Fatal("unexpected return stack")
		}

	}), []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).Export("return_add")

	if _, err := builder.Instantiate(ctx); err != nil {
		t.Fatal(err)
	}

	binpath := "./internal/testing/tracepoint/traced_add.wasm"

	compiled, err := runtime.CompileModule(ctx, wasmdata(t, binpath))
	if err != nil {
		t.Fatal(err)
	}

	instance, err := runtime.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := instance.ExportedFunction("add").Call(ctx, addParams...); err != nil {
		t.Fatal(err)
	}

	t.Logf("latency: %v", lat)
}

func wasmdata(t *testing.T, path string) []byte {
	fh, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	wasmdata, err := ioutil.ReadAll(fh)
	if err != nil {
		t.Fatal(err)
	}
	return wasmdata
}
