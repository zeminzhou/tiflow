// Copyright 2022 PingCAP, Inc.
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

package internal

import (
	"context"

	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/processor/tablepb"
)

// TableExecutor is an abstraction for "Processor".
//
// This interface is so designed that it would be the least problematic
// to adapt the current Processor implementation to it.
// TODO find a way to make the semantics easier to understand.
type TableExecutor interface {
	// AddTableSpan add a new table span with `startTs`
	// if `isPrepare` is true, the 1st phase of the 2 phase scheduling protocol.
	// if `isPrepare` is false, the 2nd phase.
	AddTableSpan(
		ctx context.Context, span tablepb.Span, startTs model.Ts, isPrepare bool,
	) (done bool, err error)

	// IsAddTableSpanFinished make sure the requested table span is in the proper status
	IsAddTableSpanFinished(span tablepb.Span, isPrepare bool) (done bool)

	// RemoveTableSpan remove the table, return true if the table is already removed
	RemoveTableSpan(span tablepb.Span) (done bool)
	// IsRemoveTableSpanFinished convince the table is fully stopped.
	// return false if table is not stopped
	// return true and corresponding checkpoint otherwise.
	IsRemoveTableSpanFinished(span tablepb.Span) (model.Ts, bool)

	// GetTableSpanCount should return the number of table spans that are being run,
	// being added and being removed.
	//
	// NOTE: two subsequent calls to the method should return the same
	// result, unless there is a call to AddTable, RemoveTable, IsAddTableFinished
	// or IsRemoveTableSpanFinished in between two calls to this method.
	GetTableSpanCount() int

	// GetCheckpoint returns the local checkpoint-ts and resolved-ts of
	// the processor. Its calculation should take into consideration all
	// tables that would have been returned if GetTableSpanCount had been
	// called immediately before.
	GetCheckpoint() (checkpointTs, resolvedTs model.Ts)

	// GetTableSpanStatus return the checkpoint and resolved ts for the given table span.
	GetTableSpanStatus(span tablepb.Span) tablepb.TableStatus
}
