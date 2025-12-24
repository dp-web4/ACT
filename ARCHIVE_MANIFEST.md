# ACT Repository Archive Manifest

**Date**: 2025-12-24
**Housekeeping Session**: Repository cleanup and organization

---

## Archived Content

### research-archive/society2/
**Original Location**: `implementation/society2/`
**Size**: ~5.4 MB
**Archived**: 2025-12-24
**Status**: Frozen reference implementation

**Reason for Archival**:
- Only 3 commits (Sep 28, 2025)
- No activity since late September
- Society4 is the active implementation (Nov-Dec 2025)
- Valuable as reference but not under active development

**Contents**:
- Complete Society2 blockchain implementation
- Early governance experiments
- Federation inbox/outbox patterns
- Original consensus mechanisms

**Historical Value**:
Shows early Web4 blockchain approach before Society4 evolution.

**Decision**: Keep for reference, archive for clarity.

---

### research-archive/society-prototype/
**Original Location**: `implementation/society-prototype/`
**Size**: Single document
**Archived**: 2025-12-24
**Status**: Unused prototype

**Reason for Archival**:
- Created Sep 18, 2025
- Never developed beyond initial document
- Superseded by Society2/4 implementations

**Contents**:
- Initial prototype design document

**Historical Value**:
Shows very early thinking before Society blockchain series.

**Decision**: Archive for historical record.

---

## Deleted Content

### /ACT/ Empty Duplicate Directory
**Location**: `/home/dp/ai-workspace/ACT/`
**Deleted**: 2025-12-24
**Status**: Empty duplicate (filesystem case sensitivity artifact)

**Reason for Deletion**:
- Completely empty (only had empty `implementation/` subdirectory)
- Case sensitivity confusion (/ACT/ vs /act/)
- /act/ is the real repository with git history
- No valuable content

**Decision**: Delete immediately to prevent confusion.

---

### Build Artifacts (Removed from git tracking)
**Files Removed**:
- `ignite_29.0.0_linux_amd64.tar.gz` (20+ MB tarball at root)
- `implementation/ledger/go.mod.backup`
- `implementation/ledger/go.sum.backup`

**Deleted**: 2025-12-24
**Status**: Build artifacts, not source code

**Reason for Removal**:
- Should not be tracked in git
- Tarballs should be downloaded as needed
- .backup files are build system artifacts
- Bloat repository without value

**Decision**: Remove from tracking, add patterns to .gitignore.

---

## .gitignore Updates

**Added patterns** (2025-12-24):
```gitignore
# Build artifacts and backups
*.backup
ignite_*.tar.gz
*.tar.gz

# Blockchain data directories (all locations)
society/data/
implementation/ledger/data/
data/
```

**Reason**: Prevent future tracking of build artifacts and runtime data.

---

## Content Kept (NOT Archived)

### federation_inbox/ and federation_outbox/
**Location**: Root level
**Size**: 41 governance files
**Status**: **KEPT** for governance history

**Reason for Keeping**:
- Valuable governance history
- Referenced by current implementations
- Frozen but not deprecated
- May inform future federation work

**Decision**: Keep for now, revisit in quarterly review.

---

## Current Active Components

Post-housekeeping, active development focuses on:

### Primary Implementation
- `implementation/ledger/` - Main blockchain (913 Go files, 5.4 MB)
- `implementation/society4/` - Production target (very active, Dec 23)
- `implementation/cbp-chain/` - Agent simulation framework

### Specifications
- `core-spec/` - Core protocol specifications
- `docs/` - Technical documentation
- `rfcs/` - RFC governance documents

### Supporting
- `philosophy/` - Foundational concepts
- `machines/` - Machine-specific configurations
- `tool/` - Development utilities
- `demo-society/` - Educational demos

---

## Repository Health

**Before Housekeeping**:
- Duplicate /ACT/ directory causing confusion
- Build artifacts tracked in git (23+ MB)
- Mixed frozen/active implementations
- No clear archival structure

**After Housekeeping**:
- Single /act/ directory (lowercase)
- Clean git tracking (build artifacts removed)
- research-archive/ for frozen implementations
- Clear separation: active vs reference vs archived

**Improvement**: Significant - repository now clearly organized.

---

## Next Housekeeping (Quarterly)

**Candidates for Future Review**:
1. federation_inbox/outbox/ - If no longer referenced, move to archive
2. philosophy/ - Consider separate archive repo for research content
3. demo-society/ - If superseded by better demos, archive

**Timeline**: Q1 2026 (March 2026)

---

## Archive Access

**All archived content remains in git history and in research-archive/.**

To access archived content:
```bash
# View current archive
ls -la research-archive/

# Access society2 code
cd research-archive/society2/

# View historical reference
cat research-archive/society-prototype/[file]
```

**Nothing was lost** - only reorganized for clarity.

---

**Housekeeping completed**: 2025-12-24
**Next review**: Q1 2026
**Repository status**: CLEAN AND ORGANIZED ✅
