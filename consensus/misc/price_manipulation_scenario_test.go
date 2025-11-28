package misc

import (
    "math/big"
    "testing"

    "github.com/kaiachain/kaia/blockchain/types"
    "github.com/kaiachain/kaia/common"
    "github.com/kaiachain/kaia/params"
)

// Scenario: attack followed by normal blocks to show delayed drop
func TestPriceManipulation_Scenario_AttackThenNormal(t *testing.T) {
    cfg := params.GetDefaultKIP71Config()
    cfg.GasTarget = 30000000
    cfg.BaseFeeDenominator = 20

    parentBase := new(big.Int).SetUint64(500000000000)

    // Attack: 20 blocks at target - 10
    attackGas := cfg.GasTarget - 10
    prev := parentBase
    // apply first attack block
    hdr := &types.Header{ Number: common.Big3, GasUsed: attackGas, BaseFee: prev }
    prev = NextMagmaBlockBaseFee(hdr, cfg)
    for i := 0; i < 19; i++ {
        hdr := &types.Header{ Number: common.Big3, GasUsed: attackGas, BaseFee: prev }
        prev = NextMagmaBlockBaseFee(hdr, cfg)
    }

    // Now normal blocks with 50% utilization
    normalGas := cfg.GasTarget / 2
    for i := 0; i < 5; i++ {
        hdr := &types.Header{ Number: common.Big3, GasUsed: normalGas, BaseFee: prev }
        prev = NextMagmaBlockBaseFee(hdr, cfg)
    }

    // Ensure baseFee decreased but not to expected 'no-attack' baseline
    if prev.Cmp(parentBase) >= 0 {
        t.Fatalf("expected eventual decrease after normal blocks, got %s >= %s", prev.String(), parentBase.String())
    }
}

