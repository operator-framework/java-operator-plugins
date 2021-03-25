package templates

import "sigs.k8s.io/kubebuilder/v3/pkg/model/file"

var _ file.Template = &PomXmlFile{}

type PomXmlFile struct {
	file.TemplateMixin

	// Package is the source files package
	Package         string
	ProjectName     string
	OperatorVersion string
}

func (f *PomXmlFile) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "pom.xml"
	}

	f.TemplateBody = pomxmlTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const pomxmlTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <groupId>{{ .Package }}</groupId>
  <artifactId>{{ .ProjectName }}</artifactId>
  <name>{{ .ProjectName }}</name>
  <version>{{ .OperatorVersion }}</version>
  <packaging>jar</packaging>
  <properties>
    <compiler-plugin.version>3.8.1</compiler-plugin.version>
    <maven.compiler.parameters>true</maven.compiler.parameters>
    <maven.compiler.source>11</maven.compiler.source>
      <maven.compiler.target>11</maven.compiler.target>
      <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
      <project.reporting.outputEncoding>UTF-8</project.reporting.outputEncoding>
      <fabric8-client.version>5.2.1</fabric8-client.version>
      <quarkus-sdk.version>1.8.0</quarkus-sdk.version>
      <java-sdk.version>1.8.2</java-sdk.version>
      <quarkus.version>1.12.2.Final</quarkus.version>
      <quarkus.native.builder-image>quay.io/quarkus/ubi-quarkus-native-image:19.3.1-java11</quarkus.native.builder-image>
  </properties>

  <dependencyManagement>
    <dependencies>
      <dependency>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-bom</artifactId>
        <version>${quarkus.version}</version>
        <type>pom</type>
	    <scope>import</scope>
      </dependency>
    </dependencies>
  </dependencyManagement>
  <dependencies>
    <dependency>
      <groupId>io.quarkiverse.operatorsdk</groupId>
      <artifactId>quarkus-operator-sdk</artifactId>
      <version>${quarkus-sdk.version}</version>
    </dependency>
    <dependency>
      <groupId>io.fabric8</groupId>
      <artifactId>crd-generator-apt</artifactId>
      <version>${fabric8-client.version}</version>
    </dependency>
    <dependency>
      <groupId>io.javaoperatorsdk</groupId>
      <artifactId>operator-framework</artifactId>
      <version>${java-sdk.version}</version>
    </dependency>
    <dependency>
      <groupId>io.quarkus</groupId>
      <artifactId>quarkus-container-image-jib</artifactId>
      <version>${quarkus.version}</version>
    </dependency>
  </dependencies>

  <build>
    <plugins>
      <plugin>
        <groupId>io.quarkus</groupId>
        <artifactId>quarkus-maven-plugin</artifactId>
        <version>${quarkus.version}</version>
        <executions>
          <execution>
            <goals>
              <goal>build</goal>
            </goals>
          </execution>
        </executions>
    </plugin>
    <plugin>
      <artifactId>maven-compiler-plugin</artifactId>
      <version>${compiler-plugin.version}</version>
    </plugin>
    </plugins>
  </build>

  <profiles>
    <profile>
      <id>native</id>
      <properties>
        <quarkus.package.type>native</quarkus.package.type>
      </properties>
    </profile>
  </profiles>

</project>
`
