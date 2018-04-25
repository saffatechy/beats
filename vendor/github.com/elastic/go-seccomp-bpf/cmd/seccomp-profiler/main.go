// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"bufio"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v2"

	"github.com/elastic/go-seccomp-bpf"
	"github.com/elastic/go-seccomp-bpf/arch"
	"github.com/elastic/go-seccomp-bpf/cmd/seccomp-profiler/disasm"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(value string) error {
	list := strings.FieldsFunc(value, func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == ';'
	})
	*s = append(*s, list...)
	return nil
}

// Flags.
var (
	debug        bool
	format       string
	templateFile string
	packageName  string
	blacklist    stringSlice
)

func init() {
	flag.StringVar(&format, "format", "code", "output format (code or config)")
	flag.StringVar(&templateFile, "t", "", "custom code template file")
	flag.StringVar(&packageName, "pkg", "main", "package name to use in source code")
	flag.BoolVar(&debug, "d", false, "add debug to the config output")
	flag.Var(&blacklist, "b", "blacklist syscalls by name")
}

func main() {
	flag.Parse()

	binary := flag.Arg(0)
	if binary == "" {
		log.Fatal("no binary specified")
	}
	log.Println("Binary file:", binary)

	archInfo, goarch, err := getBinaryArch(binary)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Detected architecture:", archInfo.Name)

	hash, err := hashBinary(binary)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SHA256:", hash)

	objDump, err := doObjdump(binary, hash)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Objdump File:", objDump)

	syscalls, err := disasm.ExtractSyscalls(archInfo, objDump)
	if err != nil {
		log.Fatal(err)
	}

	// Deduplicate syscalls.
	m := make(map[int]disasm.Syscall, len(syscalls))
	for _, s := range syscalls {
		m[s.Num] = s
	}
	var names []string
	for _, s := range m {
		names = append(names, s.Name)
	}
	sort.Strings(names)
	names = filterBlacklist(names)

	log.Printf("Found %d total syscalls", len(syscalls))
	log.Printf("Found %d unique syscalls", len(m))
	if len(blacklist) > 0 {
		log.Printf("Filtered %d blacklisted syscalls", len(m)-len(names))
	}

	// Write output.
	switch format {
	case "code":
		if err = writeGoTemplate(goarch, names); err != nil {
			log.Fatal(err)
		}
	case "config":
		if debug {
			if err = writeDebugYAML(syscalls); err != nil {
				log.Fatal(err)
			}
		}

		if err = writeProfileConfig(names); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid format=%v", format)
	}
}

func getBinaryArch(binary string) (*arch.Info, string, error) {
	f, err := os.Open(binary)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	bin, err := elf.NewFile(f)
	if err != nil {
		return nil, "", err
	}

	if section := bin.Section(".note.go.buildid"); section == nil {
		return nil, "", fmt.Errorf("%v is not a Go binary", binary)
	}

	libs, err := bin.DynString(elf.DT_NEEDED)
	if err != nil {
		return nil, "", err
	}
	if len(libs) > 0 {
		log.Println("Binary is dynamically linked with", strings.Join(libs, ", "))
	}

	switch bin.Machine {
	case elf.EM_386:
		return arch.I386, "386", nil
	case elf.EM_ARM:
		return arch.ARM, "arm", nil
	case elf.EM_X86_64:
		return arch.X86_64, "amd64", nil
	default:
		return nil, "", fmt.Errorf("%v architecture is not supported by go-seccomp-bpf", bin.Machine)
	}
}

func hashBinary(binary string) (string, error) {
	f, err := os.Open(binary)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, bufio.NewReader(f)); err != nil {
		return "", nil
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func cachedDumpFile(binary string) (string, error) {
	abs, err := filepath.Abs(binary)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write([]byte(abs)); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	hash = hash[:10]

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dumpDir := filepath.Join(usr.HomeDir, ".seccomp-profiler")
	if err := os.MkdirAll(dumpDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(dumpDir, filepath.Base(binary)+"-"+hash), nil
}

func doObjdump(binary, hash string) (string, error) {
	dumpFile, err := cachedDumpFile(binary)
	if err != nil {
		return "", err
	}

	f, err := os.Open(dumpFile)
	if err == nil {
		buf := make([]byte, 256/4)
		n, err := f.Read(buf)
		f.Close()
		if err == nil && n == len(buf) && hash == string(buf) {
			log.Println("Using cached objdump.")
			return dumpFile, nil
		}
	}

	f, err = os.Create(dumpFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	out := bufio.NewWriter(f)
	defer out.Flush()

	if _, err = out.WriteString(hash + "\n"); err != nil {
		return "", err
	}

	cmd := exec.Command("go", "tool", "objdump", binary)
	cmd.Stdout = out
	if err = cmd.Run(); err != nil {
		return "", err
	}

	log.Println("objdump written to", dumpFile)
	return dumpFile, nil
}

func filterBlacklist(syscalls []string) []string {
	filter := make(map[string]struct{}, len(blacklist))
	for _, s := range blacklist {
		filter[s] = struct{}{}
	}

	var out []string
	for _, s := range syscalls {
		if _, found := filter[s]; !found {
			out = append(out, s)
		}
	}
	return out
}

func writeDebugYAML(syscalls []disasm.Syscall) error {
	sort.Slice(syscalls, func(i, j int) bool {
		return syscalls[i].Name < syscalls[j].Name
	})

	var debug = struct {
		AllSyscalls []disasm.Syscall `yaml:"all_syscalls"`
	}{
		AllSyscalls: syscalls,
	}

	data, err := yaml.Marshal(debug)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func writeProfileConfig(syscalls []string) error {
	type Config struct {
		Seccomp seccomp.Policy `yaml:"seccomp"`
	}

	config := Config{
		Seccomp: seccomp.Policy{
			DefaultAction: seccomp.ActionErrno,
			Syscalls: []seccomp.SyscallGroup{
				{
					Action: seccomp.ActionAllow,
					Names:  syscalls,
				},
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

const defaultTemplate = `// Code generated by seccomp-profiler - DO NOT EDIT.

// {{ printf "+build linux,"}}{{.GOARCH}}

package {{.Package}}

import (
    "github.com/elastic/go-seccomp-bpf"
)

var SeccompProfile = seccomp.Policy{
    DefaultAction: seccomp.ActionErrno,
    Syscalls: []seccomp.SyscallGroup{
        {
            Action: seccomp.ActionAllow,
            Names:  []string{
{{- range $syscall := .SyscallNames}}
                "{{ $syscall}}",
{{- end}}
            },
        },
    },
}
`

var codeTemplate = template.Must(template.New("profile").Parse(defaultTemplate))

func writeGoTemplate(goarch string, syscalls []string) error {
	t := codeTemplate
	if templateFile != "" {
		var err error
		t, err = template.ParseFiles(templateFile)
		if err != nil {
			return err
		}
	}

	type Params struct {
		Package      string
		GOARCH       string
		SyscallNames []string
	}

	p := Params{
		Package:      packageName,
		GOARCH:       goarch,
		SyscallNames: syscalls,
	}
	return t.Execute(os.Stdout, p)
}
