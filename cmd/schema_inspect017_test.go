package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	cli "github.com/PastureStack/cli/internal/legacycli"
	client "github.com/PastureStack/cli/internal/rancherclient/v2"
)

const cli017SchemasPath = "/v2-beta/projects/1a7/schemas"
const cli017StackPath = "/v2-beta/projects/1a7/stacks/1st7"

// Run the real legacy-parser context and command action. Its returned error is
// captured before the outer CLI can exit the test process; it is never replaced
// inside the production action. Config is explicitly absent, not the user's.
func cli017Action(t *testing.T, command cli.Command, serverURL string, arguments ...string) (error, string) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer
	defer func() {
		os.Stdout = originalStdout
		_ = writer.Close()
		_ = reader.Close()
	}()
	output := make(chan []byte, 1)
	go func() {
		value, _ := io.ReadAll(reader)
		output <- value
	}()

	app := cli.NewApp()
	app.Name = "cli017-focused"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config", Value: filepath.Join(t.TempDir(), "absent.json")},
		cli.StringFlag{Name: "url", Value: serverURL},
		cli.StringFlag{Name: "access-key", Value: "cli017-access"},
		cli.StringFlag{Name: "secret-key", Value: "cli017-secret"},
		cli.StringFlag{Name: "environment", Value: "1a7"},
		cli.BoolFlag{Name: "wait"},
		cli.IntFlag{Name: "wait-timeout", Value: 600},
		cli.StringFlag{Name: "wait-state"},
	}
	native := command.Action
	called := false
	var actionErr error
	command.Action = func(ctx *cli.Context) error {
		called = true
		actionErr = native(ctx)
		return nil
	}
	app.Commands = []cli.Command{command}
	if err := app.Run(append([]string{app.Name, command.Name}, arguments...)); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("native action was not called")
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = originalStdout
	return actionErr, string(<-output)
}

func cli017Server(t *testing.T, schemaStatus, resourceStatus int) (*httptest.Server, chan string) {
	t.Helper()
	paths := make(chan string, 16)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		paths <- request.URL.Path
		if request.Method != http.MethodGet {
			http.Error(response, "unexpected test method", http.StatusMethodNotAllowed)
			return
		}
		user, secret, ok := request.BasicAuth()
		if !ok || user != "cli017-access" || secret != "cli017-secret" {
			http.Error(response, "unexpected test authorization", http.StatusUnauthorized)
			return
		}
		response.Header().Set("Content-Type", "application/json")
		switch request.URL.Path {
		case cli017SchemasPath:
			if schemaStatus != http.StatusOK {
				response.WriteHeader(schemaStatus)
				_, _ = io.WriteString(response, `{"type":"error","status":403,"code":"KeyPolicyDenied","message":"closed test key"}`)
				return
			}
			response.Header().Set("X-API-Schemas", server.URL+cli017SchemasPath)
			_, _ = fmt.Fprintf(response, `{"type":"collection","data":[{"id":"stack","pluralName":"stacks","collectionMethods":["GET"],"resourceMethods":["GET"],"links":{"collection":%q}}]}`,
				server.URL+"/v2-beta/projects/1a7/stacks")
		case cli017StackPath:
			if resourceStatus != http.StatusOK {
				response.WriteHeader(resourceStatus)
				_, _ = io.WriteString(response, `{"type":"error","status":403,"code":"KeyScopeDenied","message":"blocked test resource"}`)
				return
			}
			_, _ = fmt.Fprintf(response, `{"id":"1st7","type":"stack","name":"owned-test","links":{"self":%q},"actions":{"remove":%q}}`,
				server.URL+cli017StackPath, server.URL+cli017StackPath+"?action=remove")
		default:
			http.NotFound(response, request)
		}
	}))
	t.Cleanup(server.Close)
	return server, paths
}

func cli017Paths(paths chan string) []string {
	var result []string
	for {
		select {
		case path := <-paths:
			result = append(result, path)
		default:
			return result
		}
	}
}

