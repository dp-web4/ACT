# ACT Repository - Quick Reference Guide

**Last Updated**: December 24, 2025
**Analysis Scope**: Complete directory structure, file organization, deprecation status
**Key Finding**: Well-organized research prototype with ~1,400 files across 913 Go source files

---

## Health Status at a Glance

| Metric | Status | Notes |
|--------|--------|-------|
| **Overall Status** | EXPERIMENTAL | 65% complete (per README) |
| **Last Commit** | Dec 18-23, 2025 | Byzantine consensus work |
| **Code Quality** | GOOD | Structured patterns, well-documented |
| **Organization** | ADEQUATE | Cleanup opportunities exist |
| **Main Issue** | CASE SENSITIVITY | Both `/act/` and `/ACT/` directories exist |

---

## Component Quick View

### Keep Everything (TIER 1 - ACTIVE)
- **Ledger** (5.4M, 913 Go files) - Primary blockchain implementation
- **Society4** (5.4M) - Production target blockchain with Phase 1-3 compliance
- **CBP-Chain** (Python) - Autonomous agent simulation framework
- **Core-Spec** - Architecture specifications
- **Docs & RFCs** - Governance framework and standards proposals

### Keep But Review (TIER 2 - SUPPORTING)
- **Philosophy/** - Foundational research (candidate for archival)
- **Machines/** - Machine-specific setup (operational reference)
- **Tool/** - Small utility library
- **Demo-Society/** - Educational example implementation

### Archive or Delete (TIER 3 - DEPRECATED)
- **Society-Prototype/** (12KB) - Single unused document
- **Conversations/** (1MB) - Single historical session from Jan 2025
- **Federation Inbox/Outbox** (41 files) - Archive but preserve governance history
- **Society2** (4.9M) - Frozen reference (verify intent before deleting)

### Remove from Git (TIER 4 - BUILD ARTIFACTS)
- **ignite_29.0.0_linux_amd64.tar.gz** (24MB) - Replace with download script
- **go.mod.backup, go.sum.backup** - Remove from version control
- **Database files** (not in git yet) - Prevent via .gitignore

---

## The One Problem: Case Sensitivity Issue

```
PROBLEM: Two directories exist:
  /home/dp/ai-workspace/act/         ← REAL (has .git history)
  /home/dp/ai-workspace/ACT/         ← COPY (empty except /implementation/)

SOLUTION:
  rm -rf /home/dp/ai-workspace/ACT
```

---

## What's Actually Here

### By the Numbers
- **Go source files**: 913
- **Protobuf definitions**: 190
- **Python utilities**: 50
- **Documentation files**: 175+
- **Total tracked files**: ~1,400

### By Component
```
implementation/
  ├─ ledger/          (5.4M) ACTIVE - Core blockchain
  ├─ society4/        (5.4M) VERY ACTIVE - Production target
  ├─ society2/        (4.9M) FROZEN - Reference implementation
  ├─ cbp-chain/       (ACTIVE) - Agent simulation
  └─ society-prototype/ (12KB) DEPRECATED

docs/                 ACTIVE - Proposal documents
core-spec/           ACTIVE - Architecture specs
rfcs/                ACTIVE - Standards proposals
philosophy/          RESEARCH - Conceptual foundations
machines/            ACTIVE - Machine setup reference
tool/                SMALL - Utility library
demo-society/        EXAMPLE - Educational
swarm-bootstrap/     RESEARCH - Experiments
conversations/       HISTORICAL - Single session
society/             DEV - Runtime blockchain state
```

---

## Activity Timeline

```
CURRENT (Dec 18-23)  → Byzantine consensus, trust monitoring
RECENT (Nov-Dec)     → ATP/ADP integration, compliance verification  
ACTIVE (Oct)         → Federation, RFC voting, governance
SETUP (Sep)          → Initial blockchain setup, specifications
```

**Pattern**: 30-40 commits/month during active phases, session-based organization

---

## Critical Cleanup Items (Priority Order)

### IMMEDIATE
1. Delete `/home/dp/ai-workspace/ACT/` (empty duplicate)
2. Add `ignite_*.tar.gz` to .gitignore
3. Add `*.backup` to .gitignore
4. Update .gitignore for blockchain data directories

### NEXT WEEK
1. Decide on Society2 (keep as reference or delete?)
2. Archive federation governance files to `research-archive/`
3. Delete or archive `society-prototype/`
4. Create `ARCHIVE_MANIFEST.md` documenting decisions

### NEXT MONTH
1. Move philosophical/research content to `research-archive/`
2. Update README with maturity levels for each component
3. Clarify blockchain development status (active vs reference)

---

## Key Documentation Locations

| What | Where |
|------|-------|
| **Main README** | `/act/README.md` |
| **Key Discoveries** | `/act/DISCOVERIES.md` |
| **Complete Analysis** | `/act/ACT_REPOSITORY_ANALYSIS.md` (this file) |
| **Architecture Specs** | `/act/core-spec/` |
| **Society4 Details** | `/act/implementation/society4/README.md` |
| **Governance Framework** | `/act/rfcs/` and federation documents |
| **API Documentation** | `/act/implementation/ledger/api-bridge/README.md` |

---

## For New Contributors

**Start here**: `/home/dp/ai-workspace/act/README.md`

**Understand structure**: This Quick Reference Guide

**Deep dive**: `/act/ACT_REPOSITORY_ANALYSIS.md` (comprehensive 18KB analysis)

**Learn the tech**: `/act/core-spec/` and `/act/implementation/ledger/docs/`

**See it working**: `/act/demo-society/` (small, educational example)

---

## Gotchas & Notes

1. **Case Sensitivity**: Both `/act/` and `/ACT/` exist (delete `/ACT/`)
2. **Large Tarball**: 24MB ignite binary in git (should be downloaded instead)
3. **Duplicate Docs**: 9+ large markdown files duplicated across Society2/4 (intentional - independent blockchains)
4. **Database Files**: Blockchain runtime state accumulating locally (need .gitignore refresh)
5. **Research Code**: 41 federation governance files preserved as historical record

---

## At a Glance: What Gets Deleted vs Archived

### DELETE
- `/ACT/` directory (uppercase duplicate, empty)
- `society-prototype/` (single unused document)
- `.backup` files from git history

### ARCHIVE to research-archive/
- `philosophy/` (foundational research)
- `federation_inbox/` and `federation_outbox/` (governance history)
- `conversations/` (historical sessions)
- `swarm-bootstrap/` (if superseded by SAGE/HRM)

### KEEP EVERYTHING ELSE
- Implementation code (Ledger, Society4, CBP-Chain)
- Specifications and documentation
- RFC and governance framework
- Machines and setup information
- Demo and example code

---

## Directory Tree (Active Components Only)

```
/act/
├── implementation/
│   ├── ledger/              ✅ PRIMARY - Active development
│   ├── society4/            ✅ PRIMARY - Production target
│   ├── cbp-chain/           ✅ Supporting - Agent simulation
│   ├── society2/            ⚠️  Reference (frozen, 3 commits)
│   └── society-prototype/   ❌ Deprecated (1 file)
├── core-spec/               ✅ Architecture specifications
├── docs/                    ✅ Proposals and governance
├── rfcs/                    ✅ Standards evolution
├── demo-society/            ✅ Educational example
├── machines/                ✅ Setup documentation
├── tool/                    ✅ Utility library
├── philosophy/              📦 Research (archive candidate)
├── swarm-bootstrap/         📦 Research (archive candidate)
├── conversations/           📦 Historical (archive candidate)
├── README.md                ✅ Main documentation
└── DISCOVERIES.md           ✅ Key insights
```

Legend: ✅ Keep | ⚠️ Review | ❌ Delete | 📦 Archive

---

## One-Page Summary for Meetings

**What**: ACT is a research implementation of Web4 (trust-native protocol for human-AI interaction)

**Size**: ~1,400 files, 1.1K lines of code (Go), primarily blockchain implementations

**Status**: Experimental, 65% complete, actively developed (commits every 1-2 days)

**Quality**: Good code organization, well-documented, follows Cosmos SDK patterns

**Main Issues**: 
1. Case-sensitive directory duplication (ACT vs act)
2. 24MB tarball tracked in git
3. Some cleanup opportunities (old research files, federation archives)

**Recommendation**: Keep core implementation, archive research/exploratory work, fix case sensitivity issue

**Effort to Clean**: 1-2 hours for immediate fixes, 1 day for complete reorganization

---

## For Repository Maintainers

### Immediate Actions (30 minutes)
```bash
# 1. Remove case-sensitive duplicate
rm -rf /home/dp/ai-workspace/ACT

# 2. Improve .gitignore
echo "ignite_*.tar.gz" >> .gitignore
echo "*.backup" >> .gitignore
echo "implementation/*/data/" >> .gitignore

