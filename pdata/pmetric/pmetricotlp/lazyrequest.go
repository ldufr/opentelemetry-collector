// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package pmetricotlp // import "go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"

import (
	"slices"

	"go.opentelemetry.io/collector/pdata/internal"
	"go.opentelemetry.io/collector/pdata/internal/json"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// LazyExportRequest represents the request for gRPC/HTTP client/server.
// It's a wrapper for pmetric.Metrics data.
type LazyExportRequest struct {
	orig  *internal.LazyExportMetricsServiceRequest
	state *internal.State
}

// NewLazyExportRequest returns an empty LazyExportRequest.
func NewLazyExportRequest() LazyExportRequest {
	return LazyExportRequest{
		orig:  &internal.LazyExportMetricsServiceRequest{},
		state: internal.NewState(),
	}
}

// NewLazyExportRequestFromMetrics returns a LazyExportRequest from pmetric.Metrics.
// Because LazyExportRequest is a wrapper for pmetric.Metrics,
// any changes to the provided Metrics struct will be reflected in the LazyExportRequest and vice versa.
func NewLazyExportRequestFromMetrics(md pmetric.LazyMetrics) LazyExportRequest {
	return LazyExportRequest{
		orig:  internal.GetLazyMetricsOrig(internal.LazyMetricsWrapper(md)),
		state: internal.GetLazyMetricsState(internal.LazyMetricsWrapper(md)),
	}
}

// MarshalProto marshals LazyExportRequest into proto bytes.
func (ms LazyExportRequest) MarshalProto() ([]byte, error) {
	size := ms.orig.SizeProto()
	buf := make([]byte, size)
	_ = ms.orig.MarshalProto(buf)
	return buf, nil
}

// UnmarshalProto unmarshalls LazyExportRequest from proto bytes.
func (ms LazyExportRequest) UnmarshalProto(data []byte) error {
	err := ms.orig.UnmarshalProto(data)
	if err != nil {
		return err
	}
	return nil
}

// MarshalJSON marshals LazyExportRequest into JSON bytes.
func (ms LazyExportRequest) MarshalJSON() ([]byte, error) {
	dest := json.BorrowStream(nil)
	defer json.ReturnStream(dest)
	ms.orig.MarshalJSON(dest)
	if dest.Error() != nil {
		return nil, dest.Error()
	}
	return slices.Clone(dest.Buffer()), nil
}

// UnmarshalJSON unmarshalls LazyExportRequest from JSON bytes.
func (ms LazyExportRequest) UnmarshalJSON(data []byte) error {
	iter := json.BorrowIterator(data)
	defer json.ReturnIterator(iter)
	ms.orig.UnmarshalJSON(iter)
	return iter.Error()
}

func (ms LazyExportRequest) Metrics() pmetric.LazyMetrics {
	return pmetric.LazyMetrics(internal.NewLazyMetricsWrapper(ms.orig, ms.state))
}
