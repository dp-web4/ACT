#!/bin/bash
# Pending Consensus System for Society 4
# Manages decisions made while network-isolated from federation

set -e

# Configuration
PENDING_DIR="$HOME/.society4chain/pending_consensus"
PENDING_FILE="$PENDING_DIR/pending.json"
PROCESSED_FILE="$PENDING_DIR/processed.json"
NETWORK_STATE_FILE="$PENDING_DIR/network_state"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Initialize pending consensus directory
init_pending_consensus() {
    mkdir -p "$PENDING_DIR"

    if [ ! -f "$PENDING_FILE" ]; then
        echo '{"pending_decisions": [], "metadata": {}}' > "$PENDING_FILE"
    fi

    if [ ! -f "$PROCESSED_FILE" ]; then
        echo '{"processed": []}' > "$PROCESSED_FILE"
    fi

    echo -e "${GREEN}✓${NC} Pending consensus system initialized"
}

# Detect current network state
detect_network() {
    local current_ip=$(ip addr show eth0 2>/dev/null | grep "inet " | awk '{print $2}' | cut -d/ -f1)
    local network_type="unknown"

    if [[ "$current_ip" =~ ^10\.0\.0\. ]]; then
        network_type="home_federation"
        echo -e "${GREEN}✓${NC} Connected to HOME network (Federation accessible)"
    elif [[ "$current_ip" =~ ^172\.25\. ]]; then
        network_type="work_isolated"
        echo -e "${YELLOW}⚠${NC} Connected to WORK network (Federation isolated)"
    else
        network_type="unknown"
        echo -e "${RED}✗${NC} Unknown network: $current_ip"
    fi

    echo "$network_type|$current_ip|$(date -Iseconds)" > "$NETWORK_STATE_FILE"
    echo "$network_type"
}

# Add a pending decision
add_pending_decision() {
    local decision_type="$1"
    local decision_data="$2"
    local reason="$3"

    local network_state=$(detect_network)
    local timestamp=$(date -Iseconds)
    local hardware_hash=$(bash /mnt/c/projects/ai-agents/ACT/implementation/society4/blockchain/source/extract_hardware.sh 2>/dev/null | grep "Hardware Hash" | awk '{print $3}')

    # Create decision object
    local decision=$(jq -n \
        --arg type "$decision_type" \
        --arg data "$decision_data" \
        --arg reason "$reason" \
        --arg network "$network_state" \
        --arg time "$timestamp" \
        --arg hw "$hardware_hash" \
        '{
            id: ($time | split(".")[0] | split("T")[1] | gsub("[:-]"; "")),
            type: $type,
            data: $data,
            reason: $reason,
            network_state: $network,
            timestamp: $time,
            hardware_attestation: $hw,
            status: "pending"
        }')

    # Add to pending file
    jq ".pending_decisions += [$decision]" "$PENDING_FILE" > "$PENDING_FILE.tmp" && mv "$PENDING_FILE.tmp" "$PENDING_FILE"

    echo -e "${GREEN}✓${NC} Pending decision added: $decision_type"
    echo "$decision" | jq -c
}

# List pending decisions
list_pending() {
    local count=$(jq '.pending_decisions | length' "$PENDING_FILE")

    if [ "$count" -eq 0 ]; then
        echo -e "${YELLOW}No pending decisions${NC}"
        return
    fi

    echo -e "${GREEN}Pending Decisions ($count):${NC}"
    jq -r '.pending_decisions[] | "\(.timestamp | split("T")[0]) | \(.type) | \(.reason) | Status: \(.status)"' "$PENDING_FILE"
}

