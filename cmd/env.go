package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rancher/go-rancher/v2"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func EnvCommand() cli.Command {
	envLsFlags := []cli.Flag{
		listAllFlag(),
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "Only display IDs",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "'json' or Custom format: '{{.ID}} {{.Environment.Name}}'",
		},
	}

	return cli.Command{
		Name:      "environment",
		ShortName: "env",
		Usage:     "Interact with environments",
		Action:    defaultAction(envLs),
		Flags:     envLsFlags,
		Subcommands: []cli.Command{
			{
				Name:        "ls",
				Usage:       "List environments",
				Description: "\nWith an account API key, all compatible environments are listed. An environment API key lists only its own environment.\n\nExample:\n\t$ pasturestack env ls\n",
				ArgsUsage:   "None",
				Action:      envLs,
				Flags:       envLsFlags,
			},
			{
				Name:  "create",
				Usage: "Create an environment",
				Description: `
By default, an environment with the built-in compatible orchestration framework is created. This command only works with account API keys.

Example:

	$ pasturestack env create newEnv

To add an orchestration framework do TODO
	$ pasturestack env create -t kubernetes newK8sEnv
	$ pasturestack env create -t mesos newMesosEnv
	$ pasturestack env create -t swarm newSwarmEnv
`,
				ArgsUsage: "[NEWENVNAME...]",
				Action:    envCreate,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "template,t",
						Usage: "Environment template to create from",
						Value: "Cattle",
					},
				},
			},
			{
				Name:      "templates",
				ShortName: "template",
				Usage:     "Interact with environment templates",
				Action:    defaultAction(envTemplateLs),
				Flags:     envLsFlags,
				Subcommands: []cli.Command{
					{
						Name:      "export",
						Usage:     "Export an environment template to STDOUT",
						ArgsUsage: "[TEMPLATEID TEMPLATENAME...]",
						Action:    envTemplateExport,
						Flags:     []cli.Flag{},
					},
					{
						Name:      "import",
						Usage:     "Import an environment template to from file",
						ArgsUsage: "[FILE FILE...]",
						Action:    envTemplateImport,
						Flags: []cli.Flag{
							cli.BoolFlag{
								Name:  "public",
								Usage: "Make template public",
							},
						},
					},
				},
			},
			{
				Name:        "rm",
				Usage:       "Remove environment(s)",
				Description: "\nExample:\n\t$ pasturestack env rm 1a5\n\t$ pasturestack env rm newEnv\n",
				ArgsUsage:   "[ENVID ENVNAME...]",
				Action:      envRm,
				Flags:       []cli.Flag{},
			},
			{
				Name:  "deactivate",
				Usage: "Deactivate environment(s)",
				Description: `
Deactivate an environment by ID or name

Example:
	$ pasturestack env deactivate 1a5
	$ pasturestack env deactivate Default
`,
				ArgsUsage: "[ID NAME...]",
				Action:    envDeactivate,
				Flags:     []cli.Flag{},
			},
			{
				Name:  "activate",
				Usage: "Activate environment(s)",
				Description: `
Activate an environment by ID or name

Example:
	$ pasturestack env activate 1a5
	$ pasturestack env activate Default
`,
				ArgsUsage: "[ID NAME...]",
				Action:    envActivate,
				Flags:     []cli.Flag{},
			},
		},
	}
}

type EnvData struct {
	ID          string
	Environment *client.Project
}

type TemplateData struct {
	ID              string
	ProjectTemplate *client.ProjectTemplate
}

func NewEnvData(project client.Project) *EnvData {
	return &EnvData{
		ID:          project.Id,
		Environment: &project,
	}
}

func envRm(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	return forEachResourceWithClient(c, ctx, []string{"project"}, func(c *client.RancherClient, resource *client.Resource) (string, error) {
		return resource.Id, c.Delete(resource)
	})
}

