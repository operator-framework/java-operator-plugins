# Enable kubebuilder-plugin for operator-sdk


To use kubebuilder-plugin for java operators we need to clone the operator-sdk repo. 

### Updates in Operator-SDK go.mod

- Add the kubebuilder plugin to `go.mod`

```
github.com/java-operator-sdk/kubebuilder-plugin v0.0.0-20210225171707-e42ea87455e3
```

- Replace the kubebuilder-plugin path in go-mod pointing to the local dir of your kube-builder repo. Example.

```
github.com/java-operator-sdk/kubebuilder-plugin => /Users/sushah/go/src/github.com/sujil02/kubebuilder-plugin
```

### Updates in Operator-SDK `internal/cmd/operator-sdk/cli/cli.go`

- Add the java-operator-sdk import

```
javav1 "github.com/java-operator-sdk/kubebuilder-plugin/pkg/quarkus/v1"
```

- Introduce the java bundle in `GetPluginsCLIAndRoot()` method. 
```
javaBundle, _ := plugin.NewBundle("quarkus"+plugins.DefaultNameQualifier, plugin.Version{Number: 1},
		&javav1.Plugin{},
	)
```

- Add the created javaBundle to the `cli.New`

```
    cli.WithPlugins(
			ansibleBundle,
			gov2Bundle,
			gov3Bundle,
			helmBundle,
			javaBundle,
		),
```


### Build and Install the Operator-SDK
```
go mod tidy
make install
```

Now that the plugin is integrated with the `operator-sdk` you can run the `init` command to generate the sample java operator

- Use the quarkus plugin flag
- Pick the domain and project name as prefered.

```
operator-sdk init --plugins quarkus --domain xyz.com --project-name java-op
```

Once the operator is scaffolded check for the following files

```
├── PROJECT
├── pom.xml
└── src
    └── main
        ├── java
        │   └── com
        │       └── xyz
        │           └── JavaOpOperator.java
        └── resources
            └── application.properties

```