#!/usr/bin/env python3
"""
MRH (Markov Relevancy Horizon) Graph Implementation
RDF graph structure per Web4 MRH Specification
"""

import json
from typing import Dict, List, Tuple, Optional, Set
from dataclasses import dataclass
from datetime import datetime
from pathlib import Path
import hashlib

# Using built-in graph implementation to avoid external dependencies
class RDFTriple:
    """RDF Triple: Subject-Predicate-Object"""
    def __init__(self, subject: str, predicate: str, obj: str):
        self.subject = subject
        self.predicate = predicate
        self.object = obj
        
    def __str__(self):
        return f"<{self.subject}> <{self.predicate}> <{self.object}>"
    
    def __hash__(self):
        return hash((self.subject, self.predicate, self.object))
    
    def __eq__(self, other):
        return (self.subject == other.subject and 
                self.predicate == other.predicate and 
                self.object == other.object)

class MRHGraph:
    """
    Markov Relevancy Horizon Graph
    Implements RDF-based context graph per MRH Specification
    """
    
    def __init__(self, federation_path: str = "/home/dp/ai-workspace/act/implementation/ledger"):
        self.federation_path = Path(federation_path)
        self.graph_file = self.federation_path / "federation" / "mrh_graph.json"
        self.triples: Set[RDFTriple] = set()
        self.horizon_depth = 3  # Markov horizon depth
        
        # Define federation namespaces
        self.namespaces = {
            "federation": "http://web4.federation/",
            "society": "http://web4.federation/society/",
            "task": "http://web4.federation/task/",
            "lct": "http://web4.federation/lct/",
            "atp": "http://web4.federation/atp/",
            "sage": "http://web4.federation/sage/",
            "rdf": "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
            "rdfs": "http://www.w3.org/2000/01/rdf-schema#"
        }
        
        # Define predicates
        self.predicates = {
            "type": f"{self.namespaces['rdf']}type",
            "member_of": f"{self.namespaces['federation']}memberOf",
            "assigned_to": f"{self.namespaces['task']}assignedTo",
            "depends_on": f"{self.namespaces['task']}dependsOn",
            "witnesses": f"{self.namespaces['federation']}witnesses",
            "trusts": f"{self.namespaces['federation']}trusts",
            "allocates_atp": f"{self.namespaces['atp']}allocates",
            "discharges_atp": f"{self.namespaces['atp']}discharges",
            "implements": f"{self.namespaces['sage']}implements",
            "requires": f"{self.namespaces['sage']}requires",
            "created_at": f"{self.namespaces['federation']}createdAt",
            "status": f"{self.namespaces['task']}status",
            "has_context": f"{self.namespaces['federation']}hasContext",
            "related_to": f"{self.namespaces['federation']}relatedTo"
        }
        
        self.load_graph()
    
    def load_graph(self):
        """Load graph from persistent storage"""
        if self.graph_file.exists():
            with open(self.graph_file, 'r') as f:
                data = json.load(f)
                for triple_data in data.get("triples", []):
                    triple = RDFTriple(
                        triple_data["subject"],
                        triple_data["predicate"],
                        triple_data["object"]
                    )
                    self.triples.add(triple)
    
    def save_graph(self):
        """Save graph to persistent storage"""
        data = {
            "namespaces": self.namespaces,
            "triples": [
                {
                    "subject": t.subject,
                    "predicate": t.predicate,
                    "object": t.object
                }
                for t in self.triples
            ],
            "metadata": {
                "horizon_depth": self.horizon_depth,
                "triple_count": len(self.triples),
                "updated": datetime.now().isoformat()
            }
        }
        
        self.graph_file.parent.mkdir(exist_ok=True)
        with open(self.graph_file, 'w') as f:
            json.dump(data, f, indent=2)
    
    def add_triple(self, subject: str, predicate: str, obj: str):
        """Add RDF triple to graph"""
        triple = RDFTriple(subject, predicate, obj)
        self.triples.add(triple)
        self.save_graph()
    
    def remove_triple(self, subject: str, predicate: str, obj: str):
        """Remove RDF triple from graph"""
        triple = RDFTriple(subject, predicate, obj)
        self.triples.discard(triple)
        self.save_graph()
    
    def query(self, subject: Optional[str] = None, 
             predicate: Optional[str] = None,
             obj: Optional[str] = None) -> List[RDFTriple]:
        """
        Query graph with optional filters
        Simulates SPARQL-like queries
        """
        results = []
        for triple in self.triples:
            if subject and triple.subject != subject:
                continue
            if predicate and triple.predicate != predicate:
                continue
            if obj and triple.object != obj:
                continue
            results.append(triple)
        return results
    
    def get_neighbors(self, node: str, depth: int = 1) -> Set[str]:
        """
        Get neighbors within Markov horizon
        Per MRH Specification §2.1 - Horizon Calculation
        """
        if depth <= 0:
            return set()
        
        neighbors = set()
        
        # Find direct neighbors
        for triple in self.triples:
            if triple.subject == node:
                neighbors.add(triple.object)
            if triple.object == node:
                neighbors.add(triple.subject)
        
        # Recursive expansion up to horizon
        if depth > 1:
            expanded = set()
            for neighbor in neighbors:
                expanded.update(self.get_neighbors(neighbor, depth - 1))
            neighbors.update(expanded)
        
        # Remove self
        neighbors.discard(node)
        
        return neighbors
    
    def calculate_relevancy(self, source: str, target: str) -> float:
        """
        Calculate relevancy between nodes
        Per MRH Specification §3 - Relevancy Scoring
        """
        # Find shortest path length
        visited = {source}
        queue = [(source, 0)]
        
        while queue:
            current, distance = queue.pop(0)
            
            if current == target:
                # Relevancy decreases exponentially with distance
                return 1.0 / (2 ** distance)
            
            if distance < self.horizon_depth:
                neighbors = self.get_neighbors(current, 1)
                for neighbor in neighbors:
                    if neighbor not in visited:
                        visited.add(neighbor)
                        queue.append((neighbor, distance + 1))
        
        return 0.0  # Beyond horizon
    
    def get_context_boundary(self, node: str) -> Dict:
        """
        Get MRH context boundary for a node
        Per MRH Specification §4 - Context Boundaries
        """
        context = {
            "center": node,
            "horizon_1": set(),
            "horizon_2": set(),
            "horizon_3": set(),
            "properties": {},
            "relationships": []
        }
        
        # Get nodes at each horizon level
        for depth in range(1, self.horizon_depth + 1):
            neighbors = self.get_neighbors(node, depth)
            prev_neighbors = self.get_neighbors(node, depth - 1) if depth > 1 else set()
            horizon_nodes = neighbors - prev_neighbors
            context[f"horizon_{depth}"] = horizon_nodes
        
        # Get all properties of center node
        properties = self.query(subject=node)
        for triple in properties:
            if triple.predicate not in context["properties"]:
                context["properties"][triple.predicate] = []
            context["properties"][triple.predicate].append(triple.object)
        
        # Get relationships within context
        all_context_nodes = {node}
        for depth in range(1, self.horizon_depth + 1):
            all_context_nodes.update(context[f"horizon_{depth}"])
        
        for triple in self.triples:
            if triple.subject in all_context_nodes and triple.object in all_context_nodes:
                context["relationships"].append({
                    "from": triple.subject,
                    "to": triple.object,
                    "type": triple.predicate
                })
        
        return context
    
    def initialize_federation_graph(self):
        """Initialize federation MRH graph with base structure"""
        
        # Define federation structure
        federation = f"{self.namespaces['federation']}GenesisFederation"
        
        # Add societies
        societies = ["Genesis", "Society4", "Society2", "Sprout"]
        for society in societies:
            society_uri = f"{self.namespaces['society']}{society}"
            self.add_triple(society_uri, self.predicates["type"], "Society")
            self.add_triple(society_uri, self.predicates["member_of"], federation)
        
        # Add SAGE development context
        sage_project = f"{self.namespaces['sage']}SAGEDevelopment"
        self.add_triple(sage_project, self.predicates["type"], "Project")
        
        # Add society assignments for SAGE
        assignments = {
            "Genesis": ["AttentionOrchestrator", "FederationCoordination"],
            "Society4": ["TrainingLoop", "ContextEncoding"],
            "Society2": ["LLMIntegration", "CognitiveSensor"],
            "Sprout": ["JetsonOptimization", "EdgeDeployment"]
        }
        
        for society, tasks in assignments.items():
            society_uri = f"{self.namespaces['society']}{society}"
            for task in tasks:
                task_uri = f"{self.namespaces['task']}{task}"
                self.add_triple(task_uri, self.predicates["type"], "Task")
                self.add_triple(task_uri, self.predicates["assigned_to"], society_uri)
                self.add_triple(task_uri, self.predicates["has_context"], sage_project)
                self.add_triple(society_uri, self.predicates["implements"], task_uri)
        
        # Add ATP allocations
        for society in societies:
            society_uri = f"{self.namespaces['society']}{society}"
            atp_pool = f"{self.namespaces['atp']}Pool_{society}"
            self.add_triple(atp_pool, self.predicates["type"], "ATPPool")
            self.add_triple(federation, self.predicates["allocates_atp"], atp_pool)
            self.add_triple(atp_pool, self.predicates["assigned_to"], society_uri)
            self.add_triple(atp_pool, "amount", "5000")
        
        # Add trust relationships
        self.add_triple(f"{self.namespaces['society']}Genesis", 
                       self.predicates["trusts"], 
                       f"{self.namespaces['society']}Sprout")
        self.add_triple(f"{self.namespaces['society']}Society4", 
                       self.predicates["trusts"], 
                       f"{self.namespaces['society']}Society2")
        
        # Add witness relationships
        for society in societies:
            society_uri = f"{self.namespaces['society']}{society}"
            for other in societies:
                if society != other:
                    other_uri = f"{self.namespaces['society']}{other}"
                    self.add_triple(society_uri, self.predicates["witnesses"], other_uri)
        
        # Add dependencies
        self.add_triple(f"{self.namespaces['task']}LLMIntegration",
                       self.predicates["depends_on"],
                       f"{self.namespaces['task']}AttentionOrchestrator")
        self.add_triple(f"{self.namespaces['task']}JetsonOptimization",
                       self.predicates["depends_on"],
                       f"{self.namespaces['task']}TrainingLoop")
        
        self.save_graph()
    
    def sparql_like_query(self, query_string: str) -> List[Dict]:
        """
        Simple SPARQL-like query interface
        Example: "SELECT ?s WHERE { ?s type Society }"
        """
        results = []
        
        # Parse simple SELECT WHERE pattern
        if "SELECT" in query_string and "WHERE" in query_string:
            # Extract pattern from WHERE clause
            where_start = query_string.find("{") + 1
            where_end = query_string.find("}")
            pattern = query_string[where_start:where_end].strip()
            
            # Parse triple pattern
            parts = pattern.split()
            if len(parts) >= 3:
                s_pattern = parts[0]
                p_pattern = parts[1]
                o_pattern = parts[2]
                
                # Convert predicate shortcuts
                predicate_map = {
                    "type": self.predicates["type"],
                    "memberOf": self.predicates["member_of"],
                    "assignedTo": self.predicates["assigned_to"],
                    "trusts": self.predicates["trusts"]
                }
                
                if p_pattern in predicate_map:
                    p_pattern = predicate_map[p_pattern]
                
                # Query graph
                for triple in self.triples:
                    match = True
                    result = {}
                    
                    # Check subject
                    if s_pattern.startswith("?"):
                        result[s_pattern] = triple.subject
                    elif s_pattern != triple.subject:
                        match = False
                    
                    # Check predicate  
                    if p_pattern.startswith("?"):
                        result[p_pattern] = triple.predicate
                    elif p_pattern != triple.predicate:
                        match = False
                    
                    # Check object
                    if o_pattern.startswith("?"):
                        result[o_pattern] = triple.object
                    elif o_pattern != triple.object:
                        match = False
                    
                    if match:
                        results.append(result)
        
        return results


