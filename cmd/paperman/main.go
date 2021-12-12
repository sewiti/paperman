package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/sewiti/paperman/internal/server"
	"github.com/sewiti/paperman/pkg/screen"
)

func main() {
	const srvDir = "/opt/paperman/servers"

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
		if err != nil {
			break
		}
		fmt.Printf("Created %s\n", filepath.Join(srvDir, name))

	case "send":
		requireNArgs(os.Args, 3)
		name := os.Args[2]
		screen := server.Screen(name)
		stuff := strings.Join(os.Args[3:], " ") + "\x0f"
		err = screen.SendStuffContext(ctx, stuff)
		if err != nil {
			break
		}
		fmt.Println("Command sent")

	case "launch":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		var srv server.Server
		srv, err = server.Read(filepath.Join(srvDir, name))
		if err != nil {
			break
		}
		err = srv.Launch(ctx, filepath.Join(srvDir, name))

	case "backup":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		err = backup(ctx, srvDir, name)

	case "backups-purge":
		requireNArgs(os.Args, 3)
		name := os.Args[2]
		var count int
		count, err = strconv.Atoi(os.Args[3])
		if err != nil {
			break
		}
		err = purgeBackups(srvDir, name, count)

	case "restore":
		requireNArgs(os.Args, 3)
		name := os.Args[2]
		backupFile := os.Args[3]

		confirmation := fmt.Sprintf("Yes, restore %s", name)
		fmt.Printf("Are you sure you want to restore %q? This will delete current data. [%s]: ", name, confirmation)
		var cont bool
		cont, err = promptConfirm(os.Stdin, confirmation)
		if err != nil {
			break
		}
		if !cont {
			fmt.Println("Restore aborted")
			break
		}
		err = restore(ctx, srvDir, name, backupFile)

	case "delete":
		requireNArgs(os.Args, 2)
		name := os.Args[2]
		dir := filepath.Join(srvDir, name)

		confirmation := fmt.Sprintf("Yes, delete %s", name)
		fmt.Printf("Are you sure you want to delete %q? [%s]: ", dir, confirmation)
		var cont bool
		cont, err = promptConfirm(os.Stdin, confirmation)
		if err != nil {
			break
		}
		if !cont {
			fmt.Println("Deletion aborted")
			break
		}
		var enabled bool
		enabled, err = server.Server{Name: name}.IsEnabledStandalone(ctx)
		if err != nil {
			break
		}
		if enabled {
			systemControl("disable", name)
		}
		err = os.RemoveAll(dir)

	case "list":
		var running []screen.Screen
		running, err = screen.ListContext(ctx, "paperman-")
		if err != nil {
			break
		}
		err = listInstances(os.Stdout, running, srvDir)

	case "install":
		err = install()

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
		fmt.Fprintf(os.Stderr, "%s: ", os.Args[1])
		if eerr, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "%s: %s", eerr.Error(), string(eerr.Stderr))
			os.Exit(1)
		}
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
		"Instances commands:\n" +
		"  create NAME VERSION PORT\t: Create new instance\n" +
		"  send NAME COMMAND       \t: Send command to a running instance\n" +
		"  launch NAME             \t: Launch instance\n" +
		"  list                    \t: List instances\n" +
		"  delete NAME             \t: Delete instance\n" +
		"  backup NAME             \t: Backup instance (live supported)\n" +
		"  backups-purge NAME COUNT\t: Purge old backups, leaving only COUNT backups\n" +
		"  restore NAME ARCHIVE    \t: Restore instance\n" +
		"\n" +
		"Service commands:\n" +
		"  install     \t: Install service (requires root)\n" +
		"  start NAME  \t: Start instance service\n" +
		"  stop NAME   \t: Stop instance service\n" +
		"  restart NAME\t: Restart instance service\n" +
		"  enable NAME \t: Enable instance service\n" +
		"  disable NAME\t: Disable instance service\n"

	fmt.Fprintf(w, help, os.Args[0])
}

func promptConfirm(r io.Reader, confirmation string) (bool, error) {
	sc := bufio.NewScanner(r)
	sc.Scan()
	if sc.Err() != nil {
		return false, sc.Err()
	}
	return sc.Text() == confirmation, nil
}
