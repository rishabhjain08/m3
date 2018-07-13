// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package aggregator

import (
	"fmt"

	"math"

	"sync"

	"time"

	maggregation "github.com/m3db/m3metrics/aggregation"

	"github.com/m3db/m3metrics/metric/id"

	"github.com/m3db/m3metrics/metric/unaggregated"

	"github.com/m3db/m3metrics/pipeline/applied"

	"github.com/m3db/m3metrics/policy"

	"github.com/m3db/m3metrics/transformation"

	"github.com/willf/bitset"
)

type lockedCounterAggregation struct {
	sync.Mutex

	closed bool

	// sourcesReady is only used for elements receiving forwarded metrics.
	// It determines whether the current aggregation can use its source set
	// to determine whether it has received data from all forwarding sources
	// to perform eager forwarding if enabled.
	sourcesReady bool

	// expectedSources is only used for elements receiving forwarded metrics.
	// It keeps track of all the sources the current aggregation expect to receive
	// data from.
	expectedSources *bitset.BitSet

	// seenSources keeps track of all the sources the current aggregation has
	// seen so far.
	seenSources *bitset.BitSet

	// consumeState is only used for elements receiving forwarded metrics. It
	// describes whether the current aggregation is ready to be consumed or has
	// been consumed. This in turn determines whether the aggregation can be
	// eagerly consumed, or should be skipped during consumption.
	consumeState consumeState
	aggregation  counterAggregation
}

func (lockedAgg *lockedCounterAggregation) close() {
	if lockedAgg.closed {
		return
	}
	lockedAgg.closed = true
	lockedAgg.expectedSources = nil
	lockedAgg.seenSources = nil
	lockedAgg.aggregation.Close()
}

type timedCounter struct {
	startAtNanos int64 // start time of an aggregation window
	lockedAgg    *lockedCounterAggregation
}

func (ta *timedCounter) Reset() {
	ta.startAtNanos = 0
	ta.lockedAgg = nil
}

// CounterElem is an element storing time-bucketed aggregations.
type CounterElem struct {
	elemBase
	counterElemBase

	values              []timedCounter // metric aggregations sorted by time in ascending order
	toConsume           []timedCounter // small buffer to avoid memory allocations during consumption
	lastConsumedAtNanos int64          // last consumed at in Unix nanoseconds
	lastConsumedValues  []float64      // last consumed values
}

// NewCounterElem creates a new element for the given metric type.
func NewCounterElem(
	incomingMetricType IncomingMetricType,
	id id.RawID,
	sp policy.StoragePolicy,
	aggTypes maggregation.Types,
	pipeline applied.Pipeline,
	numForwardedTimes int,
	opts Options,
) (*CounterElem, error) {
	e := &CounterElem{
		elemBase: newElemBase(opts),
		values:   make([]timedCounter, 0, defaultNumAggregations), // in most cases values will have two entries
	}
	if err := e.ResetSetData(incomingMetricType, id, sp, aggTypes, pipeline, numForwardedTimes); err != nil {
		return nil, err
	}
	return e, nil
}

// MustNewCounterElem creates a new element, or panics if the input is invalid.
func MustNewCounterElem(
	incomingMetricType IncomingMetricType,
	id id.RawID,
	sp policy.StoragePolicy,
	aggTypes maggregation.Types,
	pipeline applied.Pipeline,
	numForwardedTimes int,
	opts Options,
) *CounterElem {
	elem, err := NewCounterElem(incomingMetricType, id, sp, aggTypes, pipeline, numForwardedTimes, opts)
	if err != nil {
		panic(fmt.Errorf("unable to create element: %v", err))
	}
	return elem
}

