import { computed, Directive, input } from "@angular/core";

export type ButtonVariant = "primary" | "secondary";

// Shared identity for action buttons; call sites add only size/layout classes.
const BASE =
	"inline-flex items-center justify-center gap-2 rounded-md font-semibold no-underline cursor-pointer transition-colors disabled:opacity-50 disabled:cursor-not-allowed";

const VARIANTS: Record<ButtonVariant, string> = {
	primary:
		"bg-blue-600 text-white border border-blue-500 shadow-sm hover:bg-blue-500",
	secondary:
		"bg-app-surface text-app-text-muted border border-app-border hover:text-app-text hover:bg-app-surface-hover hover:border-app-text-muted/40",
};

/**
 * Applies the shared button styling to any <button> or <a>. Use `appButton`
 * for the primary (blue) look or `appButton="secondary"` for the outline look.
 * Static classes on the element (padding, width, etc.) are merged by Angular.
 */
@Directive({
	selector: "[appButton]",
	host: { "[class]": "hostClasses()" },
})
export class ButtonDirective {
	readonly appButton = input<ButtonVariant | "">("");

	protected readonly hostClasses = computed(
		() => `${BASE} ${VARIANTS[this.appButton() || "primary"]}`,
	);
}
