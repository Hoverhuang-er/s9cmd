package cmd

import (
	"github.com/urfave/cli/v2"
	"runtime"
	"s9cmd/internal"
)

const (
	// App info
	appName = "s9cmd"
	// Global info
	globalName = "global"
	gshortDesc = "s9cmd global command"
	glongDesc  = "s9cmd global command which can be used to execute the command"
	// ls info
	lsName  = "ls"
	lsShort = "List the objects"
	lsLong  = "List the objects"
	// cp info
	cpName  = "cp"
	cpShort = "Copy the objects"
	cpLong  = "Copy the objects with the source and destination"
	// mv info
	mvName  = "mv"
	mvShort = "Move the objects"
	mvLong  = "Move the objects with the source and destination"
	// sync info
	syncName  = "sync"
	syncShort = "Sync the objects"
	syncLong  = "Sync the objects with the source and destination"
	// du info
	duName  = "du"
	duShort = "Disk usage"
	duLong  = "Disk usage of the objects"
	// get info
	getName  = "get"
	getShort = "Get the objects"
	getLong  = "Get the objects with the source and destination"
	// put info
	putName  = "put"
	putShort = "Put the objects"
	putLong  = "Put the objects with the source and destination"
	// policy info
	policyName  = "policy"
	policyShort = "Policy the objects"
	policyLong  = "Policy the objects with the source and destination"
	// multipart info
	multipartName  = "multipart"
	multipartShort = "Multipart upload"
	multipartLong  = "Multipart upload with the source and destination"
	// mkdir info
	mkdirName  = "mkdir"
	mkdirShort = "Make directory"
	mkdirLong  = "Make directory with the source and destination"
	// rm info
	rmName  = "rm"
	rmShort = "Remove the objects"
	rmLong  = "Remove the objects with the source and destination"
	// chmod info
	chmodName  = "chmod"
	chmodShort = "Change the mode"
	chmodLong  = "Change the mode with the source and destination"
	// chown info
	chownName  = "chown"
	chownShort = "Change the owner"
	chownLong  = "Change the owner with the source and destination"
	// unlink info
	unlinkName  = "unlink"
	unlinkShort = "Unlink the objects"
	unlinkLong  = "Unlink the objects with the source and destination"
)

var (
	version, commit, date = "dev", "dev", "dev"
	rootCmd               = &cli.Command{
		Name:        globalName,
		Usage:       getShort,
		Description: getLong,
	}
	duCmd = &cli.Command{
		Name:        duName,
		Usage:       duShort,
		Description: duLong,
		Action:      Du,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "human",
				Aliases: []string{"-h"},
			},
			&cli.BoolFlag{
				Name:    "summarize",
				Aliases: []string{"-s"},
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"-a"},
			},
		},
	}
	lsCmd = &cli.Command{
		Name:        lsName,
		Usage:       lsShort,
		Description: lsLong,
	}
	cpCmd = &cli.Command{
		Name:        cpName,
		Usage:       cpShort,
		Description: cpLong,
	}
	mvCmd = &cli.Command{
		Name:        mvName,
		Usage:       mvShort,
		Description: mvLong,
	}
	syncCmd = &cli.Command{
		Name:        syncName,
		Usage:       syncShort,
		Description: syncLong,
	}
	getCmd = &cli.Command{
		Name:        getName,
		Usage:       getShort,
		Description: getLong,
	}
	putCmd = &cli.Command{
		Name:        putName,
		Usage:       putShort,
		Description: putLong,
	}
	policyCmd = &cli.Command{
		Name:        policyName,
		Usage:       policyShort,
		Description: policyLong,
	}
	multipartCmd = &cli.Command{
		Name:        multipartName,
		Usage:       multipartShort,
		Description: multipartLong,
	}
	mkdirCmd = &cli.Command{
		Name:        mkdirName,
		Usage:       mkdirShort,
		Description: mkdirLong,
	}
	rmCmd = &cli.Command{
		Name:        rmName,
		Usage:       rmShort,
		Description: rmLong,
	}
	chmodCmd = &cli.Command{
		Name:        chmodName,
		Usage:       chmodShort,
		Description: chmodLong,
	}
	chownCmd = &cli.Command{
		Name:        chownName,
		Usage:       chownShort,
		Description: chownLong,
	}
	unlinkCmd = &cli.Command{
		Name:        unlinkName,
		Usage:       unlinkShort,
		Description: unlinkLong,
	}
)

type flagError struct{ err error }

func (e flagError) Error() string { return e.err.Error() }
func init() {
	rootCmd.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:   "version",
			Usage:  "print the version",
			Action: Version,
		},
		&cli.BoolFlag{
			Name:   "info",
			Usage:  "print the info",
			Action: Info,
		},
	}
}

// Cmd
func Cmd() *cli.App {
	return &cli.App{
		Name:        appName,
		Usage:       gshortDesc,
		Description: glongDesc,
		Version:     version,
		Commands: []*cli.Command{
			rootCmd, duCmd, lsCmd, cpCmd, mvCmd, syncCmd, getCmd, putCmd, policyCmd, multipartCmd, mkdirCmd, rmCmd, chmodCmd, chownCmd, unlinkCmd,
		},
	}
}

// Du
func Du(ctx *cli.Context) error {
	// TODO: Implement the root command
	if internal.Threads > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(internal.Threads)
	}
	getCfg, err := internal.NewConfig(ctx)
	if err != nil {
		return err
	}
	return internal.GetUsage(ctx.Context, getCfg, ctx)
}

// Ls
func Ls(ctx *cli.Context) error {
	// TODO: Implement the root command
	if internal.Threads > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(internal.Threads)
	}
	getCfg, err := internal.NewConfig(ctx)
	if err != nil {
		return err
	}
	return internal.ListBucket(ctx.Context, getCfg, ctx)
}

// Execute
func Execute(ctx *cli.Context) error {
	// TODO: Implement the root command
	if internal.Threads > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(internal.Threads)
	}
	return nil
}

// Version
func Version(ctx *cli.Context, enable bool) error {
	// TODO: Implement the root command
	if internal.Threads > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(internal.Threads)
	}
	return nil
}

// Info
func Info(ctx *cli.Context, enable bool) error {
	// TODO: Implement the root command
	if internal.Threads > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(internal.Threads)
	}
	return nil
}