# 3. Remove backup files from git
git rm --cached implementation/ledger/go.mod.backup
git rm --cached implementation/ledger/go.sum.backup
git commit -m "fix: Remove build artifacts from version control"
```

### Short-term Cleanup (1 week)
```bash
# 1. Archive research content
mkdir research-archive
mv philosophy/ research-archive/
mv conversations/ research-archive/
mv implementation/ledger/federation_inbox research-archive/
mv implementation/ledger/federation_outbox research-archive/

# 2. Decide on society2 (archive or reference?)
git rm -r implementation/society2/  # if deprecated
# OR keep and document as reference implementation

# 3. Delete unused prototype
git rm -r implementation/society-prototype/

# 4. Create manifest
cat > ARCHIVE_MANIFEST.md << 'MANIFEST'
# Archived Content Manifest

## Research Archive
- philosophy/ - Foundational research
- conversations/ - Historical sessions
- federation_inbox/outbox - Governance history
- swarm-bootstrap/ - Architecture experiments

## Rationale
These contain valuable intellectual history but are not active development.
Moved to research-archive/ to reduce main directory clutter while preserving content.
MANIFEST

git add -A
git commit -m "refactor: Reorganize repository structure"
```

### Long-term Improvements (ongoing)
1. Establish maturity tiers for components
2. Consider splitting into separate repositories if monorepo becomes unwieldy
3. Implement automated documentation updates
4. Create contributor guidelines reflecting structure

---

## Questions to Verify with Dennis

1. **Society2**: Keep as reference implementation or deprecate?
2. **Swarm-Bootstrap**: Still relevant or superseded by SAGE/HRM work?
3. **Federation Governance Files**: Historical value only, or still referenced?
4. **Philosophy Directory**: Keep in main repo or separate research collection?

---

## Files for Reference

**For Comprehensive Analysis**: See `/act/ACT_REPOSITORY_ANALYSIS.md` (18KB, complete details)

**For Quick Overview**: You're reading it

**For Getting Started**: `/act/README.md`

