<img src="https://raw.githubusercontent.com/operator-framework/operator-sdk/master/website/static/operator_logo_sdk_color.svg" height="125px"></img>

# Java Operator Plugin

## Overview

This project is a component of the [Operator Framework][of-home], an
open source toolkit to manage Kubernetes native applications, called
Operators, in an effective, automated, and scalable way. Read more in
the [introduction blog post][of-blog].

[Operators][operator-link] make it easy to manage complex stateful
applications on top of Kubernetes. However writing an operator today can
be difficult because of challenges such as using low level APIs, writing
boilerplate, and a lack of modularity which leads to duplication.



## License

Operator SDK is under Apache 2.0 license. See the [LICENSE][license_file] file for details.

[license_file]:./LICENSE
[of-home]: https://github.com/operator-framework
[of-blog]: https://coreos.com/blog/introducing-operator-framework
[operator-link]: https://coreos.com/operators/

## Enable java-operator-plugins for operator-sdk


To use java-operator-plugins for java operators we need to clone the operator-sdk repo. 

### Updates in Operator-SDK go.mod

- Add the kubebuilder plugin to `go.mod`

```
github.com/operator-framework/java-operator-plugins v0.0.0-20210225171707-e42ea87455e3
```

- Replace the java-operator-plugins path in go-mod pointing to the local dir of your kube-builder repo. Example.

```
github.com/operator-framework/java-operator-plugins => /Users/sushah/go/src/github.com/sujil02/java-operator-plugins
```

### Updates in Operator-SDK `internal/cmd/operator-sdk/cli/cli.go`

- Add the java-operator-sdk import

```
quarkusv1 "github.com/operator-framework/java-operator-plugins/pkg/quarkus/v1alpha"
```

- Introduce the java bundle in `GetPluginsCLIAndRoot()` method. 
```
quarkusBundle, _ := plugin.NewBundle("quarkus"+plugins.DefaultNameQualifier, plugin.Version{Number: 1},
		&quarkusv1.Plugin{},
	)
```

- Add the created quarkusBundle to the `cli.New`

```
    cli.WithPlugins(
			ansibleBundle,
			gov2Bundle,
			gov3Bundle,
			helmBundle,
			quarkusBundle,
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
.
├── Makefile
├── PROJECT
├── pom.xml
└── src
    └── main
        ├── java
        └── resources
            └── application.properties

4 directories, 4 files
```
