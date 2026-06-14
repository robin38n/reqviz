import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
} from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { ActivatedRoute, Router, RouterLink } from "@angular/router";
import type {
	EndpointNode,
	GraphNode,
	SchemaNode,
} from "../../models/graph.model";
import { ListToolbarComponent } from "../../shared/components/list-toolbar/list-toolbar.component";
import { MethodBadgeComponent } from "../../shared/components/method-badge/method-badge.component";
import { GraphCanvasComponent } from "./graph/graph-canvas.component";
import { GraphCanvasForceComponent } from "./graph/graph-canvas-force.component";
import { GraphToolbarComponent } from "./graph/graph-toolbar.component";
import { NodeDetailComponent } from "./node-detail/node-detail.component";
import { SpecGraphService } from "./services/spec-graph.service";
import { SpecTabsService } from "./services/spec-tabs.service";

@Component({
	selector: "app-spec-viewer",
	imports: [
		RouterLink,
		GraphCanvasComponent,
		GraphCanvasForceComponent,
		GraphToolbarComponent,
		MethodBadgeComponent,
		NodeDetailComponent,
		ListToolbarComponent,
	],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./spec-viewer.component.html",
})
export class SpecViewerComponent {
	protected readonly svc = inject(SpecGraphService);
	protected readonly tabs = inject(SpecTabsService);
	private readonly route = inject(ActivatedRoute);
	private readonly router = inject(Router);

	constructor() {
		// The component is reused across `/specs/:id` changes, so react to every
		// param emission rather than reading the snapshot once.
		this.route.paramMap.pipe(takeUntilDestroyed()).subscribe((params) => {
			const id = params.get("id");
			if (!id) return;
			if (!this.tabs.open(id)) {
				// Tab cap hit: revert the URL to the still-active tab.
				const active = this.tabs.activeId();
				if (active && active !== id) {
					this.router.navigate(["/specs", active]);
				}
			}
		});
	}

	switchTab(id: string): void {
		if (id !== this.tabs.activeId()) {
			this.router.navigate(["/specs", id]);
		}
	}

	closeTab(id: string, event: Event): void {
		event.stopPropagation();
		const wasActive = id === this.tabs.activeId();
		const nextActive = this.tabs.close(id);
		if (wasActive) {
			this.router.navigate(nextActive ? ["/specs", nextActive] : ["/"]);
		}
	}

	protected readonly displayGraph = computed(
		() => this.svc.filteredGraph() ?? this.svc.graph(),
	);

	protected readonly filteredEndpoints = computed(
		() =>
			this.displayGraph()?.nodes.filter(
				(n): n is EndpointNode => n.type === "endpoint",
			) ?? [],
	);

	protected readonly filteredSchemas = computed(
		() =>
			this.displayGraph()?.nodes.filter(
				(n): n is SchemaNode => n.type === "schema",
			) ?? [],
	);

	protected readonly localEndpoints = computed(() => {
		let items = this.filteredEndpoints();
		const query = this.svc.listSearch().toLowerCase().trim();
		const sort = this.svc.listSort();

		if (query) {
			items = items.filter(
				(ep) =>
					ep.path.toLowerCase().includes(query) ||
					ep.method.toLowerCase().includes(query),
			);
		}

		return [...items].sort((a, b) => {
			if (sort === "az") return a.path.localeCompare(b.path);
			if (sort === "za") return b.path.localeCompare(a.path);
			if (sort === "method") return a.method.localeCompare(b.method);
			return 0;
		});
	});

	protected readonly localSchemas = computed(() => {
		let items = this.filteredSchemas();
		const query = this.svc.listSearch().toLowerCase().trim();
		const sort = this.svc.listSort();

		if (query) {
			items = items.filter((sc) => sc.name.toLowerCase().includes(query));
		}

		return [...items].sort((a, b) => {
			if (sort === "az") return a.name.localeCompare(b.name);
			if (sort === "za") return b.name.localeCompare(a.name);
			return 0;
		});
	});

	onNodeClick(node: GraphNode): void {
		this.svc.selectNode(node);
	}
}
