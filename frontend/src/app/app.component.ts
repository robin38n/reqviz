import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
} from "@angular/core";
import { RouterLink, RouterLinkActive, RouterOutlet } from "@angular/router";
import { ThemeService } from "./core/theme.service";
import { SpecTabsService } from "./features/spec-viewer/services/spec-tabs.service";
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
	private readonly tabs = inject(SpecTabsService);

	// Explorer points at the active spec tab; falls back to the start page.
	protected readonly explorerLink = computed(() => {
		const id = this.tabs.activeId();
		return id ? ["/specs", id] : ["/"];
	});
}
