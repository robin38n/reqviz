import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
} from "@angular/core";
import { RouterLink, RouterLinkActive, RouterOutlet } from "@angular/router";
import { ThemeService } from "./core/theme.service";
import { SpecGraphService } from "./features/spec-viewer/services/spec-graph.service";
import { ReactiveBackgroundComponent } from "./shared/components/reactive-background/reactive-background.component";

@Component({
	selector: "app-root",
	imports: [
		RouterLink,
		RouterLinkActive,
		RouterOutlet,
		ReactiveBackgroundComponent,
	],
	changeDetection: ChangeDetectionStrategy.OnPush,
	templateUrl: "./app.component.html",
})
export class AppComponent {
	protected readonly theme = inject(ThemeService);
	private readonly specGraph = inject(SpecGraphService);

	// Explorer points at the loaded spec's graph; falls back to the start page.
	protected readonly explorerLink = computed(() => {
		const id = this.specGraph.specId();
		return id ? ["/specs", id] : ["/"];
	});
}
