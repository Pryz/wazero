package wasmdebug

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ianlancetaylor/demangle"
)

type PerfMap struct {
	entries []entry
	fh      *os.File
}

type entry struct {
	addr *uint64
	size uint64
	name string
}

func NewPerfMap() *PerfMap {
	return &PerfMap{}
}

func (f *PerfMap) AddEntry(addr *uint64, size uint64, name string) {
	e := entry{addr, size, name}
	if f.entries == nil {
		f.entries = []entry{e}
		return
	}
	f.entries = append(f.entries, e)
}

func (f *PerfMap) Flush() error {
	defer func() {
		_ = f.fh.Sync()
		_ = f.fh.Close()
	}()

	var err error
	f.fh, err = f.file()
	if err != nil {
		return err
	}

	for _, e := range f.entries {
		dem, err := demangle.ToString(e.name)
		if err != nil {
			dem = e.name
		}
		f.fh.WriteString(fmt.Sprintf("%x %s %s\n",
			e.addr,
			strconv.FormatUint(e.size, 16),
			dem,
		))
	}
	return nil
}

func (f *PerfMap) file() (*os.File, error) {
	pid := os.Getpid()
	filename := "/tmp/perf-" + strconv.Itoa(pid) + ".map"

	var fh *os.File
	var err error
	if _, err := os.Stat(filename); err == nil {
		fh, err = os.Open(filename)
	} else {
		fh, err = os.Create(filename)
	}
	if err != nil {
		return nil, err
	}
	return fh, nil
}
