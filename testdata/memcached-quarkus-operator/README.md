<img src="https://raw.githubusercontent.com/operator-framework/operator-sdk/master/website/static/operator_logo_sdk_color.svg" height="125px"></img>

# Java Operator SDK

## License

Operator SDK is under Apache 2.0 license. See the [LICENSE][license_file] file for details.

[license_file]:./LICENSE
[of-home]: https://github.com/operator-framework
[of-blog]: https://coreos.com/blog/introducing-operator-framework
[operator-link]: https://coreos.com/operators/

# Enable kubebuilder-plugin for operator-sdk


To use kubebuilder-plugin for java operators we need to clone the operator-sdk repo. 

### Updates in Operator-SDK go.mod

- Add the kubebuilder plugin to `go.mod`

```
github.com/operator-framework/java-operator-plugins v0.0.0-20210225171707-e42ea87455e3
```

- Replace the kubebuilder-plugin path in go-mod pointing to the local dir of your kube-builder repo. Example.

```
github.com/operator-framework/java-operator-plugins => /Users/sushah/go/src/github.com/sujil02/kubebuilder-plugin
```

### Updates in Operator-SDK `internal/cmd/operator-sdk/cli/cli.go`

- Add the java-operator-sdk import

```
javav1 "github.com/operator-framework/java-operator-plugins/pkg/quarkus/v1"
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
- Pick the domain and project name as preferred.

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

you can run the `crete api` command to generate the api's for java operator

```
operator-sdk create api --plugins quarkus --group cache --version v1 --kind Memcached
```

Once the api's added to the operator check for the following files

```
.
├── PROJECT
├── pom.xml
└── src
    └── main
        ├── java
        │   └── com
        │       └── lucky
        │           ├── JavaOpOperator.java
        │           ├── Memcached.java
        │           ├── MemcachedController.java
        │           ├── MemcachedSpec.java
        │           └── MemcachedStatus.java
        └── resources
            └── application.properties
```

Now, create kind or minikube cluster. Then, create CRD's and CR as follows from the k8s folder.

```
kubectl apply -f crd.yaml
```

```
kubectl apply -f memcached-sample.yaml
```

Change the Memcached file as follows for group and version as it is aligned with Deployment.

```
@Version("v1alpha1")
@Group("cache.example.com")
```

At the end, change the MemcachedController file as show in the ```Memcached Quarkus Operator```.