// ResetSetData resets the element and sets data.
func (e *CounterElem) ResetSetData(
	incomingMetricType IncomingMetricType,
	id id.RawID,
	sp policy.StoragePolicy,
	aggTypes maggregation.Types,
	pipeline applied.Pipeline,
	numForwardedTimes int,
) error {
	useDefaultAggregation := aggTypes.IsDefault()
	if useDefaultAggregation {
		aggTypes = e.DefaultAggregationTypes(e.aggTypesOpts)
	}
	if err := e.elemBase.resetSetData(incomingMetricType, id, sp, aggTypes, useDefaultAggregation, pipeline, numForwardedTimes); err != nil {
		return err
	}
	if err := e.counterElemBase.ResetSetData(e.aggTypesOpts, aggTypes, useDefaultAggregation); err != nil {
		return err
	}
	// If the pipeline contains derivative transformations, we need to store past
	// values in order to compute the derivatives.
	if !e.parsedPipeline.HasDerivativeTransform {
		return nil
	}
	numAggTypes := len(e.aggTypes)
	if cap(e.lastConsumedValues) < numAggTypes {
		e.lastConsumedValues = make([]float64, numAggTypes)
	}
	e.lastConsumedValues = e.lastConsumedValues[:numAggTypes]
	for i := 0; i < len(e.lastConsumedValues); i++ {
		e.lastConsumedValues[i] = nan
	}
	return nil
}

// AddUnion adds a metric value union at a given timestamp.
func (e *CounterElem) AddUnion(timestamp time.Time, mu unaggregated.MetricUnion) error {
	alignedStart := timestamp.Truncate(e.sp.Resolution().Window).UnixNano()
	lockedAgg, err := e.findOrCreate(alignedStart, sourcesOptions{})
	if err != nil {
		return err
	}
	lockedAgg.Lock()
	if lockedAgg.closed {
		lockedAgg.Unlock()
		return errAggregationClosed
	}
	lockedAgg.aggregation.AddUnion(mu)
	lockedAgg.Unlock()
	return nil
}

// AddUnique adds a metric value from a given source at a given timestamp.
// If previous values from the same source have already been added to the
// same aggregation, the incoming value is discarded.
func (e *CounterElem) AddUnique(timestamp time.Time, values []float64, sourceID uint32) error {
	alignedStart := timestamp.Truncate(e.sp.Resolution().Window).UnixNano()
	lockedAgg, err := e.findOrCreate(alignedStart, sourcesOptions{updateSources: true, source: sourceID})
	if err != nil {
		return err
	}
	lockedAgg.Lock()
	if lockedAgg.closed {
		lockedAgg.Unlock()
		return errAggregationClosed
	}
	source := uint(sourceID)
	if lockedAgg.seenSources.Test(source) {
		lockedAgg.Unlock()
		return errDuplicateForwardingSource
	}
	lockedAgg.seenSources.Set(source)
	if lockedAgg.sourcesReady {
		// If the sources are ready, the expected sources will be a pre-filled
		// bitset populated with sources the aggregation is expected to see data from.
		// As such, we need to clear the source bit in the expected sources.
		if lockedAgg.expectedSources.Test(source) {
			// This source is never seen before and is in the expected source list,
			// as a result, we need to clear the source bit.
			lockedAgg.expectedSources.Clear(source)
			if lockedAgg.expectedSources.None() {
				lockedAgg.consumeState = readyToConsume
			}
		}
		// New sources that are not in the expected source list are still allowed
		// to go through.
	}
	for _, v := range values {
		lockedAgg.aggregation.Add(v)
	}
	lockedAgg.Unlock()
	return nil
}

