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
	// Id of the open spec the loaded endpoint came from, used to scope the proxy
	// request (specId) and the approval gate. Null for manually-typed URLs.
	readonly selectedSpecId = signal<string | null>(null);
}
