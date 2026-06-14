import { DOCUMENT } from "@angular/common";
import {
	ChangeDetectionStrategy,
	Component,
	type ElementRef,
	effect,
	inject,
	input,
	model,
	output,
	viewChild,
} from "@angular/core";

/**
 * Confirmation modal shown before the proxy is allowed to forward requests for a
 * spec. Pure presentation: the parent supplies the host list and reacts to
 * `(confirm)` by calling the approval API. Shared by the Explorer's Try-It-Out
 * tab and the standalone API Client.
 *
 * The overlay is teleported to <body> while open so it always covers the full
 * viewport, even when an ancestor establishes a containing block for fixed
 * positioning (e.g. the spec-viewer's `backdrop-blur` panel).
 */
@Component({
	selector: "app-approval-dialog",
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
		@if (open()) {
			<div #overlay class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" (click)="cancel()">
				<div class="bg-app-bg border border-app-border rounded-lg shadow-xl max-w-md w-full p-5" (click)="$event.stopPropagation()">
					<div class="flex items-start gap-3 mb-3">
						<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-amber-600 dark:text-amber-400 shrink-0 mt-0.5"><title>Warning</title><path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
						<h3 class="m-0 text-base font-semibold text-app-text">Confirm API approval</h3>
					</div>

					<p class="text-xs text-app-text-muted m-0 mb-3">By approving, you allow the ReqViz backend to forward HTTP requests on your behalf to:</p>

					<ul class="m-0 mb-4 pl-4 list-disc text-xs text-app-text">
						@for (h of hosts(); track h) {
							<li><code class="font-mono">{{ h }}</code></li>
						}
					</ul>

					<div class="bg-app-surface border border-app-border rounded p-3 mb-4">
						<h4 class="m-0 mb-2 text-[0.7rem] font-semibold text-app-text-muted uppercase tracking-wide">Before you approve</h4>
						<ul class="m-0 pl-4 list-disc text-[0.7rem] text-app-text-muted space-y-1.5">
							<li><strong class="text-app-text">Requests are real.</strong> Each call hits the live API and counts against your rate limits and quota. On paid APIs (OpenAI, Stripe, Twilio, …) a request can cost real money on whatever key you provide.</li>
							<li><strong class="text-app-text">Your credentials pass through ReqViz.</strong> Headers, tokens, and API keys you enter are sent from the ReqViz server to the target — not directly from your browser — and may appear in the target's logs. Prefer test or restricted keys.</li>
							<li><strong class="text-app-text">Only the listed hosts are reachable.</strong> Approval lets ReqViz call the hosts shown above and nothing else. It applies to this spec only and is forgotten when the server restarts.</li>
							<li><strong class="text-app-text">Trust the spec's source.</strong> A spec you don't trust can point requests at endpoints you never intended to call. Only approve specs from a source you trust.</li>
						</ul>
					</div>

					<div class="flex justify-end gap-2">
						<button type="button"
							class="py-1.5 px-3 text-xs font-medium bg-app-surface hover:bg-app-surface-hover text-app-text border border-app-border rounded cursor-pointer transition-colors"
							(click)="cancel()"
						>Cancel</button>
						<button type="button"
							class="py-1.5 px-3 text-xs font-medium bg-amber-600 hover:bg-amber-700 text-white rounded border-none cursor-pointer transition-colors"
							(click)="approve()"
						>I understand, approve</button>
					</div>
				</div>
			</div>
		}
	`,
})
export class ApprovalDialogComponent {
	readonly open = model(false);
	readonly hosts = input<string[]>([]);
	readonly confirm = output<void>();

	private readonly doc = inject(DOCUMENT);
	private readonly overlay = viewChild<ElementRef<HTMLElement>>("overlay");

	constructor() {
		// Move the overlay to <body> once rendered so `fixed` resolves against the
		// viewport rather than a blurred/clipped ancestor. Angular still owns the
		// node, so bindings and teardown keep working from its new location.
		effect(() => {
			const el = this.overlay()?.nativeElement;
			if (el && el.parentElement !== this.doc.body) {
				this.doc.body.appendChild(el);
			}
		});
	}

	cancel(): void {
		this.open.set(false);
	}

	approve(): void {
		this.open.set(false);
		this.confirm.emit();
	}
}
