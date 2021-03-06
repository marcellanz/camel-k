/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/util/log"
	"github.com/spf13/cobra"
)

func newCmdLog(rootCmdOptions *RootCmdOptions) *cobra.Command {
	options := logCmdOptions{
		RootCmdOptions: rootCmdOptions,
	}

	cmd := cobra.Command{
		Use:   "log integration",
		Short: "Print the logs of an integration",
		Long:  `Print the logs of an integration.`,
		Args:  options.validate,
		RunE:  options.run,
	}

	// completion support
	configureKnownCompletions(&cmd)

	return &cmd
}

type logCmdOptions struct {
	*RootCmdOptions
}

func (o *logCmdOptions) validate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("accepts 1 arg, received %d", len(args))
	}

	return nil
}

func (o *logCmdOptions) run(cmd *cobra.Command, args []string) error {
	integration := v1alpha1.Integration{
		TypeMeta: metav1.TypeMeta{
			Kind:       v1alpha1.IntegrationKind,
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: o.Namespace,
			Name:      args[0],
		},
	}

	if err := sdk.Get(&integration); err != nil {
		return err
	}
	if err := log.Print(o.Context, &integration); err != nil {
		return err
	}

	// Let's add a wait point, otherwise the script terminates
	<-o.Context.Done()

	return nil
}
