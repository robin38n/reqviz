import {
	ChangeDetectionStrategy,
	Component,
	computed,
	input,
} from "@angular/core";

const METHOD_BG: Record<string, string> = {
	GET: "bg-green-600",
	POST: "bg-blue-600",
	PUT: "bg-amber-600",
	PATCH: "bg-purple-600",
	DELETE: "bg-red-600",
};

@Component({
	selector: "app-method-badge",
	standalone: true,
	changeDetection: ChangeDetectionStrategy.OnPush,
	template: `
    <span 
      class="inline-flex items-center justify-center font-bold text-white dark:text-zinc-950 rounded-sm leading-none shrink-0"
      [class]="classes()"
    >
      {{ method() }}
    </span>
  `,
})
export class MethodBadgeComponent {
	readonly method = input.required<string>();
	readonly size = input<"xs" | "sm" | "lg">("sm");

	protected readonly classes = computed(() => {
		const bg = METHOD_BG[this.method()] ?? "bg-zinc-500";
		const size = this.size();

		let sizeClasses = "text-[0.65rem] px-1 py-0.5 min-w-[32px]";
		if (size === "xs") sizeClasses = "text-[0.55rem] px-0.5 py-0 min-w-[24px]";
		if (size === "lg") sizeClasses = "text-[0.8rem] px-2 py-1 min-w-[48px]";

		return `${bg} ${sizeClasses}`;
	});
}