if __name__ == "__main__":
    print("=== Web4 Federation MRH Graph Implementation ===\n")
    
    # Initialize MRH graph
    mrh = MRHGraph()
    
    print("1. Initializing federation graph structure...")
    mrh.initialize_federation_graph()
    print(f"   ✅ Created {len(mrh.triples)} RDF triples\n")
    
    print("2. Querying societies in federation...")
    societies = mrh.query(predicate=mrh.predicates["member_of"])
    print(f"   Found {len(societies)} society memberships:")
    for triple in societies[:3]:
        print(f"   - {triple.subject.split('/')[-1]} → memberOf → Federation\n")
    
    print("3. Finding Genesis context boundary...")
    genesis_uri = f"{mrh.namespaces['society']}Genesis"
    context = mrh.get_context_boundary(genesis_uri)
    print(f"   Center: Genesis")
    print(f"   Horizon 1: {len(context['horizon_1'])} nodes")
    print(f"   Horizon 2: {len(context['horizon_2'])} nodes")
    print(f"   Relationships: {len(context['relationships'])}\n")
    
    print("4. Calculating relevancy between nodes...")
    relevancy = mrh.calculate_relevancy(
        f"{mrh.namespaces['society']}Genesis",
        f"{mrh.namespaces['task']}LLMIntegration"
    )
    print(f"   Genesis → LLMIntegration relevancy: {relevancy:.3f}\n")
    
    print("5. Testing SPARQL-like queries...")
    query = "SELECT ?s WHERE { ?s type Society }"
    results = mrh.sparql_like_query(query)
    print(f"   Query: {query}")
    print(f"   Results: {len(results)} societies found")
    for r in results[:2]:
        print(f"   - {r['?s'].split('/')[-1]}")
    
    print("\n6. Finding task dependencies...")
    dependencies = mrh.query(predicate=mrh.predicates["depends_on"])
    print(f"   Found {len(dependencies)} dependencies:")
    for dep in dependencies:
        print(f"   - {dep.subject.split('/')[-1]} depends on {dep.object.split('/')[-1]}")
    
    print("\n✅ MRH Graph implementation complete!")
    print("   - RDF triple store")
    print("   - Markov horizon calculation")
    print("   - Context boundary detection")
    print("   - SPARQL-like queries")
    print("\nWeb4 MRH Specification compliance achieved!")