package misc

import (
    "math/big"
    "testing"

    "github.com/kaiachain/kaia/blockchain/types"
    "github.com/kaiachain/kaia/common"
    "github.com/kaiachain/kaia/params"
)

// PoC: show 20-block attack sequence causing stickiness as described
func TestPriceManipulation_PoC_Stickiness(t *testing.T) {
    cfg := params.GetDefaultKIP71Config()
    cfg.GasTarget = 30000000
    cfg.BaseFeeDenominator = 20

    // Start baseFee 500 Gwei
    baseFee := new(big.Int).SetUint64(500000000000)
    parent := &types.Header{ Number: common.Big3, GasUsed: cfg.GasTarget, BaseFee: baseFee }

    // Attack blocks (20) with GasUsed = target - 10
    attackGas := cfg.GasTarget - 10
    seq := make([]*big.Int, 0, 30)
    seq = append(seq, parent.BaseFee)

    // first attack block
    parent.GasUsed = attackGas
    next := NextMagmaBlockBaseFee(parent, cfg)
    seq = append(seq, next)

    prev := next
    for i := 0; i < 19; i++ {
        hdr := &types.Header{ Number: common.Big3, GasUsed: attackGas, BaseFee: prev }
        nxt := NextMagmaBlockBaseFee(hdr, cfg)
        seq = append(seq, nxt)
        prev = nxt
    }

    // Evaluate cumulative drop after 20 attack blocks
    final := prev
    initial := baseFee
    diff := new(big.Int).Sub(initial, final)
    // Expect the cumulative drop to be relatively small compared to initial
    // (this demonstrates very slow decrease under small deltas)
    threshold := new(big.Int).Div(new(big.Int).Mul(initial, big.NewInt(1)), big.NewInt(1000)) // 0.1%
    if diff.Cmp(threshold) > 0 {
        t.Fatalf("expected cumulative drop after 20 attack blocks to be small (<0.1%%), got diff %s > threshold %s; sequence: %v", diff.String(), threshold.String(), seq)
    }
}

