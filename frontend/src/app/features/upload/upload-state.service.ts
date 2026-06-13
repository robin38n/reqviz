import { Injectable, signal } from "@angular/core";
import type { components } from "../../core/schema";

type SpecSummaryRaw = components["schemas"]["SpecSummary"];
export type SpecSummary = SpecSummaryRaw & {
	endpoints?: number;
	schemas?: number;
};

/**
 * Holds the landing-page upload draft (pasted/loaded spec text and last result)
 * so it survives navigation between pages within a session. Root-scoped.
 */
@Injectable({ providedIn: "root" })
export class UploadStateService {
	readonly draft = signal("");
	readonly selectedDemoSlug = signal("");
	readonly summary = signal<SpecSummary | null>(null);
}