func TestCLISchemaInspect017GetClientUsesSingleSchemasSuffix(t *testing.T) {
	server, paths := cli017Server(t, http.StatusOK, http.StatusOK)
	var value *client.RancherClient
	err, output := cli017Action(t, cli.Command{Name: "client", Action: func(ctx *cli.Context) error {
		var err error
		value, err = GetClient(ctx)
		return err
	}}, server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if value == nil || value.GetOpts().Url != server.URL+cli017SchemasPath || output != "" {
		t.Fatal("GetClient did not retain the exact EnvironmentURL schema endpoint")
	}
	if got, want := cli017Paths(paths), []string{cli017SchemasPath}; !reflect.DeepEqual(got, want) {
		t.Fatalf("wire paths=%v, want %v", got, want)
	}
}

func TestCLISchemaInspect017InspectFlushesSuccessJSON(t *testing.T) {
	server, paths := cli017Server(t, http.StatusOK, http.StatusOK)
	err, output := cli017Action(t, InspectCommand(), server.URL, "--type", "stack", "1st7")
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]interface{}
	if err := json.Unmarshal([]byte(output), &value); err != nil {
		t.Fatalf("native inspect did not flush one valid JSON resource: %v", err)
	}
	if value["id"] != "1st7" || value["type"] != "stack" || value["name"] != "owned-test" {
		t.Fatalf("unexpected inspect result: %v", value)
	}
	if _, present := value["links"]; present {
		t.Fatal("inspect --links=false changed")
	}
	if _, present := value["actions"]; present {
		t.Fatal("inspect --links=false changed")
	}
	// NewWaiter currently constructs a second client even without --wait. These
	// are valid repeated schema reads, not the erroneous /schemas/schemas URL.
	want := []string{cli017SchemasPath, cli017SchemasPath, cli017StackPath, cli017StackPath}
	if got := cli017Paths(paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("wire paths=%v, want %v", got, want)
	}
}

func TestCLISchemaInspect017InspectReturnsMissingConfigURL(t *testing.T) {
	err, output := cli017Action(t, InspectCommand(), "", "--type", "stack", "1st7")
	if !errors.Is(err, errNoURL) || output != "" {
		t.Fatalf("inspect discarded pre-HTTP setup error: err=%v output=%q", err, output)
	}
}

func TestCLISchemaInspect017InspectReturnsSchemaHTTPFailure(t *testing.T) {
	server, paths := cli017Server(t, http.StatusForbidden, http.StatusOK)
	err, output := cli017Action(t, InspectCommand(), server.URL, "--type", "stack", "1st7")
	var apiError *client.ApiError
	if !errors.As(err, &apiError) || apiError.StatusCode != http.StatusForbidden || output != "" {
		t.Fatalf("inspect discarded schema HTTP denial: err=%v output=%q", err, output)
	}
	if !strings.Contains(apiError.Body, "code=KeyPolicyDenied") {
		t.Fatal("schema denial lost the client's existing remote error representation")
	}
	if got, want := cli017Paths(paths), []string{cli017SchemasPath}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected requests after schema denial: %v", got)
	}
}

func TestCLISchemaInspect017InspectReturnsResourceHTTPFailure(t *testing.T) {
	server, paths := cli017Server(t, http.StatusOK, http.StatusForbidden)
	err, output := cli017Action(t, InspectCommand(), server.URL, "--type", "stack", "1st7")
	var exitCode interface{ ExitCode() int }
	if !errors.As(err, &exitCode) || exitCode.ExitCode() != 1 || !strings.Contains(output, "KeyScopeDenied") {
		t.Fatalf("inspect discarded resource HTTP denial: err=%v output=%q", err, output)
	}
	if strings.Contains(output, `"id":"1st7"`) || strings.Contains(output, `"type":"stack"`) {
		t.Fatal("denied resource was emitted as successful JSON")
	}
	want := []string{cli017SchemasPath, cli017SchemasPath, cli017StackPath}
	if got := cli017Paths(paths); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected requests after resource denial: %v", got)
	}
}
