package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	flag "github.com/docker/docker/pkg/mflag"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/registry"
)

type DockerCli struct {
	proto      string
	addr       string
	configFile *registry.ConfigFile
	in         io.ReadCloser
	out        io.Writer
	err        io.Writer
	isTerminal bool
	terminalFd uintptr
	tlsConfig  *tls.Config
	scheme     string
}

var funcMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
}

func (cli *DockerCli) getMethod(name string) (func(...string) error, bool) {
	if len(name) == 0 {
		return nil, false
	}
	methodName := "Cmd" + strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
	switch (methodName) {
	default:
		return nil, false
	case "CmdAttach":
		return cli.CmdAttach, true
	case "CmdAttach":
		return cli.CmdAttach, true
	case "CmdBuild":
		return cli.CmdBuild, true
	case "CmdCommit":
		return cli.CmdCommit, true
	case "CmdCp":
		return cli.CmdCp, true
	case "CmdDiff":
		return cli.CmdDiff, true
	case "CmdEvents":
		return cli.CmdEvents, true
	case "CmdExport":
		return cli.CmdExport, true
	case "CmdHelp":
		return cli.CmdHelp, true
	case "CmdHistory":
		return cli.CmdHistory, true
	case "CmdImages":
		return cli.CmdImages, true
	case "CmdImport":
		return cli.CmdImport, true
	case "CmdInfo":
		return cli.CmdInfo, true
	case "CmdInspect":
		return cli.CmdInspect, true
	case "CmdKill":
		return cli.CmdKill, true
	case "CmdLoad":
		return cli.CmdLoad, true
	case "CmdLogin":
		return cli.CmdLogin, true
	case "CmdLogout":
		return cli.CmdLogout, true
	case "CmdLogs":
		return cli.CmdLogs, true
	case "CmdPause":
		return cli.CmdPause, true
	case "CmdPort":
		return cli.CmdPort, true
	case "CmdPs":
		return cli.CmdPs, true
	case "CmdPull":
		return cli.CmdPull, true
	case "CmdPush":
		return cli.CmdPush, true
	case "CmdRestart":
		return cli.CmdRestart, true
	case "CmdRm":
		return cli.CmdRm, true
	case "CmdRmi":
		return cli.CmdRmi, true
	case "CmdRun":
		return cli.CmdRun, true
	case "CmdSave":
		return cli.CmdSave, true
	case "CmdSearch":
		return cli.CmdSearch, true
	case "CmdStart":
		return cli.CmdStart, true
	case "CmdStop":
		return cli.CmdStop, true
	case "CmdTag":
		return cli.CmdTag, true
	case "CmdTop":
		return cli.CmdTop, true
	case "CmdUnpause":
		return cli.CmdUnpause, true
	case "CmdVersion":
		return cli.CmdVersion, true
	case "CmdWait":
		return cli.CmdWait, true
	}
}

// Cmd executes the specified command
func (cli *DockerCli) Cmd(args ...string) error {
	if len(args) > 0 {
		method, exists := cli.getMethod(args[0])
		if !exists {
			fmt.Println("Error: Command not found:", args[0])
			return cli.CmdHelp(args[1:]...)
		}
		return method(args[1:]...)
	}
	return cli.CmdHelp(args...)
}

func (cli *DockerCli) Subcmd(name, signature, description string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.Usage = func() {
		options := ""
		if flags.FlagCountUndeprecated() > 0 {
			options = "[OPTIONS] "
		}
		fmt.Fprintf(cli.err, "\nUsage: docker %s %s%s\n\n%s\n\n", name, options, signature, description)
		flags.PrintDefaults()
		os.Exit(2)
	}
	return flags
}

func (cli *DockerCli) LoadConfigFile() (err error) {
	cli.configFile, err = registry.LoadConfig(os.Getenv("HOME"))
	if err != nil {
		fmt.Fprintf(cli.err, "WARNING: %s\n", err)
	}
	return err
}

func NewDockerCli(in io.ReadCloser, out, err io.Writer, proto, addr string, tlsConfig *tls.Config) *DockerCli {
	var (
		isTerminal = false
		terminalFd uintptr
		scheme     = "http"
	)

	if tlsConfig != nil {
		scheme = "https"
	}

	if in != nil {
		if file, ok := out.(*os.File); ok {
			terminalFd = file.Fd()
			isTerminal = term.IsTerminal(terminalFd)
		}
	}

	if err == nil {
		err = out
	}
	return &DockerCli{
		proto:      proto,
		addr:       addr,
		in:         in,
		out:        out,
		err:        err,
		isTerminal: isTerminal,
		terminalFd: terminalFd,
		tlsConfig:  tlsConfig,
		scheme:     scheme,
	}
}
