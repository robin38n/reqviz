import type { Selection } from "d3";
import type {
	EndpointNode,
	GraphNode,
	SchemaNode,
	SpecGraph,
} from "../../../models/graph.model";
import { METHOD_COLORS } from "../../../shared/constants/method-colors";
import { truncateLabel } from "../../../shared/utils/truncate-label";

const MAX_ENDPOINT_W = 280;
const MAX_SCHEMA_W = 200;

/** A graph node with the visual fields both layouts need; each layout extends it. */
export interface SimNode {
	id: string;
	type: "endpoint" | "schema";
	label: string;
	fullLabel: string;
	sublabel: string;
	method?: string;
	width: number;
	height: number;
	original: GraphNode;
}

export function buildSimNodes(
	graph: SpecGraph,
	schemaHeight: number,
	schemaSublabel: (sc: SchemaNode) => string,
): SimNode[] {
	return graph.nodes.map((n): SimNode => {
		if (n.type === "endpoint") {
			const ep = n as EndpointNode;
			const fullLabel = `${ep.method} ${ep.path}`;
			return {
				id: n.id,
				type: "endpoint",
				label: truncateLabel(fullLabel, MAX_ENDPOINT_W, 7.5),
				fullLabel,
				sublabel: ep.summary || "",
				method: ep.method,
				width: Math.min(
					MAX_ENDPOINT_W,
					Math.max(140, fullLabel.length * 7.5 + 24),
				),
				height: 40,
				original: n,
			};
		}
		const sc = n as SchemaNode;
		return {
			id: n.id,
			type: "schema",
			label: truncateLabel(sc.name, MAX_SCHEMA_W, 8),
			fullLabel: sc.name,
			sublabel: schemaSublabel(sc),
			width: Math.min(MAX_SCHEMA_W, Math.max(120, sc.name.length * 8 + 24)),
			height: schemaHeight,
			original: n,
		};
	});
}

export function appendArrowMarker(
	defs: Selection<SVGDefsElement, unknown, null, undefined>,
	id: string,
	color: string,
): void {
	defs
		.append("marker")
		.attr("id", id)
		.attr("viewBox", "0 0 10 6")
		.attr("refX", 10)
		.attr("refY", 3)
		.attr("markerWidth", 8)
		.attr("markerHeight", 6)
		.attr("orient", "auto")
		.append("path")
		.attr("d", "M0,0 L10,3 L0,6 Z")
		.attr("fill", color);
}

/** Renders endpoint rectangles, labels, and truncation tooltips (shared by both layouts). */
export function appendEndpointShapes<N extends SimNode>(
	nodeGroup: Selection<SVGGElement, N, SVGGElement, unknown>,
): void {
	const endpoints = nodeGroup.filter((d) => d.type === "endpoint");

	endpoints
		.append("rect")
		.attr("width", (d) => d.width)
		.attr("height", (d) => d.height)
		.attr("rx", 4)
		.attr("ry", 4)
		.attr("fill", (d) => METHOD_COLORS[d.method ?? "GET"] ?? "#6b7280")
		.attr("opacity", 0.9);

	endpoints
		.append("text")
		.text((d) => d.label)
		.attr("x", (d) => d.width / 2)
		.attr("y", (d) => d.height / 2)
		.attr("text-anchor", "middle")
		.attr("dominant-baseline", "central")
		.attr("fill", "var(--graph-endpoint-text)")
		.attr("font-size", 11)
		.attr("font-family", "monospace")
		.attr("font-weight", 600);

	endpoints
		.filter((d) => d.label !== d.fullLabel)
		.append("title")
		.text((d) => d.fullLabel);
}
