import { Injectable, signal } from "@angular/core";

export interface HeaderRow {
	key: string;
	value: string;
}

/**
 * Holds the API Client's in-progress request form so it survives navigation
 * between pages within a session. Root-scoped (singleton).
 */
@Injectable({ providedIn: "root" })
export class ApiClientStateService {
	readonly method = signal<string>("GET");
	readonly url = signal("");
	readonly headers = signal<HeaderRow[]>([]);
	readonly body = signal("");
}
