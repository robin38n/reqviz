import { DecimalPipe } from "@angular/common";
import {
	ChangeDetectionStrategy,
	Component,
	inject,
	signal,
} from "@angular/core";
import { RouterLink } from "@angular/router";
import { ApiService } from "../../core/api.service";
import type { components } from "../../core/schema";
import { ThemeService } from "../../core/theme.service";
import { AlertComponent } from "../../shared/components/alert/alert.component";
import { ButtonDirective } from "../../shared/directives/button.directive";
import { type SpecSummary, UploadStateService } from "./upload-state.service";

type DemoInfo = components["schemas"]["DemoInfo"];

@Component({
	selector: "app-upload",
	imports: [RouterLink, DecimalPipe, AlertComponent, ButtonDirective],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./upload.component.html",
})
export class UploadComponent {
	private readonly api = inject(ApiService);
	protected readonly theme = inject(ThemeService);
	protected readonly state = inject(UploadStateService);

	// Draft, demo selection, and result persist across navigation via the service.
	readonly draft = this.state.draft;
	readonly summary = this.state.summary;
	readonly selectedDemoSlug = this.state.selectedDemoSlug;

	readonly loading = signal(false);
	readonly error = signal<string | null>(null);
	readonly demos = signal<DemoInfo[]>([]);

	constructor() {
		this.api.listDemos().then(({ data }) => {
			if (data && data.length > 0) {
				this.demos.set(data);
			}
		});
	}

	async loadDemo(slug: string) {
		this.selectedDemoSlug.set(slug);
		try {
			const { data } = await this.api.getDemoSpec(slug);
			if (data) {
				this.draft.set(JSON.stringify(data, null, 2));
			}
		} catch {
			this.error.set("Could not load the demo. Please try again.");
		}
	}

	async onFileSelected(event: Event) {
		const input = event.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		this.draft.set(await file.text());
	}

	async visualize() {
		this.error.set(null);
		const text = this.draft().trim();
		if (!text) return;

		this.loading.set(true);
		try {
			const { data, error } = await this.api.uploadSpecRaw(
				text,
				text.startsWith("{") ? "application/json" : "application/x-yaml",
			);
			if (error) {
				this.error.set(error.error || "Upload failed. Please try again.");
			} else if (data) {
				const s = data as SpecSummary;
				s.endpoints = data.endpointCount;
				s.schemas = data.schemaCount;
				this.summary.set(s);
			}
		} catch {
			this.error.set("Couldn't reach the server.");
		} finally {
			this.loading.set(false);
		}
	}
}
