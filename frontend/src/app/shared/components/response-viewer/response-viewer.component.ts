import { KeyValuePipe } from "@angular/common";
import {
	ChangeDetectionStrategy,
	Component,
	computed,
	input,
	signal,
} from "@angular/core";
import type { ProxyResponse } from "../../../features/spec-viewer/services/try-it-out.service";
import { StatusBadgeComponent } from "../status-badge/status-badge.component";

const PREVIEW_LINES = 12;

@Component({
	selector: "app-response-viewer",
	imports: [StatusBadgeComponent, KeyValuePipe],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./response-viewer.component.html",
})
export class ResponseViewerComponent {
	readonly response = input.required<ProxyResponse>();
	readonly showHeaders = signal(false);
	readonly expanded = signal(false);

	readonly formattedBody = computed(() => {
		const body = this.response().body;
		if (body == null) return "";
		if (typeof body === "string") return body;
		return JSON.stringify(body, null, 2);
	});

	private readonly bodyLines = computed(() => this.formattedBody().split("\n"));

	readonly isTruncatable = computed(
		() => this.bodyLines().length > PREVIEW_LINES,
	);

	readonly hiddenLineCount = computed(() =>
		Math.max(0, this.bodyLines().length - PREVIEW_LINES),
	);

	/** Body to render: full when expanded or short enough, else a preview. */
	readonly displayBody = computed(() => {
		if (this.expanded() || !this.isTruncatable()) return this.formattedBody();
		return this.bodyLines().slice(0, PREVIEW_LINES).join("\n");
	});
}
