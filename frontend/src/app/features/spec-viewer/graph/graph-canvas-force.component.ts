import { ChangeDetectionStrategy, Component } from "@angular/core";
import * as d3 from "d3";
import type { EdgeKind } from "../../../models/graph.model";
import {
	EDGE_COLORS,
	EDGE_DASH,
	SCHEMA_FILL,
	SCHEMA_STROKE,
} from "../../../shared/constants/edge-styles";
import { GraphCanvasBase } from "./graph-canvas-base";
import { GraphControlsComponent } from "./graph-controls.component";
import { GraphLegendComponent } from "./graph-legend.component";
import {
	appendArrowMarker,
	appendEndpointShapes,
	buildSimNodes,
	type SimNode,
} from "./graph-render";

type ForceNode = SimNode & d3.SimulationNodeDatum;

interface SimLink extends d3.SimulationLinkDatum<ForceNode> {
	kind: EdgeKind;
	label?: string;
	curveOffset: number;
}

@Component({
	selector: "app-graph-canvas-force",
	imports: [GraphControlsComponent, GraphLegendComponent],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./graph-canvas-force.component.html",
	styleUrl: "./graph-canvas-force.component.css",
})
export class GraphCanvasForceComponent extends GraphCanvasBase {
	private simulation: d3.Simulation<ForceNode, SimLink> | null = null;

