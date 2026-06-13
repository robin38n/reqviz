import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
	signal,
} from "@angular/core";
import { Router } from "@angular/router";
import type { EndpointNode } from "../../models/graph.model";
import { AlertComponent } from "../../shared/components/alert/alert.component";
import { ApprovalDialogComponent } from "../../shared/components/approval-dialog/approval-dialog.component";
import { MethodBadgeComponent } from "../../shared/components/method-badge/method-badge.component";
import { RequestHistoryComponent } from "../../shared/components/request-history/request-history.component";
import { ResponseViewerComponent } from "../../shared/components/response-viewer/response-viewer.component";
import { ButtonDirective } from "../../shared/directives/button.directive";
import { SpecGraphService } from "../spec-viewer/services/spec-graph.service";
import {
	type HistoryEntry,
	type ProxyRequest,
	TryItOutService,
} from "../spec-viewer/services/try-it-out.service";
import { ApiClientStateService } from "./api-client-state.service";

const METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE"] as const;

@Component({
	selector: "app-api-client",
	imports: [
		ResponseViewerComponent,
		RequestHistoryComponent,
		MethodBadgeComponent,
		AlertComponent,
		ButtonDirective,
		ApprovalDialogComponent,
	],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./api-client.component.html",
})
export class ApiClientComponent {
	private readonly router = inject(Router);
	protected readonly tryItOut = inject(TryItOutService);
	protected readonly specGraph = inject(SpecGraphService);
	protected readonly state = inject(ApiClientStateService);

	readonly methods = METHODS;
	// Form state lives in ApiClientStateService so it survives navigation.
	readonly method = this.state.method;
	readonly url = this.state.url;
	readonly headers = this.state.headers;
	readonly body = this.state.body;
	readonly urlError = signal("");
	readonly showApprovalDialog = signal(false);

	readonly showBody = computed(() => {
		const m = this.method();
		return m === "POST" || m === "PUT" || m === "PATCH";
	});

	/** Fill the request form from a spec endpoint (method + full URL). */
	loadEndpoint(ep: EndpointNode): void {
		this.method.set(ep.method);
		this.url.set(`${this.specGraph.serverBaseUrl()}${ep.path}`);
		this.urlError.set("");
	}

	constructor() {
		const state = this.router.getCurrentNavigation()?.extras?.state;
		if (state) {
			if (state.method) this.method.set(String(state.method));
			if (state.url) this.url.set(String(state.url));
			if (state.body) this.body.set(String(state.body));
			if (Array.isArray(state.headers)) {
				this.headers.set(state.headers);
			}
		}
	}

	inputValue(event: Event): string {
		return (event.target as HTMLInputElement).value;
	}

	asMethod(event: Event): string {
		return (event.target as HTMLSelectElement).value;
	}

	addHeader(): void {
		this.headers.update((h) => [...h, { key: "", value: "" }]);
	}

	removeHeader(index: number): void {
		this.headers.update((h) => h.filter((_, i) => i !== index));
	}

	updateHeader(index: number, field: "key" | "value", event: Event): void {
		const val = (event.target as HTMLInputElement).value;
		this.headers.update((h) =>
			h.map((item, i) => (i === index ? { ...item, [field]: val } : item)),
		);
	}

	async sendRequest(): Promise<void> {
		const reqUrl = this.url().trim();
		if (!reqUrl) return;

		try {
			const parsed = new URL(reqUrl);
			if (parsed.protocol !== "http:" && parsed.protocol !== "https:") {
				this.urlError.set("Only http:// and https:// URLs are allowed.");
				return;
			}
		} catch {
			this.urlError.set("Enter a valid URL starting with http:// or https://.");
			return;
		}
		this.urlError.set("");

		// A loaded spec must be approved before the proxy will forward requests.
		// Prompt for approval instead of letting the backend reject it.
		if (this.specGraph.specId() && !this.specGraph.approved()) {
			this.showApprovalDialog.set(true);
			return;
		}

		const headers: Record<string, string> = {};
		for (const h of this.headers()) {
			if (h.key.trim()) {
				headers[h.key.trim()] = h.value;
			}
		}

		const method = this.method() as ProxyRequest["method"];
		const hasBody = ["POST", "PUT", "PATCH"].includes(method);

		if (hasBody && !headers["Content-Type"]) {
			headers["Content-Type"] = "application/json";
		}

		let bodyPayload: unknown;
		if (hasBody && this.body().trim()) {
			try {
				bodyPayload = JSON.parse(this.body());
			} catch {
				bodyPayload = this.body();
			}
		}

		await this.tryItOut.sendRequest({
			method,
			url: reqUrl,
			headers: Object.keys(headers).length > 0 ? headers : undefined,
			body: hasBody && bodyPayload != null ? bodyPayload : undefined,
			specId: this.specGraph.specId() ?? undefined,
		});
	}

	/** Approve the loaded spec, then retry the send the user originally clicked. */
	async onApprovalConfirmed(): Promise<void> {
		await this.specGraph.approve();
		if (this.specGraph.approved()) {
			await this.sendRequest();
		}
	}

	onReplayRequest(entry: HistoryEntry): void {
		this.method.set(entry.request.method);
		this.url.set(entry.request.url);

		if (entry.request.headers) {
			this.headers.set(
				Object.entries(entry.request.headers).map(([key, value]) => ({
					key,
					value,
				})),
			);
		} else {
			this.headers.set([]);
		}

		if (entry.request.body != null) {
			this.body.set(
				typeof entry.request.body === "string"
					? entry.request.body
					: JSON.stringify(entry.request.body, null, 2),
			);
		} else {
			this.body.set("");
		}

		this.tryItOut.lastResponse.set(entry.response);
	}
}
