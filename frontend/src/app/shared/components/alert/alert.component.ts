import {
	ChangeDetectionStrategy,
	Component,
	input,
	output,
} from "@angular/core";

/**
 * Red error/alert panel. Use `compact` for the smaller inline variant (icon +
 * message only) or the default for the full variant (icon, title, message, and
 * an optional dismiss button).
 */
@Component({
	selector: "app-alert",
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
		@if (compact()) {
			<div class="p-2.5 bg-app-surface/80 backdrop-blur-md border border-red-500/30 rounded-lg text-[0.8rem] flex items-center gap-2.5 shadow-sm animate-in fade-in slide-in-from-top-2 duration-300">
				<div class="flex items-center justify-center size-6 rounded-full bg-red-500/10 shrink-0">
					<svg aria-hidden="true" class="size-3.5 text-red-500 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<span class="text-app-text">{{ message() }}</span>
			</div>
		} @else {
			<div class="p-4 bg-app-surface/80 backdrop-blur-md border border-red-500/30 rounded-lg text-sm flex items-start gap-3 shadow-lg shadow-red-500/5 animate-in fade-in slide-in-from-top-2 duration-300">
				<div class="flex items-center justify-center size-8 rounded-full bg-red-500/10 shrink-0 mt-px">
					<svg aria-hidden="true" class="size-4 text-red-500 dark:text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</div>
				<div class="flex-1 min-w-0">
					<p class="text-xs font-semibold uppercase tracking-wider text-red-500 dark:text-red-400 mb-0.5">{{ title() }}</p>
					<p class="m-0 text-app-text leading-relaxed">{{ message() }}</p>
				</div>
				@if (dismissible()) {
					<button type="button" (click)="dismiss.emit()" aria-label="Dismiss"
						class="p-1 rounded-md bg-transparent border-none text-app-text-muted hover:text-app-text hover:bg-app-surface-hover cursor-pointer transition-colors shrink-0">
						<svg aria-hidden="true" class="size-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				}
			</div>
		}
	`,
})
export class AlertComponent {
	readonly title = input("Error");
	readonly message = input<string | null | undefined>(null);
	readonly dismissible = input(false);
	readonly compact = input(false);
	readonly dismiss = output<void>();
}