# Process pending decisions (when reconnected to federation)
process_pending() {
    local network_state=$(detect_network)

    if [ "$network_state" != "home_federation" ]; then
        echo -e "${RED}✗${NC} Cannot process pending decisions - not connected to federation"
        echo -e "Current network: $network_state"
        return 1
    fi

    local pending_count=$(jq '.pending_decisions | map(select(.status == "pending")) | length' "$PENDING_FILE")

    if [ "$pending_count" -eq 0 ]; then
        echo -e "${YELLOW}No pending decisions to process${NC}"
        return 0
    fi

    echo -e "${GREEN}Processing $pending_count pending decisions...${NC}"

    # Process each pending decision
    jq -c '.pending_decisions[] | select(.status == "pending")' "$PENDING_FILE" | while read -r decision; do
        local id=$(echo "$decision" | jq -r '.id')
        local type=$(echo "$decision" | jq -r '.type')
        local data=$(echo "$decision" | jq -r '.data')

        echo -e "Processing decision $id (type: $type)..."

        # Here you would submit to the actual federation
        # For now, we'll simulate processing
        sleep 1

        # Mark as processed
        jq "(.pending_decisions[] | select(.id == \"$id\")) |= . + {status: \"processed\", processed_at: \"$(date -Iseconds)\"}" \
            "$PENDING_FILE" > "$PENDING_FILE.tmp" && mv "$PENDING_FILE.tmp" "$PENDING_FILE"

        # Add to processed log
        jq ".processed += [$decision + {processed_at: \"$(date -Iseconds)\"}]" \
            "$PROCESSED_FILE" > "$PROCESSED_FILE.tmp" && mv "$PROCESSED_FILE.tmp" "$PROCESSED_FILE"

        echo -e "${GREEN}✓${NC} Decision $id processed"
    done
}

# Export pending decisions for git commit
export_for_git() {
    local export_file="/mnt/c/projects/ai-agents/ACT/implementation/society4/public/pending_consensus_$(date +%Y%m%d_%H%M%S).json"

    # Create export with metadata
    jq --arg network "$(cat $NETWORK_STATE_FILE 2>/dev/null || echo 'unknown')" \
       --arg hw "$(bash /mnt/c/projects/ai-agents/ACT/implementation/society4/blockchain/source/extract_hardware.sh 2>/dev/null | grep 'Hardware Hash' | awk '{print $3}')" \
       '. + {
            export_metadata: {
                timestamp: now | todate,
                network_state: $network,
                hardware_attestation: $hw,
                pending_count: (.pending_decisions | map(select(.status == "pending")) | length),
                processed_count: (.pending_decisions | map(select(.status == "processed")) | length)
            }
        }' "$PENDING_FILE" > "$export_file"

    echo -e "${GREEN}✓${NC} Exported pending consensus to:"
    echo "$export_file"
    echo -e "\nYou can now commit this to git for federation visibility"
}

# Show current status
show_status() {
    echo -e "\n${GREEN}=== Society 4 Pending Consensus Status ===${NC}"

    # Network state
    echo -e "\n${YELLOW}Network State:${NC}"
    detect_network

    # Hardware binding
    echo -e "\n${YELLOW}Hardware Binding:${NC}"
    bash /mnt/c/projects/ai-agents/ACT/implementation/society4/blockchain/source/extract_hardware.sh 2>/dev/null | grep "Hardware Hash" || echo "Unable to extract"

    # Pending decisions
    echo -e "\n${YELLOW}Pending Decisions:${NC}"
    local pending=$(jq '.pending_decisions | map(select(.status == "pending")) | length' "$PENDING_FILE")
    local processed=$(jq '.pending_decisions | map(select(.status == "processed")) | length' "$PENDING_FILE")
    echo "  Pending: $pending"
    echo "  Processed: $processed"

    # Last network transition
    if [ -f "$NETWORK_STATE_FILE" ]; then
        echo -e "\n${YELLOW}Last Network State:${NC}"
        cat "$NETWORK_STATE_FILE"
    fi
}

# Main command handler
case "${1:-help}" in
    init)
        init_pending_consensus
        ;;
    add)
        if [ $# -lt 4 ]; then
            echo "Usage: $0 add <type> <data> <reason>"
            exit 1
        fi
        add_pending_decision "$2" "$3" "$4"
        ;;
    list)
        list_pending
        ;;
    process)
        process_pending
        ;;
    export)
        export_for_git
        ;;
    status)
        show_status
        ;;
    network)
        detect_network
        ;;
    *)
        echo "Society 4 Pending Consensus System"
        echo ""
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  init     - Initialize pending consensus system"
        echo "  add      - Add a pending decision: add <type> <data> <reason>"
        echo "  list     - List all pending decisions"
        echo "  process  - Process pending decisions (requires federation network)"
        echo "  export   - Export pending decisions for git commit"
        echo "  status   - Show current system status"
        echo "  network  - Detect current network state"
        echo ""
        echo "Examples:"
        echo "  $0 add vote '{\"proposal\":\"123\",\"vote\":\"yes\"}' 'Support Web4 enhancement'"
        echo "  $0 add transaction '{\"to\":\"society1\",\"amount\":10}' 'ATP delegation'"
        echo "  $0 process  # When back on home network"
        ;;
esac