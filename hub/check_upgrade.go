// Copyright (c) 2017-2020 VMware, Inc. or its affiliates
// SPDX-License-Identifier: Apache-2.0

package hub

import (
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/greenplum-db/gpupgrade/step"
)

type UpgradeChecker interface {
	UpgradeMaster(args UpgradeMasterArgs) error
	UpgradePrimaries(args UpgradePrimaryArgs) error
}

type upgradeChecker struct{}

func (upgradeChecker) UpgradeMaster(args UpgradeMasterArgs) error {
	return UpgradeMaster(args)
}

func (upgradeChecker) UpgradePrimaries(args UpgradePrimaryArgs) error {
	return UpgradePrimaries(args)
}

var upgrader UpgradeChecker = upgradeChecker{}

func (s *Server) CheckUpgrade(stream step.OutStreams, conns []*Connection) error {
	var wg sync.WaitGroup
	checkErrs := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		checkErrs <- upgrader.UpgradeMaster(UpgradeMasterArgs{
			Source:      s.Source,
			Target:      s.Target,
			StateDir:    s.StateDir,
			Stream:      stream,
			CheckOnly:   true,
			UseLinkMode: s.UseLinkMode,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		dataDirPairMap, dataDirPairsErr := s.GetDataDirPairs()
		if dataDirPairsErr != nil {
			checkErrs <- errors.Wrap(dataDirPairsErr, "failed to get source and target primary data directories")
			return
		}

		checkErrs <- upgrader.UpgradePrimaries(UpgradePrimaryArgs{
			CheckOnly:       true,
			MasterBackupDir: "",
			AgentConns:      conns,
			DataDirPairMap:  dataDirPairMap,
			Source:          s.Source,
			Target:          s.Target,
			UseLinkMode:     s.UseLinkMode,
		})
	}()

	wg.Wait()
	close(checkErrs)

	var multiErr *multierror.Error
	for err := range checkErrs {
		multiErr = multierror.Append(multiErr, err)
	}

	return multiErr.ErrorOrNil()
}
