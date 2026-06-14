import { ChangeDetectionStrategy, Component } from "@angular/core";
import { layout as dagreLayout, Graph } from "@dagrejs/dagre";
import * as d3 from "d3";
import type { EdgeKind } from "../../../models/graph.model";
import {
	EDGE_DASH,
	edgeColor,
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

type StructuredNode = SimNode & { x: number; y: number };

interface SimLink {
	source: string;
	target: string;
	kind: EdgeKind;
	label?: string;
	points?: Array<{ x: number; y: number }>;
}

@Component({
	selector: "app-graph-canvas",
	imports: [GraphControlsComponent, GraphLegendComponent],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./graph-canvas.component.html",
	styleUrl: "./graph-canvas.component.css",
})
export class GraphCanvasComponent extends GraphCanvasBase {
	protected initGraph(): void {
		const graph = this.graph();
		const svgEl = this.svgRef().nativeElement;

		d3.select(svgEl).selectAll("*").remove();

		if (!graph || graph.nodes.length === 0) return;

		const container = this.containerRef().nativeElement;
		const width = container.clientWidth || 800;
		const height = container.clientHeight || 500;

		const nodes: StructuredNode[] = buildSimNodes(
			graph,
			44,
			(sc) => `${sc.properties.length} props`,
		).map((n) => ({ ...n, x: 0, y: 0 }));

		const nodeMap = new Map(nodes.map((n) => [n.id, n]));

		const links: SimLink[] = graph.edges
			.filter((e) => nodeMap.has(e.source) && nodeMap.has(e.target))
			.map((e) => ({
				source: e.source,
				target: e.target,
				kind: e.kind,
				label: e.label,
			}));

		const g2 = new Graph()
			.setGraph({ rankdir: "TB", nodesep: 40, ranksep: 120, edgesep: 20 })
			.setDefaultEdgeLabel(() => ({}));

		for (const node of nodes) {
			g2.setNode(node.id, { width: node.width, height: node.height });
		}
		for (const link of links) {
			g2.setEdge(link.source, link.target);
		}

		dagreLayout(g2);

		for (const node of nodes) {
			const pos = g2.node(node.id);
			node.x = pos.x - node.width / 2;
			node.y = pos.y - node.height / 2;
		}

		for (const link of links) {
			const edgeData = g2.edge(link.source, link.target);
			if (edgeData?.points) link.points = edgeData.points;
		}

		const svg = d3
			.select(svgEl)
			.attr("viewBox", `0 0 ${width} ${height}`)
			.attr("preserveAspectRatio", "xMidYMid meet");

		const defs = svg.append("defs");
		const uniqueColors = new Set(links.map((l) => edgeColor(l.kind, l.label)));
		for (const color of uniqueColors) {
			appendArrowMarker(defs, `arrow-${color.replace("#", "")}`, color);
		}

		const rootG = svg.append("g");
		const zoom = this.setupZoom(svg, rootG);

		const lineGen = d3
			.line<{ x: number; y: number }>()
			.x((d) => d.x)
			.y((d) => d.y)
			.curve(d3.curveBasis);

		const buildPathPoints = (d: SimLink): Array<{ x: number; y: number }> => {
			if (d.points && d.points.length > 0) return d.points;
			// biome-ignore lint/style/noNonNullAssertion: nodes guaranteed to exist in the graph
			const src = nodeMap.get(d.source)!;
			// biome-ignore lint/style/noNonNullAssertion: nodes guaranteed to exist in the graph
			const tgt = nodeMap.get(d.target)!;
			const sx = src.x + src.width / 2;
			const sy = src.y + src.height / 2;
			const tx = tgt.x + tgt.width / 2;
			const ty = tgt.y + tgt.height / 2;
			return [
				{ x: sx, y: sy },
				{ x: (sx + tx) / 2, y: (sy + ty) / 2 },
				{ x: tx, y: ty },
			];
		};

		const linkGroup = rootG
			.append("g")
			.attr("class", "links")
			.selectAll("path")
			.data(links)
			.join("path")
			.attr("d", (d) => lineGen(buildPathPoints(d)))
			.attr("stroke", (d) => edgeColor(d.kind, d.label))
			.attr("stroke-width", 1.5)
			.attr("stroke-dasharray", (d) => EDGE_DASH[d.kind])
			.attr(
				"marker-end",
				(d) => `url(#arrow-${edgeColor(d.kind, d.label).replace("#", "")})`,
			)
			.attr("fill", "none")
			.attr("opacity", 0.7);

		const edgeLabelGroup = rootG
			.append("g")
			.attr("class", "edge-labels")
			.selectAll("text")
			.data(links.filter((l) => l.label))
			.join("text")
			.text((d) => d.label ?? "")
			.attr("font-size", 9)
			.attr("fill", "var(--graph-text-sub)")
			.attr("text-anchor", "middle")
			.attr("dy", -4)
			.attr("x", (d) => {
				const pts = buildPathPoints(d);
				return pts[Math.floor(pts.length / 2)].x;
			})
			.attr("y", (d) => {
				const pts = buildPathPoints(d);
				return pts[Math.floor(pts.length / 2)].y;
			});

		const nodeGroup = rootG
			.append("g")
			.attr("class", "nodes")
			.selectAll<SVGGElement, StructuredNode>("g")
			.data(nodes)
			.join("g")
			.attr("transform", (d) => `translate(${d.x},${d.y})`)
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
			.append("line")
			.attr("x1", 0)
			.attr("y1", 24)
			.attr("x2", (d) => d.width)
			.attr("y2", 24)
			.attr("stroke", SCHEMA_STROKE)
			.attr("stroke-width", 1);

		schemaNodes
			.append("text")
			.text((d) => d.label)
			.attr("x", (d) => d.width / 2)
			.attr("y", 15)
			.attr("text-anchor", "middle")
			.attr("dominant-baseline", "central")
			.attr("fill", "var(--graph-text)")
			.attr("font-size", 12)
			.attr("font-weight", 700);

		schemaNodes
			.filter((d) => d.label !== d.fullLabel)
			.append("title")
			.text((d) => d.fullLabel);

		schemaNodes
			.append("text")
			.text((d) => d.sublabel)
			.attr("x", (d) => d.width / 2)
			.attr("y", 35)
			.attr("text-anchor", "middle")
			.attr("dominant-baseline", "central")
			.attr("fill", "var(--graph-text-sub)")
			.attr("font-size", 10);

		const drag = d3
			.drag<SVGGElement, StructuredNode>()
			.on("drag", function (event, d) {
				d.x += event.dx;
				d.y += event.dy;
				d3.select(this).attr("transform", `translate(${d.x},${d.y})`);

				linkGroup
					.filter((l) => l.source === d.id || l.target === d.id)
					.each((l) => {
						l.points = undefined;
					})
					.attr("d", (l) => lineGen(buildPathPoints(l)));

				edgeLabelGroup
					.filter((l) => l.source === d.id || l.target === d.id)
					.attr("x", (l) => {
						const pts = buildPathPoints(l);
						return pts[Math.floor(pts.length / 2)].x;
					})
					.attr("y", (l) => {
						const pts = buildPathPoints(l);
						return pts[Math.floor(pts.length / 2)].y;
					});
			});

		nodeGroup.call(drag);

		const graphInfo = g2.graph();
		const gWidth = graphInfo.width ?? width;
		const gHeight = graphInfo.height ?? height;
		const padding = 40;
		const scale = Math.min(
			(width - padding * 2) / gWidth,
			(height - padding * 2) / gHeight,
			1,
		);
		const translateX = (width - gWidth * scale) / 2;
		const translateY = (height - gHeight * scale) / 2;
		svg.call(
			zoom.transform,
			d3.zoomIdentity.translate(translateX, translateY).scale(scale),
		);
	}
}
