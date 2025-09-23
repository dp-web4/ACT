#!/bin/bash
# Federation Discussion Monitor
# Checks for society responses every 5 minutes

echo "🔍 Starting Federation Discussion Monitor"
echo "📊 Checking for responses from Sprout and Society4..."
echo ""

while true; do
    cd /home/dp/ai-workspace/act
    git pull origin main 2>/dev/null
    
    echo "════════════════════════════════════════════════════"
    echo "📅 Discussion Status Check - $(date '+%Y-%m-%d %H:%M:%S')"
    echo "════════════════════════════════════════════════════"
    
    # Check for position documents
    echo "📝 Position Statements:"
    if ls implementation/ledger/*POSITION*.md 2>/dev/null; then
        for file in implementation/ledger/*POSITION*.md; do
            echo "  ✓ Found: $(basename $file)"
        done
    else
        echo "  ⏳ No formal positions submitted yet"
    fi
    
    echo ""
    
    # Check for response documents  
    echo "💬 Society Responses:"
    if ls implementation/ledger/*RESPONSE*.md 2>/dev/null | grep -v TEMPLATE; then
        for file in implementation/ledger/*RESPONSE*.md; do
            if [[ ! "$file" =~ TEMPLATE ]]; then
                echo "  ✓ Found: $(basename $file)"
            fi
        done
    else
        echo "  ⏳ No responses received yet"
    fi
    
    echo ""
    
    # Check for questions
    echo "❓ Pending Questions:"
    if [ -f "implementation/ledger/QUESTIONS_PENDING.md" ]; then
        echo "  🔴 Questions need answers!"
        head -3 implementation/ledger/QUESTIONS_PENDING.md 2>/dev/null
    else
        echo "  ✅ No pending questions"
    fi
    
    echo ""
    
    # Check discussion log
    echo "📊 Discussion Activity:"
    if [ -f "implementation/ledger/FEDERATION_DISCUSSION_LOG.md" ]; then
        lines=$(wc -l < implementation/ledger/FEDERATION_DISCUSSION_LOG.md)
        echo "  📈 Discussion log has $lines lines"
        echo "  Last entry:"
        tail -1 implementation/ledger/FEDERATION_DISCUSSION_LOG.md
    else
        echo "  📋 No discussion log started yet"
    fi
    
    echo ""
    
    # Calculate time remaining
    start_time=$(date -d "2025-09-22 21:00:00" +%s 2>/dev/null || date +%s)
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))
    phase1_remaining=$((259200 - elapsed)) # 72 hours in seconds
    
    if [ $phase1_remaining -gt 0 ]; then
        hours=$((phase1_remaining / 3600))
        echo "⏰ Time Remaining in Education Phase: $hours hours"
    else
        echo "⏰ Education Phase Complete - Position Development Active"
    fi
    
    echo ""
    echo "💤 Next check in 5 minutes... (Press Ctrl+C to stop)"
    echo "════════════════════════════════════════════════════"
    echo ""
    
    sleep 300
done