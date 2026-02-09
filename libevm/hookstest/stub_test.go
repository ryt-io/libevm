// Copyright 2024-2025 the libevm authors.
//
// The libevm additions to go-ethereum are free software: you can redistribute
// them and/or modify them under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The libevm additions are distributed in the hope that they will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see
// <http://www.gnu.org/licenses/>.

package hookstest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ryt-io/libevm/common"
	"github.com/ryt-io/libevm/core"
	"github.com/ryt-io/libevm/core/rawdb"
	"github.com/ryt-io/libevm/libevm/ethtest"
	"github.com/ryt-io/libevm/log"
	"github.com/ryt-io/libevm/params"
)

func TestSetupGenesisBlockWithStub(t *testing.T) {
	// The original bug was due to [Stub] resulting in an error when being
	// marshalled to JSON for [rawdb.WriteChainConfig], which resulted in
	// [log.Crit] being called.
	l := log.Root()
	t.Cleanup(func() { log.SetDefault(l) })
	log.SetDefault(log.NewLogger(ethtest.NewTBLogHandler(t, log.LevelDebug)))

	stub := &Stub{}
	extras := stub.Register(t)

	config := &params.ChainConfig{}
	extras.ChainConfig.Set(config, stub)
	gen := &core.Genesis{
		Config: config,
	}

	db, cache, _ := ethtest.NewEmptyStateDB(t)

	// An eventual call to this function was the root cause.
	rawdb.WriteChainConfig(db, common.Hash{1, 2, 3, 4, 5, 6}, config)

	// Also check calls to this function, which was the desired behaviour.
	_, _, err := core.SetupGenesisBlock(db, cache.TrieDB(), gen)
	require.NoError(t, err, "core.SetupGenesisBlock([%T with %T as registered %T extra])", gen, stub, config)
}
