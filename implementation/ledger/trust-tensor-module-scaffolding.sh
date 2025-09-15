#!/bin/bash
# Web4 Race Car Demo - Trust Tensor Module Only
# Safe incremental approach: Get Trust Tensor working first
# DEPENDS ONLY ON: bank, lctmanager (both proven to work)

set -e

echo "🧠 Adding Trust Tensor module to Web4 system..."
echo "📋 Dependencies: bank (✅), lctmanager (✅)"

# Navigate to your racecar-web directory
# cd racecar-web

# Trust Tensor Module (DEPENDS ONLY ON LCTMANAGER)
echo "🧠 Creating Trust Tensor module (T3/V3 attached to relationships)..."
ignite scaffold module trusttensor --dep bank,lctmanager

# Trust tensor structures attached to LCT relationships
echo "📊 Adding Trust Tensor data structures..."
ignite scaffold type relationship-trust-tensor tensor_id:string lct_id:string tensor_type:string talent_score:string training_score:string temperament_score:string context:string created_at:int updated_at:int version:int --module trusttensor

ignite scaffold type value-tensor tensor_id:string lct_id:string operation_id:string valuation_score:string veracity_score:string validity_score:string created_at:int --module trusttensor

ignite scaffold type tensor-entry entry_id:string tensor_id:string dimension:string entry_type:string score_value:string witness_chain:string witness_confidence:string timestamp:int --module trusttensor

# Trust tensor operations on relationships  
echo "⚙️ Adding Trust Tensor operations..."
ignite scaffold message create-relationship-tensor lct_id:string tensor_type:string context:string --module trusttensor --response tensor_id:string

ignite scaffold message update-tensor-score tensor_id:string dimension:string value:string context:string witness_data:string --module trusttensor

ignite scaffold message add-tensor-witness tensor_id:string dimension:string witness_lct:string confidence:string evidence_hash:string --module trusttensor

# Trust tensor queries
echo "🔍 Adding Trust Tensor queries..."
ignite scaffold query get-relationship-tensor lct_id:string tensor_type:string --module trusttensor --response relationship_trust_tensor

ignite scaffold query calculate-relationship-trust lct_id:string context:string --module trusttensor --response trust_score:string,factors:string

ignite scaffold query get-tensor-history tensor_id:string --module trusttensor --response tensor_entries

# Build and test
echo "🔨 Building Web4 system with Trust Tensor module..."
ignite chain build

echo "🎉 Trust Tensor module scaffolded successfully!"
echo ""
echo "⚠️ NEXT STEPS:"
echo "1. Apply v29 Magic Formula fixes:"
echo "   - Copy x/trusttensor/types/keys.go"
echo "   - Copy x/trusttensor/types/expected_keepers.go"
echo "   - Copy x/trusttensor/keeper/genesis.go"
echo "   - Copy x/trusttensor/keeper/keeper.go (complete implementation)"
echo ""
echo "2. Test the 5-module system:"
echo "   ignite chain build"
echo "   ignite chain serve"
echo ""
echo "3. Confirm Trust Tensor works before adding Energy Cycle"
echo ""
echo "📊 Current System Status:"
echo "   ✅ componentregistry"  
echo "   ✅ pairingqueue"
echo "   ✅ lctmanager"
echo "   ✅ pairing"
echo "   🔧 trusttensor (needs v29 fixes)"
echo ""
echo "🎯 Trust Tensor provides T3/V3 scoring for LCT relationships"