// Consume consumes values before a given time and removes them from the element
// after they are consumed, returning whether the element can be collected after
// the consumption is completed.
// NB: Consume is not thread-safe and must be called within a single goroutine
// to avoid race conditions.
func (e *CounterElem) Consume(
	targetNanos int64,
	eagerForwardingMode eagerForwardingMode,
	isEarlierThanFn isEarlierThanFn,
	timestampNanosFn timestampNanosFn,
	flushLocalFn flushLocalMetricFn,
	flushForwardedFn flushForwardedMetricFn,
	onForwardedFlushedFn onForwardingElemFlushedFn,
) bool {
	resolution := e.sp.Resolution().Window
	e.Lock()
	if e.closed {
		e.Unlock()
		return false
	}
	idx := 0
	for range e.values {
		// Bail as soon as the timestamp is no later than the target time.
		timeNanos := timestampNanosFn(e.values[idx].startAtNanos, resolution)
		if !isEarlierThanFn(timeNanos, targetNanos) {
			break
		}
		idx++
	}
	e.toConsume = e.toConsume[:0]
	if idx > 0 {
		// Shift remaining values to the left and shrink the values slice.
		e.toConsume = append(e.toConsume, e.values[:idx]...)
		n := copy(e.values[0:], e.values[idx:])
		// Clear out the invalid items to avoid holding references to objects
		// for reduced GC overhead..
		for i := n; i < len(e.values); i++ {
			e.values[i].Reset()
		}
		e.values = e.values[:n]
	}
	canCollect := len(e.values) == 0 && e.tombstoned

	// Additionally for elements receiving forwarded metrics and sending aggregated metrics
	// to local backends, we also check if any aggregations are ready to be consumed. We however
	// do not remove the aggregations as we do for aggregations whose timestamps are old enough,
	// since for aggregations receiving forwarded metrics that are marked "consume ready", it is
	// possible that metrics still go to the such aggregation bucket after they are marked "consume
	// ready" due to delayed source re-delivery or new sources showing up, and removing such
	// aggregation prematurely would mean the values from the delayed sources and/or new sources
	// would be considered as the aggregated value for such aggregation bucket, which is incorrect.
	// By keeping such aggregation buckets and only removing them when they are considered old enough
	// (i.e., when their timestmaps are earlier than the target timestamp), we ensure no metrics may
	// go to such aggregation buckets after they are consumed and therefore avoid the aformentioned
	// problem.
	aggregationIdxToCloseUntil := len(e.toConsume)
	if e.incomingMetricType == ForwardedIncomingMetric && e.isSourcesSetReadyWithElemLock() {
		e.maybeRefreshSourcesSetWithElemLock()
		// We only attempt to consume if the outgoing metrics type is local instead of forwarded.
		// This is because forwarded metrics are sent in batches and can only be sent when all sources
		// in the same shard have been consumed, and as such is not well suited for pre-emptive consumption.
		if e.outgoingMetricType() == localOutgoingMetric && eagerForwardingMode == allowEagerForwarding {
			for i := 0; i < len(e.values); i++ {
				// NB: This makes the logic easier to understand but it would be more efficient to use
				// an atomic here to avoid locking aggregations.
				e.values[i].lockedAgg.Lock()
				if e.values[i].lockedAgg.consumeState == readyToConsume {
					e.toConsume = append(e.toConsume, e.values[i])
					e.values[i].lockedAgg.consumeState = consuming
				}
				e.values[i].lockedAgg.Unlock()
			}
		}
	}
	e.Unlock()

	// Process the aggregations that are ready for consumption.
	for i := range e.toConsume {
		timeNanos := timestampNanosFn(e.toConsume[i].startAtNanos, resolution)
		e.toConsume[i].lockedAgg.Lock()
		if e.toConsume[i].lockedAgg.consumeState != consumed {
			e.processValueWithAggregationLock(timeNanos, e.toConsume[i].lockedAgg, flushLocalFn, flushForwardedFn)
		}
		e.toConsume[i].lockedAgg.consumeState = consumed
		if i < aggregationIdxToCloseUntil {
			if e.toConsume[i].lockedAgg.seenSources != nil {
				e.sourcesLock.Lock()
				// This is to make sure there aren't too many cached source sets taking up
				// too much space.
				if len(e.cachedSourceSets) < e.opts.MaxNumCachedSourceSets() {
					e.cachedSourceSets = append(e.cachedSourceSets, e.toConsume[i].lockedAgg.seenSources)
					e.toConsume[i].lockedAgg.seenSources = nil
				}
				e.sourcesLock.Unlock()
			}
			e.toConsume[i].lockedAgg.close()
		}
		e.toConsume[i].lockedAgg.Unlock()
		e.toConsume[i].Reset()
	}

	if e.outgoingMetricType() == forwardedOutgoingMetric {
		forwardedAggregationKey, _ := e.ForwardedAggregationKey()
		onForwardedFlushedFn(e.onForwardedAggregationWrittenFn, forwardedAggregationKey)
	}

	return canCollect
}

// Close closes the element.
func (e *CounterElem) Close() {
	e.Lock()
	if e.closed {
		e.Unlock()
		return
	}
	e.closed = true
	e.id = nil
	e.parsedPipeline = parsedPipeline{}
	e.writeForwardedMetricFn = nil
	e.onForwardedAggregationWrittenFn = nil
	e.sourcesHeartbeat = nil
	e.sourcesSet = nil
	for idx := range e.cachedSourceSets {
		e.cachedSourceSets[idx] = nil
	}
	e.cachedSourceSets = nil
	for idx := range e.values {
		// Close the underlying aggregation objects.
		e.values[idx].lockedAgg.close()
		e.values[idx].Reset()
	}
	e.values = e.values[:0]
	e.toConsume = e.toConsume[:0]
	e.lastConsumedValues = e.lastConsumedValues[:0]
	e.counterElemBase.Close()
	aggTypesPool := e.aggTypesOpts.TypesPool()
	pool := e.ElemPool(e.opts)
	e.Unlock()

	if !e.useDefaultAggregation {
		aggTypesPool.Put(e.aggTypes)
	}
	pool.Put(e)
}

