import { computed, Injectable, inject } from "@angular/core";
import { ApiService } from "../../../core/api.service";
import type { GraphNode } from "../../../models/graph.model";
import { SpecTab } from "./spec-tab";
import { SpecTabsService } from "./spec-tabs.service";

/**
 * Thin facade over the currently active {@link SpecTab}. The spec-viewer renders
 * only one tab at a time, so its child components (graph toolbar, node/endpoint/
 * schema detail, try-it-out) keep injecting this service and transparently read
 * and write the active tab's state. Reading `active()` inside each getter makes
 * every binding reactive to tab switches.
 *
 * Components that need *all* open specs (e.g. the api-client) use
 * {@link SpecTabsService} directly instead.
 */
@Injectable({ providedIn: "root" })
export class SpecGraphService {
	private readonly tabs = inject(SpecTabsService);
	private readonly api = inject(ApiService);

	// Fallback so getters stay null-safe when no tab is active.
	private readonly empty = new SpecTab(this.api, "");

	private get t(): SpecTab {
		return this.tabs.active() ?? this.empty;
	}

	readonly specId = computed(() => this.tabs.active()?.id ?? null);

	get loading() {
		return this.t.loading;
	}
	get error() {
		return this.t.error;
	}
	get summary() {
		return this.t.summary;
	}
	get graph() {
		return this.t.graph;
	}
	get rawSpec() {
		return this.t.rawSpec;
	}
	get selectedNodeId() {
		return this.t.selectedNodeId;
	}
	get searchQuery() {
		return this.t.searchQuery;
	}
	get selectedTags() {
		return this.t.selectedTags;
	}
	get selectedMethods() {
		return this.t.selectedMethods;
	}
	get layout() {
		return this.t.layout;
	}
	get listSearch() {
		return this.t.listSearch;
	}
	get listSort() {
		return this.t.listSort;
	}

	get serverBaseUrl() {
		return this.t.serverBaseUrl;
	}
	get approved() {
		return this.t.approved;
	}
	get allowedHosts() {
		return this.t.allowedHosts;
	}
	get endpointNodes() {
		return this.t.endpointNodes;
	}
	get schemaNodes() {
		return this.t.schemaNodes;
	}
	get edgeCount() {
		return this.t.edgeCount;
	}
	get allTags() {
		return this.t.allTags;
	}
	get filteredGraph() {
		return this.t.filteredGraph;
	}
	get selectedNode() {
		return this.t.selectedNode;
	}
	get selectedNodeEdges() {
		return this.t.selectedNodeEdges;
	}

	selectNode(node: GraphNode): void {
		this.t.selectNode(node);
	}
	clearSelection(): void {
		this.t.clearSelection();
	}
	toggleTag(tag: string): void {
		this.t.toggleTag(tag);
	}
	toggleMethod(method: string): void {
		this.t.toggleMethod(method);
	}
	clearFilters(): void {
		this.t.clearFilters();
	}
	approve(hosts?: string[]): Promise<void> {
		return this.t.approve(hosts);
	}
}
