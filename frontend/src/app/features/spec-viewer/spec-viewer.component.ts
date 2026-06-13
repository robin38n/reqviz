import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
	type OnInit,
} from "@angular/core";
import { ActivatedRoute, RouterLink } from "@angular/router";
import type {
	EndpointNode,
	GraphNode,
	SchemaNode,
} from "../../models/graph.model";
import { ListToolbarComponent } from "../../shared/components/list-toolbar/list-toolbar.component";
import { MethodBadgeComponent } from "../../shared/components/method-badge/method-badge.component";
import { ButtonDirective } from "../../shared/directives/button.directive";
import { GraphCanvasComponent } from "./graph/graph-canvas.component";
import { GraphCanvasForceComponent } from "./graph/graph-canvas-force.component";
import { GraphToolbarComponent } from "./graph/graph-toolbar.component";
import { NodeDetailComponent } from "./node-detail/node-detail.component";
import { SpecGraphService } from "./services/spec-graph.service";

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
		ButtonDirective,
	],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./spec-viewer.component.html",
})
export class SpecViewerComponent implements OnInit {
	protected readonly svc = inject(SpecGraphService);
	private readonly route = inject(ActivatedRoute);

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

	ngOnInit() {
		const id = this.route.snapshot.paramMap.get("id");
		if (id) {
			this.svc.loadSpec(id);
		}
	}

	onNodeClick(node: GraphNode): void {
		this.svc.selectNode(node);
	}
}
