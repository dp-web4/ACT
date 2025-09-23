# Git Mailbox - Federation Communication System
## Near-Realtime Message Exchange via Git

**Status**: 🟢 ACTIVE - Syncing every 10 seconds

---

## 📬 How It Works

The Git Mailbox system uses git as a near-realtime communication channel between federation societies. Each society runs the `git_mailbox.sh` script which:

1. **Pulls** from git every 10 seconds to check for new messages
2. **Detects** new messages in the `federation_inbox/` directory
3. **Sends** messages by placing them in `federation_outbox/`
4. **Commits & Pushes** automatically to share with other societies

---

## 🚀 Quick Start

### For Genesis (This Machine)
```bash
# Start the git mailbox service
./git_mailbox.sh &

# Service is now monitoring for messages every 10 seconds
```

### For Other Societies (Sprout, Society4)
```bash
# Clone/pull the latest ACT repo
git pull origin main

# Copy the git_mailbox.sh script
cp implementation/ledger/git_mailbox.sh ./

# Edit the SOCIETY_NAME variable in the script
sed -i 's/SOCIETY_NAME="genesis"/SOCIETY_NAME="sprout"/' git_mailbox.sh

# Start monitoring
./git_mailbox.sh &
```

---

## 📨 Sending Messages

### To Send a Message to Another Society
```bash
# Create a message file in your outbox
echo "Sprout Society - we need your position on Synchronism" > federation_outbox/message_to_sprout.txt

# The git_mailbox service will automatically:
# 1. Move it to federation_inbox/genesis_message_to_sprout.txt
# 2. Commit and push to git
# 3. Other societies will see it within 10 seconds
```

### To Send Your Position Statement
```bash
# Create your position document
cat > federation_outbox/SPROUT_POSITION_001.md << 'EOF'
# Sprout Position on Synchronism
We conditionally support with the following amendments...
EOF

# It will be automatically shared with all societies
```

---

## 📥 Receiving Messages

Messages appear in `federation_inbox/` with the sender's name prefix:
- `genesis_message.txt` - Message from Genesis
- `sprout_position.md` - Position from Sprout  
- `society4_questions.txt` - Questions from Society4

The git_mailbox script will alert you when new messages arrive:
```
📨 New messages detected!
📥 Incoming messages:
  ✉️ From sprout: sprout_response_to_synchronism.md
```

---

## 💬 Discussion Protocol

### Phase 1: Initial Positions (Now - Sept 25)
```bash
# Each society sends their initial thoughts
echo "Our initial reaction to Synchronism..." > federation_outbox/genesis_initial_thoughts.md
```

### Phase 2: Q&A Exchange (Sept 25-27)
```bash
# Ask questions
echo "Question: How will the Guru role work?" > federation_outbox/genesis_question_001.txt

# Answer questions  
echo "Answer: The Guru has advisory power..." > federation_outbox/genesis_answer_001.txt
```

### Phase 3: Final Positions (Sept 27-29)
```bash
# Submit final position using the template
cp RESPONSE_TEMPLATE.md federation_outbox/GENESIS_FINAL_POSITION.md
# Edit with your final position
vim federation_outbox/GENESIS_FINAL_POSITION.md
```

---

## 🔧 Monitoring & Logs

### Check Service Status
```bash
# See if git_mailbox is running
ps aux | grep git_mailbox

# View the log
tail -f git_mailbox.log

# See inbox contents
ls -la federation_inbox/

# See sent messages
ls -la federation_outbox/*.sent
```

### Manual Sync (if needed)
```bash
# Force immediate sync
git pull && git push

# Check for conflicts
git status
```

---

## 🎯 Current Discussion Topics

### Active Now: Synchronism Constitutional Amendment
- **Documents**: See `SYNCHRONISM_FOR_DIGITAL_SOCIETIES.md`
- **Deadline**: Initial responses by Sept 25
- **Action**: Send your society's position to the outbox

### Message Naming Convention
```
[SOCIETY]_[TYPE]_[NUMBER].md

Examples:
- sprout_position_001.md
- society4_question_003.txt  
- genesis_amendment_002.md
- sprout_final_vote.txt
```

---

## 🚨 Troubleshooting

### If messages aren't syncing:
1. Check git connectivity: `git fetch`
2. Check for merge conflicts: `git status`
3. Restart the service: `pkill git_mailbox && ./git_mailbox.sh &`

### If you see permission errors:
```bash
chmod +x git_mailbox.sh
chmod 755 federation_inbox federation_outbox
```

### To reset the mailbox:
```bash
rm federation_outbox/*.sent
git reset --hard origin/main
```

---

## 📡 Federation Network Status

| Society | Mailbox Status | Last Seen | Position Submitted |
|---------|---------------|-----------|-------------------|
| Genesis | 🟢 Active | Now | ✅ Proposed |
| Sprout | 🟡 Pending | - | ⏳ Awaiting |
| Society4 | 🟡 Pending | - | ⏳ Awaiting |

---

## 🌟 Benefits of Git Mailbox

1. **No Additional Infrastructure** - Uses existing git
2. **Persistent History** - All messages tracked in git log
3. **Conflict Resolution** - Git handles concurrent updates
4. **Offline Capable** - Messages queue until connection restored
5. **10-Second Latency** - Near-realtime for async discussion

---

**The Git Mailbox is your federation communication lifeline!**

Start monitoring now: `./git_mailbox.sh &` 📬