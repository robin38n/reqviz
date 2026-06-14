import {
	afterNextRender,
	Directive,
	type ElementRef,
	effect,
	input,
	output,
	viewChild,
} from "@angular/core";
import * as d3 from "d3";
import type { GraphNode, SpecGraph } from "../../../models/graph.model";

/**
 * Shared scaffolding for the two graph canvases: the graph input/click output,
 * SVG/container refs, zoom state, the render trigger, and the zoom/fullscreen
 * controls. Subclasses implement `initGraph` with their own layout engine.
 */
@Directive()
export abstract class GraphCanvasBase {
	readonly graph = input.required<SpecGraph>();
	readonly nodeClick = output<GraphNode>();

	protected readonly svgRef =
		viewChild.required<ElementRef<SVGSVGElement>>("svg");
	protected readonly containerRef =
		viewChild.required<ElementRef<HTMLDivElement>>("container");

	protected zoom: d3.ZoomBehavior<SVGSVGElement, unknown> | null = null;
	protected svgSelection: d3.Selection<
		SVGSVGElement,
		unknown,
		null,
		undefined
	> | null = null;
	private initialized = false;

	constructor() {
		afterNextRender(() => {
			this.initialized = true;
			this.initGraph();
		});
		effect(() => {
			const g = this.graph();
			if (this.initialized && g) this.initGraph();
		});
	}

	protected abstract initGraph(): void;

	/** Wires pan/zoom on the SVG and records the handles the controls operate on. */
	protected setupZoom(
		svg: d3.Selection<SVGSVGElement, unknown, null, undefined>,
		rootG: d3.Selection<SVGGElement, unknown, null, undefined>,
	): d3.ZoomBehavior<SVGSVGElement, unknown> {
		const zoom = d3
			.zoom<SVGSVGElement, unknown>()
			.scaleExtent([0.2, 4])
			.on("zoom", (event) => {
				rootG.attr("transform", event.transform);
			});
		svg.call(zoom);
		this.zoom = zoom;
		this.svgSelection = svg;
		return zoom;
	}

	onZoomIn(): void {
		if (this.zoom && this.svgSelection) {
			this.svgSelection.transition().duration(300).call(this.zoom.scaleBy, 1.3);
		}
	}

	onZoomOut(): void {
		if (this.zoom && this.svgSelection) {
			this.svgSelection.transition().duration(300).call(this.zoom.scaleBy, 0.7);
		}
	}

	onResetZoom(): void {
		if (this.zoom && this.svgSelection) {
			this.svgSelection
				.transition()
				.duration(300)
				.call(this.zoom.transform, d3.zoomIdentity);
		}
	}

	onFullscreen(): void {
		const el = this.containerRef().nativeElement;
		if (document.fullscreenElement) {
			document.exitFullscreen();
		} else {
			el.requestFullscreen();
		}
	}
}