	protected initGraph(): void {
		const graph = this.graph();
		const svgEl = this.svgRef().nativeElement;

		if (this.simulation) {
			this.simulation.stop();
			this.simulation = null;
		}
		d3.select(svgEl).selectAll("*").remove();

		if (!graph || graph.nodes.length === 0) return;

		const container = this.containerRef().nativeElement;
		const width = container.clientWidth || 800;
		const height = container.clientHeight || 500;

		const nodes: ForceNode[] = buildSimNodes(graph, 30, () => "");

		const nodeMap = new Map(nodes.map((n) => [n.id, n]));

		const links: SimLink[] = graph.edges
			.filter((e) => nodeMap.has(e.source) && nodeMap.has(e.target))
			.map((e) => ({
				source: e.source,
				target: e.target,
				kind: e.kind,
				label: e.label,
				curveOffset: 0,
			}));

		// Spread parallel edges (same source+target pair) symmetrically so they
		// fan out instead of overlapping: -25,+25 for two; -25,0,+25 for three.
		const pairCount = new Map<string, number>();
		const pairIndex = new Map<string, number>();
		for (const l of links) {
			const key = [String(l.source), String(l.target)].sort().join("||");
			pairCount.set(key, (pairCount.get(key) ?? 0) + 1);
		}
		for (const l of links) {
			const key = [String(l.source), String(l.target)].sort().join("||");
			const total = pairCount.get(key) ?? 1;
			if (total > 1) {
				const idx = pairIndex.get(key) ?? 0;
				pairIndex.set(key, idx + 1);
				l.curveOffset = (idx - (total - 1) / 2) * 25;
			}
		}

		const svg = d3
			.select(svgEl)
			.attr("viewBox", `0 0 ${width} ${height}`)
			.attr("preserveAspectRatio", "xMidYMid meet");

		const defs = svg.append("defs");
		for (const [kind, color] of Object.entries(EDGE_COLORS)) {
			appendArrowMarker(defs, `arrow-force-${kind}`, color);
		}

		const g = svg.append("g");
		this.setupZoom(svg, g);

		const linkGroup = g
			.append("g")
			.attr("class", "links")
			.selectAll("path")
			.data(links)
			.join("path")
			.attr("stroke", (d) => EDGE_COLORS[d.kind])
			.attr("stroke-width", 1.5)
			.attr("stroke-dasharray", (d) => EDGE_DASH[d.kind])
			.attr("marker-end", (d) => `url(#arrow-force-${d.kind})`)
			.attr("opacity", 0.7)
			.attr("fill", "none");

		const edgeLabelGroup = g
			.append("g")
			.attr("class", "edge-labels")
			.selectAll("text")
			.data(links.filter((l) => l.label))
			.join("text")
			.text((d) => d.label ?? "")
			.attr("font-size", 9)
			.attr("fill", "var(--graph-text-sub)")
			.attr("text-anchor", "middle")
			.attr("dy", -4);

		const nodeGroup = g
			.append("g")
			.attr("class", "nodes")
			.selectAll<SVGGElement, ForceNode>("g")
			.data(nodes)
			.join("g")
			.attr("cursor", "pointer")
			.on("click", (_event, d) => {
				this.nodeClick.emit(d.original);
			});

		appendEndpointShapes(nodeGroup);

		const schemaNodes = nodeGroup.filter((d) => d.type === "schema");

		schemaNodes
			.append("rect")
			.attr("width", (d) => d.width)
			.attr("height", (d) => d.height)
			.attr("rx", 3)
			.attr("ry", 3)
			.attr("fill", SCHEMA_FILL)
			.attr("stroke", SCHEMA_STROKE)
			.attr("stroke-width", 1.5);

		schemaNodes
			.append("text")
			.text((d) => d.label)
			.attr("x", (d) => d.width / 2)
			.attr("y", (d) => d.height / 2)
			.attr("text-anchor", "middle")
			.attr("dominant-baseline", "central")
			.attr("fill", "var(--graph-text)")
			.attr("font-size", 12)
			.attr("font-weight", 700);

		schemaNodes
			.filter((d) => d.label !== d.fullLabel)
			.append("title")
			.text((d) => d.fullLabel);

		// Low alphaTarget while dragging minimizes drift of unconnected subgraphs.
		const drag = d3
			.drag<SVGGElement, ForceNode>()
			.on("start", (event, d) => {
				if (!event.active) simulation.alphaTarget(0.05).restart();
				d.fx = d.x;
				d.fy = d.y;
			})
			.on("drag", (event, d) => {
				d.fx = event.x;
				d.fy = event.y;
			})
			.on("end", (event, d) => {
				if (!event.active) simulation.alphaTarget(0);
				d.fx = null;
				d.fy = null;
			});

		nodeGroup.call(drag);

		// High velocityDecay dampens drift of unconnected subgraphs; forceX/forceY
		// at low strength replace forceCenter so disconnected clusters aren't pulled
		// toward the center during drag.
		const simulation = d3
			.forceSimulation<ForceNode, SimLink>(nodes)
			.velocityDecay(0.7)
			.force(
				"link",
				d3
					.forceLink<ForceNode, SimLink>(links)
					.id((d) => d.id)
					.distance(160),
			)
			.force("charge", d3.forceManyBody().strength(-350))
			.force("x", d3.forceX(width / 2).strength(0.03))
			.force("y", d3.forceY(height / 2).strength(0.03))
			.force(
				"collision",
				d3
					.forceCollide<ForceNode>()
					.radius((d) => Math.max(d.width, d.height) / 2 + 10),
			)
			.on("tick", () => {
				linkGroup.attr("d", (d) => {
					const src = d.source as ForceNode;
					const tgt = d.target as ForceNode;
					const x1 = (src.x ?? 0) + src.width / 2;
					const y1 = (src.y ?? 0) + src.height / 2;
					const x2 = (tgt.x ?? 0) + tgt.width / 2;
					const y2 = (tgt.y ?? 0) + tgt.height / 2;

					if (d.curveOffset === 0) {
						return `M${x1},${y1} L${x2},${y2}`;
					}

					const dx = x2 - x1;
					const dy = y2 - y1;
					const len = Math.sqrt(dx * dx + dy * dy) || 1;
					const nx = -dy / len;
					const ny = dx / len;
					const cx = (x1 + x2) / 2 + nx * d.curveOffset;
					const cy = (y1 + y2) / 2 + ny * d.curveOffset;
					return `M${x1},${y1} Q${cx},${cy} ${x2},${y2}`;
				});

				edgeLabelGroup
					.attr("x", (d) => {
						const src = d.source as ForceNode;
						const tgt = d.target as ForceNode;
						const x1 = (src.x ?? 0) + src.width / 2;
						const x2 = (tgt.x ?? 0) + tgt.width / 2;
						if (d.curveOffset === 0) return (x1 + x2) / 2;
						const dx = x2 - x1;
						const dy =
							(tgt.y ?? 0) + tgt.height / 2 - ((src.y ?? 0) + src.height / 2);
						const len = Math.sqrt(dx * dx + dy * dy) || 1;
						return (x1 + x2) / 2 + (-dy / len) * d.curveOffset;
					})
					.attr("y", (d) => {
						const src = d.source as ForceNode;
						const tgt = d.target as ForceNode;
						const y1 = (src.y ?? 0) + src.height / 2;
						const y2 = (tgt.y ?? 0) + tgt.height / 2;
						if (d.curveOffset === 0) return (y1 + y2) / 2;
						const dx =
							(tgt.x ?? 0) + tgt.width / 2 - ((src.x ?? 0) + src.width / 2);
						const dy = y2 - y1;
						const len = Math.sqrt(dx * dx + dy * dy) || 1;
						return (y1 + y2) / 2 + (dx / len) * d.curveOffset;
					});

				nodeGroup.attr(
					"transform",
					(d) => `translate(${d.x ?? 0},${d.y ?? 0})`,
				);
			});

		this.simulation = simulation;
	}
}
