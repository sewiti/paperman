package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/sewiti/paperman/internal/instances"
	"github.com/sewiti/paperman/internal/server"
	"github.com/sewiti/paperman/pkg/tmux"
)

const srvDir = "servers"

func main() {
	requireNArgs(os.Args, 1)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var err error
	switch cmd := os.Args[1]; cmd {
	case "create":
		requireNArgs(os.Args, 4)
		name := os.Args[2]
		version := os.Args[3]
		port := os.Args[4]
		err = server.Create(srvDir, name, version, port)

	case "launch":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		srv, err := server.Read(filepath.Join(srvDir, name))
		if err != nil {
			break
		}
		err = srv.Launch(ctx, filepath.Join(srvDir, name))

	case "backup":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		err = backup(srvDir, name)

	case "restore":
		fmt.Println("not implemented yet")
		os.Exit(1)
		// requireNArgs(os.Args, 3)
		// name := os.Args[2]
		// backupFile := os.Args[3]
		// err = restoer(srvDir, name, backupFile)

	case "delete":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		err = os.RemoveAll(filepath.Join(srvDir, name))

	case "list":
		var ss []tmux.Session
		ss, err = instances.ListAll(ctx)
		if err != nil {
			break
		}
		err = listServers(os.Stdout, ss, srvDir)

	case "start", "stop", "restart", "enable", "disable":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		err = systemControl(cmd, name)

	case "h", "help", "-h", "-help", "--help":
		printHelp(os.Stdout)
		os.Exit(0)

	default:
		printHelp(os.Stderr)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func requireNArgs(args []string, n int) {
	if len(args) <= n {
		printHelp(os.Stderr)
		os.Exit(1)
	}
}

func printHelp(w io.Writer) {
	const help = "usage: %s COMMAND [ARGS]\n" +
		"\n" +
		"Server commands:\n" +
		"  create NAME VERSION PORT\t: Create new server\n" +
		"  launch NAME             \t: Launch server\n" +
		"  backup NAME             \t: Backup server (live supported)\n" +
		"  restore NAME            \t: Backup server (live supported)\n" +
		"  delete NAME             \t: Delete server\n" +
		"  list                    \t: List servers\n" +
		"\n" +
		"Service commands:\n" +
		"  start NAME  \t: Start server's service\n" +
		"  stop NAME   \t: Stop server's service\n" +
		"  restart NAME\t: Restart server's service\n" +
		"  enable NAME \t: Enable server's service\n" +
		"  disable NAME\t: Disable server's service\n"

	fmt.Fprintf(w, help, os.Args[0])
}