func envCreate(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	name := RandomName()
	if ctx.NArg() > 0 {
		name = ctx.Args()[0]
	}

	data := map[string]interface{}{
		"name": name,
	}

	template := ctx.String("template")
	if template != "" {
		template, err := Lookup(c, template, "projectTemplate")
		if err != nil {
			return err
		}
		data["projectTemplateId"] = template.Id
	}

	var newEnv client.Project
	if err := c.Create("project", data, &newEnv); err != nil {
		return err
	}

	fmt.Println(newEnv.Id)
	return nil
}

func envLs(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	writer := NewTableWriter([][]string{
		{"ID", "ID"},
		{"NAME", "Environment.Name"},
		{"ORCHESTRATION", "Environment.Orchestration"},
		{"STATE", "Environment.State"},
		{"CREATED", "Environment.Created"},
	}, ctx)
	defer writer.Close()

	collection, err := c.Project.List(defaultListOpts(ctx))
	if err != nil {
		return err
	}

	collectiondata := collection.Data

	for {
		collection, _ = collection.Next()
		if collection == nil {
			break
		}
		collectiondata = append(collectiondata, collection.Data...)
		if !collection.Pagination.Partial {
			break
		}
	}

	for _, item := range collectiondata {
		writer.Write(NewEnvData(item))
	}

	return writer.Err()
}

func envDeactivate(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	return forEachResourceWithClient(c, ctx, []string{"project"}, func(c *client.RancherClient, resource *client.Resource) (string, error) {
		action, err := pickAction(resource, "deactivate")
		if err != nil {
			return "", err
		}
		return resource.Id, c.Action(resource.Type, action, resource, nil, resource)
	})
}

func envActivate(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	return forEachResourceWithClient(c, ctx, []string{"project"}, func(c *client.RancherClient, resource *client.Resource) (string, error) {
		action, err := pickAction(resource, "activate")
		if err != nil {
			return "", err
		}
		return resource.Id, c.Action(resource.Type, action, resource, nil, resource)
	})
}

func envTemplateLs(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	writer := NewTableWriter([][]string{
		{"ID", "ID"},
		{"NAME", "ProjectTemplate.Name"},
		{"DESC", "ProjectTemplate.Description"},
	}, ctx)
	defer writer.Close()

	collection, err := c.ProjectTemplate.List(defaultListOpts(ctx))
	if err != nil {
		return err
	}

	for _, item := range collection.Data {
		writer.Write(TemplateData{
			ID:              item.Id,
			ProjectTemplate: &item,
		})
	}

	return writer.Err()
}

func envTemplateImport(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	w, err := NewWaiter(ctx)
	if err != nil {
		return err
	}

	for _, file := range ctx.Args() {
		template := client.ProjectTemplate{
			IsPublic: ctx.Bool("public"),
		}
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(content, &template); err != nil {
			return err
		}

		created, err := c.ProjectTemplate.Create(&template)
		if err != nil {
			return err
		}

		w.Add(created.Id)
	}

	return w.Wait()
}

func envTemplateExport(ctx *cli.Context) error {
	c, err := GetRawClient(ctx)
	if err != nil {
		return err
	}

	for _, name := range ctx.Args() {
		r, err := Lookup(c, name, "projectTemplate")
		if err != nil {
			return err
		}

		template, err := c.ProjectTemplate.ById(r.Id)
		if err != nil {
			return err
		}

		stacks := []map[string]interface{}{}
		for _, s := range template.Stacks {
			data := map[string]interface{}{
				"name": s.Name,
			}
			if s.TemplateId != "" {
				data["template_id"] = s.TemplateId
			}
			if s.TemplateVersionId != "" {
				data["template_version_id"] = s.TemplateVersionId
			}
			if len(s.Answers) > 0 {
				data["answers"] = s.Answers
			}
			stacks = append(stacks, data)
		}

		data := map[string]interface{}{
			"name":        template.Name,
			"description": template.Description,
			"stacks":      stacks,
		}

		content, err := yaml.Marshal(&data)
		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(content)
		if err != nil {
			return err
		}
	}

	return nil
}
