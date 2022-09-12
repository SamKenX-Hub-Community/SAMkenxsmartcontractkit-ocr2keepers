package keepers

import (
	"context"
	"fmt"
	"sort"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	ktypes "github.com/smartcontractkit/ocr2keepers/pkg/types"
)

// Query implements the types.ReportingPlugin interface in OCR2. The query produced from this
// method is intended to be empty.
func (k *keepers) Query(_ context.Context, _ types.ReportTimestamp) (types.Query, error) {
	return types.Query{}, nil
}

// Observation implements the types.ReportingPlugin interface in OCR2. This method samples a set
// of upkeeps available in and UpkeepService and produces an observation containing upkeeps that
// need to be executed.
func (k *keepers) Observation(ctx context.Context, _ types.ReportTimestamp, _ types.Query) (types.Observation, error) {
	results, err := k.service.SampleUpkeeps(ctx)
	if err != nil {
		return types.Observation{}, err
	}

	keys := keyList(filterUpkeeps(results, Perform))

	b, err := Encode(keys)
	if err != nil {
		return types.Observation{}, err
	}

	return types.Observation(b), err
}

// Report implements the types.ReportingPlugin inteface in OC2. This method chooses a single upkeep
// from the provided observations by the earliest block number, checks the upkeep, and builds a
// report. Multiple upkeeps in a single report is supported by how the data is abi encoded, but
// no gas estimations exist yet.
func (k *keepers) Report(ctx context.Context, _ types.ReportTimestamp, _ types.Query, attributed []types.AttributedObservation) (bool, types.Report, error) {
	var err error

	// collect all observations
	sets := make([][]ktypes.UpkeepKey, len(attributed))
	for i, a := range attributed {
		var values []ktypes.UpkeepKey
		err = Decode([]byte(a.Observation), &values)
		if err != nil {
			// TODO: handle errors better; this currently results in a hard failure on bad encoding
			return false, nil, err
		}
		sets[i] = values
	}

	// dedupe, flatten, and sort
	allKeys, err := dedupe(sets)
	if err != nil {
		return false, nil, fmt.Errorf("%w: observation dedupe", err)
	}
	sort.Sort(sortUpkeepKeys(allKeys))

	// select, verify, and build report
	toPerform := []ktypes.UpkeepResult{}
	for _, key := range allKeys {
		upkeep, err := k.service.CheckUpkeep(ctx, key)
		if err != nil {
			return false, nil, fmt.Errorf("%w: check upkeep failure in report", err)
		}

		if upkeep.State == Perform {
			// only build a report from a single upkeep for now
			toPerform = append(toPerform, upkeep)
			break
		}
	}

	// if nothing to report, return false with no error
	if len(toPerform) == 0 {
		return false, nil, nil
	}

	b, err := k.encoder.EncodeReport(toPerform)
	if err != nil {
		// TODO: handle errors better
		return false, nil, fmt.Errorf("%w: report encoding", err)
	}

	return true, types.Report(b), err
}

// ShouldAcceptFinalizedReport implements the types.ReportingPlugin interface in OCR2. The implementation
// is the most basic possible at this point and assumes all reports should be accepted.
func (k *keepers) ShouldAcceptFinalizedReport(_ context.Context, _ types.ReportTimestamp, _ types.Report) (bool, error) {
	// TODO: add some checks to verify that a report should be accepted to transmit
	return true, nil
}

// ShouldTransmitAcceptedReport implements the types.ReportingPlugin interface in OCR2. The implementation
// is the most basic possible at this point and assumes all reports should be accepted.
func (k *keepers) ShouldTransmitAcceptedReport(c context.Context, t types.ReportTimestamp, r types.Report) (bool, error) {
	// TODO: add some checks to verify that a report should be accepted to transmit
	return true, nil
}

// Close implements the types.ReportingPlugin interface in OCR2.
func (k *keepers) Close() error {
	return nil
}