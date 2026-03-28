import { Component } from "@angular/core";
import { RouterLink, RouterOutlet } from "@angular/router";

@Component({
	selector: "app-root",
	imports: [RouterLink, RouterOutlet],
	template: `
		<nav class="app-header">
			<a routerLink="/" class="brand">
				<img src="assets/icons/restatlas.svg" alt="" class="brand-icon" />
				<span class="brand-name">RestAtlas</span>
			</a>
		</nav>
		<router-outlet />
	`,
	styles: `
		.app-header {
			display: flex;
			align-items: center;
			padding: 0.5rem 1rem;
			border-bottom: 1px solid #e5e7eb;
			background: #fff;
			font-family: system-ui, sans-serif;
		}
		.brand {
			display: flex;
			align-items: center;
			gap: 0.5rem;
			text-decoration: none;
			color: #111;
		}
		.brand-icon {
			width: 24px;
			height: 24px;
		}
		.brand-name {
			font-weight: 700;
			font-size: 1.1rem;
		}
	`,
})
export class App {}
