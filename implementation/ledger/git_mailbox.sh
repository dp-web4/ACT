#!/bin/bash
# Git Mailbox - Near-realtime federation communication via git
# Monitors and pulls every 10 seconds for fast message exchange

REPO_DIR="/home/dp/ai-workspace/act"
LOG_FILE="$REPO_DIR/implementation/ledger/git_mailbox.log"
INBOX_DIR="$REPO_DIR/implementation/ledger/federation_inbox"
OUTBOX_DIR="$REPO_DIR/implementation/ledger/federation_outbox"
SOCIETY_NAME="genesis"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Git Mailbox Service Starting${NC}"
echo "📬 Society: $SOCIETY_NAME"
echo "🔄 Sync interval: 10 seconds"
echo "📁 Repository: $REPO_DIR"
echo ""

# Create inbox/outbox directories if they don't exist
mkdir -p "$INBOX_DIR"
mkdir -p "$OUTBOX_DIR"

# Initialize log
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Git Mailbox service started" > "$LOG_FILE"

# Track last seen commits to detect new messages
LAST_COMMIT=""

# Function to check for new messages
check_new_messages() {
    local current_commit=$(git rev-parse HEAD 2>/dev/null)
    
    if [ "$current_commit" != "$LAST_COMMIT" ] && [ -n "$LAST_COMMIT" ]; then
        echo -e "${YELLOW}📨 New messages detected!${NC}"
        
        # Check for new files in inbox
        local new_files=$(git diff --name-only "$LAST_COMMIT" HEAD | grep "federation_inbox/" | grep -v "$SOCIETY_NAME")
        
        if [ -n "$new_files" ]; then
            echo -e "${BLUE}📥 Incoming messages:${NC}"
            echo "$new_files" | while read file; do
                if [ -f "$REPO_DIR/$file" ]; then
                    local sender=$(echo "$file" | sed 's/.*inbox\/\([^_]*\).*/\1/')
                    echo -e "  ${GREEN}✉️${NC} From $sender: $(basename $file)"
                    
                    # Log the message
                    echo "[$(date '+%Y-%m-%d %H:%M:%S')] Message from $sender: $file" >> "$LOG_FILE"
                fi
            done
        fi
        
        # Check for discussion updates
        local discussion_updates=$(git diff --name-only "$LAST_COMMIT" HEAD | grep -E "(DISCUSSION|POSITION|RESPONSE)" | grep -v TEMPLATE)
        if [ -n "$discussion_updates" ]; then
            echo -e "${BLUE}💬 Discussion updates:${NC}"
            echo "$discussion_updates" | while read file; do
                echo -e "  ${GREEN}📋${NC} $(basename $file)"
            done
        fi
    fi
    
    LAST_COMMIT="$current_commit"
}

# Function to send outgoing messages
send_messages() {
    local outgoing=$(ls "$OUTBOX_DIR" 2>/dev/null | grep -v ".sent")
    
    if [ -n "$outgoing" ]; then
        echo -e "${YELLOW}📤 Sending messages...${NC}"
        
        for file in $outgoing; do
            if [ -f "$OUTBOX_DIR/$file" ]; then
                # Move to inbox for other societies to see
                cp "$OUTBOX_DIR/$file" "$INBOX_DIR/${SOCIETY_NAME}_$(basename $file)"
                mv "$OUTBOX_DIR/$file" "$OUTBOX_DIR/$file.sent"
                
                echo -e "  ${GREEN}✓${NC} Sent: $file"
                echo "[$(date '+%Y-%m-%d %H:%M:%S')] Sent message: $file" >> "$LOG_FILE"
            fi
        done
        
        # Commit and push
        git add "$INBOX_DIR"
        git commit -m "[$SOCIETY_NAME] Outgoing messages - $(date '+%H:%M:%S')" 2>/dev/null
        git push origin main 2>/dev/null
    fi
}

# Main monitoring loop
echo -e "${GREEN}✅ Monitoring active - Press Ctrl+C to stop${NC}"
echo "════════════════════════════════════════════════════"

# Get initial commit
cd "$REPO_DIR"
git pull origin main --quiet 2>/dev/null
LAST_COMMIT=$(git rev-parse HEAD 2>/dev/null)

while true; do
    # Silent pull to check for updates
    git pull origin main --quiet 2>/dev/null
    
    # Check for new messages
    check_new_messages
    
    # Send any pending messages
    send_messages
    
    # Show status indicator (subtle heartbeat)
    printf "💓"
    
    # Sleep for 10 seconds
    sleep 10
    
    # Clear the heartbeat
    printf "\b "
    printf "\b"
done