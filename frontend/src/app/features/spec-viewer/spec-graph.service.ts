import { computed, Injectable, inject, signal } from "@angular/core";
import { ApiService } from "../../api/api.service";
import type { components } from "../../api/schema";
import type {
	EndpointNode,
	SchemaNode,
	SpecGraph,
} from "../../models/graph.model";
import { buildSpecGraph } from "../../models/spec-to-graph";

type SpecSummary = components["schemas"]["SpecSummary"];

@Injectable({ providedIn: "root" })
export class SpecGraphService {
	private readonly api = inject(ApiService);

	readonly loading = signal(false);
	readonly error = signal<string | null>(null);
	readonly specId = signal<string | null>(null);
	readonly summary = signal<SpecSummary | null>(null);
	readonly graph = signal<SpecGraph | null>(null);

	readonly endpointNodes = computed(
		() =>
			this.graph()?.nodes.filter(
				(n): n is EndpointNode => n.type === "endpoint",
			) ?? [],
	);

	readonly schemaNodes = computed(
		() =>
			this.graph()?.nodes.filter((n): n is SchemaNode => n.type === "schema") ??
			[],
	);

	readonly edgeCount = computed(() => this.graph()?.edges.length ?? 0);

	async loadSpec(id: string): Promise<void> {
		this.loading.set(true);
		this.error.set(null);
		this.specId.set(id);

		try {
			const { data, error } = await this.api.getSpec(id);
			if (error) {
				this.error.set("Failed to load spec");
				return;
			}
			if (data) {
				this.summary.set(data.summary);
				this.graph.set(buildSpecGraph(data.raw));
			}
		} catch {
			this.error.set("Network error — is the backend running?");
		} finally {
			this.loading.set(false);
		}
	}
}
