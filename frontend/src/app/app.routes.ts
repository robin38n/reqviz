import type { Routes } from "@angular/router";
import { UploadComponent } from "./features/upload/upload";

export const routes: Routes = [
	{ path: "", component: UploadComponent },
	{
		path: "specs/:id",
		loadComponent: () =>
			import("./features/spec-viewer/spec-viewer").then(
				(m) => m.SpecViewerComponent,
			),
	},
	{ path: "**", redirectTo: "" },
];
