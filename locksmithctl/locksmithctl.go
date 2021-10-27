// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"text/tabwriter"

	"github.com/flatcar-linux/fleetlock/pkg/client"
	"github.com/flatcar-linux/locksmith/version"
)

const (
	cliName        = "locksmithctl"
	cliDescription = `Manage the cluster wide reboot lock.`
)

var (
	out *tabwriter.Writer

	commands      []*Command
	globalFlagSet = flag.NewFlagSet("locksmithctl", flag.ExitOnError)

	globalFlags = struct {
		Debug        bool
		Endpoints    endpoints
		EtcdKeyFile  string
		EtcdCertFile string
		EtcdCAFile   string
		EtcdUsername string
		EtcdPassword string
		Group        string
		ID           string
		Version      bool
	}{}

	defaultEndpoints = []string{
		"http://127.0.0.1:2379",
		"http://127.0.0.1:4001",
	}
)

type endpoints []string

func (e *endpoints) String() string {
	if len(*e) == 0 {
		return strings.Join(defaultEndpoints, ",")
	}

	return strings.Join(*e, ",")
}

func (e *endpoints) Set(value string) error {
	*e = []string{}
	for _, url := range strings.Split(value, ",") {
		*e = append(*e, strings.TrimSpace(url))
	}

	return nil
}

func init() {
	out = new(tabwriter.Writer)
	out.Init(os.Stdout, 0, 8, 1, '\t', 0)

	globalFlagSet.BoolVar(&globalFlags.Debug, "debug", false, "Print out debug information to stderr.")
	globalFlagSet.Var(&globalFlags.Endpoints, "endpoint", "etcd endpoint for locksmith. Specify multiple times to use multiple endpoints.")
	globalFlagSet.StringVar(&globalFlags.EtcdKeyFile, "etcd-keyfile", "", "etcd key file authentication")
	globalFlagSet.StringVar(&globalFlags.EtcdCertFile, "etcd-certfile", "", "etcd cert file authentication")
	globalFlagSet.StringVar(&globalFlags.EtcdCAFile, "etcd-cafile", "", "etcd CA file authentication")
	globalFlagSet.StringVar(&globalFlags.EtcdUsername, "etcd-username", "", "username for secure etcd communication")
	globalFlagSet.StringVar(&globalFlags.EtcdPassword, "etcd-password", "", "password for secure etcd communication")
	globalFlagSet.StringVar(&globalFlags.Group, "group", "", "locksmith group")
	globalFlagSet.StringVar(&globalFlags.ID, "id", "", "locksmith id")
	globalFlagSet.BoolVar(&globalFlags.Version, "version", false, "Print the version and exit.")

	commands = []*Command{
		cmdHelp,
		cmdLock,
		cmdReboot,
		cmdSendNeedReboot,
		cmdUnlock,
	}
}

// Command is the struct representation of a subcommand for a cli.
type Command struct {
	Name        string                  // Name of the Command and the string to use to invoke it
	Summary     string                  // One-sentence summary of what the Command does
	Usage       string                  // Usage options/arguments
	Description string                  // Detailed description of command
	Flags       flag.FlagSet            // Set of flags associated with this command
	Run         func(args []string) int // Run a command with the given arguments, return exit status
}

func getAllFlags() (flags []*flag.Flag) {
	return getFlags(globalFlagSet)
}

func getFlags(flagSet *flag.FlagSet) (flags []*flag.Flag) {
	flags = make([]*flag.Flag, 0)
	flagSet.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f)
	})
	return
}

func main() {
	globalFlagSet.Parse(os.Args[1:])
	var args = globalFlagSet.Args()

	if len(globalFlags.Endpoints) == 0 {
		globalFlags.Endpoints = defaultEndpoints
	}

	progName := path.Base(os.Args[0])

	if globalFlags.Version {
		fmt.Printf("%s version %s\n", progName, version.Version)
		os.Exit(0)
	}

	if progName == "locksmithd" {
		flagsFromEnv("LOCKSMITHD", globalFlagSet)
		os.Exit(runDaemon())
	}

	// no command specified - trigger help
	if len(args) < 1 {
		args = append(args, "help")
	}

	flagsFromEnv("LOCKSMITHCTL", globalFlagSet)

	var cmd *Command

	// determine which Command should be run
	for _, c := range commands {
		if c.Name == args[0] {
			cmd = c
			if err := c.Flags.Parse(args[1:]); err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}
			break
		}
	}

	if cmd == nil {
		fmt.Printf("%v: unknown subcommand: %q\n", cliName, args[0])
		fmt.Printf("Run '%v help' for usage.\n", cliName)
		os.Exit(2)
	}

	os.Exit(cmd.Run(cmd.Flags.Args()))
}

// getLockClient returns an initialized fleetlockClient,
// using the global flags and http client
func getClient() (*client.Client, error) {
	// copy of github.com/coreos/etcd/client.DefaultTransport so that
	// TLSClientConfig can be overridden.
	transport := http.DefaultClient
	if globalFlags.EtcdCAFile != "" || globalFlags.EtcdCertFile != "" || globalFlags.EtcdKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(globalFlags.EtcdCertFile, globalFlags.EtcdKeyFile)
		if err != nil {
			return nil, err
		}

		ca, err := ioutil.ReadFile(globalFlags.EtcdCAFile)
		if err != nil {
			return nil, err
		}

		capool := x509.NewCertPool()
		capool.AppendCertsFromPEM(ca)

		tlsconf := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      capool,
		}

		tlsconf.BuildNameToCertificate()

		transport.TLSClientConfig = tlsconf
	}

	for _, ep := range globalFlags.Endpoints {
		cfg := client.Client{
			URL:   ep,
		}
		flc, err := client.New(cfg.URL, globalFlags.Group, globalFlags.ID, transport)
		if err != nil {
			return nil, err
		}
		return flc, nil
	}

	return nil, fmt.Errorf("no url available, tried: %s", strings.Join(globalFlags.Endpoints, ","))
}

// flagsFromEnv parses all registered flags in the given flagSet,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, have the given prefix, and any dashes are replaced by
// underscores - for example: some-flag => PREFIX_SOME_FLAG
func flagsFromEnv(prefix string, fs *flag.FlagSet) {
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})
	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := strings.ToUpper(prefix + "_" + strings.Replace(f.Name, "-", "_", -1))
			val := os.Getenv(key)
			if val != "" {
				fs.Set(f.Name, val)
			}
		}
	})
}
