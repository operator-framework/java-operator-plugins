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