// findOrCreate finds the aggregation for a given time, or creates one
// if it doesn't exist.
func (e *CounterElem) findOrCreate(
	alignedStart int64,
	sourcesOpts sourcesOptions,
) (*lockedCounterAggregation, error) {
	e.RLock()
	if e.closed {
		e.RUnlock()
		return nil, errElemClosed
	}
	if sourcesOpts.updateSources {
		e.updateSources(sourcesOpts.source)
	}
	idx, found := e.indexOfWithLock(alignedStart)
	if found {
		agg := e.values[idx].lockedAgg
		e.RUnlock()
		return agg, nil
	}
	e.RUnlock()

	e.Lock()
	if e.closed {
		e.Unlock()
		return nil, errElemClosed
	}
	idx, found = e.indexOfWithLock(alignedStart)
	if found {
		agg := e.values[idx].lockedAgg
		e.Unlock()
		return agg, nil
	}

	// If not found, create a new aggregation.
	numValues := len(e.values)
	e.values = append(e.values, timedCounter{})
	copy(e.values[idx+1:numValues+1], e.values[idx:numValues])

	var (
		sourcesReady    = e.isSourcesSetReadyWithElemLock()
		expectedSources *bitset.BitSet
		seenSources     *bitset.BitSet
	)
	if sourcesOpts.updateSources {
		e.sourcesLock.Lock()
		// If the sources set is ready, we clone it ane use the clone to
		// determine when we have received from all the expected sources.
		if sourcesReady {
			expectedSources = e.sourcesSet.Clone()
		}
		if numCachedSourceSets := len(e.cachedSourceSets); numCachedSourceSets > 0 {
			seenSources = e.cachedSourceSets[numCachedSourceSets-1]
			e.cachedSourceSets[numCachedSourceSets-1] = nil
			e.cachedSourceSets = e.cachedSourceSets[:numCachedSourceSets-1]
			seenSources.ClearAll()
		} else {
			seenSources = bitset.New(defaultNumSources)
		}
		e.sourcesLock.Unlock()
	}

	e.values[idx] = timedCounter{
		startAtNanos: alignedStart,
		lockedAgg: &lockedCounterAggregation{
			sourcesReady:    sourcesReady,
			expectedSources: expectedSources,
			seenSources:     seenSources,
			aggregation:     e.NewAggregation(e.opts, e.aggOpts),
		},
	}
	agg := e.values[idx].lockedAgg
	e.Unlock()
	return agg, nil
}

// indexOfWithLock finds the smallest element index whose timestamp
// is no smaller than the start time passed in, and true if it's an
// exact match, false otherwise.
func (e *CounterElem) indexOfWithLock(alignedStart int64) (int, bool) {
	numValues := len(e.values)
	// Optimize for the common case.
	if numValues > 0 && e.values[numValues-1].startAtNanos == alignedStart {
		return numValues - 1, true
	}
	// Binary search for the unusual case. We intentionally do not
	// use the sort.Search() function because it requires passing
	// in a closure.
	left, right := 0, numValues
	for left < right {
		mid := left + (right-left)/2 // avoid overflow
		if e.values[mid].startAtNanos < alignedStart {
			left = mid + 1
		} else {
			right = mid
		}
	}
	// If the current timestamp is equal to or larger than the target time,
	// return the index as is.
	if left < numValues && e.values[left].startAtNanos == alignedStart {
		return left, true
	}
	return left, false
}

