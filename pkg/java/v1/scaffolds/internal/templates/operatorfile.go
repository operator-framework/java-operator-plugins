package templates

import (
	"fmt"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kubebuilder/v3/pkg/model/file"
)

var _ file.Template = &OperatorFile{}

type OperatorFile struct {
	file.TemplateMixin

	// Package is the source files package
	Package string

	// Name of the operator used for the main file.
	OperatorName string
}

const (
	FilePathSep = string(filepath.Separator)
	javaPath    = "src" + FilePathSep + "main" + FilePathSep + "java"
)

func prependJavaPath(filename string, pkg string) string {
	return javaPath + FilePathSep + pkg + FilePathSep + filename
}

func asPath(s string) string {
	return strings.ReplaceAll(s, ".", FilePathSep)
}

func (f *OperatorFile) SetTemplateDefaults() error {
	if f.OperatorName == "" {
		return fmt.Errorf("invalid operator name")
	}

	if f.Path == "" {
		if strings.HasSuffix(strings.ToLower(f.OperatorName), "operator") {
			f.Path = prependJavaPath(f.OperatorName+".java", asPath(f.Package))
		} else {
			f.Path = prependJavaPath(f.OperatorName+"Operator.java", asPath(f.Package))
		}
	}

	f.TemplateBody = operatorTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const operatorTemplate = `
package {{ .Package }};

import io.javaoperatorsdk.operator.Operator;
import io.quarkus.runtime.Quarkus;
import io.quarkus.runtime.QuarkusApplication;
import io.quarkus.runtime.annotations.QuarkusMain;
import javax.inject.Inject;

@QuarkusMain
public class {{ .OperatorName }}Operator implements QuarkusApplication {

  @Inject Operator operator;

  public static void main(String... args) {
    Quarkus.run({{ .OperatorName }}Operator.class, args);
  }

  @Override
  public int run(String... args) throws Exception {
    operator.start();

    Quarkus.waitForExit();
    return 0;
  }
}
`
