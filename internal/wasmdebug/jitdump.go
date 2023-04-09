package wasmdebug

// https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/tools/perf/Documentation/jitdump-specification.txt

type FileHeader struct {
	// * uint32_t magic     : a magic number tagging the file type. The value is 4-byte long and represents the string "JiTD" in ASCII form. It written is as 0x4A695444. The reader will detect an endian mismatch when it reads 0x4454694a.
	Magic uint32
	// * uint32_t version   : a 4-byte value representing the format version. It is currently set to 1
	Version uint32
	// * uint32_t total_size: size in bytes of file header
	TotalSize uint32
	// * uint32_t elf_mach  : ELF architecture encoding (ELF e_machine value as specified in /usr/include/elf.h)
	ElfMach uint32
	// * uint32_t pad1      : padding. Reserved for future use
	Pad1 uint32
	// * uint32_t pid       : JIT runtime process identification (OS specific)
	Pid uint32
	// * uint64_t timestamp : timestamp of when the file was created
	Timestamp uint64
	// * uint64_t flags     : a bitmask of flags
	Flags uint64
}

type JitDumpFile struct {
	Header FileHeader
}

type RecordHeader struct {
	// * uint32_t id        : a value identifying the record type (see below)
	Id uint32
	// * uint32_t total_size: the size in bytes of the record including the header.
	TotalSize uint32
	// * uint64_t timestamp : a timestamp of when the record was created.
	Timestamp uint64
}

type RecordType int

const (
	JITCodeLoad RecordType = iota
	JITCodeMove
	JITCodeDebugInfo
	JITCodeClose
	JITCodeUnwindingInfo
)

type CodeLoadRecord struct {
	RecordHeader

	// * uint32_t pid: OS process id of the runtime generating the jitted code
	Pid uint32
	// * uint32_t tid: OS thread identification of the runtime thread generating the jitted code
	Tid uint32
	// * uint64_t vma: virtual address of jitted code start
	Vma uint64
	// * uint64_t code_addr: code start address for the jitted code. By default vma = code_addr
	CodeAddr uint64
	// * uint64_t code_size: size in bytes of the generated jitted code
	CodeSize uint64
	// * uint64_t code_index: unique identifier for the jitted code (see below)
	CodeIndex uint64
	// * char[n]: function name in ASCII including the null termination
	FunctionName string
	// * native code: raw byte encoding of the jitted code
	Code []byte
}

// The record contains source lines debug information, i.e., a way to map a code
// address back to a source line. This information may be used by the performance tool.
type DebugInfoRecord struct {
	RecordHeader

	// * uint64_t code_addr: address of function for which the debug information is generated
	CodeAddr uint64
	// * uint64_t nr_entry : number of debug entries for the function
	NrEntry uint64
	// * debug_entry[n]: array of nr_entry debug entries for the function
	DebugEntry []DebugEntry
}

type DebugEntry struct {
	// * uint64_t code_addr: address of function for which the debug information is generated
	CodeAddr uint64
	// * uint32_t line     : source file line number (starting at 1)
	Line uint32
	// * uint32_t discrim  : column discriminator, 0 is default
	Discrim uint32
	// * char name[n]      : source file name in ASCII, including null termination
	Name string
}
