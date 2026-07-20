package main

import (
	"os"

	"github.com/PastureStack/cli/cmd"
	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var VERSION = "dev"

var AppHelpTemplate = `{{.Usage}}

Usage: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]

Version: {{.Version}}
{{if .Flags}}
Options:
  {{range .Flags}}{{if .Hidden}}{{else}}{{.}}
  {{end}}{{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`

var CommandHelpTemplate = `{{.Usage}}
{{if .Description}}{{.Description}}{{end}}
Usage: pasturestack [global options] {{.Name}} {{if .Flags}}[OPTIONS] {{end}}{{if ne "None" .ArgsUsage}}{{if ne "" .ArgsUsage}}{{.ArgsUsage}}{{else}}[arg...]{{end}}{{end}}

{{if .Flags}}Options:{{range .Flags}}
	 {{.}}{{end}}{{end}}
`

func main() {
	if err := mainErr(); err != nil {
		logrus.Fatal(err)
	}
}

func mainErr() error {
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = CommandHelpTemplate

	app := cli.NewApp()
	app.Name = "pasturestack"
	app.Usage = cmd.T("app.usage")
	app.Before = func(ctx *cli.Context) error {
		cmd.SetLocale(ctx.GlobalString("locale"))
		if ctx.GlobalBool("debug") {
			logrus.SetLevel(logrus.DebugLevel)
		}
		return nil
	}
	app.Version = VERSION
	app.Author = "Original contributors and PastureStack maintainers"
	app.Email = ""
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debug logging",
		},
		cli.StringFlag{
			Name:   "config,c",
			Usage:  cmd.T("flag.config"),
			EnvVar: "PASTURESTACK_CLIENT_CONFIG,PLATFORM_CLIENT_CONFIG,RANCHER_CLIENT_CONFIG",
		},
		cli.StringFlag{
			Name:   "environment,env",
			Usage:  "Environment name or ID",
			EnvVar: "PASTURESTACK_ENVIRONMENT,PLATFORM_ENVIRONMENT,RANCHER_ENVIRONMENT",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  cmd.T("flag.url"),
			EnvVar: "PASTURESTACK_URL,PLATFORM_URL,RANCHER_URL",
		},
		cli.StringFlag{
			Name:   "access-key",
			Usage:  cmd.T("flag.accessKey"),
			EnvVar: "PASTURESTACK_ACCESS_KEY,PLATFORM_ACCESS_KEY,RANCHER_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "secret-key",
			Usage:  cmd.T("flag.secretKey"),
			EnvVar: "PASTURESTACK_SECRET_KEY,PLATFORM_SECRET_KEY,RANCHER_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "host",
			Usage:  "Host used for docker command",
			EnvVar: "PASTURESTACK_DOCKER_HOST,PLATFORM_DOCKER_HOST,RANCHER_DOCKER_HOST",
		},
		cli.BoolFlag{
			Name:  "wait,w",
			Usage: "Wait for resource to reach resting state",
		},
		cli.IntFlag{
			Name:  "wait-timeout",
			Usage: "Timeout in seconds to wait",
			Value: 600,
		},
		cli.StringFlag{
			Name:  "wait-state",
			Usage: "State to wait for (active, healthy, etc)",
		},
		// These hidden flags bridge the preserved Compose compatibility library.
		cli.StringFlag{
			Name:   "rancher-file",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "env-file",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "locale",
			Usage:  "Operator message locale: en-US or zh-TW",
			EnvVar: "PASTURESTACK_LOCALE",
			Value:  "en-US",
		},
		cli.StringSliceFlag{
			Name:   "file,f",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "project-name",
			Hidden: true,
		},
	}
	app.Commands = []cli.Command{
		cmd.CatalogCommand(),
		cmd.ConfigCommand(),
		cmd.DockerCommand(),
		cmd.EnvCommand(),
		cmd.EventsCommand(),
		cmd.ExecCommand(),
		cmd.ExportCommand(),
		cmd.HostCommand(),
		cmd.LogsCommand(),
		cmd.PsCommand(),
		cmd.RestartCommand(),
		cmd.RmCommand(),
		cmd.RunCommand(),
		cmd.ScaleCommand(),
		cmd.SecretCommand(),
		cmd.SSHCommand(),
		cmd.StackCommand(),
		cmd.StartCommand(),
		cmd.StopCommand(),
		cmd.UpCommand(),
		cmd.VolumeCommand(),
		cmd.InspectCommand(),
		cmd.WaitCommand(),
	}

	return app.Run(os.Args)
}