func (e *CounterElem) processValueWithAggregationLock(
	timeNanos int64,
	lockedAgg *lockedCounterAggregation,
	flushLocalFn flushLocalMetricFn,
	flushForwardedFn flushForwardedMetricFn,
) {
	if lockedAgg.aggregation.Count() == 0 {
		return
	}
	var (
		fullPrefix       = e.FullPrefix(e.opts)
		transformations  = e.parsedPipeline.Transformations
		discardNaNValues = e.opts.DiscardNaNAggregatedValues()
	)
	for aggTypeIdx, aggType := range e.aggTypes {
		value := lockedAgg.aggregation.ValueOf(aggType)
		for i := 0; i < transformations.Len(); i++ {
			transformType := transformations.At(i).Transformation.Type
			if transformType.IsUnaryTransform() {
				fn := transformType.MustUnaryTransform()
				res := fn(transformation.Datapoint{TimeNanos: timeNanos, Value: value})
				value = res.Value
			} else {
				fn := transformType.MustBinaryTransform()
				prev := transformation.Datapoint{TimeNanos: e.lastConsumedAtNanos, Value: e.lastConsumedValues[aggTypeIdx]}
				curr := transformation.Datapoint{TimeNanos: timeNanos, Value: value}
				res := fn(prev, curr)
				// NB: we only need to record the value needed for derivative transformations.
				// We currently only support first-order derivative transformations so we only
				// need to keep one value. In the future if we need to support higher-order
				// derivative transformations, we need to store an array of values here.
				e.lastConsumedValues[aggTypeIdx] = value
				value = res.Value
			}
		}
		if discardNaNValues && math.IsNaN(value) {
			continue
		}
		if e.outgoingMetricType() == localOutgoingMetric {
			flushLocalFn(fullPrefix, e.id, e.TypeStringFor(e.aggTypesOpts, aggType), timeNanos, value, e.sp)
		} else {
			forwardedAggregationKey, _ := e.ForwardedAggregationKey()
			flushForwardedFn(e.writeForwardedMetricFn, forwardedAggregationKey, timeNanos, value)
		}
	}
	e.lastConsumedAtNanos = timeNanos

	// Emit latency metrics for forwarded metrics.
	if e.outgoingMetricType() == localOutgoingMetric && e.incomingMetricType == ForwardedIncomingMetric {
		e.opts.FullForwardingLatencyHistograms().RecordDuration(
			e.sp.Resolution().Window,
			e.numForwardedTimes,
			time.Duration(e.nowFn().UnixNano()-timeNanos),
		)
	}
}

func (e *CounterElem) outgoingMetricType() outgoingMetricType {
	if !e.parsedPipeline.HasRollup {
		return localOutgoingMetric
	}
	return forwardedOutgoingMetric
}

func (e *CounterElem) isSourcesSetReadyWithElemLock() bool {
	if !e.opts.EnableEagerForwarding() {
		return false
	}
	if e.buildingSourcesAtNanos == 0 {
		return false
	}
	// NB: Allow TTL for the source set to build up.
	return e.nowFn().UnixNano() >= e.buildingSourcesAtNanos+e.sourcesTTLNanos
}

func (e *CounterElem) maybeRefreshSourcesSetWithElemLock() {
	if !e.opts.EnableEagerForwarding() {
		return
	}
	nowNanos := e.nowFn().UnixNano()
	if nowNanos-e.lastSourcesRefreshNanos < e.sourcesTTLNanos {
		return
	}
	e.sourcesLock.Lock()
	for sourceID, lastHeartbeatNanos := range e.sourcesHeartbeat {
		if nowNanos-lastHeartbeatNanos >= e.sourcesTTLNanos {
			delete(e.sourcesHeartbeat, sourceID)
			e.sourcesSet.Clear(uint(sourceID))
		}
	}
	e.lastSourcesRefreshNanos = nowNanos
	e.sourcesLock.Unlock()
}

func (e *CounterElem) updateSources(source uint32) {
	if !e.opts.EnableEagerForwarding() {
		return
	}
	nowNanos := e.nowFn().UnixNano()
	e.sourcesLock.Lock()
	// First time a source is received.
	if e.sourcesHeartbeat == nil {
		e.sourcesHeartbeat = make(map[uint32]int64, defaultNumSources)
		e.sourcesSet = bitset.New(defaultNumSources)
		e.buildingSourcesAtNanos = nowNanos
		e.lastSourcesRefreshNanos = nowNanos
	}
	if v, exists := e.sourcesHeartbeat[source]; !exists || v < nowNanos {
		e.sourcesHeartbeat[source] = nowNanos
	}
	e.sourcesSet.Set(uint(source))
	e.sourcesLock.Unlock()
}
