package templates

import "sigs.k8s.io/kubebuilder/v3/pkg/model/file"

var _ file.Template = &OperatorFile{}

type OperatorFile struct {
	file.TemplateMixin
}

func (f *OperatorFile) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = "Main.java"
	}

	f.TemplateBody = operatorTemplate

	return nil
}

// TODO: pass in the name of the operator i.e. replace Memcached
const operatorTemplate = `
package io.javaoperatorsdk.sample.memcached;

import io.javaoperatorsdk.operator.Operator;
import io.quarkus.runtime.Quarkus;
import io.quarkus.runtime.QuarkusApplication;
import io.quarkus.runtime.annotations.QuarkusMain;
import javax.inject.Inject;

@QuarkusMain
public class MemcachedOperator implements QuarkusApplication {

  @Inject Operator operator;

  public static void main(String... args) {
    Quarkus.run(MemcachedOperator.class, args);
  }

  @Override
  public int run(String... args) throws Exception {
    operator.start();

    Quarkus.waitForExit();
    return 0;
  }
}
`
