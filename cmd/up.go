package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	cli "github.com/PastureStack/cli/internal/legacycli"
)

const composeBinaryEnvironment = "PASTURESTACK_COMPOSE_BINARY"

// UpCommand delegates Compose reconciliation to the dedicated maintained
// binary. Keeping a second embedded executor here previously duplicated its
// dependency and security lifecycle.
func UpCommand() cli.Command {
	return cli.Command{
		Name:   "up",
		Usage:  "Bring all services up with pasturestack-compose",
		Action: runComposeUp,
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "pull,p", Usage: "Pull images before upgrade"},
			cli.BoolFlag{Name: "prune", Usage: "Remove services absent from the Compose file"},
			cli.BoolFlag{Name: "d", Usage: "Do not block and log"},
			cli.BoolFlag{Name: "render", Usage: "Render processed Compose files and exit"},
			cli.BoolFlag{Name: "upgrade,u,recreate", Usage: "Upgrade changed services"},
			cli.BoolFlag{Name: "force-upgrade,force-recreate", Usage: "Upgrade regardless of changes"},
			cli.BoolFlag{Name: "confirm-upgrade", Usage: "Confirm a successful upgrade"},
			cli.BoolFlag{Name: "rollback", Usage: "Roll back the previous deployment"},
			cli.IntFlag{Name: "batch-size", Usage: "Containers to upgrade at once", Value: 2},
			cli.IntFlag{Name: "interval", Usage: "Update interval in milliseconds", Value: 1000},
			cli.StringFlag{Name: "description", Usage: "Stack description"},
			cli.StringFlag{Name: "rancher-file", Usage: "Compatibility Compose settings file"},
			cli.StringFlag{Name: "env-file,e", Usage: "Environment file"},
			cli.StringSliceFlag{Name: "file,f", Usage: "Compose file"},
			cli.StringFlag{Name: "stack,s", Usage: "Project name"},
		},
	}
}

func runComposeUp(ctx *cli.Context) error {
	config, err := lookupConfig(ctx)
	if err != nil {
		return err
	}
	environmentURL, err := config.EnvironmentURL()
	if err != nil {
		return err
	}

	binary := os.Getenv(composeBinaryEnvironment)
	if binary == "" {
		binary, err = exec.LookPath("pasturestack-compose")
		if err != nil {
			return fmt.Errorf("pasturestack-compose is required for 'up'; install the matching maintained Compose CLI or set %s: %w", composeBinaryEnvironment, err)
		}
	}

	args := make([]string, 0, 32)
	for _, file := range ctx.StringSlice("file") {
		args = append(args, "--file", file)
	}
	if value := ctx.String("stack"); value != "" {
		args = append(args, "--project-name", value)
	}
	if value := ctx.String("rancher-file"); value != "" {
		args = append(args, "--platform-file", value)
	}
	if value := ctx.String("env-file"); value != "" {
		args = append(args, "--env-file", value)
	}
	args = append(args, "up")
	for _, name := range []string{"pull", "prune", "d", "render", "upgrade", "force-upgrade", "confirm-upgrade", "rollback"} {
		if ctx.Bool(name) {
			args = append(args, "--"+name)
		}
	}
	args = append(args, "--batch-size", fmt.Sprint(ctx.Int("batch-size")))
	args = append(args, "--interval", fmt.Sprint(ctx.Int("interval")))
	if value := ctx.String("description"); value != "" {
		args = append(args, "--description", value)
	}
	args = append(args, ctx.Args()...)

	command := exec.Command(binary, args...)
	command.Env = withEnvironment(os.Environ(), map[string]string{
		"PLATFORM_URL":        environmentURL,
		"PLATFORM_ACCESS_KEY": config.AccessKey,
		"PLATFORM_SECRET_KEY": config.SecretKey,
		"PASTURESTACK_LOCALE": ctx.GlobalString("locale"),
	})
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func withEnvironment(base []string, overrides map[string]string) []string {
	result := make([]string, 0, len(base)+len(overrides))
	for _, entry := range base {
		name, _, ok := strings.Cut(entry, "=")
		if !ok {
			result = append(result, entry)
			continue
		}
		replaced := false
		for key := range overrides {
			if strings.EqualFold(name, key) {
				replaced = true
				break
			}
		}
		if !replaced {
			result = append(result, entry)
		}
	}

	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		result = append(result, key+"="+overrides[key])
	}
	return result
}
