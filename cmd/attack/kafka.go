// Copyright 2020 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package attack

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/chaos-mesh/chaosd/cmd/server"
	"github.com/chaos-mesh/chaosd/pkg/core"
	"github.com/chaos-mesh/chaosd/pkg/server/chaosd"
	"github.com/chaos-mesh/chaosd/pkg/utils"
)

func NewKafkaAttackCommand(uid *string) *cobra.Command {
	options := core.NewKafkaCommand()
	dep := fx.Options(
		server.Module,
		fx.Provide(func() *core.KafkaCommand {
			options.UID = *uid
			return options
		}),
	)

	cmd := &cobra.Command{
		Use:   "kafka <subcommand>",
		Short: "Kafka attack related commands",
	}

	cmd.AddCommand(
		NewKafkaFloodCommand(dep, options),
		NewKafkaIOCommand(dep, options),
	)

	return cmd
}

func NewKafkaFloodCommand(dep fx.Option, options *core.KafkaCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "flood [options]",
		Short:             "Flood kafka cluster with messages",
		ValidArgsFunction: cobra.NoFileCompletions,
		Run: func(*cobra.Command, []string) {
			options.Action = core.KafkaFloodAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(kafkaCommandFunc)).Run()
		},
	}

	return cmd
}

func NewKafkaIOCommand(dep fx.Option, options *core.KafkaCommand) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "io [options]",
		Short:             "Make kafka cluster non-writable/non-readable",
		ValidArgsFunction: cobra.NoFileCompletions,
		Run: func(*cobra.Command, []string) {
			options.Action = core.KafkaIOAction
			utils.FxNewAppWithoutLog(dep, fx.Invoke(kafkaCommandFunc)).Run()
		},
	}

	return cmd
}

func kafkaCommandFunc(options *core.KafkaCommand, chaos *chaosd.Server) {
	options.CompleteDefaults()

	if err := options.Validate(); err != nil {
		utils.ExitWithError(utils.ExitBadArgs, err)
	}

	uid, err := chaos.ExecuteAttack(chaosd.KafkaAttack, options, core.CommandMode)
	if err != nil {
		utils.ExitWithError(utils.ExitError, err)
	}

	utils.NormalExit(fmt.Sprintf("Attack jvm successfully, uid: %s", uid))
}