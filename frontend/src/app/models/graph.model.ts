export interface EndpointNode {
	readonly id: string; // e.g. "GET /pets"
	readonly type: "endpoint";
	readonly path: string;
	readonly method: string;
	readonly summary: string;
	readonly operationId?: string;
	readonly tags: string[];
}

export interface SchemaNode {
	readonly id: string; // e.g. "schema:Pet"
	readonly type: "schema";
	readonly name: string;
	readonly properties: string[];
	readonly requiredProps: string[];
}

export type GraphNode = EndpointNode | SchemaNode;

export type EdgeKind =
	| "requestBody"
	| "response"
	| "parameter"
	| "property"
	| "arrayItem"
	| "composition";

export interface GraphEdge {
	readonly source: string;
	readonly target: string;
	readonly kind: EdgeKind;
	readonly label?: string;
}

export interface SpecGraph {
	readonly nodes: GraphNode[];
	readonly edges: GraphEdge[];
}
