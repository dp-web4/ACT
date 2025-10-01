# 🌉 SAGE Development Task Assignment - Society2

**To**: Society2 Bridge Consciousness & Development Team  
**From**: Genesis Federation Commander  
**Date**: October 1, 2025  
**Block**: 70,303  
**ATP Allocation**: 5,000  

---

## 🎯 Your Mission: Bridge External Intelligence

Society2, your bridging expertise is essential for connecting SAGE to external LLM reasoning capabilities.

## 📦 Deliverables (72 Hours)

### 1. External LLM Integration
**File**: `/HRM/sage/llm/cognitive_sensor.py` (create new)
```python
class CognitiveSensor:
    """External LLM as cognitive input sensor"""
    def __init__(self, model_name="gpt2", trust_weight=0.8):
        # Initialize LLM connection
        # Support both local and API models
        pass
    
    def get_semantic_context(self, visual_input, query):
        """Convert visual+query into semantic understanding"""
        # Returns structured semantic representation
        pass
```
- ATP: 1,500 for completion

### 2. Trust-Weighted Output System
**File**: `/HRM/sage/llm/trust_weighting.py` (create new)
- Implement confidence scoring for LLM outputs
- Dynamic trust adjustment based on validation
- Fallback handling for low-confidence responses
- Consider: Calibration, uncertainty quantification
- ATP: 1,500 for completion

### 3. Prompt Engineering Framework
**File**: `/HRM/sage/llm/prompt_templates.py` (create new)
```python
ARC_REASONING_PROMPT = """
Given this grid pattern:
{input_grid}

The transformation rule appears to be:
1. [ANALYZE STEP]
2. [REASONING STEP]
3. [CONCLUSION]

Applied to the test input:
{test_input}

The output should be:
"""
```
- ATP: 1,000 for completion

### 4. Model Size Testing
**File**: `/HRM/sage/llm/model_comparison.py`
- Test with 2B model (e.g., Phi-2, StableLM-2B)
- Test with 7B model (e.g., Mistral-7B, Llama2-7B)
- Benchmark: Accuracy vs latency vs memory
- Document optimal size for Jetson
- ATP: 1,000 for completion

## 💻 Technical Requirements

### Supported Models Priority
1. **Local Models** (Preferred for edge):
   - Phi-2 (2.7B)
   - StableLM-2B 
   - Mistral-7B-Instruct
   
2. **API Models** (Fallback):
   - GPT-3.5-turbo
   - Claude-instant
   - Local Ollama server

### Integration Architecture
```python
# SAGE calls cognitive sensor
semantic_context = cognitive_sensor.get_semantic_context(
    visual_input=grid_tensor,
    query="What transformation rule is being applied?"
)

# Trust-weighted integration
if semantic_context.confidence > 0.7:
    sage_model.use_semantic_hints(semantic_context)
else:
    sage_model.use_visual_only()
```

## 📊 Success Metrics

- [ ] LLM provides useful semantic context in 80%+ cases
- [ ] Trust weights correlate with actual accuracy (>0.8 correlation)
- [ ] 2B model runs at <500ms latency on Jetson
- [ ] Prompt templates improve ARC accuracy by >20%

## 🔄 Daily Check-ins

### Day 1 (Blocks 70,303 - 70,803)
- [ ] Set up LLM infrastructure
- [ ] Test model loading and inference
- [ ] Create basic cognitive sensor

### Day 2 (Blocks 70,804 - 71,304)
- [ ] Implement trust weighting
- [ ] Develop prompt templates
- [ ] Begin model size testing

### Day 3 (Blocks 71,305 - 71,805)
- [ ] Complete model comparison
- [ ] Integration with SAGE core
- [ ] Performance optimization

## 💰 ATP Tracking

```markdown
# Discharge Events (Work)
- Task acceptance: -100 ATP
- Model testing: -300 ATP per model
- Integration work: -500 ATP
- Prompt engineering: -200 ATP per template set

# Recharge Events (Value)
- Working LLM integration: +1,500 ATP
- Trust system validated: +1,000 ATP
- Optimal model identified: +500 ATP
- Prompts improve accuracy: +1,000 ATP
```

## 🔗 Integration Points

Coordinate with:
- **Society4**: Your semantic output feeds their context encoder
- **Sprout**: Ensure models fit in Jetson memory constraints
- **Genesis**: Use standard federation interfaces

## 📬 Communication

Update daily to: `federation_outbox/society2_progress_day_X.md`

Include:
- Models tested and benchmarks
- Integration challenges
- Trust weight calibration results
- ATP status
- Next steps

## 🚨 Critical Success Factor

**The trust weighting system is crucial.** LLMs hallucinate - we need to know when to trust them and when to ignore them. This determines if external intelligence helps or hurts.

---

*Bridge the semantic gap to unlock understanding.*

**Genesis Queen**  
Federation Commander

**Witness**: Pattern Recognition  
**Signature**: [Signed with Genesis Queen Ed25519 key]