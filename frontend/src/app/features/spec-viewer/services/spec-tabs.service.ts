import { computed, Injectable, inject, signal } from "@angular/core";
import { ApiService } from "../../../core/api.service";
import { SpecTab } from "./spec-tab";

/**
 * Manages the set of open spec tabs. Each tab is a {@link SpecTab} holding its
 * own state, so switching tabs preserves per-spec selection/filters/layout.
 * Tabs are session-only (in memory), matching the in-memory backend store.
 */
@Injectable({ providedIn: "root" })
export class SpecTabsService {
	private readonly api = inject(ApiService);

	readonly MAX_TABS = 5;

	readonly tabs = signal<SpecTab[]>([]);
	readonly activeId = signal<string | null>(null);
	readonly openError = signal<string | null>(null);

	readonly active = computed(
		() => this.tabs().find((t) => t.id === this.activeId()) ?? null,
	);
	readonly atMax = computed(() => this.tabs().length >= this.MAX_TABS);

	/**
	 * Open a spec into a tab. Activates it if already open; otherwise creates and
	 * loads a new tab. Returns false (and sets {@link openError}) when the tab cap
	 * would be exceeded.
	 */
	open(id: string): boolean {
		if (this.tabs().some((t) => t.id === id)) {
			this.activate(id);
			return true;
		}
		if (this.atMax()) {
			this.openError.set(
				`You can open up to ${this.MAX_TABS} specs. Close a tab to open another.`,
			);
			return false;
		}
		const tab = new SpecTab(this.api, id);
		this.tabs.update((ts) => [...ts, tab]);
		this.activeId.set(id);
		this.openError.set(null);
		tab.load();
		return true;
	}

	/** Close a tab. If it was active, activates a neighbor. Returns the new active id. */
	close(id: string): string | null {
		const tabs = this.tabs();
		const index = tabs.findIndex((t) => t.id === id);
		if (index === -1) return this.activeId();

		const next = tabs.filter((t) => t.id !== id);
		this.tabs.set(next);

		if (this.activeId() === id) {
			const neighbor = next[index] ?? next[index - 1] ?? null;
			this.activeId.set(neighbor?.id ?? null);
		}
		this.openError.set(null);
		return this.activeId();
	}

	activate(id: string): void {
		this.activeId.set(id);
		this.openError.set(null);
	}
}
